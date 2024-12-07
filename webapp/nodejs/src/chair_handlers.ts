import type { Context } from "hono";
import { setCookie } from "hono/cookie";
import type { RowDataPacket } from "mysql2";
import { ulid } from "ulid";
import { getLatestRideStatus } from "./common.js";
import type { Environment } from "./types/hono.js";
import type {
  ChairLocation,
  Coordinate,
  Owner,
  Ride,
  RideStatus,
  User,
} from "./types/models.js";
import { secureRandomStr } from "./utils/random.js";

export const chairPostChairs = async (ctx: Context<Environment>) => {
  const reqJson = await ctx.req.json<{
    name: string;
    model: string;
    chair_register_token: string;
  }>();
  const { name, model, chair_register_token } = reqJson;
  if (!name || !model || !chair_register_token) {
    return ctx.text(
      "some of required fields(name, model, chair_register_token) are empty",
      400,
    );
  }
  const [[owner]] = await ctx.var.dbConn.query<Array<Owner & RowDataPacket>>(
    "SELECT * FROM owners WHERE chair_register_token = ?",
    [chair_register_token],
  );
  if (!owner) {
    return ctx.text("invalid chair_register_token", 401);
  }
  const chairID = ulid();
  const accessToken = secureRandomStr(32);
  await ctx.var.dbConn.query(
    "INSERT INTO chairs (id, owner_id, name, model, is_active, access_token) VALUES (?, ?, ?, ?, ?, ?)",
    [chairID, owner.id, name, model, false, accessToken],
  );

  setCookie(ctx, "chair_session", accessToken, { path: "/" });

  return ctx.json({ id: chairID, owner_id: owner.id }, 201);
};

export const chairPostActivity = async (ctx: Context<Environment>) => {
  const chair = ctx.var.chair;
  const reqJson = await ctx.req.json<{ is_active: boolean }>();
  try {
    await ctx.var.dbConn.query("UPDATE chairs SET is_active = ? WHERE id = ?", [
      reqJson.is_active,
      chair.id,
    ]);
  } catch (e) {
    return ctx.text(`${e}`, 500);
  }
  return ctx.body(null, 204);
};

export const chairPostCoordinate = async (ctx: Context<Environment>) => {
  const reqJson = await ctx.req.json<Coordinate>();
  const chair = ctx.var.chair;
  const chairLocationID = ulid();
  await ctx.var.dbConn.beginTransaction();
  try {
    await ctx.var.dbConn.query(
      "INSERT INTO chair_locations (id, chair_id, latitude, longitude) VALUES (?, ?, ?, ?)",
      [chairLocationID, chair.id, reqJson.latitude, reqJson.longitude],
    );
    const [[location]] = await ctx.var.dbConn.query<
      Array<ChairLocation & RowDataPacket>
    >("SELECT * FROM chair_locations WHERE id = ?", [chairLocationID]);
    const [[ride]] = await ctx.var.dbConn.query<Array<Ride & RowDataPacket>>(
      "SELECT * FROM rides WHERE chair_id = ? ORDER BY updated_at DESC LIMIT 1",
      [chair.id],
    );
    if (ride) {
      const status = await getLatestRideStatus(ctx.var.dbConn, ride.id);
      if (status !== "COMPLETED" && status !== "CANCELED") {
        if (
          reqJson.latitude === ride.pickup_latitude &&
          reqJson.longitude === ride.pickup_longitude &&
          status === "ENROUTE"
        ) {
          await ctx.var.dbConn.query(
            "INSERT INTO ride_statuses (id, ride_id, status) VALUES (?, ?, ?)",
            [ulid(), ride.id, "PICKUP"],
          );
        }
        if (
          reqJson.latitude === ride.destination_latitude &&
          reqJson.longitude === ride.destination_longitude &&
          status === "CARRYING"
        ) {
          await ctx.var.dbConn.query(
            "INSERT INTO ride_statuses (id, ride_id, status) VALUES (?, ?, ?)",
            [ulid(), ride.id, "ARRIVED"],
          );
        }
      }
    }
    await ctx.var.dbConn.commit();
    return ctx.json({ recorded_at: location.created_at.getTime() }, 200);
  } catch (e) {
    await ctx.var.dbConn.rollback();
    return ctx.text(`${e}`, 500);
  }
};

