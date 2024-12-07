import type { MetaFunction } from "@remix-run/node";
import { useNavigate } from "@remix-run/react";
import { useEffect, useState } from "react";
import colors from "tailwindcss/colors";
import { AccountSwitchIcon } from "~/components/icon/account-switch";
import { Button } from "~/components/primitives/button/button";
import { Text } from "~/components/primitives/text/text";
import { getCookieValue } from "~/utils/get-cookie-value";

export const meta: MetaFunction = () => {
  return [
    { title: "User | ISURIDE" },
    { name: "description", content: "ユーザーページ" },
  ];
};

export default function Index() {
  const navigate = useNavigate();
  const [sessionToken, setSessionToken] = useState<string>();

  useEffect(() => {
    const token = getCookieValue(document.cookie, "app_session");
    if (token) {
      setSessionToken(token);
    }
  }, []);

  return (
    <section className="mx-8 flex-1">
      <h2 className="text-xl my-6">ユーザー</h2>
      <div className="mb-4 border-t pt-4">
        <Text bold className="mb-1">
          セッショントークン
        </Text>
        <Text size="sm" className="text-neutral-500">
          {sessionToken}
        </Text>
      </div>
      <div className="mb-4 border-t pt-4">
        <Button
          className="w-full flex items-center justify-center "
          onClick={() => navigate("/client/login")}
        >
          <AccountSwitchIcon className="me-1" fill={colors.neutral[600]} />
          ユーザーを切り替える
        </Button>
      </div>
    </section>
  );
}
