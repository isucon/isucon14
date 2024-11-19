type TextProps<T extends string> = {
  tabs: readonly { key: T; label: string }[];
  activeTab?: T;
  onTabClick?: (tab: T) => void;
};

export const Tab = <T extends string>({
  tabs,
  activeTab,
  onTabClick,
}: TextProps<T>) => {
  return (
    <nav className="border-b">
      <ul className="flex">
        {tabs.map((tab) => (
          <li
            key={tab.key}
            className={tab.key === activeTab ? "border-b-4 border-black" : ""}
          >
            <button className="px-4 py-2" onClick={() => onTabClick?.(tab.key)}>
              {tab.label}
            </button>
          </li>
        ))}
      </ul>
    </nav>
  );
};
