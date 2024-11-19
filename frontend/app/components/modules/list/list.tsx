import { ComponentProps, PropsWithChildren } from "react";

export function List({
  children,
  className,
  ...props
}: PropsWithChildren<ComponentProps<"ul">>) {
  return (
    <ul {...props} className={className}>
      {children}
    </ul>
  );
}
