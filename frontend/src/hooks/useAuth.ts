import useAuthUser from 'react-auth-kit/hooks/useAuthUser';
import useSignOut from 'react-auth-kit/hooks/useSignOut';
import type { AuthUser } from '../lib/auth';

/**
 * Custom hook that provides a similar interface to our previous useAuth hook
 * but uses react-auth-kit under the hood
 */
export const useAuth = () => {
  const user = useAuthUser<AuthUser>();
  const signOut = useSignOut();

  return {
    user,
    login: async (_email: string, _password: string) => {
      // This will be implemented in the LoginPage using React Query mutations
      // For now, we'll keep the signature but the actual implementation
      // will be in the LoginPage component
      throw new Error('Login should be implemented using React Query mutations in LoginPage');
    },
    logout: async () => {
      signOut();
    },
    isAuthenticated: !!user,
    isAdmin: user?.role === 'admin',
    loading: false, // react-auth-kit handles loading state internally
  };
};