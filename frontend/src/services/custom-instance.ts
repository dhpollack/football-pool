import Axios, { AxiosRequestConfig } from "axios";

const API_BASE_URL = import.meta.env.VITE_BACKEND_URL;
export const AXIOS_INSTANCE = Axios.create({ baseURL: API_BASE_URL });

export const customInstance = async <T>(
  config: AxiosRequestConfig,
  options?: AxiosRequestConfig,
): Promise<T> => {
  const source = Axios.CancelToken.source();
  const promise = AXIOS_INSTANCE({
    ...config,
    ...options,
    cancelToken: source.token,
  }).then(({ data }) => data);
  // @ts-ignore
  promise.cancel = () => {
    source.cancel("Query was cancelled");
  };
  return promise;
};
