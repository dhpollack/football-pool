import {
  createContext,
  useContext,
  useState,
  useEffect,
  type ReactNode,
} from "react";
import { api } from "../services/api";
import { ApiError } from "../services/api";

interface User {
  id: number;
  name: string;
  email: string;
}

interface AuthContextType {
  user: User | null;
  login: (email: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
  isAuthenticated: boolean;
  loading: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};

interface AuthProviderProps {
  children: ReactNode;
}

export const AuthProvider = ({ children }: AuthProviderProps) => {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let isMounted = true;

    const initializeAuth = async () => {
      try {
        // Verify if user is authenticated by making a protected request
        const userData = await api.get<User>("/api/profile");
        if (isMounted) {
          setUser(userData);
        }
      } catch (error) {
        // Not authenticated or token expired, clear any stale data
        if (isMounted) {
          setUser(null);
        }
        // Handle auth errors (auto-logout on 401)
        if (error instanceof ApiError && error.status === 401) {
          handleAuthError(error);
        }
      } finally {
        if (isMounted) {
          setLoading(false);
        }
      }
    };

    initializeAuth();

    // Set up periodic token check (every 5 minutes)
    const tokenCheckInterval = setInterval(() => {
      if (user) {
        api.get("/api/profile").catch((error) => {
          if (error instanceof ApiError && error.status === 401) {
            handleAuthError(error);
          }
        });
      }
    }, 5 * 60 * 1000);

    return () => {
      isMounted = false;
      clearInterval(tokenCheckInterval);
    };
  }, [user]);

  const login = async (email: string, password: string) => {
    try {
      const response = await api.post<{ user: User }>("/api/login", {
        email,
        password,
      });

      setUser(response.user);
    } catch (error) {
      throw new Error(error instanceof Error ? error.message : "Login failed");
    }
  };

  const handleAuthError = (error: unknown) => {
    if (error instanceof ApiError && error.status === 401) {
      // Auto-logout on authentication errors
      setUser(null);
      window.location.href = "/login";
    }
    throw error;
  };

  const logout = async () => {
    try {
      await api.post("/api/logout");
    } catch (error) {
      console.error("Logout error:", error);
    } finally {
      setUser(null);
      window.location.href = "/login";
    }
  };

  const value: AuthContextType = {
    user,
    login,
    logout,
    isAuthenticated: !!user,
    loading,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};
