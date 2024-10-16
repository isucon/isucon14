/**
 * Generated by @openapi-codegen
 *
 * @version 1.0
 */
import * as reactQuery from "@tanstack/react-query";
import { useApiContext, ApiContext } from "./apiContext";
import type * as Fetcher from "./apiFetcher";
import { apiFetch } from "./apiFetcher";
import type * as Schemas from "./apiSchemas";

export type PostInitializeError = Fetcher.ErrorWrapper<undefined>;

export type PostInitializeResponse = {
  /**
   * 実装言語
   */
  language: "go" | "perl" | "php" | "python" | "ruby" | "rust" | "node";
};

export type PostInitializeRequestBody = {
  /**
   * 決済サーバーアドレス
   */
  payment_server: string;
};

export type PostInitializeVariables = {
  body: PostInitializeRequestBody;
} & ApiContext["fetcherOptions"];

export const fetchPostInitialize = (
  variables: PostInitializeVariables,
  signal?: AbortSignal,
) =>
  apiFetch<
    PostInitializeResponse,
    PostInitializeError,
    PostInitializeRequestBody,
    {},
    {},
    {}
  >({ url: "/initialize", method: "post", ...variables, signal });

export const usePostInitialize = (
  options?: Omit<
    reactQuery.UseMutationOptions<
      PostInitializeResponse,
      PostInitializeError,
      PostInitializeVariables
    >,
    "mutationFn"
  >,
) => {
  const { fetcherOptions } = useApiContext();
  return reactQuery.useMutation<
    PostInitializeResponse,
    PostInitializeError,
    PostInitializeVariables
  >({
    mutationFn: (variables: PostInitializeVariables) =>
      fetchPostInitialize({ ...fetcherOptions, ...variables }),
    ...options,
  });
};

export type AppPostRegisterError = Fetcher.ErrorWrapper<{
  status: 400;
  payload: Schemas.Error;
}>;

export type AppPostRegisterResponse = {
  /**
   * アクセストークン
   */
  access_token: string;
  /**
   * ユーザーID
   */
  id: string;
};

export type AppPostRegisterRequestBody = {
  /**
   * ユーザー名
   */
  username: string;
  /**
   * 名前
   */
  firstname: string;
  /**
   * 名字
   */
  lastname: string;
  /**
   * 生年月日
   */
  date_of_birth: string;
};

export type AppPostRegisterVariables = {
  body: AppPostRegisterRequestBody;
} & ApiContext["fetcherOptions"];

export const fetchAppPostRegister = (
  variables: AppPostRegisterVariables,
  signal?: AbortSignal,
) =>
  apiFetch<
    AppPostRegisterResponse,
    AppPostRegisterError,
    AppPostRegisterRequestBody,
    {},
    {},
    {}
  >({ url: "/app/register", method: "post", ...variables, signal });

export const useAppPostRegister = (
  options?: Omit<
    reactQuery.UseMutationOptions<
      AppPostRegisterResponse,
      AppPostRegisterError,
      AppPostRegisterVariables
    >,
    "mutationFn"
  >,
) => {
  const { fetcherOptions } = useApiContext();
  return reactQuery.useMutation<
    AppPostRegisterResponse,
    AppPostRegisterError,
    AppPostRegisterVariables
  >({
    mutationFn: (variables: AppPostRegisterVariables) =>
      fetchAppPostRegister({ ...fetcherOptions, ...variables }),
    ...options,
  });
};

export type AppPostPaymentMethodsError = Fetcher.ErrorWrapper<undefined>;

export type AppPostPaymentMethodsRequestBody = {
  /**
   * 決済トークン
   */
  token: string;
};

export type AppPostPaymentMethodsVariables = {
  body: AppPostPaymentMethodsRequestBody;
} & ApiContext["fetcherOptions"];

export const fetchAppPostPaymentMethods = (
  variables: AppPostPaymentMethodsVariables,
  signal?: AbortSignal,
) =>
  apiFetch<
    undefined,
    AppPostPaymentMethodsError,
    AppPostPaymentMethodsRequestBody,
    {},
    {},
    {}
  >({ url: "/app/payment-methods", method: "post", ...variables, signal });

