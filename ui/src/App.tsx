import { Routes, Route } from 'react-router';
import HomePage from './pages/home';
import LoginPage from './pages/login';
import SignupPage from './pages/signup';
import { AuthProvider } from './context/AuthContext';
import { MessagingProvider } from './context/MessagingContext';

function App() {
    return (
        <AuthProvider>
            <MessagingProvider>
                <Routes>
                    <Route path="/" element={<HomePage />} />
                    <Route path="/login" element={<LoginPage />} />
                    <Route path="/signup" element={<SignupPage />} />
                </Routes>
            </MessagingProvider>
        </AuthProvider>
    );
}

export default App;

