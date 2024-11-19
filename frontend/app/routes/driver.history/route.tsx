import type { MetaFunction } from "@remix-run/node";
import { List } from "~/components/modules/list/list";
import { ListItem } from "~/components/modules/list/list-item";
import { PriceText } from "~/components/modules/price-text/price-text";

export const meta: MetaFunction = () => {
  return [
    { title: "椅子履歴 | ISURIDE" },
    { name: "description", content: "配椅子履歴" },
  ];
};

export default function Index() {
  const items = [
    {
      date: "2024/08/24",
      from: "xxx",
      to: "yyy",
      price: 1234,
    },
  ];

  return (
    <section className="flex-1 mx-4">
      <h2 className="text-2xl my-4">履歴</h2>
      <List>
        {items.map((item) => (
          <ListItem key={item.date} className="flex justify-between">
            <span>
              <span>{item.date}</span>
              <span className="ms-4">
                {item.from} → {item.to}
              </span>
            </span>
            <PriceText value={item.price} />
          </ListItem>
        ))}
      </List>
    </section>
  );
}
