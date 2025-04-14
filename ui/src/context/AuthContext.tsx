import {
    createContext,
    useState,
    useContext,
    useEffect,
    ReactNode,
    useCallback,
} from 'react';
import { User } from '../lib/definitions_v2';
import { ErrorResponse } from '../lib/definitions_v2';

type LoginResponse = {
    id: number;
    firstName: string;
    lastName: string;
    email: string;
    token: string;
}

type ValidateResponse = {
    id: number;
    firstName: string;
    lastName: string;
    email: string;
    token: string;
}

type IAuthContext = {
    user: User | null;
    login: (email: string, password: string) => Promise<void>;
    logout: () => void;
    isLoading: boolean;
};

const AuthContext = createContext<IAuthContext | null>(null);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
    const [user, setUser] = useState<User | null>(null);
    const [isLoading, setIsLoading] = useState(true);

    const login = useCallback(async (email: string, password: string) => {
        const response = await fetch("http://localhost:3000/login", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            credentials: "include",
            body: JSON.stringify({ email, password }),
        });

        if (response.ok) {
            const loginResponse: LoginResponse = await response.json();
            const user: User = {
                id: loginResponse.id,
                firstName: loginResponse.firstName,
                lastName: loginResponse.lastName,
                email: loginResponse.email,
            }
            setUser(user);
        } else {
            const errorResponse: ErrorResponse = await response.json();
            throw new Error(errorResponse.details);
        }
    }, []);

    const logout = () => {
        setUser(null);
    };

    const validate = async () => {
        setIsLoading(true);
        const response = await fetch("http://localhost:3000/user/me", {
            method: "GET",
            credentials: "include",
        });

        if (response.ok) {
            const validateResponse: ValidateResponse = await response.json();
            const user: User = {
                id: validateResponse.id,
                firstName: validateResponse.firstName,
                lastName: validateResponse.lastName,
                email: validateResponse.email,
            }
            setUser(user);
        } else {
            setUser(null);
        }
        setIsLoading(false);
    }

    useEffect(() => {
        validate();
    }, []);

    return (
        <AuthContext.Provider value={{
            user,
            login,
            logout,
            isLoading,
        }}>
            {children}
        </AuthContext.Provider>
    );
};

export const useAuth = () => {
    const context = useContext(AuthContext);
    if (!context) {
        throw new Error('useAuth must be used within an AuthProvider');
    }
    return context;
}
