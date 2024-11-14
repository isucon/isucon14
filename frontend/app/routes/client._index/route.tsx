import type { MetaFunction } from "@remix-run/node";
import { useCallback, useEffect, useRef, useState } from "react";
import {
  fetchAppPostRides,
  fetchAppPostRidesEstimatedFare,
} from "~/apiClient/apiComponents";
import { Coordinate } from "~/apiClient/apiSchemas";
import { useOnClickOutside } from "~/components/hooks/use-on-click-outside";
import { LocationButton } from "~/components/modules/location-button/location-button";
import { Map } from "~/components/modules/map/map";
import { PriceText } from "~/components/modules/price-text/price-text";
import { Button } from "~/components/primitives/button/button";
import { Modal } from "~/components/primitives/modal/modal";
import { Text } from "~/components/primitives/text/text";
import { useClientAppRequestContext } from "~/contexts/user-context";
import { Arrived } from "./driving-state/arrived";
import { Carrying } from "./driving-state/carrying";
import { Dispatched } from "./driving-state/dispatched";

export const meta: MetaFunction = () => {
  return [
    { title: "Top | ISURIDE" },
    { name: "description", content: "目的地まで椅子で快適に移動しましょう" },
  ];
};

type Action = "from" | "to";
type EstimatePrice = { fare: number; discount: number };

export default function Index() {
  const data = useClientAppRequestContext();

  const [action, setAction] = useState<Action>();
  const [selectedLocation, setSelectedLocation] = useState<Coordinate>();
  const [currentLocation, setCurrentLocation] = useState<Coordinate>();
  const [destLocation, setDestLocation] = useState<Coordinate>();
  const [estimatePrice, setEstimatePrice] = useState<EstimatePrice>();

  const [isSelectorModalOpen, setIsSelectorModalOpen] = useState(false);
  const selectorModalRef = useRef<HTMLElement & { close: () => void }>(null);
  const handleSelectorModalClose = useCallback(() => {
    if (selectorModalRef.current) {
      selectorModalRef.current.close();
    }
  }, []);

  const drivingStateModalRef = useRef(null);

  const onClose = useCallback(() => {
    if (action === "from") setCurrentLocation(selectedLocation);
    if (action === "to") setDestLocation(selectedLocation);
    setIsSelectorModalOpen(false);
  }, [action, selectedLocation]);

  const onMove = useCallback((coordinate: Coordinate) => {
    setSelectedLocation(coordinate);
  }, []);

  const handleOpenModal = useCallback((action: Action) => {
    setIsSelectorModalOpen(true);
    setAction(action);
  }, []);

  // TODO: requestId をベースに配車キャンセルしたい
  /* eslint-disable-next-line @typescript-eslint/no-unused-vars */
  const [requestId, setRequestId] = useState<string>("");
  const handleRideRequest = useCallback(async () => {
    if (!currentLocation || !destLocation) {
      return;
    }
    await fetchAppPostRides({
      body: {
        pickup_coordinate: currentLocation,
        destination_coordinate: destLocation,
      },
      headers: {
        Authorization: `Bearer ${data.auth?.accessToken}`,
      },
    }).then((res) => setRequestId(res.ride_id));
  }, [data, currentLocation, destLocation]);

  useEffect(() => {
    if (!currentLocation || !destLocation) {
      return;
    }

    fetchAppPostRidesEstimatedFare({
      body: {
        pickup_coordinate: currentLocation,
        destination_coordinate: destLocation,
      },
    })
      .then((res) =>
        setEstimatePrice({ fare: res.fare, discount: res.discount }),
      )
      .catch((err) => {
        console.error(err);
        setEstimatePrice(undefined);
      });
  }, [currentLocation, destLocation]);

  useOnClickOutside(selectorModalRef, handleSelectorModalClose);

  return (
    <>
      <Map
        from={currentLocation}
        to={destLocation}
        initialCoordinate={selectedLocation}
      />
      <div className="w-full px-8 py-8 flex flex-col items-center justify-center">
        <LocationButton
          className="w-full"
          location={currentLocation}
          onClick={() => {
            handleOpenModal("from");
          }}
          placeholder="現在地を選択する"
          label="現在地"
        />
        <Text size="xl">↓</Text>
        <LocationButton
          location={destLocation}
          className="w-full"
          onClick={() => {
            handleOpenModal("to");
          }}
          placeholder="目的地を選択する"
          label="目的地"
        />
        {estimatePrice && (
          <div className="flex mt-4">
            <Text>推定運賃: </Text>
            <PriceText className="px-4" value={estimatePrice.fare} />
            <Text>(割引額: </Text>
            <PriceText value={estimatePrice.discount} />
            <Text>)</Text>
          </div>
        )}
        <Button
          variant="primary"
          className="w-full mt-6 font-bold"
          onClick={() => void handleRideRequest()}
          disabled={!(Boolean(currentLocation) && Boolean(destLocation))}
        >
          ISURIDE
        </Button>
      </div>
      {isSelectorModalOpen && (
        <Modal ref={selectorModalRef} onClose={onClose}>
          <div className="flex flex-col items-center mt-4 h-full">
            <div className="flex-grow w-full max-h-[75%] mb-6">
              <Map
                onMove={onMove}
                from={currentLocation}
                to={destLocation}
                initialCoordinate={
                  action === "from" ? currentLocation : destLocation
                }
                selectable
                className="rounded-2xl"
              />
            </div>
            <p className="font-bold mb-4 text-base">
              {action === "from" ? "現在地" : "目的地"}を選択してください
            </p>
            <Button onClick={handleSelectorModalClose}>
              {action === "from"
                ? "この場所から移動する"
                : "この場所に移動する"}
            </Button>
          </div>
        </Modal>
      )}
      {data?.status && (
        <Modal ref={drivingStateModalRef}>
          {data.status === "PICKUP" && (
            <Dispatched destLocation={data?.payload?.coordinate?.destination} />
          )}
          {data.status === "CARRYING" && (
            <Carrying destLocation={data?.payload?.coordinate?.destination} />
          )}
          {data.status === "ARRIVED" && <Arrived />}
        </Modal>
      )}
    </>
  );
}
