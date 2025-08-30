import { vi } from "vitest";
import { api, ApiError } from "./api";

// Mock fetch globally
const mockFetch = vi.fn();
global.fetch = mockFetch;

describe("API Service", () => {
  beforeEach(() => {
    mockFetch.mockReset();
    vi.stubEnv("VITE_BACKEND_URL", "http://localhost:8080");
  });

  afterEach(() => {
    vi.unstubAllEnvs();
  });

  describe("CSRF Token Management", () => {
    it("fetches CSRF token successfully", async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({ token: "test-csrf-token" }),
      });

      await api.initialize();

      expect(mockFetch).toHaveBeenCalledWith(
        "http://localhost:8080/api/csrf-token",
        {
          credentials: "include",
        },
      );
    });

    it("handles CSRF token fetch failure", async () => {
      mockFetch.mockRejectedValueOnce(new Error("Network error"));

      await expect(api.initialize()).resolves.not.toThrow();
    });
  });

  describe("Request Handling", () => {
    it("includes CSRF token in non-GET requests", async () => {
      // First mock CSRF token fetch
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({ token: "test-csrf-token" }),
      });

      await api.initialize();

      // Mock successful POST request
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: () => Promise.resolve({ success: true }),
      });

      await api.post("/api/test", { data: "test" });

      expect(mockFetch).toHaveBeenCalledWith(
        "http://localhost:8080/api/test",
        expect.objectContaining({
          method: "POST",
          headers: expect.objectContaining({
            "X-CSRF-Token": "test-csrf-token",
          }),
        }),
      );
    });

    it("handles CSRF token errors by retrying", async () => {
      // Mock CSRF token fetch
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({ token: "test-csrf-token" }),
      });

      await api.initialize();

      // First request fails with CSRF error
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 403,
        json: () =>
          Promise.resolve({
            code: "csrf_token_invalid",
            error: "CSRF token invalid",
          }),
      });

      // Second CSRF token fetch
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({ token: "new-csrf-token" }),
      });

      // Second request succeeds
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: () => Promise.resolve({ success: true }),
      });

      const result = await api.post("/api/test", { data: "test" });

      expect(result).toEqual({ success: true });
      expect(mockFetch).toHaveBeenCalledTimes(4);
    });

    it("handles token expiration by refreshing", async () => {
      // Mock initial request failing with token expiration
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 401,
        json: () =>
          Promise.resolve({
            code: "token_expired",
            error: "Token expired",
          }),
      });

      // Mock refresh token request
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
      });

      // Mock CSRF token refresh
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({ token: "refreshed-csrf-token" }),
      });

      // Mock retried request succeeding
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: () => Promise.resolve({ success: true }),
      });

      const result = await api.post("/api/test", { data: "test" });

      expect(result).toEqual({ success: true });
      expect(mockFetch).toHaveBeenCalledTimes(4);
    });
  });

  describe("Error Handling", () => {
    it("throws ApiError for HTTP errors", async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 404,
        statusText: "Not Found",
        json: () => Promise.reject(new Error("JSON parse error")),
      });

      await expect(api.get("/api/not-found")).rejects.toThrow(ApiError);
    });

    it("includes status code in ApiError", async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 500,
        json: () => Promise.resolve({ error: "Internal Server Error" }),
      });

      try {
        await api.get("/api/error");
        expect.fail("Should have thrown");
      } catch (error) {
        expect(error).toBeInstanceOf(ApiError);
        expect((error as ApiError).status).toBe(500);
      }
    });
  });

  describe("HTTP Methods", () => {
    it("supports GET requests", async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({ data: "test" }),
      });

      const result = await api.get("/api/data");
      expect(result).toEqual({ data: "test" });
    });

    it("supports POST requests", async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({ created: true }),
      });

      const result = await api.post("/api/data", { name: "test" });
      expect(result).toEqual({ created: true });
    });

    it("supports PUT requests", async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({ updated: true }),
      });

      const result = await api.put("/api/data/1", { name: "updated" });
      expect(result).toEqual({ updated: true });
    });

    it("supports DELETE requests", async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 204,
      });

      const result = await api.delete("/api/data/1");
      expect(result).toBeUndefined();
    });

    it("supports PATCH requests", async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({ patched: true }),
      });

      const result = await api.patch("/api/data/1", { name: "patched" });
      expect(result).toEqual({ patched: true });
    });
  });
});
