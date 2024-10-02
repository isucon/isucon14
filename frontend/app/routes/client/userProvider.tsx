import { ReactNode, createContext, useContext } from "react";

export type AccessToken = string;

/**
 * フロント側で利用するクライアント情報
 */
type User = {
  id: string;
  name: string;
  accessToken: AccessToken;
};

const userContext = createContext<Partial<User>>({});

export const UserProvider = ({
  children,
  accessToken,
}: {
  children: ReactNode;
  accessToken: string;
}) => {
  /**
   * openapi上にfetchするものがないので一旦仮置き
   * 想定では、ここで通信を行い子供に流す。
   */
  const fetchedValue = { id: "fetched-id", name: "fetched-name", accessToken };

  return (
    <userContext.Provider value={fetchedValue}>{children}</userContext.Provider>
  );
};

export const useClient = () => useContext(userContext);
