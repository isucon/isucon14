import { ComponentPropsWithoutRef, FC } from "react";
import { Text } from "~/components/primitives/text/text";

type PriceTextProps = Omit<
  ComponentPropsWithoutRef<typeof Text>,
  "children"
> & {
  value: number;
};

export const PriceText: FC<PriceTextProps> = ({ value, ...rest }) => {
  return (
    <Text {...rest}>{new Intl.NumberFormat("ja-JP").format(value)} 円</Text>
  );
};
