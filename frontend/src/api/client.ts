export type ApiError = {
  code: string;
  category: string;
  safeMessage: string;
  fieldErrors?: Array<{ field: string; code: string; safeMessage: string }>;
};

export async function apiRequest<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(path, {
    ...init,
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
      ...(init?.headers ?? {})
    }
  });
  const body = await response.json().catch(() => ({}));
  if (!response.ok) {
    const error = body.error as ApiError | undefined;
    throw error ?? { code: 'REQUEST_FAILED', category: 'system', safeMessage: 'Request failed.' };
  }
  return body as T;
}
