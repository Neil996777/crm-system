import { apiRequest } from './client';
import { UserRole } from './auth';

export type UserStatus = 'Active' | 'Disabled';

export type ManagedUser = {
  id: string;
  email: string;
  displayName: string;
  role: UserRole;
  status: UserStatus;
};

export type UserListResponse = {
  items: ManagedUser[];
  activeAdministratorCount: number;
};

export async function listUsers() {
  return apiRequest<UserListResponse>('/admin/users');
}

export async function createUser(input: { email: string; displayName: string; password: string; role: UserRole }) {
  const response = await apiRequest<{ user: ManagedUser }>('/admin/users', {
    method: 'POST',
    body: JSON.stringify(input)
  });
  return response.user;
}

export async function changeUserRole(id: string, role: UserRole) {
  const response = await apiRequest<{ user: ManagedUser }>(`/admin/users/${id}/role`, {
    method: 'PATCH',
    body: JSON.stringify({ role })
  });
  return response.user;
}

export async function changeUserStatus(id: string, status: UserStatus) {
  const response = await apiRequest<{ user: ManagedUser }>(`/admin/users/${id}/status`, {
    method: 'PATCH',
    body: JSON.stringify({ status })
  });
  return response.user;
}
