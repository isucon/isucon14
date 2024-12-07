import { ComponentProps, FC } from "react";
import { twMerge } from "tailwind-merge";
import { ChairIcon } from "~/components/icon/chair";
import { Text } from "~/components/primitives/text/text";
import { ClientAppChair } from "~/types";

export const ChairInformation: FC<
  { chair: ClientAppChair } & ComponentProps<"div">
> = ({ chair, className, ...props }) => {
  return (
    <div
      className={twMerge(
        "flex flex-row items-center space-x-4 px-6 py-4 bg-neutral-100 rounded-md w-full",
        className,
      )}
      {...props}
    >
      <div className="rounded-full bg-neutral-200 p-5">
        <ChairIcon model={chair.model} width={40} height={40}></ChairIcon>
      </div>
      <div className="flex flex-col space-y-0.5">
        <div className="flex flex-col space-y-0.5 mb-1">
          <Text tagName="span" bold>
            {chair.name}
          </Text>
          <Text tagName="span" size="xs" className="text-neutral-500">
            {chair.model}
          </Text>
        </div>
        {chair.stats?.total_evaluation_avg && (
          <Text tagName="span" size="xs" className="text-neutral-600">
            <Text tagName="span" className="pr-1">
              評価:
            </Text>
            {chair.stats.total_evaluation_avg.toFixed(1)}
          </Text>
        )}
        {chair.stats?.total_rides_count && (
          <Text tagName="span" size="xs" className="text-neutral-600">
            <Text tagName="span" className="pr-1">
              これまでの配車数:
            </Text>
            {chair.stats?.total_rides_count}
          </Text>
        )}
      </div>
    </div>
  );
};