export const useAppPostPaymentMethods = (
  options?: Omit<
    reactQuery.UseMutationOptions<
      undefined,
      AppPostPaymentMethodsError,
      AppPostPaymentMethodsVariables
    >,
    "mutationFn"
  >,
) => {
  const { fetcherOptions } = useApiContext();
  return reactQuery.useMutation<
    undefined,
    AppPostPaymentMethodsError,
    AppPostPaymentMethodsVariables
  >({
    mutationFn: (variables: AppPostPaymentMethodsVariables) =>
      fetchAppPostPaymentMethods({ ...fetcherOptions, ...variables }),
    ...options,
  });
};

export type AppPostRequestError = Fetcher.ErrorWrapper<undefined>;

export type AppPostRequestResponse = {
  /**
   * 配車要求ID
   */
  request_id: string;
};

export type AppPostRequestRequestBody = {
  /**
   * 配車位置
   */
  pickup_coordinate: Schemas.Coordinate;
  /**
   * 目的地
   */
  destination_coordinate: Schemas.Coordinate;
};

export type AppPostRequestVariables = {
  body: AppPostRequestRequestBody;
} & ApiContext["fetcherOptions"];

export const fetchAppPostRequest = (
  variables: AppPostRequestVariables,
  signal?: AbortSignal,
) =>
  apiFetch<
    AppPostRequestResponse,
    AppPostRequestError,
    AppPostRequestRequestBody,
    {},
    {},
    {}
  >({ url: "/app/requests", method: "post", ...variables, signal });

export const useAppPostRequest = (
  options?: Omit<
    reactQuery.UseMutationOptions<
      AppPostRequestResponse,
      AppPostRequestError,
      AppPostRequestVariables
    >,
    "mutationFn"
  >,
) => {
  const { fetcherOptions } = useApiContext();
  return reactQuery.useMutation<
    AppPostRequestResponse,
    AppPostRequestError,
    AppPostRequestVariables
  >({
    mutationFn: (variables: AppPostRequestVariables) =>
      fetchAppPostRequest({ ...fetcherOptions, ...variables }),
    ...options,
  });
};

export type AppGetRequestPathParams = {
  /**
   * 配車要求ID
   */
  requestId: string;
};

export type AppGetRequestError = Fetcher.ErrorWrapper<{
  status: 404;
  payload: Schemas.Error;
}>;

export type AppGetRequestVariables = {
  pathParams: AppGetRequestPathParams;
} & ApiContext["fetcherOptions"];

export const fetchAppGetRequest = (
  variables: AppGetRequestVariables,
  signal?: AbortSignal,
) =>
  apiFetch<
    Schemas.AppRequest,
    AppGetRequestError,
    undefined,
    {},
    {},
    AppGetRequestPathParams
  >({ url: "/app/requests/{requestId}", method: "get", ...variables, signal });

export const useAppGetRequest = <TData = Schemas.AppRequest,>(
  variables: AppGetRequestVariables,
  options?: Omit<
    reactQuery.UseQueryOptions<Schemas.AppRequest, AppGetRequestError, TData>,
    "queryKey" | "queryFn" | "initialData"
  >,
) => {
  const { fetcherOptions, queryOptions, queryKeyFn } = useApiContext(options);
  return reactQuery.useQuery<Schemas.AppRequest, AppGetRequestError, TData>({
    queryKey: queryKeyFn({
      path: "/app/requests/{requestId}",
      operationId: "appGetRequest",
      variables,
    }),
    queryFn: ({ signal }) =>
      fetchAppGetRequest({ ...fetcherOptions, ...variables }, signal),
    ...options,
    ...queryOptions,
  });
};

export type AppPostRequestEvaluatePathParams = {
  /**
   * 配車要求ID
   */
  requestId: string;
};

export type AppPostRequestEvaluateError = Fetcher.ErrorWrapper<
  | {
      status: 400;
      payload: Schemas.Error;
    }
  | {
      status: 404;
      payload: Schemas.Error;
    }
