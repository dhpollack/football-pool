
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Layout from './components/Layout';
import HomePage from './pages/HomePage';
import LoginPage from './pages/LoginPage';
import RegisterPage from './pages/RegisterPage';
import ProfilePage from './pages/ProfilePage';
import PickEntryPage from './pages/PickEntryPage';
import ResultEntryPage from './pages/ResultEntryPage';
import WeeklyResultsPage from './pages/WeeklyResultsPage';
import OverallResultsPage from './pages/OverallResultsPage';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Layout />}>
          <Route index element={<HomePage />} />
          <Route path="login" element={<LoginPage />} />
          <Route path="register" element={<RegisterPage />} />
          <Route path="profile" element={<ProfilePage />} />
          <Route path="picks" element={<PickEntryPage />} />
          <Route path="results" element={<ResultEntryPage />} />
          <Route path="weekly-results" element={<WeeklyResultsPage />} />
          <Route path="overall-results" element={<OverallResultsPage />} />
        </Route>
      </Routes>
    </Router>
  );
}

export default App;