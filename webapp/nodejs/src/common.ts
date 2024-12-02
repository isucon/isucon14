import type { Connection, RowDataPacket } from "mysql2/promise";
import type { Ride, RideStatus } from "./types/models.js";
import type { Context } from "hono";
import type { StatusCode } from "hono/utils/http-status";

export const INITIAL_FARE = 500;
export const FARE_PER_DISTANCE = 100;

// マンハッタン距離を求める
export const calculateDistance = (
  aLatitude: number,
  aLongitude: number,
  bLatitude: number,
  bLongitude: number,
): number => {
  return Math.abs(aLatitude - bLatitude) + Math.abs(aLongitude - bLongitude);
};

export const calculateFare = (
  pickupLatitude: number,
  pickupLongitude: number,
  destLatitude: number,
  destLongitude: number,
): number => {
  const meterdFare =
    FARE_PER_DISTANCE *
    calculateDistance(
      pickupLatitude,
      pickupLongitude,
      destLatitude,
      destLongitude,
    );
  return INITIAL_FARE + meterdFare;
};

export const calculateSale = (ride: Ride): number => {
  return calculateFare(
    ride.pickup_latitude,
    ride.pickup_longitude,
    ride.destination_latitude,
    ride.destination_longitude,
  );
};

export const getLatestRideStatus = async (
  dbConn: Connection,
  rideId: string,
): Promise<string> => {
  const [[{ status }]] = await dbConn.query<
    Array<Pick<RideStatus, "status"> & RowDataPacket>
  >(
    "SELECT status FROM ride_statuses WHERE ride_id = ? ORDER BY created_at DESC LIMIT 1",
    [rideId],
  );
  return status;
};

export class ErroredUpstream extends Error {
  constructor(message: string) {
    super(message);
    this.name = "ErroredUpstream";
  }
}

export const responseError = (
  ctx: Context,
  error: unknown,
  statusCode: StatusCode,
) => {
  if (error instanceof Error) {
    return ctx.json({ message: error.message }, statusCode);
  }
  throw new Error("error is not an instance of Error");
};