>;

export type AppPostRequestEvaluateRequestBody = {
  /**
   * 椅子の評価
   *
   * @minimum 1
   * @maximum 5
   */
  evaluation: number;
};

export type AppPostRequestEvaluateVariables = {
  body: AppPostRequestEvaluateRequestBody;
  pathParams: AppPostRequestEvaluatePathParams;
} & ApiContext["fetcherOptions"];

export const fetchAppPostRequestEvaluate = (
  variables: AppPostRequestEvaluateVariables,
  signal?: AbortSignal,
) =>
  apiFetch<
    undefined,
    AppPostRequestEvaluateError,
    AppPostRequestEvaluateRequestBody,
    {},
    {},
    AppPostRequestEvaluatePathParams
  >({
    url: "/app/requests/{requestId}/evaluate",
    method: "post",
    ...variables,
    signal,
  });

export const useAppPostRequestEvaluate = (
  options?: Omit<
    reactQuery.UseMutationOptions<
      undefined,
      AppPostRequestEvaluateError,
      AppPostRequestEvaluateVariables
    >,
    "mutationFn"
  >,
) => {
  const { fetcherOptions } = useApiContext();
  return reactQuery.useMutation<
    undefined,
    AppPostRequestEvaluateError,
    AppPostRequestEvaluateVariables
  >({
    mutationFn: (variables: AppPostRequestEvaluateVariables) =>
      fetchAppPostRequestEvaluate({ ...fetcherOptions, ...variables }),
    ...options,
  });
};

export type AppGetNotificationError = Fetcher.ErrorWrapper<undefined>;

export type AppGetNotificationVariables = ApiContext["fetcherOptions"];

/**
 * 最新の自分の配車要求を取得します。
 */
export const fetchAppGetNotification = (
  variables: AppGetNotificationVariables,
  signal?: AbortSignal,
) =>
  apiFetch<Schemas.AppRequest, AppGetNotificationError, undefined, {}, {}, {}>({
    url: "/app/notification",
    method: "get",
    ...variables,
    signal,
  });

/**
 * 最新の自分の配車要求を取得します。
 */
export const useAppGetNotification = <TData = Schemas.AppRequest,>(
  variables: AppGetNotificationVariables,
  options?: Omit<
    reactQuery.UseQueryOptions<
      Schemas.AppRequest,
      AppGetNotificationError,
      TData
    >,
    "queryKey" | "queryFn" | "initialData"
  >,
) => {
  const { fetcherOptions, queryOptions, queryKeyFn } = useApiContext(options);
  return reactQuery.useQuery<
    Schemas.AppRequest,
    AppGetNotificationError,
    TData
  >({
    queryKey: queryKeyFn({
      path: "/app/notification",
      operationId: "appGetNotification",
      variables,
    }),
    queryFn: ({ signal }) =>
      fetchAppGetNotification({ ...fetcherOptions, ...variables }, signal),
    ...options,
    ...queryOptions,
  });
};

export type ProviderPostRegisterError = Fetcher.ErrorWrapper<undefined>;

export type ProviderPostRegisterResponse = {
  /**
   * アクセストークン
   */
  access_token: string;
  /**
   * プロバイダーID
   */
  id: string;
};

export type ProviderPostRegisterRequestBody = {
  /**
   * プロバイダー名
   */
  name: string;
};

export type ProviderPostRegisterVariables = {
  body: ProviderPostRegisterRequestBody;
} & ApiContext["fetcherOptions"];

export const fetchProviderPostRegister = (
  variables: ProviderPostRegisterVariables,
  signal?: AbortSignal,
) =>
  apiFetch<
    ProviderPostRegisterResponse,
    ProviderPostRegisterError,
    ProviderPostRegisterRequestBody,
    {},
    {},
    {}
  >({ url: "/provider/register", method: "post", ...variables, signal });

