import "@testing-library/jest-dom";

// Only override for unit tests (Jest), not for E2E tests (Playwright)
if (typeof process !== "undefined" && process.env.JEST_WORKER_ID) {
  const { server } = await import("./mocks/server");
  beforeAll(() => server.listen());
  afterEach(() => server.resetHandlers());
  afterAll(() => server.close());

  import.meta.env.VITE_BACKEND_URL = "http://localhost:18080";
}

window.alert = () => {};
