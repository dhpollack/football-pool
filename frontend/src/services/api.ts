const API_BASE_URL = import.meta.env.VITE_BACKEND_URL;

// CSRF token management
let csrfToken: string | null = null;

async function fetchCsrfToken(): Promise<void> {
  try {
    const response = await fetch(`${API_BASE_URL}/api/csrf-token`, {
      credentials: "include",
    });

    if (response.ok) {
      const data = await response.json();
      csrfToken = data.token || data.csrfToken;
    }
  } catch (error) {
    console.warn("Failed to fetch CSRF token:", error);
  }
}

export interface ApiErrorDetails {
  error?: string;
  message?: string;
  code?: string;
  [key: string]: unknown;
}

export type ApiData = Record<string, unknown> | unknown[];

export class ApiError extends Error {
  status?: number;
  details?: ApiErrorDetails;

  constructor(message: string, status?: number, details?: ApiErrorDetails) {
    super(message);
    this.status = status;
    this.details = details;
    this.name = "ApiError";
  }
}

export const api = {
  async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    const url = `${API_BASE_URL}${endpoint}`;

    // Add CSRF token for non-GET requests
    const headers: Record<string, string> = {
      "Content-Type": "application/json",
    };

    // Add any custom headers from options
    if (options.headers) {
      if (options.headers instanceof Headers) {
        options.headers.forEach((value, key) => {
          headers[key] = value;
        });
      } else if (Array.isArray(options.headers)) {
        options.headers.forEach(([key, value]) => {
          headers[key] = value;
        });
      } else {
        Object.entries(options.headers).forEach(([key, value]) => {
          headers[key] = value;
        });
      }
    }

    const method = options.method?.toUpperCase() || "GET";
    if (method !== "GET" && csrfToken) {
      headers["X-CSRF-Token"] = csrfToken;
    }

    const config: RequestInit = {
      ...options,
      credentials: "include", // Include cookies in requests
      headers,
    };

    try {
      const response = await fetch(url, config);

      if (!response.ok) {
        let errorMessage = `HTTP error ${response.status}`;
        let errorData: ApiErrorDetails;

        try {
          errorData = await response.json();
          errorMessage = errorData.error || errorData.message || errorMessage;
        } catch {
          // If we can't parse JSON, use the status text
          errorMessage = response.statusText || errorMessage;
        }

        // Handle CSRF token errors by fetching a new token and retrying
        if (
          response.status === 403 &&
          errorData?.code === "csrf_token_invalid"
        ) {
          await fetchCsrfToken();
          // Retry the request with new CSRF token
          return this.request<T>(endpoint, options);
        }

        // Handle token expiration by refreshing and retrying
        if (response.status === 401 && errorData?.code === "token_expired") {
          try {
            // Attempt to refresh the token
            await this.refreshToken();
            // Retry the original request
            return this.request<T>(endpoint, options);
          } catch (refreshError) {
            // If refresh fails, throw the original error
            throw new ApiError(errorMessage, response.status);
          }
        }

        throw new ApiError(errorMessage, response.status);
      }

      // For responses with no content
      if (response.status === 204) {
        return undefined as T;
      }

      return await response.json();
    } catch (error) {
      if (error instanceof ApiError) {
        throw error;
      }
      throw new ApiError(
        error instanceof Error ? error.message : "Network error",
      );
    }
  },

  async get<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, { method: "GET" });
  },

  async post<T>(endpoint: string, data?: ApiData): Promise<T> {
    return this.request<T>(endpoint, {
      method: "POST",
      body: data ? JSON.stringify(data) : undefined,
    });
  },

  async put<T>(endpoint: string, data?: ApiData): Promise<T> {
    return this.request<T>(endpoint, {
      method: "PUT",
      body: data ? JSON.stringify(data) : undefined,
    });
  },

  async delete<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, { method: "DELETE" });
  },

  async patch<T>(endpoint: string, data?: ApiData): Promise<T> {
    return this.request<T>(endpoint, {
      method: "PATCH",
      body: data ? JSON.stringify(data) : undefined,
    });
  },

  // Initialize CSRF token
  async initialize(): Promise<void> {
    await fetchCsrfToken();
  },

  // Refresh authentication token
  async refreshToken(): Promise<void> {
    try {
      const response = await fetch(`${API_BASE_URL}/api/refresh`, {
        method: "POST",
        credentials: "include",
      });

      if (!response.ok) {
        throw new Error("Token refresh failed");
      }

      // Token refreshed successfully, CSRF token might also need refreshing
      await fetchCsrfToken();
    } catch (error) {
      console.error("Token refresh failed:", error);
      throw error;
    }
  },
};

export default api;
