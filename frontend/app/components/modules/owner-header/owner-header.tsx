import { FC } from "react";
import { useClientProviderContext } from "~/contexts/provider-context";

export const OwnerHeader: FC = () => {
  const { provider } = useClientProviderContext();

  return (
    <div className="flex items-center my-8 px-6">
      {/* TODO: UI */}
      <div className="border rounded-full size-16"></div>
      <span className="text-2xl ms-4">{provider?.name}</span>
    </div>
  );
};
