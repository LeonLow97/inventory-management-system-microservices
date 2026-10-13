import { ApiError, request } from '../../../shared/api/http-client';
import { userSchema, type LoginCredentials, type User } from '../model/auth-schema';

async function readUser(response: Response): Promise<User> {
  let data: unknown;
  try {
    data = await response.json();
  } catch {
    throw new ApiError(
      response.status,
      'We could not read your account details. Please try again.',
    );
  }
  const result = userSchema.safeParse(data);
  if (!result.success) {
    throw new ApiError(
      response.status,
      'We could not read your account details. Please try again.',
    );
  }
  return result.data;
}

export async function login(credentials: LoginCredentials): Promise<User> {
  const response = await request('/login', { method: 'POST', body: JSON.stringify(credentials) });
  return readUser(response);
}

export async function getCurrentUser(signal?: AbortSignal): Promise<User | null> {
  try {
    return await readUser(await request('/users/me', { signal }));
  } catch (error) {
    if (error instanceof ApiError && error.status === 401) return null;
    throw error;
  }
}

export async function logout(): Promise<void> {
  await request('/logout', { method: 'POST' });
}
