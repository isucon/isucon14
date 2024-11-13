export const fetchWithCustomCookie = <T extends () => any>(
  fetcher: T,
  cookieInfo: {
    key: string;
    value: string;
    path: string;
  },
): ReturnType<T> => {
  const { key, value, path } = cookieInfo;
  document.cookie = `${key}=${value}; path=${path}`;
  return fetcher();
};
