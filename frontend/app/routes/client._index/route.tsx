import type { MetaFunction } from "@remix-run/node";
import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import colors from "tailwindcss/colors";
import {
  fetchAppGetNearbyChairs,
  fetchAppPostRides,
  fetchAppPostRidesEstimatedFare,
} from "~/api/api-components";
import { Coordinate, RideStatus } from "~/api/api-schemas";
import { useGhostChairs } from "~/components/hooks/use-ghost-chairs";
import { CampaignBanner } from "~/components/modules/campaign-banner/campaign-banner";
import { LocationButton } from "~/components/modules/location-button/location-button";
import { Map } from "~/components/modules/map/map";
import { Price } from "~/components/modules/price/price";
import { Button } from "~/components/primitives/button/button";
import { Modal } from "~/components/primitives/modal/modal";
import { Text } from "~/components/primitives/text/text";
import { useClientContext } from "~/contexts/client-context";
import type { Distance, NearByChair } from "~/types";
import { sendClientReady, sendClientRideRequested } from "~/utils/post-message";
import { Arrived } from "./driving-state/arrived";
import { Carrying } from "./driving-state/carrying";
import { Enroute } from "./driving-state/enroute";
import { Matching } from "./driving-state/matching";
import { Pickup } from "./driving-state/pickup";

export const meta: MetaFunction = () => {
  return [
    { title: "ISURIDE" },
    { name: "description", content: "目的地まで椅子で快適に移動しましょう" },
  ];
};

type Direction = "from" | "to";
type EstimatePrice = { fare: number; discount: number };

