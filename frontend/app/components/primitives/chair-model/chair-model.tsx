import { FC } from "react";
import { twMerge } from "tailwind-merge";
import { CarGreenIcon } from "~/components/icon/car-green";
import { CarRedIcon } from "~/components/icon/car-red";
import { CarYellowIcon } from "~/components/icon/car-yellow";

<CarRedIcon className="size-[76px] mb-4" />;

export const ChairModel: FC<{ model: string; className?: string }> = (
  props,
) => {
  const Chair = (() => {
    // TODO: 仮実装
    const model = props.model;
    if (/^[a-n]/i.test(model)) return CarGreenIcon;
    if (/^[o-z]/i.test(model)) return CarYellowIcon;
    return CarRedIcon;
  })();

  return <Chair className={twMerge(["size-[1.5rem]", props.className])} />;
};
