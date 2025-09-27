import useAuthUser from "react-auth-kit/hooks/useAuthUser";
import useSignOut from "react-auth-kit/hooks/useSignOut";
import useIsAuthenticated from "react-auth-kit/hooks/useIsAuthenticated";
import { useNavigate } from "react-router-dom";
import type { AuthUser } from "../lib/auth";

/**
 * Custom hook that provides a similar interface to our previous useAuth hook
 * but uses react-auth-kit under the hood
 */
export const useAuth = () => {
  const user = useAuthUser<AuthUser>();
  const isAuthenticated = useIsAuthenticated();
  const signOut = useSignOut();
  const navigate = useNavigate();

  // react-auth-kit v4 uses a different approach for loading state
  // We need to check if the authentication state is still initializing
  const loading = user === undefined && !isAuthenticated;

  return {
    user,
    login: async (_email: string, _password: string) => {
      // This will be implemented in the LoginPage using React Query mutations
      // For now, we'll keep the signature but the actual implementation
      // will be in the LoginPage component
      throw new Error(
        "Login should be implemented using React Query mutations in LoginPage",
      );
    },
    logout: async () => {
      signOut();
      // Redirect to login page after logout
      navigate("/login");
    },
    isAuthenticated,
    isAdmin: user?.role === "admin",
    loading,
  };
};
