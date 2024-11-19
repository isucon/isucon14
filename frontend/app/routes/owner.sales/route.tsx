import type { MetaFunction } from "@remix-run/node";
import { useSearchParams } from "@remix-run/react";
import { useMemo, useState } from "react";
import { List } from "~/components/modules/list/list";
import { ListItem } from "~/components/modules/list/list-item";
import { PriceText } from "~/components/modules/price-text/price-text";
import { ChairModel } from "~/components/primitives/chair-model/chair-model";
import { DateInput } from "~/components/primitives/form/date";
import { Tab } from "~/components/primitives/tab/tab";
import { useClientProviderContext } from "~/contexts/provider-context";

export const meta: MetaFunction = () => {
  return [{ title: "ISUCON14" }, { name: "description", content: "isucon14" }];
};

export default function Index() {
  const [, setSearchParams] = useSearchParams();

  const tabs = [
    { key: "chair", label: "椅子別" },
    { key: "model", label: "モデル別" },
  ] as const;

  type Tab = (typeof tabs)[number]["key"];
  const [tab, setTab] = useState<Tab>("chair");

  const { sales, chairs } = useClientProviderContext();

  const items = useMemo(() => {
    if (!sales || !chairs) {
      return [];
    }
    const chairModelMap = new Map(chairs.map((c) => [c.id, c.model]));
    return tab === "chair"
      ? sales.chairs.map((item) => ({
          key: item.id,
          name: item.name,
          model: chairModelMap.get(item.id) ?? "",
          sales: item.sales,
        }))
      : sales.models.map((item) => ({
          key: item.model,
          name: item.model,
          model: item.model,
          sales: item.sales,
        }));
  }, [sales, chairs, tab]);

  const updateDate = (key: "since" | "until", value: string) => {
    setSearchParams((prev) => {
      prev.set(key, value);
      return prev;
    });
  };

  const switchTab = (tab: Tab) => {
    setTab(tab);
  };

  return (
    <section className="flex-1 overflow-hidden flex flex-col mx-4">
      <h1 className="text-2xl my-4">売上</h1>
      <div className="flex items-baseline gap-2 mb-2">
        <DateInput
          id="sales-since"
          name="since"
          className="w-48"
          onChange={(e) => updateDate("since", e.target.value)}
        />
        →
        <DateInput
          id="sales-until"
          name="until"
          className="w-48"
          onChange={(e) => updateDate("until", e.target.value)}
        />
      </div>
      {sales ? (
        <>
          <div className="flex">
            <PriceText
              value={sales.total_sales}
              size="2xl"
              bold
              className="ms-auto px-4"
            />
          </div>
          <Tab tabs={tabs} activeTab={tab} onTabClick={switchTab} />
          <List className="overflow-auto">
            {items.map((item) => (
              <ListItem key={item.key} className="flex">
                <ChairModel model={item.model} />
                <span className="ms-4">{item.name}</span>
                <PriceText
                  tagName="span"
                  value={item.sales}
                  className="ms-auto"
                />
              </ListItem>
            ))}
          </List>
        </>
      ) : null}
    </section>
  );
}
