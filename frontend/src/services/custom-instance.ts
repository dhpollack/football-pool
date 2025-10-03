import Axios, {
  type AxiosRequestConfig,
  type InternalAxiosRequestConfig,
  type AxiosResponse,
} from "axios";

export const AXIOS_INSTANCE = Axios.create({
  baseURL: import.meta.env.VITE_BACKEND_URL || "/",
});

// Add a request interceptor
AXIOS_INSTANCE.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    try {
      // Read from react-auth-kit v4 alpha internal storage format
      const authToken = localStorage.getItem("_auth_auth");
      if (authToken) {
        // Extract token from format: "timestamp^&*&^token"
        const token = authToken.split("^&*&^")[1];
        if (token) {
          config.headers.Authorization = `Bearer ${token}`;
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

// Add a response interceptor to handle 401 responses
AXIOS_INSTANCE.interceptors.response.use(
  (response: AxiosResponse) => {
    return response;
  },
  (error) => {
    if (error.response?.status === 401) {
      // Authentication expired - let react-auth-kit handle the logout
      console.log("Authentication expired, redirecting to login");

      // Redirect to login page if we're in a browser environment
      // The react-auth-kit will automatically clear the auth state
      // when the user navigates to the login page
      if (
        typeof window !== "undefined" &&
        !window.location.pathname.includes("/login")
      ) {
        window.location.href = "/login";
      }
    }
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
