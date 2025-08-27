// E2E Test Configuration
export const E2E_CONFIG = {
  // Base URLs
  FRONTEND_URL: "http://localhost:5173",
  BACKEND_URL: "http://localhost:18080",
  
  // Test credentials
  ADMIN_CREDENTIALS: {
    email: "admin@test.com",
    password: "adminpassword",
  },
  
  // Test user template
  TEST_USER_TEMPLATE: {
    name: "E2E Test User",
    password: "SecurePassword123!",
  },
  
  // Timeouts
  TIMEOUTS: {
    PAGE_LOAD: 10000,
    ACTION: 5000,
    NAVIGATION: 3000,
    ASSERTION: 3000,
  },
  
  // Selectors
  SELECTORS: {
    LOGIN: {
      EMAIL: "input[name='email']",
      PASSWORD: "input[name='password']",
      SUBMIT: "button[type='submit']",
      ERROR: "[data-testid='error-message']",
    },
    REGISTER: {
      NAME: "input[name='name']",
      EMAIL: "input[name='email']",
      PASSWORD: "input[name='password']",
      SUBMIT: "button[type='submit']",
      ERROR: ".MuiFormHelperText-root",
    },
    COMMON: {
      LOADING: "[data-testid='loading']",
      SUCCESS: "[data-testid='success-message']",
      HEADING: "h4",
    },
  },
  
  // Routes
  ROUTES: {
    HOME: "/",
    LOGIN: "/login",
    REGISTER: "/register",
    PROFILE: "/profile",
    PICKS: "/picks",
    PROTECTED: ["/profile", "/picks", "/results"],
    PUBLIC: ["/login", "/register", "/"],
  },
};

export type E2EConfig = typeof E2E_CONFIG;