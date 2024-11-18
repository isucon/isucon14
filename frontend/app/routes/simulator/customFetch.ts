export const fetchWithCustomCookie = async <T>(
  fetcher: () => Promise<T>,
  cookieInfo: {
    key: string;
    value: string;
    path: string;
  },
): Promise<T> => {
  const { key, value, path } = cookieInfo;
  document.cookie = `${key}=${value}; path=${path}`;
  return fetcher();
};
