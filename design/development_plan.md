Based on the design documents, here is a plan to create the Football Pool web application:

### Technology Stack

*   [x] **Backend:** Go with the standard `net/http` package for the API and GORM with SQLite for the database.
*   [ ] **Frontend:** React.

### Backend Development Plan

1.  [x] **Database Schema:** I will design the database schema based on the requirements. This will include tables for:
    *   `users` (id, name, email, password_hash, role)
    *   `players` (id, user_id, name, address)
    *   `games` (id, week, season, favorite_team, underdog_team, spread)
    *   `picks` (id, user_id, game_id, picked_team, rank, quick_pick)
    *   `results` (id, game_id, favorite_score, underdog_score, outcome)
    *   `survivor_picks` (id, user_id, week, team)

2.  [ ] **API Endpoints:** I will create the following RESTful API endpoints:
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

1.  [x] **Project Setup:** I will set up a new React project using `create-react-app`.

2.  [x] **Component Library:** I will use a component library like Material-UI to build the user interface.

3.  [x] **Routing:** I will use React Router for navigation between pages.

4.  [x] **Page Implementation:** I will create the following pages:
    *   [x] **Login Page:** A form to log in.
    *   [x] **Home Page:** The main landing page with a navigation drawer.
    *   [x] **Pick Entry/View Page:** A form to enter and view weekly picks. This will include the "Quick Pick" functionality.
    *   [x] **Result Entry Page (Admin):** A form for administrators to enter game results.
    *   [x] **Weekly Results Page:** A table displaying the weekly leaderboard.
    *   [ ] **Overall Results Page:** A table displaying the overall season leaderboard.
    *   [x] **Profile Page:** A page for users to view and edit their profile.

### Implementation Steps

1.  [x] I will start by setting up the Go backend project structure.
2.  [x] Then, I will define the database models and create the database schema.
3.  [ ] Next, I will implement the API endpoints one by one, starting with authentication.
4.  [ ] Once the backend is partially complete, I will set up the React frontend project.
5.  [ ] I will then build the frontend pages and components, connecting them to the backend API as I go.
6.  [ ] I will add the "Quick Pick" and survivor pool functionality.
7.  [ ] Finally, I will implement the admin-only features for entering results.