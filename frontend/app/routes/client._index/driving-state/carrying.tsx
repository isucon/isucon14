import { FC } from "react";
import { CarGreenIcon } from "~/components/icon/car-green";
import { LocationButton } from "~/components/modules/location-button/location-button";
import { PriceText } from "~/components/modules/price-text/price-text";
import { Text } from "~/components/primitives/text/text";
import { Coordinate } from "~/types";

export const Carrying: FC<{ destLocation?: Coordinate; fare?: number }> = ({
  destLocation,
  fare,
}) => {
  return (
    <div className="w-full h-full px-8 flex flex-col items-center justify-center">
      <CarGreenIcon className="size-[76px] mb-4" />
      <Text size="xl" className="mb-6">
        快適なドライビングをお楽しみください
      </Text>
      <LocationButton label="現在地" className="w-80" />
      <Text size="xl">↓</Text>
      <LocationButton label="目的地" location={destLocation} className="w-80" />
      <p className="mt-8">
        {typeof fare === "number" ? (
          <>
            予定運賃: <PriceText tagName="span" value={fare} />
          </>
        ) : null}
      </p>
    </div>
  );
};
