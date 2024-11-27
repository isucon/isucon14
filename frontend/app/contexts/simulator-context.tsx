import {
  ReactNode,
  createContext,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";
import type { Coordinate } from "~/apiClient/apiSchemas";
import { getSimulateChair } from "~/utils/get-initial-data";

import { apiBaseURL } from "~/apiClient/APIBaseURL";
import {
  ChairGetNotificationResponse,
  fetchChairGetNotification,
} from "~/apiClient/apiComponents";
import type { ClientChairRide, SimulatorChair } from "~/types";

type ClientSimulatorContextProps = {
  targetChair?: SimulatorChair;
};

const ClientSimulatorContext = createContext<ClientSimulatorContextProps>({});

function jsonFromSseResult<T>(value: string) {
  const data = value.slice("data:".length).trim();
  return JSON.parse(data) as T;
}

export const useClientChairNotification = (id?: string) => {
  const [notification, setNotification] = useState<
    ChairGetNotificationResponse & { contentType: "event-stream" | "json" }
  >();

  useEffect(() => {
    let abortController: AbortController | undefined;
    const run = async () => {
      abortController = new AbortController();
      try {
        const notification = await fetch(`${apiBaseURL}/chair/notification`);
        const isEventStream = !!notification?.headers
          .get("Content-type")
          ?.split(";")?.[0]
          .includes("text/event-stream");
        if (isEventStream) {
          const reader = notification.body?.getReader();
          const decoder = new TextDecoder();
          const readed = (await reader?.read())?.value;
          const decoded = decoder.decode(readed);
          const json =
            jsonFromSseResult<ChairGetNotificationResponse["data"]>(decoded);
          setNotification(
            json
              ? {
                  data: json,
                  contentType: "event-stream",
                }
              : undefined,
          );
        } else {
          const json = (await notification.json()) as
            | ChairGetNotificationResponse
            | undefined;
          setNotification(json ? { ...json, contentType: "json" } : undefined);
        }
      } catch (error) {
        console.error(error);
      }
    };
    void run();
    return () => {
      abortController?.abort();
    };
  }, [setNotification]);

  const clientChairPayloadWithStatus = useMemo(
    () =>
      notification
        ? {
            status: notification.data?.status,
            payload: {
              ride_id: notification.data?.ride_id,
              coordinate: {
                pickup: notification.data?.pickup_coordinate,
                destination: notification.data?.destination_coordinate,
              },
            },
          }
        : undefined,
    [notification],
  );

  const retryAfterMs = notification?.retry_after_ms ?? 10000;
  const isSSE = notification?.contentType === "event-stream";

  useEffect(() => {
    if (!isSSE) return;
    const eventSource = new EventSource(`${apiBaseURL}/chair/notification`);
    const onMessage = (event: { data: unknown } | undefined) => {
      if (typeof event?.data === "string") {
        const eventData = JSON.parse(
          event?.data,
        ) as ChairGetNotificationResponse["data"];
        setNotification((preRequest) => {
          if (
            preRequest === undefined ||
            eventData?.status !== preRequest.data?.status ||
            eventData?.ride_id !== preRequest.data?.ride_id
          ) {
            return {
              data: eventData,
              contentType: "event-stream",
            };
          } else {
            return preRequest;
          }
        });
      }
      return () => {
        eventSource.close();
      };
    };
    eventSource.addEventListener("message", onMessage);
    return () => {
      eventSource.close();
    };
  }, [isSSE]);

  useEffect(() => {
    if (isSSE) return;
    let timeoutId: ReturnType<typeof setTimeout>;
    let abortController: AbortController | undefined;
    const polling = async () => {
      try {
        abortController = new AbortController();
        const currentNotification = await fetchChairGetNotification(
          {},
          abortController.signal,
        );
        setNotification((preRequest) => {
          if (
            preRequest === undefined ||
            currentNotification?.data?.status !== preRequest.data?.status ||
            currentNotification?.data?.ride_id !== preRequest.data?.ride_id
          ) {
            return {
              data: currentNotification.data,
              retry_after_ms: currentNotification.retry_after_ms,
              contentType: "json",
            };
          } else {
            return preRequest;
          }
        });
        timeoutId = setTimeout(() => void polling(), retryAfterMs);
      } catch (error) {
        console.error(error);
      }
    };

    timeoutId = setTimeout(() => void polling(), retryAfterMs);

    return () => {
      abortController?.abort();
      clearTimeout(timeoutId);
    };
  }, [isSSE, retryAfterMs]);

  const responseClientAppRequest = useMemo<ClientChairRide | undefined>(() => {
    const candidateAppRequest = clientChairPayloadWithStatus;
    return {
      ...candidateAppRequest,
      status: candidateAppRequest?.status,
      user: {
        id,
        name: "ISUCON太郎",
      },
    };
  }, [clientChairPayloadWithStatus, id]);

  return responseClientAppRequest;
};

const simulateChairData = getSimulateChair();

export const SimulatorProvider = ({ children }: { children: ReactNode }) => {
  useEffect(() => {
    if (simulateChairData?.token) {
      document.cookie = `chair_session=${simulateChairData.token}; path=/`;
    }
  }, []);

  const request = useClientChairNotification(simulateChairData?.id);
  const [currentCoodinate, setCurrentCoordinate] = useState<Coordinate>();

  return (
    <ClientSimulatorContext.Provider
      value={{
        targetChair: simulateChairData
          ? {
              ...simulateChairData,
              chairNotification: request,
              coordinateState: {
                setter: setCurrentCoordinate,
                coordinate: currentCoodinate,
              },
              token: simulateChairData?.token,
            }
          : undefined,
      }}
    >
      {children}
    </ClientSimulatorContext.Provider>
  );
};

export const useSimulatorContext = () => useContext(ClientSimulatorContext);
