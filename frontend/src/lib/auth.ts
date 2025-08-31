import createAuthStore from 'react-auth-kit/store/createAuthStore';
import type { UserResponse } from '../services/model';

// Define the user state type that will be stored in auth kit
export type AuthUser = UserResponse;

// Create the auth store
export const authStore = createAuthStore<AuthUser>('localstorage', {
  authName: '_auth',
  debug: process.env.NODE_ENV === 'development',
});