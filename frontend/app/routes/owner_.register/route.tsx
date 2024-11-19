import type { MetaFunction } from "@remix-run/node";
import OwnerLoginForm from "./form/login";
import OwnerRegisterForm from "./form/register";

export const meta: MetaFunction = () => {
  return [
    { title: "Regiter | ISURIDE" },
    { name: "description", content: "オーナー登録" },
  ];
};

export default function ProviderRegister() {
  return (
    <>
      <OwnerLoginForm />
      <OwnerRegisterForm />
    </>
  );
}
