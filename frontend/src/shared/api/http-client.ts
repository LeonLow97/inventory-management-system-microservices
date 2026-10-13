export class ApiError extends Error {
  readonly status: number | null;

  constructor(status: number | null, message: string) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
  }
}

function errorMessage(status: number): string {
  switch (status) {
    case 400:
      return 'Please check your email and password and try again.';
    case 401:
      return 'The email or password is incorrect.';
    case 429:
      return 'Too many login attempts. Wait a moment and try again.';
    default:
      return 'We could not complete your request. Please try again shortly.';
  }
}

// Requests remain on the frontend origin by default; Vite proxies them in development.
const baseUrl = (import.meta.env.VITE_API_BASE_URL || '/api/v1').replace(/\/$/, '');

export async function request(path: string, options: RequestInit = {}): Promise<Response> {
  const headers = new Headers(options.headers);
  headers.set('Accept', 'application/json');
  if (options.body) headers.set('Content-Type', 'application/json');

  let response: Response;
  try {
    response = await fetch(new URL(`${baseUrl}${path}`, window.location.origin), {
      ...options,
      headers,
      credentials: 'include',
      cache: 'no-store',
    });
  } catch (error) {
    if (options.signal?.aborted) throw error;
    throw new ApiError(null, 'Unable to reach Authify. Check your connection and try again.');
  }

  if (!response.ok) throw new ApiError(response.status, errorMessage(response.status));
  return response;
}
