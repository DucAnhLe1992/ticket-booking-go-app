import axios from 'axios';

export const apiClient = axios.create({
  baseURL: '/api',
  withCredentials: true,
  headers: {
    'Content-Type': 'application/json',
  },
});

export class ApiError extends Error {
  constructor(
    message: string,
    public status?: number,
    public errors?: Array<{ message: string; field?: string }>
  ) {
    super(message);
    this.name = 'ApiError';
  }
}

apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.data?.errors) {
      throw new ApiError(
        error.response.data.errors[0]?.message || 'An error occurred',
        error.response.status,
        error.response.data.errors
      );
    }
    throw new ApiError(
      error.message || 'An error occurred',
      error.response?.status
    );
  }
);
