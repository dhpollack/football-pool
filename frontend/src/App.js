import React from "react";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import HomePage from "./pages/HomePage";
import LoginPage from "./pages/LoginPage";
import PickEntryPage from "./pages/PickEntryPage";
import ProfilePage from "./pages/ProfilePage";
import ResultEntryPage from "./pages/ResultEntryPage";
import WeeklyResultsPage from "./pages/WeeklyResultsPage";
import OverallResultsPage from "./pages/OverallResultsPage";
import Layout from "./components/Layout";

function App() {
  return (
    <Router>
      <Layout>
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/login" element={<LoginPage />} />
          <Route path="/picks" element={<PickEntryPage />} />
          <Route path="/profile" element={<ProfilePage />} />
          <Route path="/admin/results" element={<ResultEntryPage />} />
          <Route path="/results/week" element={<WeeklyResultsPage />} />
          <Route path="/results/season" element={<OverallResultsPage />} />
        </Routes>
      </Layout>
    </Router>
  );
}

export default App;