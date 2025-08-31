import "@testing-library/jest-dom";

// Only override for unit tests (Jest), not for E2E tests (Playwright)
if (typeof process !== 'undefined' && process.env.JEST_WORKER_ID) {
  import.meta.env.VITE_BACKEND_URL = "http://localhost:8080";
}

window.alert = () => {};