export const chairGetNotification = async (ctx: Context<Environment>) => {
  const chair = ctx.var.chair;

  await ctx.var.dbConn.beginTransaction();
  try {
    const [[ride]] = await ctx.var.dbConn.query<Array<Ride & RowDataPacket>>(
      "SELECT * FROM rides WHERE chair_id = ? ORDER BY updated_at DESC LIMIT 1",
      [chair.id],
    );
    if (!ride) {
      return ctx.json({ retry_after_ms: 30 }, 200);
    }

    const [[yetSentRideStatus]] = await ctx.var.dbConn.query<
      Array<RideStatus & RowDataPacket>
    >(
      "SELECT * FROM ride_statuses WHERE ride_id = ? AND chair_sent_at IS NULL ORDER BY created_at ASC LIMIT 1",
      [ride.id],
    );
    const status = yetSentRideStatus
      ? yetSentRideStatus.status
      : await getLatestRideStatus(ctx.var.dbConn, ride.id);

    const [[user]] = await ctx.var.dbConn.query<Array<User & RowDataPacket>>(
      "SELECT * FROM users WHERE id = ? FOR SHARE",
      [ride.user_id],
    );

    if (yetSentRideStatus?.id) {
      await ctx.var.dbConn.query(
        "UPDATE ride_statuses SET chair_sent_at = CURRENT_TIMESTAMP(6) WHERE id = ?",
        [yetSentRideStatus.id],
      );
    }

    await ctx.var.dbConn.commit();
    return ctx.json(
      {
        data: {
          ride_id: ride.id,
          user: {
            id: user.id,
            name: `${user.firstname} ${user.lastname}`,
          },
          pickup_coordinate: {
            latitude: ride.pickup_latitude,
            longitude: ride.pickup_longitude,
          },
          destination_coordinate: {
            latitude: ride.destination_latitude,
            longitude: ride.destination_longitude,
          },
          status,
        },
        retry_after_ms: 30,
      },
      200,
    );
  } catch (e) {
    await ctx.var.dbConn.rollback();
    return ctx.text(`${e}`, 500);
  }
};

export const chairPostRideStatus = async (ctx: Context<Environment>) => {
  const rideID = ctx.req.param("ride_id");
  const chair = ctx.var.chair;
  const reqJson = await ctx.req.json<{ status: string }>();
  await ctx.var.dbConn.beginTransaction();
  try {
    const [[ride]] = await ctx.var.dbConn.query<Array<Ride & RowDataPacket>>(
      "SELECT * FROM rides WHERE id = ? FOR UPDATE",
      [rideID],
    );
    if (!ride) {
      return ctx.text("ride not found", 404);
    }
    if (ride.chair_id !== chair.id) {
      return ctx.text("not assigned to this ride", 400);
    }
    switch (reqJson.status) {
      // Acknowledge the ride
      case "ENROUTE":
        await ctx.var.dbConn.query(
          "INSERT INTO ride_statuses (id, ride_id, status) VALUES (?, ?, ?)",
          [ulid(), ride.id, "ENROUTE"],
        );
        break;
      // After Picking up user
      case "CARRYING": {
        const status = await getLatestRideStatus(ctx.var.dbConn, ride.id);
        if (status !== "PICKUP") {
          return ctx.text("chair has not arrived yet", 400);
        }
        await ctx.var.dbConn.query(
          "INSERT INTO ride_statuses (id, ride_id, status) VALUES (?, ?, ?)",
          [ulid(), ride.id, "CARRYING"],
        );
        break;
      }
      default:
        return ctx.text("invalid status", 400);
    }
    await ctx.var.dbConn.commit();
    return ctx.body(null, 204);
  } catch (e) {
    await ctx.var.dbConn.rollback();
    return ctx.text(`${e}`, 500);
  }
};
