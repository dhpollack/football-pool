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

  // Timeouts - optimized for speed
  TIMEOUTS: {
    PAGE_LOAD: 5000,
    ACTION: 3000,
    NAVIGATION: 2000,
    ASSERTION: 2000,
  },

  // Selectors
  SELECTORS: {
    LOGIN: {
      EMAIL: "input[name='email']",
      PASSWORD: "input[name='password']",
      SUBMIT: "button[type='submit']",
      ERROR: "[data-testid='error-message']",
      BUTTON: "[data-testid='login-button']", // Login button in header
    },
    REGISTER: {
      NAME: "input[name='name']",
      EMAIL: "input[name='email']",
      PASSWORD: "input[name='password']",
      SUBMIT: "button[type='submit']",
      ERROR: ".MuiFormHelperText-root",
    },
    USER_MENU: {
      BUTTON: "[data-testid='user-menu-button']",
      MENU: "[data-testid='user-menu']",
      PROFILE_LINK: "[data-testid='profile-link']",
      LOGOUT_BUTTON: "[data-testid='logout-button']",
    },
    COMMON: {
      LOADING: "[data-testid='loading']",
      SUCCESS: "[data-testid='success-message']",
      HEADING: "h4",
    },
    ADMIN: {
      NAVIGATION: {
        USERS_LINK: "a[href*='/users']",
        GAMES_LINK: "a[href*='/games']",
        PICKS_LINK: "a[href*='/picks']",
      },
      USERS: {
        SEARCH_INPUT: "input[placeholder*='Search']",
        USER_ROW: ".MuiTableRow-root",
        USER_ID: "td:nth-child(1)", // ID column
        USER_EMAIL: "td:nth-child(2)", // Email column
        USER_NAME: "td:nth-child(3)", // Name column
        USER_ROLE: "td:nth-child(4) .MuiChip-root", // Role column with chip
        USER_PICKS: "td:nth-child(5)", // Picks column
        USER_WINS: "td:nth-child(6)", // Wins column
        USER_JOINED: "td:nth-child(7)", // Joined column
        USER_ACTIONS: "td:nth-child(8)", // Actions column
        EDIT_BUTTON: "button svg[data-testid='EditIcon']",
        DELETE_BUTTON: "button svg[data-testid='DeleteIcon']",
        DETAILS_LINK: "a[href*='/users/']",
      },
      USER_DETAILS: {
        EMAIL: "h6:has(+ p:contains('Email')) + p",
        NAME: "h6:has(+ p:contains('Name')) + p",
        ROLE: "h6:has(+ p:contains('Role')) + p .MuiChip-root",
        PLAYER_NAME: "h6:has(+ p:contains('Player Name')) + p",
        PLAYER_ADDRESS: "h6:has(+ p:contains('Address')) + p",
        MEMBER_SINCE: "h6:has(+ p:contains('Member Since')) + p",
        PICKS_TAB: "button[role='tab']:has-text('Picks History')",
        STATS_TAB: "button[role='tab']:has-text('Statistics')",
      },
      GAMES: {
        SEARCH_INPUT: "input[placeholder*='Search']",
        GAME_ROW: ".MuiTableRow-root",
        GAME_ID: "td:nth-child(1)", // ID column
        GAME_WEEK: "td:nth-child(2)", // Week column
        GAME_SEASON: "td:nth-child(3)", // Season column
        GAME_MATCHUP: "td:nth-child(4)", // Matchup column
        GAME_SPREAD: "td:nth-child(5)", // Spread column
        GAME_START_TIME: "td:nth-child(6)", // Start Time column
        GAME_ACTIONS: "td:nth-child(7)", // Actions column
        CREATE_BUTTON: "button:has-text('Add Game')",
        EDIT_BUTTON: "button svg[data-testid='EditIcon']",
        DELETE_BUTTON: "button svg[data-testid='DeleteIcon']",
      },
      GAME_FORM: {
        FAVORITE_TEAM: "input[data-testid='favorite-team-input']",
        UNDERDOG_TEAM: "input[data-testid='underdog-team-input']",
        SPREAD: "input[data-testid='spread-input']",
        WEEK: "input[data-testid='week-input']",
        SEASON: "input[data-testid='season-input']",
        START_TIME: "input[data-testid='start-time-input']",
        SUBMIT_BUTTON: "button:has-text('Submit')",
        CANCEL_BUTTON: "button:has-text('Cancel')",
      },
      RESULT_FORM: {
        OUTCOME: "#outcome",
        SUBMIT_BUTTON: "button:has-text('Submit')",
        CANCEL_BUTTON: "button:has-text('Cancel')",
      },
      PICKS: {
        SEARCH_INPUT: "input[placeholder*='Search']",
        PICK_ROW: ".MuiTableRow-root:not(thead .MuiTableRow-root):not(:has-text('No picks available'))",
        PICK_ID: "td:nth-child(1)", // ID column
        PICK_USER: "td:nth-child(2)", // User column
        PICK_GAME: "td:nth-child(3)", // Game column
        PICK_WEEK: "td:nth-child(4)", // Week column
        PICK_SEASON: "td:nth-child(5)", // Season column
        PICK_CHOICE: "td:nth-child(6) .MuiChip-root", // Pick choice column with chip
        PICK_SUBMITTED: "td:nth-child(7)", // Submitted column
        PICK_ACTIONS: "td:nth-child(8)", // Actions column
        DELETE_BUTTON: "button svg[data-testid='DeleteIcon']",
        FILTER_USER: "#user-filter",
        FILTER_WEEK: "#week-filter",
        FILTER_SEASON: "#season-filter",
      },
    },
  },

  // Routes
  ROUTES: {
    HOME: "/",
    LOGIN: "/login",
    REGISTER: "/register",
    PROFILE: "/profile",
    PICKS: "/picks",
    ADMIN: "/admin",
    ADMIN_USERS: "/admin/users",
    ADMIN_GAMES: "/admin/games",
    ADMIN_PICKS: "/admin/picks",
    ADMIN_USER_DETAILS: "/admin/users/",
    PROTECTED: [
      "/profile",
      "/picks",
      "/results",
      "/admin",
      "/admin/users",
      "/admin/games",
      "/admin/picks",
    ],
    PUBLIC: ["/login", "/register", "/"],
  },
};

export type E2EConfig = typeof E2E_CONFIG;