export const useProviderPostRegister = (
  options?: Omit<
    reactQuery.UseMutationOptions<
      ProviderPostRegisterResponse,
      ProviderPostRegisterError,
      ProviderPostRegisterVariables
    >,
    "mutationFn"
  >,
) => {
  const { fetcherOptions } = useApiContext();
  return reactQuery.useMutation<
    ProviderPostRegisterResponse,
    ProviderPostRegisterError,
    ProviderPostRegisterVariables
  >({
    mutationFn: (variables: ProviderPostRegisterVariables) =>
      fetchProviderPostRegister({ ...fetcherOptions, ...variables }),
    ...options,
  });
};

export type ProviderGetSalesQueryParams = {
  /**
   * 開始日（含む）
   */
  since: string;
  /**
   * 終了日（含む）
   */
  until: string;
};

export type ProviderGetSalesError = Fetcher.ErrorWrapper<undefined>;

export type ProviderGetSalesResponse = {
  /**
   * プロバイダー全体の売上
   */
  total_sales?: number;
  /**
   * 椅子ごとの売上情報
   */
  chairs?: {
    /**
     * 椅子ID
     */
    id?: string;
    /**
     * 椅子名
     */
    name?: string;
    /**
     * 椅子ごとの売上
     */
    sales?: number;
  }[];
  /**
   * モデルごとの売上情報
   */
  models?: {
    /**
     * 椅子モデル
     */
    model?: string;
    /**
     * モデルごとの売上
     */
    sales?: number;
  }[];
};

export type ProviderGetSalesVariables = {
  queryParams: ProviderGetSalesQueryParams;
} & ApiContext["fetcherOptions"];

export const fetchProviderGetSales = (
  variables: ProviderGetSalesVariables,
  signal?: AbortSignal,
) =>
  apiFetch<
    ProviderGetSalesResponse,
    ProviderGetSalesError,
    undefined,
    {},
    ProviderGetSalesQueryParams,
    {}
  >({ url: "/provider/sales", method: "get", ...variables, signal });

export const useProviderGetSales = <TData = ProviderGetSalesResponse,>(
  variables: ProviderGetSalesVariables,
  options?: Omit<
    reactQuery.UseQueryOptions<
      ProviderGetSalesResponse,
      ProviderGetSalesError,
      TData
    >,
    "queryKey" | "queryFn" | "initialData"
  >,
) => {
  const { fetcherOptions, queryOptions, queryKeyFn } = useApiContext(options);
  return reactQuery.useQuery<
    ProviderGetSalesResponse,
    ProviderGetSalesError,
    TData
  >({
    queryKey: queryKeyFn({
      path: "/provider/sales",
      operationId: "providerGetSales",
      variables,
    }),
    queryFn: ({ signal }) =>
      fetchProviderGetSales({ ...fetcherOptions, ...variables }, signal),
    ...options,
    ...queryOptions,
  });
};

export type ChairPostRegisterError = Fetcher.ErrorWrapper<undefined>;

export type ChairPostRegisterResponse = {
  /**
   * アクセストークン
   */
  access_token: string;
  /**
   * 椅子ID
   */
  id: string;
};

export type ChairPostRegisterRequestBody = {
  /**
   * 椅子の名前
   */
  name: string;
  /**
   * 椅子のモデル
   */
  model: string;
};

export type ChairPostRegisterVariables = {
  body: ChairPostRegisterRequestBody;
} & ApiContext["fetcherOptions"];

export const fetchChairPostRegister = (
  variables: ChairPostRegisterVariables,
  signal?: AbortSignal,
) =>
  apiFetch<
    ChairPostRegisterResponse,
    ChairPostRegisterError,
    ChairPostRegisterRequestBody,
    {},
    {},
    {}
  >({ url: "/chair/register", method: "post", ...variables, signal });

export const useChairPostRegister = (
  options?: Omit<
    reactQuery.UseMutationOptions<
      ChairPostRegisterResponse,
      ChairPostRegisterError,
      ChairPostRegisterVariables
    >,
    "mutationFn"
  >,
) => {
  const { fetcherOptions } = useApiContext();
  return reactQuery.useMutation<
    ChairPostRegisterResponse,
    ChairPostRegisterError,
    ChairPostRegisterVariables
  >({
    mutationFn: (variables: ChairPostRegisterVariables) =>
      fetchChairPostRegister({ ...fetcherOptions, ...variables }),
    ...options,
  });
};

