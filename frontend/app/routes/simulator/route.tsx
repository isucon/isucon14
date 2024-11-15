import { Outlet } from "@remix-run/react";
import { useCallback, useMemo, useState } from "react";
import { SimulatorProvider } from "~/contexts/simulator-context";

type DropDownItem = {
  targetId: string;
  label: string;
};

export function DropdownMenu({
  items,
  onSelect,
}: {
  items: DropDownItem[];
  onSelect: (targetId: string) => void;
}) {
  const [isOpen, setIsOpen] = useState(false);
  const [selectedItem, setSelectedItem] = useState<DropDownItem>(items[0]);

  const targetItems = useMemo(() => [...items], [items]);
  const selectLabel = useMemo(() => selectedItem.label, [selectedItem]);

  const handleToggle = useCallback(() => {
    setIsOpen((prev) => !prev);
  }, []);

  const handleClick = useCallback(
    (item: DropDownItem) => {
      onSelect(item.targetId);
      setSelectedItem(item);
      setIsOpen(false); // 項目を選択したらメニューを閉じる
    },
    [onSelect, setSelectedItem, setIsOpen],
  );

  return (
    <div className="relative inline-block text-left">
      <button
        onClick={handleToggle}
        className="bg-blue-500 text-white px-4 py-2 rounded-md shadow-md hover:bg-blue-600 focus:outline-none"
      >
        {selectLabel}
      </button>
      {isOpen && (
        <div className="absolute left-0 mt-2 w-48 bg-white border border-gray-200 rounded-md shadow-lg z-10">
          {targetItems.map((item, index) => (
            <button
              key={index}
              onClick={() => handleClick(item)}
              className="block w-full text-left px-4 py-2 text-gray-700 hover:bg-gray-100"
            >
              {item.label}
            </button>
          ))}
        </div>
      )}
    </div>
  );
}

export default function SimulatorLayout() {
  // TODO: オーナーを動的に取得するようにする
  const list = [
    { label: "テストオーナー1", targetId: "test1" },
    { label: "テストオーナー2", targetId: "test2" },
  ] satisfies DropDownItem[];
  const [selected, SetSelected] = useState<string>(list[0].targetId);
  const onSelect = useCallback((item: string) => {
    SetSelected(item);
  }, []);

  const mainContent = (() => {
    return (
      <SimulatorProvider providerId={selected}>
        <Outlet />
      </SimulatorProvider>
    );
  })();

  return (
    <div>
      <DropdownMenu items={list} onSelect={onSelect} />
      {mainContent}
    </div>
  );
}
