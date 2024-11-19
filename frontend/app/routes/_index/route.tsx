import type { MetaFunction } from "@remix-run/node";
import { Link } from "@remix-run/react";

export const meta: MetaFunction = () => {
  return [
    { title: "ISURIDE" },
    {
      name: "description",
      content: "ISURIDEは椅子でユーザーを運ぶ新感覚のサービスです",
    },
  ];
};

export default function Index() {
  return (
    <main className="font-sans p-8">
      <h1 className="text-3xl font-bold">ISURIDE</h1>
      <ul className="mt-4 list-disc ps-8">
        <li>
          <Link to="/client" className="text-blue-600 hover:underline">
            Client Application
          </Link>
        </li>
        <li>
          <Link to="/simulator" className="text-blue-600 hover:underline">
            Chair Simulator
          </Link>
        </li>
        <li>
          <Link to="/owner" className="text-blue-600 hover:underline">
            Owner Application
          </Link>
        </li>
      </ul>
    </main>
  );
}