export type ChairPostActivateError = Fetcher.ErrorWrapper<undefined>;

export type ChairPostActivateVariables = {
  body?: Record<string, any>;
} & ApiContext["fetcherOptions"];

export const fetchChairPostActivate = (
  variables: ChairPostActivateVariables,
  signal?: AbortSignal,
) =>
  apiFetch<undefined, ChairPostActivateError, Record<string, any>, {}, {}, {}>({
    url: "/chair/activate",
    method: "post",
    ...variables,
    signal,
  });

export const useChairPostActivate = (
  options?: Omit<
    reactQuery.UseMutationOptions<
      undefined,
      ChairPostActivateError,
      ChairPostActivateVariables
    >,
    "mutationFn"
  >,
) => {
  const { fetcherOptions } = useApiContext();
  return reactQuery.useMutation<
    undefined,
    ChairPostActivateError,
    ChairPostActivateVariables
  >({
    mutationFn: (variables: ChairPostActivateVariables) =>
      fetchChairPostActivate({ ...fetcherOptions, ...variables }),
    ...options,
  });
};

export type ChairPostDeactivateError = Fetcher.ErrorWrapper<undefined>;

export type ChairPostDeactivateVariables = {
  body?: Record<string, any>;
} & ApiContext["fetcherOptions"];

export const fetchChairPostDeactivate = (
  variables: ChairPostDeactivateVariables,
  signal?: AbortSignal,
) =>
  apiFetch<
    undefined,
    ChairPostDeactivateError,
    Record<string, any>,
    {},
    {},
    {}
  >({ url: "/chair/deactivate", method: "post", ...variables, signal });

export const useChairPostDeactivate = (
  options?: Omit<
    reactQuery.UseMutationOptions<
      undefined,
      ChairPostDeactivateError,
      ChairPostDeactivateVariables
    >,
    "mutationFn"
  >,
) => {
  const { fetcherOptions } = useApiContext();
  return reactQuery.useMutation<
    undefined,
    ChairPostDeactivateError,
    ChairPostDeactivateVariables
  >({
    mutationFn: (variables: ChairPostDeactivateVariables) =>
      fetchChairPostDeactivate({ ...fetcherOptions, ...variables }),
    ...options,
  });
};

export type ChairPostCoordinateError = Fetcher.ErrorWrapper<undefined>;

export type ChairPostCoordinateVariables = {
  body: Schemas.Coordinate;
} & ApiContext["fetcherOptions"];

export const fetchChairPostCoordinate = (
  variables: ChairPostCoordinateVariables,
  signal?: AbortSignal,
) =>
  apiFetch<undefined, ChairPostCoordinateError, Schemas.Coordinate, {}, {}, {}>(
    { url: "/chair/coordinate", method: "post", ...variables, signal },
  );

export const useChairPostCoordinate = (
  options?: Omit<
    reactQuery.UseMutationOptions<
      undefined,
      ChairPostCoordinateError,
      ChairPostCoordinateVariables
    >,
    "mutationFn"
  >,
) => {
  const { fetcherOptions } = useApiContext();
  return reactQuery.useMutation<
    undefined,
    ChairPostCoordinateError,
    ChairPostCoordinateVariables
  >({
    mutationFn: (variables: ChairPostCoordinateVariables) =>
      fetchChairPostCoordinate({ ...fetcherOptions, ...variables }),
    ...options,
  });
};

export type ChairGetRequestPathParams = {
  /**
   * 配車要求ID
   */
  requestId: string;
};

export type ChairGetRequestError = Fetcher.ErrorWrapper<{
  status: 404;
  payload: Schemas.Error;
}>;

export type ChairGetRequestVariables = {
  pathParams: ChairGetRequestPathParams;
} & ApiContext["fetcherOptions"];

