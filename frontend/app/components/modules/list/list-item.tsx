import { ComponentProps, PropsWithChildren } from "react";
import { twMerge } from "tailwind-merge";

export function ListItem({
  children,
  className,
  ...props
}: PropsWithChildren<ComponentProps<"li">>) {
  return (
    <li {...props} className={twMerge("px-4 py-3 border-b", className)}>
      {children}
    </li>
  );
}