export default function Index() {
  const { data } = useClientContext();
  const emulateChairs = useGhostChairs();
  const [internalRideStatus, setInternalRideStatus] = useState<RideStatus>();
  const [currentLocation, setCurrentLocation] = useState<Coordinate>();
  const [destLocation, setDestLocation] = useState<Coordinate>();
  const [direction, setDirection] = useState<Direction | null>(null);
  const [selectedLocation, setSelectedLocation] = useState<Coordinate>();
  const [displayedChairs, setDisplayedChairs] = useState<NearByChair[]>([]);
  const [centerCoordinate, setCenterCoodirnate] = useState<Coordinate>();
  const [distance, setDistance] = useState<Distance>();
  const onMove = useCallback(
    (coordinate: Coordinate) => setCenterCoodirnate(coordinate),
    [],
  );
  const onSelectMove = useCallback(
    (coordinate: Coordinate) => setSelectedLocation(coordinate),
    [],
  );
  const onUpdateViewSize = useCallback(
    (distance: Distance) => setDistance(distance),
    [],
  );
  const [isLocationSelectorModalOpen, setLocationSelectorModalOpen] =
    useState(false);
  const locationSelectorModalRef = useRef<HTMLElement & { close: () => void }>(
    null,
  );
  const statusModalRef = useRef<HTMLElement & { close: () => void }>(null);
  const [estimatePrice, setEstimatePrice] = useState<EstimatePrice>();
  const handleConfirmLocation = useCallback(() => {
    if (direction === "from") {
      setCurrentLocation(selectedLocation);
    } else if (direction === "to") {
      setDestLocation(selectedLocation);
    }
    locationSelectorModalRef.current?.close();
  }, [direction, selectedLocation]);

  const isStatusModalOpen = useMemo(() => {
    return (
      internalRideStatus &&
      ["MATCHING", "ENROUTE", "PICKUP", "CARRYING", "ARRIVED"].includes(
        internalRideStatus,
      )
    );
  }, [internalRideStatus]);

  useEffect(() => {
    setInternalRideStatus(data?.status);
  }, [data?.status]);

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
      .catch((error) => {
        console.error(error);
        setEstimatePrice(undefined);
      });
  }, [currentLocation, destLocation]);

  const handleRideRequest = useCallback(async () => {
    if (!currentLocation || !destLocation) {
      return;
    }
    setInternalRideStatus("MATCHING");
    try {
      const { ride_id } = await fetchAppPostRides({
        body: {
          pickup_coordinate: currentLocation,
          destination_coordinate: destLocation,
        },
      });
      sendClientRideRequested(window.parent, { rideId: ride_id });
    } catch (error) {
      console.error(error);
    }
  }, [currentLocation, destLocation]);

  useEffect(() => {
    if (!centerCoordinate) return;
    if (isStatusModalOpen) return;
    let abortController: AbortController | undefined;
    let timeoutId: NodeJS.Timeout | undefined;

    const updateNearByChairs = async ({ latitude, longitude }: Coordinate) => {
      try {
        abortController?.abort();
        abortController = new AbortController();
        const { chairs } = await fetchAppGetNearbyChairs({
          queryParams: {
            latitude,
            longitude,
            distance: distance
              ? Math.max(
                  distance.horizontalDistance + 10,
                  distance.verticalDistance + 10,
                )
              : 150,
          },
        });
        setDisplayedChairs(chairs);
      } catch (error) {
        if (error instanceof DOMException && error.name === "AbortError") {
          return;
        }
        console.error(error);
      }
    };

    const polling = () => {
      void updateNearByChairs(centerCoordinate);
      timeoutId = setTimeout(polling, 10_000);
    };

    timeoutId = setTimeout(polling, 300);

    return () => {
      clearTimeout(timeoutId);
      abortController?.abort();
    };
  }, [centerCoordinate, isStatusModalOpen, distance]);

  useEffect(() => {
    sendClientReady(window.parent, { ready: true });
    return () => {
      sendClientReady(window.parent, { ready: false });
    };
  }, []);

  return (
    <>
      <CampaignBanner />
      <Map
        from={currentLocation}
        to={destLocation}
        onMove={onMove}
        onUpdateViewSize={onUpdateViewSize}
        initialCoordinate={selectedLocation}
        chairs={[...displayedChairs, ...emulateChairs]}
        className="flex-1"
      />
      <div className="w-full px-8 py-8 flex flex-col items-center justify-center">
        <LocationButton
          className="w-full"
          location={currentLocation}
          onClick={() => {
            setDirection("from");
            setLocationSelectorModalOpen(true);
          }}
          placeholder="現在地を選択する"
          label="現在地"
        />
        <Text size="xl">↓</Text>
        <LocationButton
          location={destLocation}
          className="w-full"
          onClick={() => {
            setDirection("to");
            setLocationSelectorModalOpen(true);
          }}
          placeholder="目的地を選択する"
          label="目的地"
        />
        {estimatePrice && (
          <Price
            value={estimatePrice.fare}
            pre="推定運賃"
            discount={estimatePrice.discount}
            className="mt-6 mb-4"
          ></Price>
        )}
        {currentLocation && destLocation && (
          <Button
            variant="primary"
            className="w-full font-bold"
            onClick={() => void handleRideRequest()}
            disabled={!(Boolean(currentLocation) && Boolean(destLocation))}
          >
            ISURIDE
          </Button>
        )}
      </div>
      {isLocationSelectorModalOpen && (
        <Modal
          ref={locationSelectorModalRef}
          onClose={() => setLocationSelectorModalOpen(false)}
        >
          <div className="flex flex-col items-center h-full">
            <div className="flex-grow w-full max-h-[75%] mb-6">
              <Map
                onMove={onSelectMove}
                from={currentLocation}
                to={destLocation}
                selectorPinColor={
                  direction === "from" ? colors.black : colors.red[500]
                }
                initialCoordinate={
                  direction === "from" ? currentLocation : destLocation
                }
                selectable
                className="rounded-2xl"
              />
            </div>
            <p className="font-bold mb-4 text-base">
              {direction === "from" ? "現在地" : "目的地"}
              を選択してください
            </p>
            <Button onClick={handleConfirmLocation}>
              {direction === "from"
                ? "この場所から移動する"
                : "この場所に移動する"}
            </Button>
          </div>
        </Modal>
      )}
      {isStatusModalOpen && (
        <Modal
          ref={statusModalRef}
          onClose={() => setInternalRideStatus("COMPLETED")}
        >
          {internalRideStatus === "MATCHING" && (
            <Matching
              optimistic={{
                destLocation,
                pickup: currentLocation,
                fare: estimatePrice?.fare,
              }}
            />
          )}
          {internalRideStatus === "ENROUTE" && <Enroute />}
          {internalRideStatus === "PICKUP" && <Pickup />}
          {internalRideStatus === "CARRYING" && <Carrying />}
          {internalRideStatus === "ARRIVED" && (
            <Arrived
              onEvaluated={() => {
                statusModalRef.current?.close();
                setCurrentLocation(destLocation);
                setDestLocation(undefined);
                setEstimatePrice(undefined);
              }}
            />
          )}
        </Modal>
      )}
    </>
  );
}
