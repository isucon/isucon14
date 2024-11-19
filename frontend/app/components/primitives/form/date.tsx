import { FC, PropsWithoutRef } from "react";
import { twMerge } from "tailwind-merge";

type DateInputProps = PropsWithoutRef<{
  id: string;
  name: string;
  className?: string;
  required?: boolean;
  onChange?: React.ChangeEventHandler<HTMLInputElement>;
}>;

export const DateInput: FC<DateInputProps> = (props) => {
  return (
    <input
      type="date"
      id={props.id}
      name={props.name}
      className={twMerge(
        "mt-1 p-2 w-full border border-neutral-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500",
        props?.className,
      )}
      required={props.required}
      onChange={props.onChange}
    ></input>
  );
};
