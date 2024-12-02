import type { Context } from "hono";
import type { Environment } from "./types/hono.js";
import { secureRandomStr } from "./utils/random.js";
import { setCookie } from "hono/cookie";
import type { RowDataPacket } from "mysql2";
import type { Chair, Ride } from "./types/models.js";
import { calculateSale, responseError } from "./common.js";
import { ulid } from "ulid";

export const ownerPostOwners = async (ctx: Context<Environment>) => {
  const reqJson = await ctx.req.json<{ name: string }>();
  const { name } = reqJson;
  if (!name) {
    return responseError(
      ctx,
      new Error("some of required fields(name) are empty"),
      400,
    );
  }
  const ownerId = ulid();
  const accessToken = secureRandomStr(32);
  const chairRegisterToken = secureRandomStr(32);
  await ctx.var.dbConn.query(
    "INSERT INTO owners (id, name, access_token, chair_register_token) VALUES (?, ?, ?, ?)",
    [ownerId, name, accessToken, chairRegisterToken],
  );

  setCookie(ctx, "owner_session", accessToken, { path: "/" });

  return ctx.json(
    { id: ownerId, chair_register_token: chairRegisterToken },
    201,
  );
};

export const ownerGetSales = async (ctx: Context<Environment>) => {
  const since = new Date(0);
  const until = new Date("9999-12-31T23:59:59Z");
  if (ctx.req.query("since")) {
    since.setTime(Number(ctx.req.query("since")));
  }
  if (ctx.req.query("until")) {
    until.setTime(Number(ctx.req.query("until")));
  }

  const owner = ctx.var.owner;
  await ctx.var.dbConn.beginTransaction();
  try {
    const [chairs] = await ctx.var.dbConn.query<Array<Chair & RowDataPacket>>(
      "SELECT * FROM chairs WHERE owner_id = ?",
      [owner.id],
    );

    let totalSales = 0;
    const chairSales = [];
    const modelSalesByModel: { [key: string]: number } = {};
    for (const chair of chairs) {
      const [rides] = await ctx.var.dbConn.query<Array<Ride & RowDataPacket>>(
        "SELECT rides.* FROM rides JOIN ride_statuses ON rides.id = ride_statuses.ride_id WHERE chair_id = ? AND status = 'COMPLETED' AND updated_at BETWEEN ? AND ? + INTERVAL 999 MICROSECOND",
        [chair.id, since, until],
      );
      const sales = sumSales(rides);
      totalSales += sales;
      chairSales.push({ id: chair.id, name: chair.name, sales });
      modelSalesByModel[chair.model] =
        (modelSalesByModel[chair.model] ?? 0) + sales;
    }

    const models = Object.entries(modelSalesByModel).map(([model, sales]) => ({
      model,
      sales,
    }));

    return ctx.json({
      total_sales: totalSales,
      chairs: chairSales,
      models,
    });
  } catch (e) {
    await ctx.var.dbConn.rollback();
    return responseError(ctx, e, 500);
  }
};

type ChairWithDetail = {
  id: string;
  owner_id: string;
  name: string;
  access_token: string;
  model: string;
  is_active: boolean;
  created_at: Date;
  updated_at: Date;
  total_distance: number;
  total_distance_updated_at: Date | null;
};

type OwnerGetChairsResponseChair = {
  id: string;
  name: string;
  model: string;
  active: boolean;
  registered_at: number;
  total_distance: number;
  total_distance_updated_at?: number;
};

export const ownerGetChairs = async (ctx: Context<Environment>) => {
  const owner = ctx.var.owner;
  const [chairs] = await ctx.var.dbConn.query<
    Array<ChairWithDetail & RowDataPacket>
  >(
    `SELECT id,
       owner_id,
       name,
       access_token,
       model,
       is_active,
       created_at,
       updated_at,
       IFNULL(total_distance, 0) AS total_distance,
       total_distance_updated_at
FROM chairs
       LEFT JOIN (SELECT chair_id,
                          SUM(IFNULL(distance, 0)) AS total_distance,
                          MAX(created_at)          AS total_distance_updated_at
                   FROM (SELECT chair_id,
                                created_at,
                                ABS(latitude - LAG(latitude) OVER (PARTITION BY chair_id ORDER BY created_at)) +
                                ABS(longitude - LAG(longitude) OVER (PARTITION BY chair_id ORDER BY created_at)) AS distance
                         FROM chair_locations) tmp
                   GROUP BY chair_id) distance_table ON distance_table.chair_id = chairs.id
WHERE owner_id = ?`,
    [owner.id],
  );

  const chairResponse = chairs.map((chair) => {
    const c: OwnerGetChairsResponseChair = {
      id: chair.id,
      name: chair.name,
      model: chair.model,
      active: !!chair.is_active,
      registered_at: chair.created_at.getTime(),
      total_distance: Number(chair.total_distance),
    };
    if (chair.total_distance_updated_at) {
      c.total_distance_updated_at = chair.total_distance_updated_at.getTime();
    }
    return c;
  });

  return ctx.json({ chairs: chairResponse });
};

function sumSales(rides: Ride[]) {
  return rides.reduce((acc, ride) => acc + calculateSale(ride), 0);
}
