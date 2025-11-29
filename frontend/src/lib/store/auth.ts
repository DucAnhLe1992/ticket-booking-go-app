import { create } from 'zustand';
import { authApi } from '../api/auth';
import type { User, SignupInput, SigninInput } from '../types/auth';

interface AuthState {
  user: User | null;
  isLoading: boolean;
  error: string | null;
  signin: (input: SigninInput) => Promise<void>;
  signup: (input: SignupInput) => Promise<void>;
  signout: () => Promise<void>;
  checkAuth: () => Promise<void>;
  clearError: () => void;
}

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  isLoading: true,
  error: null,

  signin: async (input: SigninInput) => {
    try {
      set({ error: null });
      const { user } = await authApi.signin(input);
      set({ user });
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Failed to sign in';
      set({ error: message });
      throw error;
    }
  },

  signup: async (input: SignupInput) => {
    try {
      set({ error: null });
      const { user } = await authApi.signup(input);
      set({ user });
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Failed to sign up';
      set({ error: message });
      throw error;
    }
  },

  signout: async () => {
    try {
      await authApi.signout();
      set({ user: null });
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Failed to sign out';
      set({ error: message });
    }
  },

  checkAuth: async () => {
    try {
      const { currentUser } = await authApi.getCurrentUser();
      set({ user: currentUser, isLoading: false });
    } catch {
      set({ user: null, isLoading: false });
    }
  },

  clearError: () => set({ error: null }),
}));
