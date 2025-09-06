import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import AuthProvider from "react-auth-kit/AuthProvider";
import { authStore } from "./lib/auth";
import Layout from "./components/Layout";
import AdminLayout from "./components/admin/AdminLayout";
import ErrorBoundary from "./components/ErrorBoundary";
import QueryErrorBoundary from "./components/QueryErrorBoundary";
import ProtectedRoute from "./components/ProtectedRoute";
import HomePage from "./pages/HomePage";
import LoginPage from "./pages/LoginPage";
import RegisterPage from "./pages/RegisterPage";
import ProfilePage from "./pages/ProfilePage";
import PickEntryPage from "./pages/PickEntryPage";
import ResultEntryPage from "./pages/ResultEntryPage";
import WeeklyResultsPage from "./pages/WeeklyResultsPage";
import OverallResultsPage from "./pages/OverallResultsPage";
import SurvivorPoolPage from "./pages/SurvivorPoolPage";
import AdminDashboard from "./pages/admin/AdminDashboard";
import AdminGamesPage from "./pages/admin/AdminGamesPage";
import AdminUsersPage from "./pages/admin/AdminUsersPage";
import AdminPicksPage from "./pages/admin/AdminPicksPage";
import UserDetailsPage from "./pages/admin/UserDetailsPage";
import WeeklyPicksPage from "./pages/admin/WeeklyPicksPage";
import UserPicksPage from "./pages/admin/UserPicksPage";

function App() {
  return (
    <ErrorBoundary>
      <QueryErrorBoundary>
        <AuthProvider store={authStore} fallbackPath="/">
          <Router>
            <Routes>
              <Route path="/" element={<Layout />}>
                <Route index element={<HomePage />} />
                <Route path="login" element={<LoginPage />} />
                <Route path="register" element={<RegisterPage />} />
                <Route
                  path="profile"
                  element={
                    <QueryErrorBoundary>
                      <ProfilePage />
                    </QueryErrorBoundary>
                  }
                />
                <Route
                  path="picks"
                  element={
                    <QueryErrorBoundary>
                      <PickEntryPage />
                    </QueryErrorBoundary>
                  }
                />
                <Route
                  path="results"
                  element={
                    <QueryErrorBoundary>
                      <ResultEntryPage />
                    </QueryErrorBoundary>
                  }
                />
                <Route
                  path="weekly-results"
                  element={
                    <QueryErrorBoundary>
                      <WeeklyResultsPage />
                    </QueryErrorBoundary>
                  }
                />
                <Route
                  path="overall-results"
                  element={
                    <QueryErrorBoundary>
                      <OverallResultsPage />
                    </QueryErrorBoundary>
                  }
                />
                <Route
                  path="survivor-pool"
                  element={
                    <QueryErrorBoundary>
                      <SurvivorPoolPage />
                    </QueryErrorBoundary>
                  }
                />
              </Route>
              <Route
                path="/admin"
                element={
                  <ProtectedRoute requireAdmin>
                    <AdminLayout />
                  </ProtectedRoute>
                }
              >
                <Route
                  index
                  element={
                    <QueryErrorBoundary>
                      <AdminDashboard />
                    </QueryErrorBoundary>
                  }
                />
                <Route
                  path="games"
                  element={
                    <QueryErrorBoundary>
                      <AdminGamesPage />
                    </QueryErrorBoundary>
                  }
                />
                <Route
                  path="users"
                  element={
                    <QueryErrorBoundary>
                      <AdminUsersPage />
                    </QueryErrorBoundary>
                  }
                />
                <Route
                  path="users/:userId"
                  element={
                    <QueryErrorBoundary>
                      <UserDetailsPage />
                    </QueryErrorBoundary>
                  }
                />
                <Route
                  path="picks"
                  element={
                    <QueryErrorBoundary>
                      <AdminPicksPage />
                    </QueryErrorBoundary>
                  }
                />
                <Route
                  path="picks/week/:week"
                  element={
                    <QueryErrorBoundary>
                      <WeeklyPicksPage />
                    </QueryErrorBoundary>
                  }
                />
                <Route
                  path="picks/user/:userId"
                  element={
                    <QueryErrorBoundary>
                      <UserPicksPage />
                    </QueryErrorBoundary>
                  }
                />
              </Route>
            </Routes>
          </Router>
        </AuthProvider>
      </QueryErrorBoundary>
    </ErrorBoundary>
  );
}

export default App;
