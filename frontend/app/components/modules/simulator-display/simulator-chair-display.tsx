import { FC, useCallback, useEffect, useMemo, useRef, useState } from "react";
import colors from "tailwindcss/colors";
import { fetchChairPostActivity } from "~/apiClient/apiComponents";
import { RideStatus } from "~/apiClient/apiSchemas";
import { useEmulator } from "~/components/hooks/use-emulate";
import { ChairIcon } from "~/components/icon/chair";
import { PinIcon } from "~/components/icon/pin";
import { LocationButton } from "~/components/modules/location-button/location-button";
import { Map } from "~/components/modules/map/map";
import { Button } from "~/components/primitives/button/button";
import { Toggle } from "~/components/primitives/form/toggle";
import { Modal } from "~/components/primitives/modal/modal";
import { Text } from "~/components/primitives/text/text";
import { useSimulatorContext } from "~/contexts/simulator-context";
import { Coordinate, SimulatorChair } from "~/types";
import { isArrayIncludes } from "~/utils/includes";
import { SimulatorChairRideStatus } from "../simulator-chair-status/simulator-chair-status";

const CoordinatePickup: FC<{
  coordinateState: SimulatorChair["coordinateState"];
}> = ({ coordinateState }) => {
  const [initialMapLocation, setInitialMapLocation] = useState<Coordinate>();
  const [mapLocation, setMapLocation] = useState<Coordinate>();
  const [visibleModal, setVisibleModal] = useState<boolean>(false);
  const modalRef = useRef<HTMLElement & { close: () => void }>(null);

  const handleOpenModal = useCallback(() => {
    setInitialMapLocation(coordinateState.coordinate);
    setVisibleModal(true);
  }, [coordinateState]);

  const handleCloseModal = useCallback(() => {
    if (mapLocation) {
      coordinateState.setter(mapLocation);
    }

    modalRef.current?.close();
    setVisibleModal(false);
  }, [mapLocation, coordinateState]);

  return (
    <>
      <LocationButton
        className="w-full text-right"
        location={coordinateState.coordinate}
        label="椅子位置"
        placeholder="現在位置を設定"
        onClick={handleOpenModal}
      />
      {visibleModal && (
        <div className="fixed inset-0 z-10">
          <Modal
            ref={modalRef}
            center
            onClose={handleCloseModal}
            className="absolute w-full max-w-[800px] max-h-none h-[700px]"
          >
            <div className="w-full h-full flex flex-col items-center">
              <Map
                className="flex-1"
                initialCoordinate={initialMapLocation}
                from={initialMapLocation}
                onMove={(c) => setMapLocation(c)}
                selectable
              />
              <Button
                className="w-full mt-6"
                onClick={handleCloseModal}
                variant="primary"
              >
                この位置で確定する
              </Button>
            </div>
          </Modal>
        </div>
      )}
    </>
  );
};

const SimulatorProgress: FC<{progress: {
  pickup: number
  destlocation: number,
}, chair: {
  model: string,
  rideStatus?: RideStatus,
}}> = ({
  progress, chair
}) => {
  return (
    <div className="flex items-center mt-8">
    {/* Progress */}
    <div className="flex border-b ms-6 pb-1 w-full">
      {/* PICKUP -> ARRIVED */}
      <div className="flex w-1/2">
        <PinIcon color={colors.red[500]} width={20} height={20} />
        {/* road */}
        <div className="relative w-full ms-6">
          {isArrayIncludes(
            [
              "CARRYING",
              "ARRIVED",
              "COMPLETED",
            ] as const satisfies RideStatus[],
            chair.rideStatus,
          ) && (
            <ChairIcon
              model={chair.model}
              className={`size-6 absolute top-[-2px] ${chair.rideStatus === "CARRYING" ? "animate-shake" : ""}`}
              style={{ right: `${progress.destlocation * 100}%` }}
            />
          )}
        </div>
      </div>
      {/* ENROUTE -> PICKUP */}
      <div className="flex w-1/2">
        <PinIcon color={colors.black} width={20} height={20} />
        {/* road */}
        <div className="relative w-full ms-6">
          {isArrayIncludes(
            [
              "MATCHING",
              "ENROUTE",
              "PICKUP",
            ] as const satisfies RideStatus[],
            chair.rideStatus,
          ) && (
            <ChairIcon
              model={chair.model}
              className={`size-6 absolute top-[-2px] ${chair.rideStatus === "ENROUTE" ? "animate-shake" : ""}`}
              style={{ right: `${progress.pickup * 100}%` }}
            />
          )}
        </div>
      </div>
    </div>
  </div>
  )
}

export const SimulatorChairDisplay: FC = () => {
  const { targetChair: chair } = useSimulatorContext();
  const [activate, setActivate] = useState<boolean>(true);
  const [progress, setProgress] = useState<{
    pickup: number;
    destlocation: number;
  }>({
    pickup: 0,
    destlocation: 0,
  });

  // TODO: 仮実装
  useEffect(() => {
    let _progress = 0;
    setInterval(() => {
      _progress = (_progress + 0.1) % 2;
      setProgress({
        pickup: Math.max(_progress - 1, 0),
        destlocation: Math.max(_progress - 1, 0),
      });
    }, 1000);
  }, []);

  const toggleActivate = useCallback(
    (activity: boolean) => {
      try {
        void fetchChairPostActivity({ body: { is_active: activity } });
        setActivate(activity);
      } catch (error) {
        console.error(error);
      }
    },
    [setActivate],
  );

  const rideStatus = useMemo(
    () => chair?.chairNotification?.status ?? "MATCHING",
    [chair],
  );

  useEmulator(chair);

  return (
    <>
      <div className="bg-white rounded shadow px-6 py-4 w-full">
        {chair ? (
          <div className="space-y-4">
            <div className="flex items-center space-x-4">
              <ChairIcon model={chair.model} className="size-12 shrink-0" />
              <div className="space-y-0.5 w-full">
                <Text bold>{chair.name}</Text>
                <Text className="text-xs text-neutral-500">{chair.model}</Text>
                <SimulatorChairRideStatus currentStatus={rideStatus} />
              </div>
            </div>
            <CoordinatePickup coordinateState={chair.coordinateState} />
            <SimulatorProgress progress={progress} chair={{rideStatus: chair.chairNotification?.status, model: chair.model}}/>
          </div>
        ) : (
          <Text className="m-4" size="sm">
            椅子のデータがありません
          </Text>
        )}
      </div>
      {chair && (
        <div className="bg-white rounded shadow px-6 py-4 w-full">
          <div className="flex justify-between items-center">
            <Text size="sm" className="text-neutral-500" bold>
              配車を受け付ける
            </Text>
            <Toggle
              checked={activate}
              onUpdate={(v) => toggleActivate(v)}
              id="chair-activity"
            />
          </div>
        </div>
      )}
    </>
  );
};