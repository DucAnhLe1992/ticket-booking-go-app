import { apiClient } from './client';
import type { SignupInput, SigninInput, AuthResponse, CurrentUserResponse } from '../types/auth';

export const authApi = {
  signup: async (input: SignupInput): Promise<AuthResponse> => {
    const { data } = await apiClient.post<AuthResponse>('/auth/signup', input);
    return data;
  },

  signin: async (input: SigninInput): Promise<AuthResponse> => {
    const { data } = await apiClient.post<AuthResponse>('/auth/signin', input);
    return data;
  },

  signout: async (): Promise<void> => {
    await apiClient.post('/auth/signout');
  },

  getCurrentUser: async (): Promise<CurrentUserResponse> => {
    const { data } = await apiClient.get<CurrentUserResponse>('/auth/currentuser');
    return data;
  },
};
