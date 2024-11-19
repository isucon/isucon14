import { ComponentProps, useCallback, useRef, useState } from "react";
import { twMerge } from "tailwind-merge";

import { RideStatus } from "~/apiClient/apiSchemas";
import { LocationButton } from "~/components/modules/location-button/location-button";
import { Map } from "~/components/modules/map/map";
import { Avatar } from "~/components/primitives/avatar/avatar";
import { Button } from "~/components/primitives/button/button";
import { Toggle } from "~/components/primitives/form/toggle";
import { Modal } from "~/components/primitives/modal/modal";
import { SimulatorChair } from "~/contexts/simulator-context";
import { Coordinate } from "~/types";

type Props = {
  chair: SimulatorChair;
};

function Statuses(
  props: ComponentProps<"div"> & {
    currentStatus: RideStatus;
  },
) {
  const labelByStatus: Record<RideStatus, [label: string, colorClass: string]> =
    {
      MATCHING: ["空車", "text-sky-600"],
      ENROUTE: ["迎車", "text-amber-600"],
      PICKUP: ["乗車待ち", "text-amber-600"],
      CARRYING: ["賃走", "text-red-600"],
      ARRIVED: ["到着", "text-green-600"],
      COMPLETED: ["完了", "text-green-600"],
    } as const;

  const { currentStatus, className, ...rest } = props;
  const [label, colorClass] = labelByStatus[currentStatus];
  return (
    <div className={twMerge(`font-bold ${colorClass}`, className)} {...rest}>
      <span className="before:content-['●'] before:mr-2">{label}</span>
    </div>
  );
}

function CoordinatePickup(props: {
  location: ReturnType<typeof useState<Coordinate>>;
}) {
  const [location, setLocation] = props.location;
  const [currentLocation, setCurrentLocation] = useState<Coordinate>();

  const [visibleModal, setVisibleModal] = useState<boolean>(false);
  const modalRef = useRef<HTMLElement & { close: () => void }>(null);

  const handleCloseModal = useCallback(() => {
    setLocation(currentLocation);
    modalRef.current?.close();
    setVisibleModal(false);
  }, [setLocation, currentLocation]);

  return (
    <>
      <LocationButton
        className="w-full"
        location={location}
        label="現在位置"
        onClick={() => setVisibleModal(true)}
      />
      {visibleModal && (
        <Modal ref={modalRef} onClose={handleCloseModal}>
          <div className="w-full h-full flex flex-col items-center">
            <Map
              className="max-h-[80%]"
              initialCoordinate={location}
              from={location}
              onMove={(c) => setCurrentLocation(c)}
              selectable
            />
            <Button
              className="w-full my-6"
              onClick={handleCloseModal}
              variant="primary"
            >
              この座標で確定する
            </Button>
          </div>
        </Modal>
      )}
    </>
  );
}

export function ChairInfo(props: Props) {
  const location = useState<Coordinate>();
  const [activate, setActivate] = useState<boolean>(false);
  const [rideStatus] = useState<RideStatus>("MATCHING"); // TODO fetch処理

  return (
    <div
      className="
        border-t
        flex
      "
    >
      <Avatar className="left-container mx-3 my-auto" />
      <div className="right-container m-3 flex-grow">
        <div className="right-top flex">
          <div className="right-top-left flex-grow">
            <div className="chair-name font-bold">
              <span>{props.chair.name}</span>
              <span className="ml-1 text-xs font-normal text-neutral-500">
                {props.chair.model}
              </span>
            </div>
            <Statuses className="my-2" currentStatus={rideStatus} />
          </div>
          <div className="right-top-right flex items-center">
            <span className="text-xs font-bold text-neutral-500 mr-3">
              配車受付
            </span>
            <Toggle value={activate} onUpdate={(v) => setActivate(v)} />
          </div>
        </div>
        <div className="right-bottom">
          <CoordinatePickup location={location} />
        </div>
      </div>
    </div>
  );
}
