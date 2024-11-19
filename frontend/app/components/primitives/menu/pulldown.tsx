type Props = {
  className?: string;
  id: string;
  label?: string;
  items: { id: string; name: string }[];
  onChange: (selected: string) => void;
};
export function PulldownSelector(props: Props) {
  return (
    <div className={`w-full flex items-center ${props.className}`}>
      {props.label !== undefined ? (
        <label className="pr-3 font-bold" htmlFor={props.id}>
          {props.label}
        </label>
      ) : null}
      <select
        className="
          bg-neutral-300 
          border-2 rounded border-neutral-500
          p-1 flex-grow
        "
        name={props.id}
        id={props.id}
        onChange={(e) => props.onChange(e.target.value)}
      >
        {props.items.map((value) => (
          <option key={value.name} value={value.id}>
            {value.name}
          </option>
        ))}
      </select>
    </div>
  );
}
