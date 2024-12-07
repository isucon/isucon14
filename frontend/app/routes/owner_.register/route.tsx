import type { MetaFunction } from "@remix-run/node";
import { Link, useNavigate } from "@remix-run/react";
import { useState } from "react";
import { fetchOwnerPostOwners } from "~/api/api-components";
import { Button } from "~/components/primitives/button/button";
import { TextInput } from "~/components/primitives/form/text-input";
import { FormFrame } from "~/components/primitives/frame/form-frame";
import { Text } from "~/components/primitives/text/text";

export const meta: MetaFunction = () => {
  return [
    { title: "Regiter | Owner | ISURIDE" },
    { name: "description", content: "オーナー登録" },
  ];
};

export default function OwnerRegister() {
  const [ownerName, setOwnerName] = useState<string>();
  const [errorMessage, setErrorMessage] = useState<string>();
  const navigate = useNavigate();

  const handleOnClick = async () => {
    try {
      if (!ownerName) {
        setErrorMessage("オーナー名を入力してください");
        return;
      }
      if (ownerName.length > 30) {
        setErrorMessage("オーナー名は30文字以内で入力してください");
        return;
      }
      await fetchOwnerPostOwners({
        body: {
          name: ownerName ?? "",
        },
      });
      navigate("/owner");
    } catch (e) {
      console.error(`ERROR: ${JSON.stringify(e)}`);
      setErrorMessage(
        "オーナーの登録に失敗しました。接続に問題があるか、ユーザー名が登録済みの可能性があります。",
      );
    }
  };

  return (
    <FormFrame>
      <div className="mb-8">
        <h1 className="text-2xl font-semibold">オーナー登録</h1>
        {errorMessage && (
          <Text variant="danger" className="mt-2">
            {errorMessage}
          </Text>
        )}
      </div>
      <div className="flex flex-col gap-8">
        <div>
          <TextInput
            id="ownerName"
            name="ownerName"
            label="オーナー名"
            onChange={(e) => setOwnerName(e.target.value)}
          />
        </div>
        <Button
          variant="primary"
          className="text-lg mt-6"
          onClick={() => void handleOnClick()}
        >
          登録
        </Button>
        <p className="text-center">
          <Link to="/owner/login" className="text-blue-600 hover:underline">
            ログイン
          </Link>
        </p>
      </div>
    </FormFrame>
  );
}
