import { Routes, Route } from 'react-router';
import HomePage from './pages/home';
import LoginPage from './pages/login';
import SignupPage from './pages/signup';

function App() {
  return (
    <Routes>
      <Route path="/" element={<HomePage />} />
      <Route path="/login" element={<LoginPage />} />
      <Route path="/signup" element={<SignupPage />} />
    </Routes>
  );
}

export default App;

