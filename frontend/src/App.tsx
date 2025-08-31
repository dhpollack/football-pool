import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import AuthProvider from 'react-auth-kit/AuthProvider';
import { authStore } from './lib/auth';
import Layout from "./components/Layout";
import ErrorBoundary from "./components/ErrorBoundary";
import QueryErrorBoundary from "./components/QueryErrorBoundary";
import HomePage from "./pages/HomePage";
import LoginPage from "./pages/LoginPage";
import RegisterPage from "./pages/RegisterPage";
import ProfilePage from "./pages/ProfilePage";
import PickEntryPage from "./pages/PickEntryPage";
import ResultEntryPage from "./pages/ResultEntryPage";
import WeeklyResultsPage from "./pages/WeeklyResultsPage";
import OverallResultsPage from "./pages/OverallResultsPage";
import SurvivorPoolPage from "./pages/SurvivorPoolPage";

function App() {
  return (
    <ErrorBoundary>
      <QueryErrorBoundary>
        <AuthProvider store={authStore} fallbackPath="/login">
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
            </Routes>
          </Router>
        </AuthProvider>
      </QueryErrorBoundary>
    </ErrorBoundary>
  );
}

export default App;
