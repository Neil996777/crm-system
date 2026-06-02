import { apiRequest } from './client';

export type UserRole = 'Administrator' | 'Sales Manager' | 'Sales';

export type CurrentUser = {
  id: string;
  email: string;
  displayName: string;
  role: UserRole;
  status: 'Active' | 'Disabled';
};

type AuthResponse = {
  user: CurrentUser;
};

export function signIn(email: string, password: string) {
  return apiRequest<AuthResponse>('/auth/sign-in', {
    method: 'POST',
    body: JSON.stringify({ email, password })
  });
}

export function currentUser() {
  return apiRequest<AuthResponse>('/auth/current');
}

export async function signOut() {
  await fetch('/auth/sign-out', { method: 'POST', credentials: 'include' });
}
