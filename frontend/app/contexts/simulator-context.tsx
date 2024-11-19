import { ReactNode, createContext, useContext } from "react";
import type { Coordinate } from "~/apiClient/apiSchemas";
import { getOwners } from "~/initialDataClient/getter";

export type SimulatorChair = {
  id: string;
  name: string;
  model: string;
  token: string;
  coordinateState: {
    coordinate?: Coordinate;
    setter: (coordinate: Coordinate) => void;
  };
};

export type SimulatorOwner = {
  id: string;
  name: string;
  token: string;
  chairs: SimulatorChair[];
};

type ClientSimulatorContextType = { owners: SimulatorOwner[] };

const ClientSimulatorContext = createContext<ClientSimulatorContextType>({
  owners: [],
});

export const SimulatorProvider = ({ children }: { children: ReactNode }) => {
  const owners = getOwners().map(
    (owner) =>
      ({
        ...owner,
        chairs: owner.chairs.map(
          (chair) =>
            ({
              ...chair,
              coordinateState: {
                setter(coordinate) {
                  this.coordinate = coordinate;
                },
              },
            }) satisfies SimulatorChair,
        ),
      }) satisfies SimulatorOwner,
  );

  return (
    <ClientSimulatorContext.Provider value={{ owners }}>
      {children}
    </ClientSimulatorContext.Provider>
  );
};

export const useSimulatorContext = () => useContext(ClientSimulatorContext);
