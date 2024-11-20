/**
 * Generated by @openapi-codegen
 *
 * @version 1.0
 */
import * as reactQuery from "@tanstack/react-query";
import { ApiContext, useApiContext } from "./apiContext";
import type * as Fetcher from "./apiFetcher";
import { apiFetch } from "./apiFetcher";
import type * as Schemas from "./apiSchemas";

export type PostInitializeError = Fetcher.ErrorWrapper<undefined>;

export type PostInitializeResponse = {
  /**
   * 実装言語
   * - go
   * - perl
   * - php
   * - python
   * - ruby
   * - rust
   * - node
   */
  language: string;
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

export type AppPostUsersError = Fetcher.ErrorWrapper<{
  status: 400;
  payload: Schemas.Error;
}>;

export type AppPostUsersResponse = {
  /**
   * ユーザーID
   */
  id: string;
  /**
   * 自分の招待コード
   */
  invitation_code: string;
};

export type AppPostUsersRequestBody = {
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
  /**
   * 他の人の招待コード
   */
  invitation_code?: string;
};

export type AppPostUsersVariables = {
  body: AppPostUsersRequestBody;
} & ApiContext["fetcherOptions"];

export const fetchAppPostUsers = (
  variables: AppPostUsersVariables,
  signal?: AbortSignal,
) =>
  apiFetch<
    AppPostUsersResponse,
    AppPostUsersError,
    AppPostUsersRequestBody,
    {},
    {},
    {}
  >({ url: "/app/users", method: "post", ...variables, signal });

export const useAppPostUsers = (
  options?: Omit<
    reactQuery.UseMutationOptions<
      AppPostUsersResponse,
      AppPostUsersError,
      AppPostUsersVariables
    >,
    "mutationFn"
  >,
) => {
  const { fetcherOptions } = useApiContext();
  return reactQuery.useMutation<
    AppPostUsersResponse,
    AppPostUsersError,
    AppPostUsersVariables
  >({
    mutationFn: (variables: AppPostUsersVariables) =>
      fetchAppPostUsers({ ...fetcherOptions, ...variables }),
    ...options,
  });
};

export type AppPostPaymentMethodsError = Fetcher.ErrorWrapper<{
  status: 400;
  payload: Schemas.Error;
}>;

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

export type AppGetRidesError = Fetcher.ErrorWrapper<undefined>;

export type AppGetRidesResponse = {
  rides: {
    /**
     * ライドID
     */
    id: string;
    pickup_coordinate: Schemas.Coordinate;
    destination_coordinate: Schemas.Coordinate;
    chair: {
      /**
       * 椅子ID
       */
      id: string;
      /**
       * オーナー名
       */
      owner: string;
      /**
       * 椅子の名前
       */
      name: string;
      /**
       * 椅子のモデル
       */
      model: string;
    };
    /**
     * 運賃
     */
    fare: number;
    /**
     * 椅子の評価
     */
    evaluation: number;
    /**
     * 配車要求日時
     *
     * @format int64
     */
    requested_at: number;
    /**
     * 評価まで完了した日時
     *
     * @format int64
     */
    completed_at: number;
  }[];
};

export type AppGetRidesVariables = ApiContext["fetcherOptions"];

export const fetchAppGetRides = (
  variables: AppGetRidesVariables,
  signal?: AbortSignal,
) =>
  apiFetch<AppGetRidesResponse, AppGetRidesError, undefined, {}, {}, {}>({
    url: "/app/rides",
    method: "get",
    ...variables,
    signal,
  });

export const useAppGetRides = <TData = AppGetRidesResponse,>(
  variables: AppGetRidesVariables,
  options?: Omit<
    reactQuery.UseQueryOptions<AppGetRidesResponse, AppGetRidesError, TData>,
    "queryKey" | "queryFn" | "initialData"
  >,
) => {
  const { fetcherOptions, queryOptions, queryKeyFn } = useApiContext(options);
  return reactQuery.useQuery<AppGetRidesResponse, AppGetRidesError, TData>({
    queryKey: queryKeyFn({
      path: "/app/rides",
      operationId: "appGetRides",
      variables,
    }),
    queryFn: ({ signal }) =>
      fetchAppGetRides({ ...fetcherOptions, ...variables }, signal),
    ...options,
    ...queryOptions,
  });
};

export type AppPostRidesError = Fetcher.ErrorWrapper<
  | {
      status: 400;
      payload: Schemas.Error;
    }
  | {
      status: 409;
      payload: Schemas.Error;
    }
>;

export type AppPostRidesResponse = {
  /**
   * ライドID
   */
  ride_id: string;
  /**
   * 割引後運賃
   */
  fare: number;
};

export type AppPostRidesRequestBody = {
  /**
   * 配車位置
   */
  pickup_coordinate: Schemas.Coordinate;
  /**
   * 目的地
   */
  destination_coordinate: Schemas.Coordinate;
};

export type AppPostRidesVariables = {
  body: AppPostRidesRequestBody;
} & ApiContext["fetcherOptions"];

export const fetchAppPostRides = (
  variables: AppPostRidesVariables,
  signal?: AbortSignal,
) =>
  apiFetch<
    AppPostRidesResponse,
    AppPostRidesError,
    AppPostRidesRequestBody,
    {},
    {},
    {}
  >({ url: "/app/rides", method: "post", ...variables, signal });

export const useAppPostRides = (
  options?: Omit<
    reactQuery.UseMutationOptions<
      AppPostRidesResponse,
      AppPostRidesError,
      AppPostRidesVariables
    >,
    "mutationFn"
  >,
) => {
  const { fetcherOptions } = useApiContext();
  return reactQuery.useMutation<
    AppPostRidesResponse,
    AppPostRidesError,
    AppPostRidesVariables
  >({
    mutationFn: (variables: AppPostRidesVariables) =>
      fetchAppPostRides({ ...fetcherOptions, ...variables }),
    ...options,
  });
};

export type AppPostRidesEstimatedFareError = Fetcher.ErrorWrapper<{
  status: 400;
  payload: Schemas.Error;
}>;

export type AppPostRidesEstimatedFareResponse = {
  /**
   * 割引後運賃
   */
  fare: number;
  /**
   * 割引額
   */
  discount: number;
};

export type AppPostRidesEstimatedFareRequestBody = {
  /**
   * 配車位置
   */
  pickup_coordinate: Schemas.Coordinate;
  /**
   * 目的地
   */
  destination_coordinate: Schemas.Coordinate;
};

export type AppPostRidesEstimatedFareVariables = {
  body: AppPostRidesEstimatedFareRequestBody;
} & ApiContext["fetcherOptions"];

export const fetchAppPostRidesEstimatedFare = (
  variables: AppPostRidesEstimatedFareVariables,
  signal?: AbortSignal,
) =>
  apiFetch<
    AppPostRidesEstimatedFareResponse,
    AppPostRidesEstimatedFareError,
    AppPostRidesEstimatedFareRequestBody,
    {},
    {},
    {}
  >({ url: "/app/rides/estimated-fare", method: "post", ...variables, signal });

export const useAppPostRidesEstimatedFare = (
  options?: Omit<
    reactQuery.UseMutationOptions<
      AppPostRidesEstimatedFareResponse,
      AppPostRidesEstimatedFareError,
      AppPostRidesEstimatedFareVariables
    >,
    "mutationFn"
  >,
) => {
  const { fetcherOptions } = useApiContext();
  return reactQuery.useMutation<
    AppPostRidesEstimatedFareResponse,
    AppPostRidesEstimatedFareError,
    AppPostRidesEstimatedFareVariables
  >({
    mutationFn: (variables: AppPostRidesEstimatedFareVariables) =>
      fetchAppPostRidesEstimatedFare({ ...fetcherOptions, ...variables }),
    ...options,
  });
};

export type AppGetRidePathParams = {
  /**
   * ライドID
   */
  rideId: string;
};

export type AppGetRideError = Fetcher.ErrorWrapper<{
  status: 404;
  payload: Schemas.Error;
}>;

export type AppGetRideVariables = {
  pathParams: AppGetRidePathParams;
} & ApiContext["fetcherOptions"];

export const fetchAppGetRide = (
  variables: AppGetRideVariables,
  signal?: AbortSignal,
) =>
  apiFetch<
    Schemas.AppRide,
    AppGetRideError,
    undefined,
    {},
    {},
    AppGetRidePathParams
  >({ url: "/app/rides/{rideId}", method: "get", ...variables, signal });

export const useAppGetRide = <TData = Schemas.AppRide,>(
  variables: AppGetRideVariables,
  options?: Omit<
    reactQuery.UseQueryOptions<Schemas.AppRide, AppGetRideError, TData>,
    "queryKey" | "queryFn" | "initialData"
  >,
) => {
  const { fetcherOptions, queryOptions, queryKeyFn } = useApiContext(options);
  return reactQuery.useQuery<Schemas.AppRide, AppGetRideError, TData>({
    queryKey: queryKeyFn({
      path: "/app/rides/{rideId}",
      operationId: "appGetRide",
      variables,
    }),
    queryFn: ({ signal }) =>
      fetchAppGetRide({ ...fetcherOptions, ...variables }, signal),
    ...options,
    ...queryOptions,
  });
};

export type AppPostRideEvaluationPathParams = {
  /**
   * ライドID
   */
  rideId: string;
};

export type AppPostRideEvaluationError = Fetcher.ErrorWrapper<
  | {
      status: 400;
      payload: Schemas.Error;
    }
  | {
      status: 404;
      payload: Schemas.Error;
    }
>;

export type AppPostRideEvaluationResponse = {
  /**
   * 割引後運賃
   */
  fare: number;
  /**
   * 完了日時
   *
   * @format int64
   */
  completed_at: number;
};

export type AppPostRideEvaluationRequestBody = {
  /**
   * ライドの評価
   *
   * @minimum 1
   * @maximum 5
   */
  evaluation: number;
};

export type AppPostRideEvaluationVariables = {
  body: AppPostRideEvaluationRequestBody;
  pathParams: AppPostRideEvaluationPathParams;
} & ApiContext["fetcherOptions"];

export const fetchAppPostRideEvaluation = (
  variables: AppPostRideEvaluationVariables,
  signal?: AbortSignal,
) =>
  apiFetch<
    AppPostRideEvaluationResponse,
    AppPostRideEvaluationError,
    AppPostRideEvaluationRequestBody,
    {},
    {},
    AppPostRideEvaluationPathParams
  >({
    url: "/app/rides/{rideId}/evaluation",
    method: "post",
    ...variables,
    signal,
  });

export const useAppPostRideEvaluation = (
  options?: Omit<
    reactQuery.UseMutationOptions<
      AppPostRideEvaluationResponse,
      AppPostRideEvaluationError,
      AppPostRideEvaluationVariables
    >,
    "mutationFn"
  >,
) => {
  const { fetcherOptions } = useApiContext();
  return reactQuery.useMutation<
    AppPostRideEvaluationResponse,
    AppPostRideEvaluationError,
    AppPostRideEvaluationVariables
  >({
    mutationFn: (variables: AppPostRideEvaluationVariables) =>
      fetchAppPostRideEvaluation({ ...fetcherOptions, ...variables }),
    ...options,
  });
};

export type AppGetNotificationError = Fetcher.ErrorWrapper<undefined>;

export type AppGetNotificationResponse = {
  /**
   * ライドID
   */
  ride_id: string;
  pickup_coordinate: Schemas.Coordinate;
  destination_coordinate: Schemas.Coordinate;
  status: Schemas.RideStatus;
  chair?: Schemas.AppChair;
  /**
   * 配車要求日時
   *
   * @format int64
   */
  created_at: number;
  /**
   * 配車要求更新日時
   *
   * @format int64
   */
  updated_at: number;
  /**
   * 次回の通知ポーリングまでの待機時間(ミリ秒単位)
   */
  retry_after_ms?: number;
};

export type AppGetNotificationVariables = ApiContext["fetcherOptions"];

/**
 * 最新の自分のライドを取得する
 */
export const fetchAppGetNotification = (
  variables: AppGetNotificationVariables,
  signal?: AbortSignal,
) =>
  apiFetch<
    AppGetNotificationResponse,
    AppGetNotificationError,
    undefined,
    {},
    {},
    {}
  >({ url: "/app/notification", method: "get", ...variables, signal });

/**
 * 最新の自分のライドを取得する
 */
export const useAppGetNotification = <TData = AppGetNotificationResponse,>(
  variables: AppGetNotificationVariables,
  options?: Omit<
    reactQuery.UseQueryOptions<
      AppGetNotificationResponse,
      AppGetNotificationError,
      TData
    >,
    "queryKey" | "queryFn" | "initialData"
  >,
) => {
  const { fetcherOptions, queryOptions, queryKeyFn } = useApiContext(options);
  return reactQuery.useQuery<
    AppGetNotificationResponse,
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

export type AppGetNearbyChairsQueryParams = {
  /**
   * 緯度
   */
  latitude: number;
  /**
   * 経度
   */
  longitude: number;
  /**
   * 検索距離
   */
  distance?: number;
};

export type AppGetNearbyChairsError = Fetcher.ErrorWrapper<undefined>;

export type AppGetNearbyChairsResponse = {
  chairs: Schemas.AppChair[];
  /**
   * 取得日時
   *
   * @format int64
   */
  retrieved_at: number;
};

export type AppGetNearbyChairsVariables = {
  queryParams: AppGetNearbyChairsQueryParams;
} & ApiContext["fetcherOptions"];

export const fetchAppGetNearbyChairs = (
  variables: AppGetNearbyChairsVariables,
  signal?: AbortSignal,
) =>
  apiFetch<
    AppGetNearbyChairsResponse,
    AppGetNearbyChairsError,
    undefined,
    {},
    AppGetNearbyChairsQueryParams,
    {}
  >({ url: "/app/nearby-chairs", method: "get", ...variables, signal });

export const useAppGetNearbyChairs = <TData = AppGetNearbyChairsResponse,>(
  variables: AppGetNearbyChairsVariables,
  options?: Omit<
    reactQuery.UseQueryOptions<
      AppGetNearbyChairsResponse,
      AppGetNearbyChairsError,
      TData
    >,
    "queryKey" | "queryFn" | "initialData"
  >,
) => {
  const { fetcherOptions, queryOptions, queryKeyFn } = useApiContext(options);
  return reactQuery.useQuery<
    AppGetNearbyChairsResponse,
    AppGetNearbyChairsError,
    TData
  >({
    queryKey: queryKeyFn({
      path: "/app/nearby-chairs",
      operationId: "appGetNearbyChairs",
      variables,
    }),
    queryFn: ({ signal }) =>
      fetchAppGetNearbyChairs({ ...fetcherOptions, ...variables }, signal),
    ...options,
    ...queryOptions,
  });
};

export type OwnerPostOwnersError = Fetcher.ErrorWrapper<{
  status: 400;
  payload: Schemas.Error;
}>;

export type OwnerPostOwnersResponse = {
  /**
   * オーナーID
   */
  id: string;
  /**
   * 椅子をオーナーに紐づけるための椅子登録用トークン
   */
  chair_register_token: string;
};

export type OwnerPostOwnersRequestBody = {
  /**
   * オーナー名
   */
  name: string;
};

export type OwnerPostOwnersVariables = {
  body: OwnerPostOwnersRequestBody;
} & ApiContext["fetcherOptions"];

export const fetchOwnerPostOwners = (
  variables: OwnerPostOwnersVariables,
  signal?: AbortSignal,
) =>
  apiFetch<
    OwnerPostOwnersResponse,
    OwnerPostOwnersError,
    OwnerPostOwnersRequestBody,
    {},
    {},
    {}
  >({ url: "/owner/owners", method: "post", ...variables, signal });

export const useOwnerPostOwners = (
  options?: Omit<
    reactQuery.UseMutationOptions<
      OwnerPostOwnersResponse,
      OwnerPostOwnersError,
      OwnerPostOwnersVariables
    >,
    "mutationFn"
  >,
) => {
  const { fetcherOptions } = useApiContext();
  return reactQuery.useMutation<
    OwnerPostOwnersResponse,
    OwnerPostOwnersError,
    OwnerPostOwnersVariables
  >({
    mutationFn: (variables: OwnerPostOwnersVariables) =>
      fetchOwnerPostOwners({ ...fetcherOptions, ...variables }),
    ...options,
  });
};

export type OwnerGetSalesQueryParams = {
  /**
   * 開始日時（含む）
   *
   * @format int64
   */
  since?: number;
  /**
   * 終了日時（含む）
   *
   * @format int64
   */
  until?: number;
};

export type OwnerGetSalesError = Fetcher.ErrorWrapper<undefined>;

export type OwnerGetSalesResponse = {
  /**
   * オーナーが管理する椅子全体の売上
   */
  total_sales: number;
  /**
   * 椅子ごとの売上情報
   */
  chairs: {
    /**
     * 椅子ID
     */
    id: string;
    /**
     * 椅子名
     */
    name: string;
    /**
     * 椅子ごとの売上
     */
    sales: number;
  }[];
  /**
   * モデルごとの売上情報
   */
  models: {
    /**
     * モデル
     */
    model: string;
    /**
     * モデルごとの売上
     */
    sales: number;
  }[];
};

export type OwnerGetSalesVariables = {
  queryParams?: OwnerGetSalesQueryParams;
} & ApiContext["fetcherOptions"];

export const fetchOwnerGetSales = (
  variables: OwnerGetSalesVariables,
  signal?: AbortSignal,
) =>
  apiFetch<
    OwnerGetSalesResponse,
    OwnerGetSalesError,
    undefined,
    {},
    OwnerGetSalesQueryParams,
    {}
  >({ url: "/owner/sales", method: "get", ...variables, signal });

export const useOwnerGetSales = <TData = OwnerGetSalesResponse,>(
  variables: OwnerGetSalesVariables,
  options?: Omit<
    reactQuery.UseQueryOptions<
      OwnerGetSalesResponse,
      OwnerGetSalesError,
      TData
    >,
    "queryKey" | "queryFn" | "initialData"
  >,
) => {
  const { fetcherOptions, queryOptions, queryKeyFn } = useApiContext(options);
  return reactQuery.useQuery<OwnerGetSalesResponse, OwnerGetSalesError, TData>({
    queryKey: queryKeyFn({
      path: "/owner/sales",
      operationId: "ownerGetSales",
      variables,
    }),
    queryFn: ({ signal }) =>
      fetchOwnerGetSales({ ...fetcherOptions, ...variables }, signal),
    ...options,
    ...queryOptions,
  });
};

export type OwnerGetChairsError = Fetcher.ErrorWrapper<undefined>;

export type OwnerGetChairsResponse = {
  chairs: {
    /**
     * 椅子ID
     */
    id: string;
    /**
     * 椅子の名前
     */
    name: string;
    /**
     * 椅子のモデル
     */
    model: string;
    /**
     * 稼働中かどうか
     */
    active: boolean;
    /**
     * 登録日時
     *
     * @format int64
     */
    registered_at: number;
    /**
     * 総移動距離
     */
    total_distance: number;
    /**
     * 総移動距離の更新日時
     *
     * @format int64
     */
    total_distance_updated_at?: number;
  }[];
};

export type OwnerGetChairsVariables = ApiContext["fetcherOptions"];

export const fetchOwnerGetChairs = (
  variables: OwnerGetChairsVariables,
  signal?: AbortSignal,
) =>
  apiFetch<OwnerGetChairsResponse, OwnerGetChairsError, undefined, {}, {}, {}>({
    url: "/owner/chairs",
    method: "get",
    ...variables,
    signal,
  });

export const useOwnerGetChairs = <TData = OwnerGetChairsResponse,>(
  variables: OwnerGetChairsVariables,
  options?: Omit<
    reactQuery.UseQueryOptions<
      OwnerGetChairsResponse,
      OwnerGetChairsError,
      TData
    >,
    "queryKey" | "queryFn" | "initialData"
  >,
) => {
  const { fetcherOptions, queryOptions, queryKeyFn } = useApiContext(options);
  return reactQuery.useQuery<
    OwnerGetChairsResponse,
    OwnerGetChairsError,
    TData
  >({
    queryKey: queryKeyFn({
      path: "/owner/chairs",
      operationId: "ownerGetChairs",
      variables,
    }),
    queryFn: ({ signal }) =>
      fetchOwnerGetChairs({ ...fetcherOptions, ...variables }, signal),
    ...options,
    ...queryOptions,
  });
};

export type OwnerGetChairPathParams = {
  /**
   * 椅子ID
   */
  chairId: string;
};

export type OwnerGetChairError = Fetcher.ErrorWrapper<undefined>;

export type OwnerGetChairResponse = {
  /**
   * 椅子ID
   */
  id: string;
  /**
   * 椅子の名前
   */
  name: string;
  /**
   * 椅子のモデル
   */
  model: string;
  /**
   * 稼働中かどうか
   */
  active: boolean;
  /**
   * 登録日時
   *
   * @format int64
   */
  registered_at: number;
  /**
   * 総移動距離
   */
  total_distance: number;
  /**
   * 総移動距離の更新日時
   *
   * @format int64
   */
  total_distance_updated_at?: number;
};

export type OwnerGetChairVariables = {
  pathParams: OwnerGetChairPathParams;
} & ApiContext["fetcherOptions"];

export const fetchOwnerGetChair = (
  variables: OwnerGetChairVariables,
  signal?: AbortSignal,
) =>
  apiFetch<
    OwnerGetChairResponse,
    OwnerGetChairError,
    undefined,
    {},
    {},
    OwnerGetChairPathParams
  >({ url: "/owner/chairs/{chairId}", method: "get", ...variables, signal });

export const useOwnerGetChair = <TData = OwnerGetChairResponse,>(
  variables: OwnerGetChairVariables,
  options?: Omit<
    reactQuery.UseQueryOptions<
      OwnerGetChairResponse,
      OwnerGetChairError,
      TData
    >,
    "queryKey" | "queryFn" | "initialData"
  >,
) => {
  const { fetcherOptions, queryOptions, queryKeyFn } = useApiContext(options);
  return reactQuery.useQuery<OwnerGetChairResponse, OwnerGetChairError, TData>({
    queryKey: queryKeyFn({
      path: "/owner/chairs/{chairId}",
      operationId: "ownerGetChair",
      variables,
    }),
    queryFn: ({ signal }) =>
      fetchOwnerGetChair({ ...fetcherOptions, ...variables }, signal),
    ...options,
    ...queryOptions,
  });
};

export type ChairPostChairsError = Fetcher.ErrorWrapper<undefined>;

export type ChairPostChairsResponse = {
  /**
   * 椅子ID
   */
  id: string;
  /**
   * オーナーID
   */
  owner_id: string;
};

export type ChairPostChairsRequestBody = {
  /**
   * 椅子の名前
   */
  name: string;
  /**
   * 椅子のモデル
   */
  model: string;
  /**
   * 椅子をオーナーに紐づけるための椅子登録用トークン
   */
  chair_register_token: string;
};

export type ChairPostChairsVariables = {
  body: ChairPostChairsRequestBody;
} & ApiContext["fetcherOptions"];

export const fetchChairPostChairs = (
  variables: ChairPostChairsVariables,
  signal?: AbortSignal,
) =>
  apiFetch<
    ChairPostChairsResponse,
    ChairPostChairsError,
    ChairPostChairsRequestBody,
    {},
    {},
    {}
  >({ url: "/chair/chairs", method: "post", ...variables, signal });

export const useChairPostChairs = (
  options?: Omit<
    reactQuery.UseMutationOptions<
      ChairPostChairsResponse,
      ChairPostChairsError,
      ChairPostChairsVariables
    >,
    "mutationFn"
  >,
) => {
  const { fetcherOptions } = useApiContext();
  return reactQuery.useMutation<
    ChairPostChairsResponse,
    ChairPostChairsError,
    ChairPostChairsVariables
  >({
    mutationFn: (variables: ChairPostChairsVariables) =>
      fetchChairPostChairs({ ...fetcherOptions, ...variables }),
    ...options,
  });
};

export type ChairPostActivityError = Fetcher.ErrorWrapper<undefined>;

export type ChairPostActivityRequestBody = {
  /**
   * 配車受付を開始するか停止するか
   */
  is_active: boolean;
};

export type ChairPostActivityVariables = {
  body: ChairPostActivityRequestBody;
} & ApiContext["fetcherOptions"];

export const fetchChairPostActivity = (
  variables: ChairPostActivityVariables,
  signal?: AbortSignal,
) =>
  apiFetch<
    undefined,
    ChairPostActivityError,
    ChairPostActivityRequestBody,
    {},
    {},
    {}
  >({ url: "/chair/activity", method: "post", ...variables, signal });

export const useChairPostActivity = (
  options?: Omit<
    reactQuery.UseMutationOptions<
      undefined,
      ChairPostActivityError,
      ChairPostActivityVariables
    >,
    "mutationFn"
  >,
) => {
  const { fetcherOptions } = useApiContext();
  return reactQuery.useMutation<
    undefined,
    ChairPostActivityError,
    ChairPostActivityVariables
  >({
    mutationFn: (variables: ChairPostActivityVariables) =>
      fetchChairPostActivity({ ...fetcherOptions, ...variables }),
    ...options,
  });
};

export type ChairPostCoordinateError = Fetcher.ErrorWrapper<undefined>;

export type ChairPostCoordinateResponse = {
  /**
   * 記録日時
   *
   * @format int64
   */
  recorded_at: number;
};

export type ChairPostCoordinateVariables = {
  body: Schemas.Coordinate;
} & ApiContext["fetcherOptions"];

export const fetchChairPostCoordinate = (
  variables: ChairPostCoordinateVariables,
  signal?: AbortSignal,
) =>
  apiFetch<
    ChairPostCoordinateResponse,
    ChairPostCoordinateError,
    Schemas.Coordinate,
    {},
    {},
    {}
  >({ url: "/chair/coordinate", method: "post", ...variables, signal });

export const useChairPostCoordinate = (
  options?: Omit<
    reactQuery.UseMutationOptions<
      ChairPostCoordinateResponse,
      ChairPostCoordinateError,
      ChairPostCoordinateVariables
    >,
    "mutationFn"
  >,
) => {
  const { fetcherOptions } = useApiContext();
  return reactQuery.useMutation<
    ChairPostCoordinateResponse,
    ChairPostCoordinateError,
    ChairPostCoordinateVariables
  >({
    mutationFn: (variables: ChairPostCoordinateVariables) =>
      fetchChairPostCoordinate({ ...fetcherOptions, ...variables }),
    ...options,
  });
};

export type ChairGetNotificationError = Fetcher.ErrorWrapper<undefined>;

export type ChairGetNotificationResponse = {
  /**
   * ライドID
   */
  ride_id: string;
  user: Schemas.User;
  pickup_coordinate: Schemas.Coordinate;
  destination_coordinate: Schemas.Coordinate;
  status: Schemas.RideStatus;
  /**
   * 次回の通知ポーリングまでの待機時間(ミリ秒単位)
   */
  retry_after_ms?: number;
};

export type ChairGetNotificationVariables = ApiContext["fetcherOptions"];

/**
 * 椅子に配車要求を通知するなどで使う想定
 */
export const fetchChairGetNotification = (
  variables: ChairGetNotificationVariables,
  signal?: AbortSignal,
) =>
  apiFetch<
    ChairGetNotificationResponse,
    ChairGetNotificationError,
    undefined,
    {},
    {},
    {}
  >({ url: "/chair/notification", method: "get", ...variables, signal });

/**
 * 椅子に配車要求を通知するなどで使う想定
 */
export const useChairGetNotification = <TData = ChairGetNotificationResponse,>(
  variables: ChairGetNotificationVariables,
  options?: Omit<
    reactQuery.UseQueryOptions<
      ChairGetNotificationResponse,
      ChairGetNotificationError,
      TData
    >,
    "queryKey" | "queryFn" | "initialData"
  >,
) => {
  const { fetcherOptions, queryOptions, queryKeyFn } = useApiContext(options);
  return reactQuery.useQuery<
    ChairGetNotificationResponse,
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

export type ChairGetRidePathParams = {
  /**
   * ライドID
   */
  rideId: string;
};

export type ChairGetRideError = Fetcher.ErrorWrapper<{
  status: 404;
  payload: Schemas.Error;
}>;

export type ChairGetRideVariables = {
  pathParams: ChairGetRidePathParams;
} & ApiContext["fetcherOptions"];

/**
 * 椅子向け通知エンドポイントから通知されたidの情報を取得する想定
 */
export const fetchChairGetRide = (
  variables: ChairGetRideVariables,
  signal?: AbortSignal,
) =>
  apiFetch<
    Schemas.ChairRide,
    ChairGetRideError,
    undefined,
    {},
    {},
    ChairGetRidePathParams
  >({ url: "/chair/rides/{rideId}", method: "get", ...variables, signal });

/**
 * 椅子向け通知エンドポイントから通知されたidの情報を取得する想定
 */
export const useChairGetRide = <TData = Schemas.ChairRide,>(
  variables: ChairGetRideVariables,
  options?: Omit<
    reactQuery.UseQueryOptions<Schemas.ChairRide, ChairGetRideError, TData>,
    "queryKey" | "queryFn" | "initialData"
  >,
) => {
  const { fetcherOptions, queryOptions, queryKeyFn } = useApiContext(options);
  return reactQuery.useQuery<Schemas.ChairRide, ChairGetRideError, TData>({
    queryKey: queryKeyFn({
      path: "/chair/rides/{rideId}",
      operationId: "chairGetRide",
      variables,
    }),
    queryFn: ({ signal }) =>
      fetchChairGetRide({ ...fetcherOptions, ...variables }, signal),
    ...options,
    ...queryOptions,
  });
};

export type ChairPostRideStatusPathParams = {
  /**
   * ライドID
   */
  rideId: string;
};

export type ChairPostRideStatusError = Fetcher.ErrorWrapper<{
  status: 404;
  payload: Schemas.Error;
}>;

export type ChairPostRideStatusRequestBody = {
  /**
   * ライドの状態
   * MATCHING: マッチングを拒否し、再度マッチング状態に戻す
   * ENROUTE: マッチングを承認し、乗車位置に向かう
   * CARRYING: ユーザーが乗車し、椅子が目的地に向かう
   */
  status: "MATCHING" | "ENROUTE" | "CARRYING";
};

export type ChairPostRideStatusVariables = {
  body: ChairPostRideStatusRequestBody;
  pathParams: ChairPostRideStatusPathParams;
} & ApiContext["fetcherOptions"];

export const fetchChairPostRideStatus = (
  variables: ChairPostRideStatusVariables,
  signal?: AbortSignal,
) =>
  apiFetch<
    undefined,
    ChairPostRideStatusError,
    ChairPostRideStatusRequestBody,
    {},
    {},
    ChairPostRideStatusPathParams
  >({
    url: "/chair/rides/{rideId}/status",
    method: "post",
    ...variables,
    signal,
  });

export const useChairPostRideStatus = (
  options?: Omit<
    reactQuery.UseMutationOptions<
      undefined,
      ChairPostRideStatusError,
      ChairPostRideStatusVariables
    >,
    "mutationFn"
  >,
) => {
  const { fetcherOptions } = useApiContext();
  return reactQuery.useMutation<
    undefined,
    ChairPostRideStatusError,
    ChairPostRideStatusVariables
  >({
    mutationFn: (variables: ChairPostRideStatusVariables) =>
      fetchChairPostRideStatus({ ...fetcherOptions, ...variables }),
    ...options,
  });
};

export type QueryOperation =
  | {
      path: "/app/rides";
      operationId: "appGetRides";
      variables: AppGetRidesVariables;
    }
  | {
      path: "/app/rides/{rideId}";
      operationId: "appGetRide";
      variables: AppGetRideVariables;
    }
  | {
      path: "/app/notification";
      operationId: "appGetNotification";
      variables: AppGetNotificationVariables;
    }
  | {
      path: "/app/nearby-chairs";
      operationId: "appGetNearbyChairs";
      variables: AppGetNearbyChairsVariables;
    }
  | {
      path: "/owner/sales";
      operationId: "ownerGetSales";
      variables: OwnerGetSalesVariables;
    }
  | {
      path: "/owner/chairs";
      operationId: "ownerGetChairs";
      variables: OwnerGetChairsVariables;
    }
  | {
      path: "/owner/chairs/{chairId}";
      operationId: "ownerGetChair";
      variables: OwnerGetChairVariables;
    }
  | {
      path: "/chair/notification";
      operationId: "chairGetNotification";
      variables: ChairGetNotificationVariables;
    }
  | {
      path: "/chair/rides/{rideId}";
      operationId: "chairGetRide";
      variables: ChairGetRideVariables;
    };
