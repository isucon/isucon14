import type {
  OwnerGetChairsResponse,
  OwnerGetSalesResponse,
} from "./api/api-components";

export type AccessToken = string;

export type ClientAppChair = {
  id: string;
  name: string;
  model: string;
  stats: Partial<{
    total_rides_count: number;
    total_evaluation_avg: number;
  }>;
};

export type Coordinate = { latitude: number; longitude: number };

export type DisplayPos = { x: number; y: number };

export type Distance = { horizontalDistance: number; verticalDistance: number };

export type NearByChair = {
  id: string;
  name: string;
  model: string;
  current_coordinate: Coordinate;
};

export type CampaignData = {
  invitationCode: string;
  registedAt: string;
  used: boolean;
};

export type OwnerChairs = OwnerGetChairsResponse["chairs"];
export type OwnerSales = OwnerGetSalesResponse;

export type SimulatorChair = {
  id: string;
  name: string;
  model: string;
  token: string;
  coordinate: Coordinate;
};
