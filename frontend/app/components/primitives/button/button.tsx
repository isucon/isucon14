import { Link } from "@remix-run/react";
import { ComponentProps, FC, PropsWithChildren, useMemo } from "react";
import { twMerge } from "tailwind-merge";

type Variant = "light" | "primary" | "danger";
type Size = "sm" | "md";

export const ButtonLink: FC<PropsWithChildren<ComponentProps<typeof Link>>> = ({
  to,
  className,
  children,
  ...props
}) => {
  return (
    <Link
      {...props}
      to={to}
      className={twMerge(
        "w-full py-2 text-center border border-neutral-500 rounded-md bg-neutral-200",
        className,
      )}
    >
      {children}
    </Link>
  );
};

export const Button: FC<
  PropsWithChildren<
    ComponentProps<"button"> & { variant?: Variant; size?: Size }
  >
> = ({ children, className, variant, size = "md", ...props }) => {
  const variantClasses = useMemo(() => {
    switch (variant) {
      case "primary":
        return "text-white bg-sky-700 active:brightness-90 hover:brightness-90 focus:brightness-90";
      case "danger":
        return "text-white bg-rose-600 active:brightness-90 hover:brightness-90 focus:brightness-90";
      case "light":
      default:
        return "bg-[#F0EFED] active:brightness-90 hover:brightness-90 focus:brightness-90";
    }
  }, [variant]);

  const sizeClasses = useMemo(() => {
    switch (size) {
      case "sm":
        return "py-2 px-3";
      case "md":
        return "py-3 px-4";
    }
  }, [size]);

  return (
    <button
      type="button"
      className={twMerge(
        "rounded-md bg-neutral-800 border border-transparent text-center text-sm transition-all shadow-md hover:shadow-lg disabled:opacity-50 disabled:shadow-none ml-2",
        variantClasses,
        sizeClasses,
        className,
      )}
      {...props}
    >
      {children}
    </button>
  );
};
