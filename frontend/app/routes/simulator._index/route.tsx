import { useState } from "react";

import { PulldownSelector } from "~/components/primitives/menu/pulldown";
import { useSimulatorContext } from "~/contexts/simulator-context";
import { ChairInfo } from "./ChairInfo";

export default function Index() {
  const {owners} = useSimulatorContext();
  const ownerNames = [...owners].map((o) => o.name);
  const getOwnerByName = (name: string) => {
    return owners.find(o => o.name === name)
  };

  const [targetOwner, setTargetOwner] = useState(
    getOwnerByName(ownerNames[0]),
  );

  return (
    <div className="p-6">
      <PulldownSelector
        className="mb-3"
        id="ownerNames"
        label="オーナー"
        items={ownerNames}
        onChange={(name) => setTargetOwner(getOwnerByName(name))}
      />
      {targetOwner !== undefined
        ? targetOwner.chairs?.map((c) => <ChairInfo key={c.id} chair={c} />)
        : null}
    </div>
  );
}
