export interface User {
  id: string;
  email: string;
  createdAt: string;
}

export interface SignupInput {
  email: string;
  password: string;
}

export interface SigninInput {
  email: string;
  password: string;
}

export interface AuthResponse {
  user: User;
}

export interface CurrentUserResponse {
  currentUser: User | null;
}
