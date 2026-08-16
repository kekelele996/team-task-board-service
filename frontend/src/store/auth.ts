import { create } from 'zustand';
import type { User } from '../types';

interface AuthState {
  token: string | null;
  user: User | null;
  setAuth: (token: string, user: User) => void;
  logout: () => void;
}

function readUser(): User | null {
  const raw = localStorage.getItem('gbkanban_user');
  if (!raw) return null;
  try {
    return JSON.parse(raw) as User;
  } catch {
    return null;
  }
}

export const useAuthStore = create<AuthState>((set) => ({
  token: localStorage.getItem('gbkanban_token'),
  user: readUser(),
  setAuth: (token, user) => {
    localStorage.setItem('gbkanban_token', token);
    localStorage.setItem('gbkanban_user', JSON.stringify(user));
    set({ token, user });
  },
  logout: () => {
    localStorage.removeItem('gbkanban_token');
    localStorage.removeItem('gbkanban_user');
    set({ token: null, user: null });
  },
}));
