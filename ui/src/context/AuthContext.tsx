import {
    createContext,
    useState,
    useContext,
    useEffect,
    ReactNode,
    useCallback,
} from 'react';
import { User } from '../lib/definitions_v2';

// Context shape
type IAuthContext = {
    user: User | null;
    login: (email: string, password: string) => Promise<void>;
    logout: () => void;
    loading: boolean;
    error: string | null;
};

const AuthContext = createContext<IAuthContext | null>(null);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
    const [user, setUser] = useState<User | null>(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const login = useCallback(async (email: string, password: string) => {
        setLoading(true);
        setError(null);

        const response: Response = await fetch("http://localhost:3000/login", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            credentials: "include",
            body: JSON.stringify({ email, password }),
        });

        const body = await response.json();
        if (response.ok) {
            setUser(body["user"]);
        } else {
            setError(body["error"]);
        }

        setLoading(false);
    }, []);

    const logout = () => {
        setUser(null);
    };

    // Optionally check local storage or session on mount
    useEffect(() => {
        // Load user if already logged in (not implemented here)
    }, []);

    return (
        <AuthContext.Provider value={{ user, login, logout, loading, error }}>
            {children}
        </AuthContext.Provider>
    );
};

// Custom hook
export const useAuth = () => {
    const context = useContext(AuthContext);
    if (!context) {
        throw new Error('useAuth must be used within an AuthProvider');
    }
    return context;
};

