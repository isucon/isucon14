import { ComponentPropsWithoutRef, FC } from "react";
import { Text } from "~/components/primitives/text/text";

type PriceTextProps = Omit<
  ComponentPropsWithoutRef<typeof Text>,
  "children"
> & {
  value: number;
};

const formatter = new Intl.NumberFormat("ja-JP");

export const PriceText: FC<PriceTextProps> = ({ value, ...rest }) => {
  return <Text {...rest}>{formatter.format(value)} 円</Text>;
};
