import type { MetaFunction } from "@remix-run/node";
import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import colors from "tailwindcss/colors";
import {
  fetchAppGetNearbyChairs,
  fetchAppPostRides,
  fetchAppPostRidesEstimatedFare,
} from "~/apiClient/apiComponents";
import { Coordinate, RideStatus } from "~/apiClient/apiSchemas";
import { LocationButton } from "~/components/modules/location-button/location-button";
import { Map } from "~/components/modules/map/map";
import { Price } from "~/components/modules/price/price";
import { Button } from "~/components/primitives/button/button";
import { Modal } from "~/components/primitives/modal/modal";
import { Text } from "~/components/primitives/text/text";
import { useClientAppRequestContext } from "~/contexts/user-context";
import { NearByChair, isClientApiError } from "~/types";
import { Arrived } from "./driving-state/arrived";
import { Carrying } from "./driving-state/carrying";
import { Enroute } from "./driving-state/enroute";
import { Matching } from "./driving-state/matching";
import { Pickup } from "./driving-state/pickup";

export const meta: MetaFunction = () => {
  return [
    { title: "Top | ISURIDE" },
    { name: "description", content: "目的地まで椅子で快適に移動しましょう" },
  ];
};

type Direction = "from" | "to";
type EstimatePrice = { fare: number; discount: number };

export default function Index() {
  const { status, payload: payload } = useClientAppRequestContext();
  const [internalRideStatus, setInternalRideStatus] = useState<RideStatus>();

  useEffect(() => {
    setInternalRideStatus(status);
  }, [status]);

  const [currentLocation, setCurrentLocation] = useState<Coordinate>();
  const [destLocation, setDestLocation] = useState<Coordinate>();

  const [direction, setDirection] = useState<Direction | null>(null);

  const [selectedLocation, setSelectedLocation] = useState<Coordinate>();

  const [fare, setFare] = useState<number>();

  const onMove = useCallback((coordinate: Coordinate) => {
    setSelectedLocation(coordinate);
  }, []);
  const [isLocationSelectorModalOpen, setLocationSelectorModalOpen] =
    useState(false);

  const locationSelectorModalRef = useRef<HTMLElement & { close: () => void }>(
    null,
  );
  const handleConfirmLocation = useCallback(() => {
    if (direction === "from") {
      setCurrentLocation(selectedLocation);
    } else if (direction === "to") {
      setDestLocation(selectedLocation);
    }
    if (locationSelectorModalRef.current) {
      locationSelectorModalRef.current.close();
    }
  }, [direction, selectedLocation]);

  const isStatusModalOpen = useMemo(() => {
    return (
      internalRideStatus &&
      ["MATCHING", "ENROUTE", "PICKUP", "CARRYING", "ARRIVED"].includes(
        internalRideStatus,
      )
    );
  }, [internalRideStatus]);

  const statusModalRef = useRef<HTMLElement & { close: () => void }>(null);

  const [estimatePrice, setEstimatePrice] = useState<EstimatePrice>();
  useEffect(() => {
    if (!currentLocation || !destLocation) {
      return;
    }
    const abortController = new AbortController();
    fetchAppPostRidesEstimatedFare(
      {
        body: {
          pickup_coordinate: currentLocation,
          destination_coordinate: destLocation,
        },
      },
      abortController.signal,
    )
      .then((res) =>
        setEstimatePrice({ fare: res.fare, discount: res.discount }),
      )
      .catch((err) => {
        console.error(err);
        setEstimatePrice(undefined);
      });
    return () => {
      abortController.abort();
    };
  }, [currentLocation, destLocation]);

  const handleRideRequest = useCallback(async () => {
    if (!currentLocation || !destLocation) {
      return;
    }
    setInternalRideStatus("MATCHING");
    try {
      const rides = await fetchAppPostRides({
        body: {
          pickup_coordinate: currentLocation,
          destination_coordinate: destLocation,
        },
      });
      setFare(rides.fare);
    } catch (error) {
      if (isClientApiError(error)) {
        console.error(error);
      }
    }
  }, [currentLocation, destLocation]);

  // TODO: NearByChairのつなぎこみは後ほど行う
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const [nearByChairs, setNearByChairs] = useState<NearByChair[]>();
  useEffect(() => {
    if (!currentLocation) {
      return;
    }
    const abortController = new AbortController();
    void (async () => {
      try {
        const { chairs } = await fetchAppGetNearbyChairs(
          {
            queryParams: {
              latitude: currentLocation?.latitude,
              longitude: currentLocation?.longitude,
            },
          },
          abortController.signal,
        );
        setNearByChairs(chairs);
      } catch (error) {
        console.error(error);
      }
    })();
    return () => abortController.abort();
  }, [setNearByChairs, currentLocation]);

  // TODO: 以下は上記が正常に返ったあとに削除する
  // const [data, setData] = useState<NearByChair[]>([
  //   {
  //     id: "hoge",
  //     current_coordinate: { latitude: 100, longitude: 100 },
  //     model: "a",
  //     name: "hoge",
  //   },
  //   {
  //     id: "1",
  //     current_coordinate: { latitude: 20, longitude: 20 },
  //     model: "b",
  //     name: "hoge",
  //   },
  //   {
  //     id: "2",
  //     current_coordinate: { latitude: -100, longitude: -100 },
  //     model: "c",
  //     name: "hoge",
  //   },
  //   {
  //     id: "3",
  //     current_coordinate: { latitude: -160, longitude: -100 },
  //     model: "d",
  //     name: "hoge",
  //   },
  //   {
  //     id: "4",
  //     current_coordinate: { latitude: -10, longitude: 100 },
  //     model: "e",
  //     name: "hoge",
  //   },
  // ]);

  // useEffect(() => {
  //   const randomInt = (min: number, max: number) => {
  //     return Math.floor(Math.random() * (max - min + 1)) + min;
  //   };
  //   const update = () => {
  //     setData((data) =>
  //       data.map((chair) => ({
  //         ...chair,
  //         current_coordinate: {
  //           latitude: chair.current_coordinate.latitude + randomInt(-2, 2),
  //           longitude: chair.current_coordinate.longitude + randomInt(-2, 2),
  //         },
  //       })),
  //     );
  //     setTimeout(update, 1000);
  //   };
  //   update();
  // }, []);

  return (
    <>
      <Map
        from={currentLocation}
        to={destLocation}
        initialCoordinate={selectedLocation}
        chairs={nearByChairs}
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
          <div className="flex flex-col items-center mt-4 h-full">
            <div className="flex-grow w-full max-h-[75%] mb-6">
              <Map
                onMove={onMove}
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
              destLocation={payload?.coordinate?.destination}
              pickup={payload?.coordinate?.pickup}
              optimisticFare={fare}
            />
          )}
          {internalRideStatus === "ENROUTE" && (
            <Enroute
              destLocation={payload?.coordinate?.destination}
              pickup={payload?.coordinate?.pickup}
            />
          )}
          {internalRideStatus === "PICKUP" && (
            <Pickup
              destLocation={payload?.coordinate?.destination}
              pickup={payload?.coordinate?.pickup}
            />
          )}
          {internalRideStatus === "CARRYING" && (
            <Carrying
              destLocation={payload?.coordinate?.destination}
              pickup={payload?.coordinate?.pickup}
            />
          )}
          {internalRideStatus === "ARRIVED" && (
            <Arrived
              onEvaluated={() => {
                statusModalRef.current?.close();
              }}
            />
          )}
        </Modal>
      )}
    </>
  );
}
