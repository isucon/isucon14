import type { MetaFunction } from "@remix-run/node";
import { ClientActionFunctionArgs, Form, redirect } from "@remix-run/react";
import { fetchChairPostChairs } from "~/apiClient/apiComponents";
import { TextInput } from "~/components/primitives/form/text";

export const meta: MetaFunction = () => {
  return [
    { title: "Regiter | ISURIDE" },
    { name: "description", content: "チェア登録" },
  ];
};
export const clientAction = async ({ request }: ClientActionFunctionArgs) => {
  const formData = await request.formData();
  // const provider = await fetchProviderPostRegister({
  //   body: {
  //     name: String(formData.get("provider_name")) ?? "",
  //   },
  // });
  const chair = await fetchChairPostChairs({
    body: {
      model: String(formData.get("model")) ?? "",
      name: String(formData.get("name")) ?? "",
      chair_register_token: String(formData.get("chair_register_token")) ?? "",
    },
  });
  return redirect(`/driver?id=${chair.id}`);
};

export default function DriverRegister() {
  return (
    <>
      <Form
        className="w-full h-full p-4 bg-white rounded-lg shadow-md flex flex-col gap-4"
        method="POST"
      >
        <div>
          <TextInput
            id="chair_register_token"
            name="chair_register_token"
            label="chair_register_token:"
            required
          />
          <TextInput id="chair_name" name="name" label="Chair name:" required />
          <TextInput
            id="chair_model"
            name="model"
            label="Chair model:"
            required
          />
        </div>
        <button type="submit">登録</button>
      </Form>
    </>
  );
}
