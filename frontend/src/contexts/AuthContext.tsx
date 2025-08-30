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

export const AuthContext = createContext<AuthContextType | undefined>(undefined);

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

    // Check authentication status on mount
    const checkAuthStatus = async () => {
      try {
        const userData = await api.get<User>("/api/users/me");
        if (isMounted) {
          setUser(userData);
        }
      } catch (error) {
        if (error instanceof ApiError && error.status === 401) {
          // Not authenticated, clear user state
          if (isMounted) {
            setUser(null);
          }
        }
        // Other errors can be ignored for initialization
      } finally {
        if (isMounted) {
          setLoading(false);
        }
      }
    };

    checkAuthStatus();

    // Set up periodic token check (every 5 minutes)
    const tokenCheckInterval = setInterval(() => {
      // Only check if we think we have a user (based on current state at time of setup)
      api.get("/api/users/me").catch((error) => {
        if (error instanceof ApiError && error.status === 401) {
          handleAuthError(error);
        }
      });
    }, 5 * 60 * 1000);

    return () => {
      isMounted = false;
      clearInterval(tokenCheckInterval);
    };
  }, []); // Empty dependency array - only run once on mount

  const login = async (email: string, password: string) => {
    try {
      // Login and get JWT token
      await api.post("/api/login", {
        email,
        password,
      });

      // Try to get user data after login
      try {
        const userData = await api.get<User>("/api/users/me");
        setUser(userData);
      } catch (profileError) {
        console.log("Profile fetch failed after login:", profileError);
        // If profile fetch fails, still set user as authenticated
        // This handles the case where user exists but no player profile
        setUser({
          id: 0,
          name: "", 
          email: email
        });
      }
    } catch (error) {
      throw new Error(error instanceof Error ? error.message : "Login failed");
    }
  };

  const handleAuthError = (error: unknown) => {
    if (error instanceof ApiError && error.status === 401) {
      // Auto-logout on authentication errors
      setUser(null);
      // Only redirect if we're not already on the login page
      if (window.location.pathname !== "/login") {
        window.location.href = "/login";
      }
    }
    // Don't throw error during initialization
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
