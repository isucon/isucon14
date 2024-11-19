import { useState } from "react";

import { PulldownSelector } from "~/components/primitives/menu/pulldown";
import {
  SimulatorChair,
  useSimulatorContext,
} from "~/contexts/simulator-context";
import { ChairInfo } from "./ChairInfo";

export default function Index() {
  const { owners } = useSimulatorContext();
  const ownerNames = [...owners].map((o) => ({ id: o.id, name: o.name }));
  const getOwnerById = (id: string) => {
    return owners.find((o) => o.id === id);
  };

  const [targetOwner, setTargetOwner] = useState<{ chairs: SimulatorChair[] }>(
    owners[0],
  );

  return (
    <div className="p-6">
      <PulldownSelector
        className="mb-3"
        id="ownerNames"
        label="オーナー"
        items={ownerNames}
        onChange={(id) => setTargetOwner(getOwnerById(id) ?? { chairs: [] })}
      />
      {targetOwner !== undefined
        ? targetOwner.chairs?.map((c) => <ChairInfo key={c.id} chair={c} />)
        : null}
    </div>
  );
}
