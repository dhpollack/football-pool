import React from "react";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import HomePage from "./pages/HomePage";
import LoginPage from "./pages/LoginPage";
import PickEntryPage from "./pages/PickEntryPage";
import ProfilePage from "./pages/ProfilePage";
import ResultEntryPage from "./pages/ResultEntryPage";
import WeeklyResultsPage from "./pages/WeeklyResultsPage";

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/login" element={<LoginPage />} />
        <Route path="/picks" element={<PickEntryPage />} />
        <Route path="/profile" element={<ProfilePage />} />
        <Route path="/admin/results" element={<ResultEntryPage />} />
        <Route path="/results/week" element={<WeeklyResultsPage />} />
      </Routes>
    </Router>
  );
}

export default App;
