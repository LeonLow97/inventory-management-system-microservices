import { http, HttpResponse } from 'msw';
import { setupServer } from 'msw/node';

export const profile = { first_name: 'Ada', last_name: 'Lovelace', email: 'ada@example.com' };

export const server = setupServer(
  http.get('*/api/v1/users/me', () => new HttpResponse(null, { status: 401 })),
);
