import Axios, {
  type AxiosRequestConfig,
  type InternalAxiosRequestConfig,
} from "axios";

export const AXIOS_INSTANCE = Axios.create({
  baseURL: import.meta.env.VITE_BACKEND_URL || "/",
});

// Add a request interceptor
AXIOS_INSTANCE.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    try {
      const authData = localStorage.getItem("_auth_auth");
      if (authData) {
        const authToken = authData.split("^&*&^")[1];
        if (authToken) {
          config.headers.Authorization = `Bearer ${authToken}`;
        }
      }
    } catch (error) {
      console.error("Error reading auth token from localStorage", error);
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  },
);

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

  // @ts-expect-error
  promise.cancel = () => {
    source.cancel("Query was cancelled");
  };
  return promise;
};