/**
 * 椅子向け通知エンドポイントから通知されたidの情報を取得する想定
 */
export const fetchChairGetRequest = (
  variables: ChairGetRequestVariables,
  signal?: AbortSignal,
) =>
  apiFetch<
    Schemas.ChairRequest,
    ChairGetRequestError,
    undefined,
    {},
    {},
    ChairGetRequestPathParams
  >({
    url: "/chair/requests/{requestId}",
    method: "get",
    ...variables,
    signal,
  });

/**
 * 椅子向け通知エンドポイントから通知されたidの情報を取得する想定
 */
export const useChairGetRequest = <TData = Schemas.ChairRequest,>(
  variables: ChairGetRequestVariables,
  options?: Omit<
    reactQuery.UseQueryOptions<
      Schemas.ChairRequest,
      ChairGetRequestError,
      TData
    >,
    "queryKey" | "queryFn" | "initialData"
  >,
) => {
  const { fetcherOptions, queryOptions, queryKeyFn } = useApiContext(options);
  return reactQuery.useQuery<Schemas.ChairRequest, ChairGetRequestError, TData>(
    {
      queryKey: queryKeyFn({
        path: "/chair/requests/{requestId}",
        operationId: "chairGetRequest",
        variables,
      }),
      queryFn: ({ signal }) =>
        fetchChairGetRequest({ ...fetcherOptions, ...variables }, signal),
      ...options,
      ...queryOptions,
    },
  );
};

export type ChairPostRequestAcceptPathParams = {
  /**
   * 配車要求ID
   */
  requestId: string;
};

export type ChairPostRequestAcceptError = Fetcher.ErrorWrapper<{
  status: 404;
  payload: Schemas.Error;
}>;

export type ChairPostRequestAcceptVariables = {
  pathParams: ChairPostRequestAcceptPathParams;
} & ApiContext["fetcherOptions"];

export const fetchChairPostRequestAccept = (
  variables: ChairPostRequestAcceptVariables,
  signal?: AbortSignal,
) =>
  apiFetch<
    undefined,
    ChairPostRequestAcceptError,
    undefined,
    {},
    {},
    ChairPostRequestAcceptPathParams
  >({
    url: "/chair/requests/{requestId}/accept",
    method: "post",
    ...variables,
    signal,
  });

export const useChairPostRequestAccept = (
  options?: Omit<
    reactQuery.UseMutationOptions<
      undefined,
      ChairPostRequestAcceptError,
      ChairPostRequestAcceptVariables
    >,
    "mutationFn"
  >,
) => {
  const { fetcherOptions } = useApiContext();
  return reactQuery.useMutation<
    undefined,
    ChairPostRequestAcceptError,
    ChairPostRequestAcceptVariables
  >({
    mutationFn: (variables: ChairPostRequestAcceptVariables) =>
      fetchChairPostRequestAccept({ ...fetcherOptions, ...variables }),
    ...options,
  });
};

export type ChairPostRequestDenyPathParams = {
  /**
   * 配車要求ID
   */
  requestId: string;
};

export type ChairPostRequestDenyError = Fetcher.ErrorWrapper<{
  status: 404;
  payload: Schemas.Error;
}>;

export type ChairPostRequestDenyVariables = {
  pathParams: ChairPostRequestDenyPathParams;
} & ApiContext["fetcherOptions"];

export const fetchChairPostRequestDeny = (
  variables: ChairPostRequestDenyVariables,
  signal?: AbortSignal,
) =>
  apiFetch<
    undefined,
    ChairPostRequestDenyError,
    undefined,
    {},
    {},
    ChairPostRequestDenyPathParams
  >({
    url: "/chair/requests/{requestId}/deny",
    method: "post",
    ...variables,
    signal,
  });

