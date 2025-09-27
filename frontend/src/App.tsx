import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import { lazy, Suspense } from "react";
import AuthProvider from "react-auth-kit/AuthProvider";
import { authStore } from "./lib/auth";
import Layout from "./components/Layout";
import AdminLayout from "./components/admin/AdminLayout";
import ErrorBoundary from "./components/ErrorBoundary";
import QueryErrorBoundary from "./components/QueryErrorBoundary";
import ProtectedRoute from "./components/ProtectedRoute";
import LoadingSpinner from "./components/LoadingSpinner";

// Non-admin pages (lazy loaded)
const HomePage = lazy(() => import("./pages/HomePage"));
const LoginPage = lazy(() => import("./pages/LoginPage"));
const RegisterPage = lazy(() => import("./pages/RegisterPage"));
const ProfilePage = lazy(() => import("./pages/ProfilePage"));
const PickEntryPage = lazy(() => import("./pages/PickEntryPage"));
const ResultEntryPage = lazy(() => import("./pages/ResultEntryPage"));
const WeeklyResultsPage = lazy(() => import("./pages/WeeklyResultsPage"));
const OverallResultsPage = lazy(() => import("./pages/OverallResultsPage"));
const SurvivorPoolPage = lazy(() => import("./pages/SurvivorPoolPage"));

// Admin pages (lazy loaded)
const AdminDashboard = lazy(() => import("./pages/admin/AdminDashboard"));
const AdminGamesPage = lazy(() => import("./pages/admin/AdminGamesPage"));
const AdminUsersPage = lazy(() => import("./pages/admin/AdminUsersPage"));
const AdminPicksPage = lazy(() => import("./pages/admin/AdminPicksPage"));
const UserDetailsPage = lazy(() => import("./pages/admin/UserDetailsPage"));
const WeeklyPicksPage = lazy(() => import("./pages/admin/WeeklyPicksPage"));
const UserPicksPage = lazy(() => import("./pages/admin/UserPicksPage"));

function App() {
  return (
    <ErrorBoundary>
      <QueryErrorBoundary>
        <AuthProvider store={authStore} fallbackPath="/">
          <Router>
            <Routes>
              <Route path="/" element={<Layout />}>
                <Route
                  index
                  element={
                    <Suspense fallback={<LoadingSpinner />}>
                      <HomePage />
                    </Suspense>
                  }
                />
                <Route
                  path="login"
                  element={
                    <Suspense fallback={<LoadingSpinner />}>
                      <LoginPage />
                    </Suspense>
                  }
                />
                <Route
                  path="register"
                  element={
                    <Suspense fallback={<LoadingSpinner />}>
                      <RegisterPage />
                    </Suspense>
                  }
                />
                <Route
                  path="profile"
                  element={
                    <ProtectedRoute>
                      <QueryErrorBoundary>
                        <Suspense fallback={<LoadingSpinner />}>
                          <ProfilePage />
                        </Suspense>
                      </QueryErrorBoundary>
                    </ProtectedRoute>
                  }
                />
                <Route
                  path="picks"
                  element={
                    <ProtectedRoute>
                      <QueryErrorBoundary>
                        <Suspense fallback={<LoadingSpinner />}>
                          <PickEntryPage />
                        </Suspense>
                      </QueryErrorBoundary>
                    </ProtectedRoute>
                  }
                />
                <Route
                  path="results"
                  element={
                    <ProtectedRoute requireAdmin>
                      <QueryErrorBoundary>
                        <Suspense fallback={<LoadingSpinner />}>
                          <ResultEntryPage />
                        </Suspense>
                      </QueryErrorBoundary>
                    </ProtectedRoute>
                  }
                />
                <Route
                  path="weekly-results"
                  element={
                    <ProtectedRoute>
                      <QueryErrorBoundary>
                        <Suspense fallback={<LoadingSpinner />}>
                          <WeeklyResultsPage />
                        </Suspense>
                      </QueryErrorBoundary>
                    </ProtectedRoute>
                  }
                />
                <Route
                  path="overall-results"
                  element={
                    <ProtectedRoute>
                      <QueryErrorBoundary>
                        <Suspense fallback={<LoadingSpinner />}>
                          <OverallResultsPage />
                        </Suspense>
                      </QueryErrorBoundary>
                    </ProtectedRoute>
                  }
                />
                <Route
                  path="survivor-pool"
                  element={
                    <ProtectedRoute>
                      <QueryErrorBoundary>
                        <Suspense fallback={<LoadingSpinner />}>
                          <SurvivorPoolPage />
                        </Suspense>
                      </QueryErrorBoundary>
                    </ProtectedRoute>
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
                      <Suspense fallback={<LoadingSpinner />}>
                        <AdminDashboard />
                      </Suspense>
                    </QueryErrorBoundary>
                  }
                />
                <Route
                  path="games"
                  element={
                    <QueryErrorBoundary>
                      <Suspense fallback={<LoadingSpinner />}>
                        <AdminGamesPage />
                      </Suspense>
                    </QueryErrorBoundary>
                  }
                />
                <Route
                  path="users"
                  element={
                    <QueryErrorBoundary>
                      <Suspense fallback={<LoadingSpinner />}>
                        <AdminUsersPage />
                      </Suspense>
                    </QueryErrorBoundary>
                  }
                />
                <Route
                  path="users/:userId"
                  element={
                    <QueryErrorBoundary>
                      <Suspense fallback={<LoadingSpinner />}>
                        <UserDetailsPage />
                      </Suspense>
                    </QueryErrorBoundary>
                  }
                />
                <Route
                  path="picks"
                  element={
                    <QueryErrorBoundary>
                      <Suspense fallback={<LoadingSpinner />}>
                        <AdminPicksPage />
                      </Suspense>
                    </QueryErrorBoundary>
                  }
                />
                <Route
                  path="picks/week/:week"
                  element={
                    <QueryErrorBoundary>
                      <Suspense fallback={<LoadingSpinner />}>
                        <WeeklyPicksPage />
                      </Suspense>
                    </QueryErrorBoundary>
                  }
                />
                <Route
                  path="picks/user/:userId"
                  element={
                    <QueryErrorBoundary>
                      <Suspense fallback={<LoadingSpinner />}>
                        <UserPicksPage />
                      </Suspense>
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
