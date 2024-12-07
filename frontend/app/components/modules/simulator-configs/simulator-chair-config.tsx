import { FC, memo, useMemo } from "react";
import { twMerge } from "tailwind-merge";
import colors from "tailwindcss/colors";
import { RideStatus } from "~/api/api-schemas";
import { ChairIcon } from "~/components/icon/chair";
import { PinIcon } from "~/components/icon/pin";
import { ConfigFrame } from "~/components/primitives/frame/config-frame";
import { Text } from "~/components/primitives/text/text";
import { useSimulatorContext } from "~/contexts/simulator-context";
import type { Coordinate } from "~/types";
import { getSimulatorStartCoordinate } from "~/utils/storage";
import { SimulatorChairLocationButton } from "../simulator-parts/simulator-chair-location-button";
import { SimulatorChairRideStatus } from "../simulator-parts/simulator-chair-status-label";

const progress = (
  start: Coordinate,
  current: Coordinate,
  end: Coordinate,
): number => {
  const distance =
    Math.abs(end.latitude - start.latitude) +
    Math.abs(end.longitude - start.longitude);
  const progress =
    Math.abs(end.latitude - current.latitude) +
    Math.abs(end.longitude - current.longitude);
  return Math.floor(((distance - progress) / distance) * 100);
};

const ChairProgress: FC<{
  model: string;
  rideStatus: RideStatus | undefined;
  currentLoc: Coordinate | undefined;
  pickupLoc?: Coordinate;
  destLoc?: Coordinate;
}> = ({ model, rideStatus, pickupLoc, destLoc, currentLoc }) => {
  const startLoc = useMemo(() => {
    return typeof rideStatus !== "undefined"
      ? getSimulatorStartCoordinate()
      : null;
  }, [rideStatus]);

  const progressToPickup: number = useMemo(() => {
    if (!rideStatus || !pickupLoc || !startLoc || !currentLoc) {
      return 0;
    }
    switch (rideStatus) {
      case "MATCHING":
      case "COMPLETED":
        return 0;
      case "PICKUP":
      case "ARRIVED":
      case "CARRYING":
        return 100;
      default:
        return progress(startLoc, currentLoc, pickupLoc);
    }
  }, [rideStatus, pickupLoc, startLoc, currentLoc]);

  const progressToDestination: number = useMemo(() => {
    if (!rideStatus || !destLoc || !pickupLoc || !currentLoc) {
      return 0;
    }
    switch (rideStatus) {
      case "MATCHING":
      case "COMPLETED":
      case "PICKUP":
      case "ENROUTE":
        return 0;
      case "ARRIVED":
        return 100;
      default:
        return progress(pickupLoc, currentLoc, destLoc);
    }
  }, [rideStatus, destLoc, pickupLoc, currentLoc]);

  return (
    <div className="flex items-center border-b pb-0.5 w-full">
      <div className="flex w-1/2">
        <div className="w-full me-6">
          {rideStatus &&
            ["MATCHING", "COMPLETED", "ENROUTE"].includes(rideStatus) && (
              <div
                className="transition-transform duration-300"
                style={{ transform: `translateX(${progressToPickup}%)` }}
              >
                <div
                  className={twMerge(
                    rideStatus === "ENROUTE" && "animate-shake",
                  )}
                >
                  <ChairIcon model={model} className={"scale-x-[-1] size-6"} />
                </div>
              </div>
            )}
        </div>
        <PinIcon color={colors.black} width={20} />
      </div>
      <div className="flex w-1/2">
        <div className="w-full me-6">
          {rideStatus &&
            ["PICKUP", "CARRYING", "ARRIVED"].includes(rideStatus) && (
              <div
                className="transition-transform duration-300"
                style={{ transform: `translateX(${progressToDestination}%)` }}
              >
                <div
                  className={twMerge(
                    rideStatus === "CARRYING" && "animate-shake",
                  )}
                >
                  <ChairIcon model={model} className={"scale-x-[-1] size-6"} />
                </div>
              </div>
            )}
        </div>
        <PinIcon color={colors.red[500]} width={20} />
      </div>
    </div>
  );
};

const ChairDetailInfo = memo(
  function ChairDetailInfo({
    chairModel,
    chairName,
    rideStatus,
  }: {
    chairModel: string;
    chairName: string;
    rideStatus: RideStatus;
  }) {
    return chairModel && chairName && rideStatus ? (
      <div className="flex items-center space-x-4">
        <ChairIcon model={chairModel} className="size-12 shrink-0" />
        <div className="space-y-0.5 w-full">
          <Text bold>{chairName}</Text>
          <Text className="text-xs text-neutral-500">{chairModel}</Text>
          <SimulatorChairRideStatus currentStatus={rideStatus} />
        </div>
      </div>
    ) : null;
  },
  (prev, next) =>
    prev.chairModel === next.chairModel &&
    prev.chairName === next.chairName &&
    prev.rideStatus === next.rideStatus,
);

export const SimulatorChairConfig: FC = () => {
  const { data, chair, setCoordinate, isAnotherSimulatorBeingUsed } =
    useSimulatorContext();
  const rideStatus = useMemo(() => data?.status ?? "MATCHING", [data]);
  return (
    <ConfigFrame aria-disabled={isAnotherSimulatorBeingUsed}>
      {!chair && (
        <Text className="m-4" size="sm">
          椅子のデータがありません
        </Text>
      )}
      {chair && (
        <>
          <div className="space-y-4">
            <ChairDetailInfo
              chairModel={chair.model}
              chairName={chair.name}
              rideStatus={rideStatus}
            />
            <SimulatorChairLocationButton
              coordinate={chair.coordinate}
              setCoordinate={setCoordinate}
            />
            <ChairProgress
              model={chair.model}
              rideStatus={rideStatus}
              currentLoc={chair.coordinate}
              pickupLoc={data?.pickup_coordinate}
              destLoc={data?.destination_coordinate}
            />
          </div>
          {isAnotherSimulatorBeingUsed && (
            <div className="absolute top-0 left-0 w-full h-full bg-neutral-500 bg-opacity-60 flex items-center justify-center cursor-not-allowed">
              <Text className="text-white" bold size="sm">
                現在、他のシミュレーターが使用中です
              </Text>
            </div>
          )}
        </>
      )}
    </ConfigFrame>
  );
};