export const useChairPostRequestDeny = (
  options?: Omit<
    reactQuery.UseMutationOptions<
      undefined,
      ChairPostRequestDenyError,
      ChairPostRequestDenyVariables
    >,
    "mutationFn"
  >,
) => {
  const { fetcherOptions } = useApiContext();
  return reactQuery.useMutation<
    undefined,
    ChairPostRequestDenyError,
    ChairPostRequestDenyVariables
  >({
    mutationFn: (variables: ChairPostRequestDenyVariables) =>
      fetchChairPostRequestDeny({ ...fetcherOptions, ...variables }),
    ...options,
  });
};

export type ChairPostRequestDepartPathParams = {
  /**
   * 配車要求ID
   */
  requestId: string;
};

export type ChairPostRequestDepartError = Fetcher.ErrorWrapper<
  | {
      status: 400;
      payload: Schemas.Error;
    }
  | {
      status: 404;
      payload: Schemas.Error;
    }
>;

export type ChairPostRequestDepartVariables = {
  pathParams: ChairPostRequestDepartPathParams;
} & ApiContext["fetcherOptions"];

export const fetchChairPostRequestDepart = (
  variables: ChairPostRequestDepartVariables,
  signal?: AbortSignal,
) =>
  apiFetch<
    undefined,
    ChairPostRequestDepartError,
    undefined,
    {},
    {},
    ChairPostRequestDepartPathParams
  >({
    url: "/chair/requests/{requestId}/depart",
    method: "post",
    ...variables,
    signal,
  });

export const useChairPostRequestDepart = (
  options?: Omit<
    reactQuery.UseMutationOptions<
      undefined,
      ChairPostRequestDepartError,
      ChairPostRequestDepartVariables
    >,
    "mutationFn"
  >,
) => {
  const { fetcherOptions } = useApiContext();
  return reactQuery.useMutation<
    undefined,
    ChairPostRequestDepartError,
    ChairPostRequestDepartVariables
  >({
    mutationFn: (variables: ChairPostRequestDepartVariables) =>
      fetchChairPostRequestDepart({ ...fetcherOptions, ...variables }),
    ...options,
  });
};

export type ChairGetNotificationError = Fetcher.ErrorWrapper<undefined>;

export type ChairGetNotificationVariables = ApiContext["fetcherOptions"];

/**
 * 椅子に配車要求を通知するなどで使う想定
 */
export const fetchChairGetNotification = (
  variables: ChairGetNotificationVariables,
  signal?: AbortSignal,
) =>
  apiFetch<
    Schemas.ChairRequest,
    ChairGetNotificationError,
    undefined,
    {},
    {},
    {}
  >({ url: "/chair/notification", method: "get", ...variables, signal });

/**
 * 椅子に配車要求を通知するなどで使う想定
 */
export const useChairGetNotification = <TData = Schemas.ChairRequest,>(
  variables: ChairGetNotificationVariables,
  options?: Omit<
    reactQuery.UseQueryOptions<
      Schemas.ChairRequest,
      ChairGetNotificationError,
      TData
    >,
    "queryKey" | "queryFn" | "initialData"
  >,
) => {
  const { fetcherOptions, queryOptions, queryKeyFn } = useApiContext(options);
  return reactQuery.useQuery<
    Schemas.ChairRequest,
    ChairGetNotificationError,
    TData
  >({
    queryKey: queryKeyFn({
      path: "/chair/notification",
      operationId: "chairGetNotification",
      variables,
    }),
    queryFn: ({ signal }) =>
      fetchChairGetNotification({ ...fetcherOptions, ...variables }, signal),
    ...options,
    ...queryOptions,
  });
};

export type QueryOperation =
  | {
      path: "/app/requests/{requestId}";
      operationId: "appGetRequest";
      variables: AppGetRequestVariables;
    }
  | {
      path: "/app/notification";
      operationId: "appGetNotification";
      variables: AppGetNotificationVariables;
    }
  | {
      path: "/provider/sales";
      operationId: "providerGetSales";
      variables: ProviderGetSalesVariables;
    }
  | {
      path: "/chair/requests/{requestId}";
      operationId: "chairGetRequest";
      variables: ChairGetRequestVariables;
    }
  | {
      path: "/chair/notification";
      operationId: "chairGetNotification";
      variables: ChairGetNotificationVariables;
    };
