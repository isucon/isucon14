import { PropsWithChildren } from "react";

export const Badge = ({ children }: PropsWithChildren) => {
  return (
    <span className="border rounded-md border-gray-300 text-sm px-2 py-1">
      {children}
    </span>
  );
};
