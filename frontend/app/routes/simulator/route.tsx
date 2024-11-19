import { Outlet } from "@remix-run/react";

import { SimulatorProvider } from "~/contexts/simulator-context";

export default function Layout() {
  return (
    <SimulatorProvider>
      <Outlet />
    </SimulatorProvider>
  );
}
