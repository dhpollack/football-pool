Based on the design documents, here is a plan to create the Football Pool web application:

### Technology Stack

*   [x] **Backend:** Go with the standard `net/http` package for the API and GORM with SQLite for the database.
*   [ ] **Frontend:** Vite.js + Tailwind CSS.

### Backend Development Plan

1.  [x] **Database Schema:** I will design the database schema based on the requirements. This will include tables for:
    *   `users` (id, name, email, password_hash, role)
    *   `players` (id, user_id, name, address)
    *   `games` (id, week, season, favorite_team, underdog_team, spread)
    *   `picks` (id, user_id, game_id, picked_team, rank, quick_pick)
    *   `results` (id, game_id, favorite_score, underdog_score, outcome)
    *   `survivor_picks` (id, user_id, week, team)

2.  [x] **API Endpoints:** I will create the following RESTful API endpoints:
    *   [x] **Authentication:**
        *   `POST /api/login`: Authenticate a user and return a JWT.
        *   `POST /api/logout`: Log out a user.
    *   [x] **Users:**
        *   `GET /api/users/me`: Get the current user's profile.
        *   `PUT /api/users/me`: Update the current user's profile.
    *   [x] **Games:**
        *   `GET /api/games`: Get a list of games for a given week and season.
    *   [x] **Picks:**
        *   `GET /api/picks`: Get the current user's picks for a week.
        *   `POST /api/picks`: Submit picks for the week.
    *   [x] **Results (Admin):**
        *   `POST /api/results`: Enter game results.
    *   **Leaderboards:**
        *   [x] `GET /api/results/week`: Get the weekly leaderboard.
        *   `GET /api/results/season`: Get the overall season leaderboard.
    *   **Survivor Pool:**
        *   `GET /api/survivor/picks`: Get the current user's survivor picks.
        *   `POST /api/survivor/picks`: Submit a survivor pick.

### Frontend Development Plan

1.  [x] **Project Setup:** I will set up a new Vite.js project with React and Tailwind CSS.
2.  [x] **Linting and Formatting:** I will use Biome for linting and auto-formatting.
3.  [x] **Component Library:** I will use a component library like Material-UI to build the user interface.
4.  [x] **Routing:** I will use React Router for navigation between pages.
5.  [x] **Layout:** I will implement a single column layout with a collapsible navbar on the left and a header with the site name and login/profile link.
6.  [x] **Page Implementation:** I will create the following pages:
    *   [x] **Homepage:** landing page of the site with a navigation drawer on the left and a login/profile at the top right
    *   [x] **Login Page:** If site visitor is not logged in, then they should be redirected to a login page
    *   [x] **Player Profile Page:** There should be a page for the player to add profile information such as name and home address
    *   [x] **Pick Entry / Pick View Page:** Players will enter their weekly picks for the weeks game.  There will be an favorite and an underdog and a spread for each game.  The player must pick either the favorite or the underdog and assign a rank value for the pick between 1 and N, where N is the total number of games for that week.  The rank for each game must be unique.  These entries should go into a database table with the game id, the player id, and the pick and the rank value. There should also be a "Quick Pick" option that selects a favorite or underdog at random and assigns a random rank.  The "Quick Pick" option needs to conform to the constraints of the entry, i.e. rank values must be unique.
    *   [x] **Result Entry Page (admin):** Administrators should be able to enter the scores for each game.  These should go into a database table with the game id, the favorite's score, the underdog's score and a result that is either the favorite wins (favorite - spread > underdog), the underdog wins (underdog > favorite - spread), or a tie (favorite - spread = underdog).
    *   [x] **Week's Results Page:** We should be able to view the results for all players for a single week.  This should show a table of the sum of the rank of the correctly selected results of the games of the players' in descending order.  The player should receive the rank value they assigned for correctly picked results and 0 points for incorrect picks.  For ties, all players should receive half the amount of points of the rank that they assigned to the game.  There should be an icon to indicate if the player used the "Quick Pick" option to make their entry.
    *   [x] **Overall Results Page:** Similar to the Week's Results, we want to see the overall results.  This should be the sum of the weekly results for a single season for all players.  The order of the results should be in descending order.

### Implementation Steps

1.  [x] I will start by setting up the Go backend project structure.
2.  [x] Then, I will define the database models and create the database schema.
3.  [x] Next, I will implement the API endpoints one by one, starting with authentication.
4.  [x] Once the backend is partially complete, I will set up the Vite.js + Tailwind CSS frontend project.
5.  [x] I will then build the frontend pages and components, connecting them to the backend API as I go.
6.  [ ] I will add the "Quick Pick" and survivor pool functionality.
7.  [ ] Finally, I will implement the admin-only features for entering results.

### Test Coverage

We will aim for high test coverage for both the backend and frontend to ensure code quality and reliability.

*   **Backend Test Coverage:**
    *   [x] Use Go's built-in testing tools to generate coverage reports.
    *   [x] Achieve at least 80% statement coverage for all `internal` packages.
    *   [x] Prioritize testing critical business logic and API endpoints.
        *   [x] Add tests for error handling in `auth` package (e.g., invalid JSON, database errors).
        *   [x] Add tests for error handling in `handlers` package (e.g., invalid parameters, database errors).
        *   [x] Add tests for `internal/database` package (e.g., connection errors, migration errors).
            *   [x] Test the `ConnectDB` function to ensure it can connect to the database.
            *   [x] Test the `ConnectTestDB` function to ensure it can connect to the test database.
            *   [x] Test the `Migrate` function to ensure it can create the database schema.
        *   [x] Add integration tests for `internal/server` package.
            *   [x] Test that the server can be started and that it listens on the correct port.
            *   [x] Test that the server can handle requests to the API endpoints.
*   **Frontend Test Coverage:**
    *   [ ] Use a testing framework like Jest and React Testing Library.
    *   [ ] Achieve at least 80% statement coverage for all React components and utility functions.
    *   [ ] Focus on testing user interactions and component rendering.