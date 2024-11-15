export const fetchWithCustomCookie = async <T extends () => Promise<any>>(
  fetcher: T,
  cookieInfo: {
    key: string;
    value: string;
    path: string;
  },
): Promise<ReturnType<T>> => {
  const { key, value, path } = cookieInfo;
  document.cookie = `${key}=${value}; path=${path}`;
  return fetcher();
};
