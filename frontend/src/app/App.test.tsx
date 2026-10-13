import { act, render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { QueryClient } from '@tanstack/react-query';
import { MemoryRouter } from 'react-router-dom';
import { http, HttpResponse } from 'msw';
import { describe, expect, it } from 'vitest';
import { App } from './App';
import { AppProviders } from './AppProviders';
import { profile, server } from '../test/server';
import { storage } from '../test/storage';

function renderApp(path = '/') {
  const client = new QueryClient({
    defaultOptions: { queries: { retry: false, gcTime: 0 }, mutations: { retry: false } },
  });
  render(
    <AppProviders client={client}>
      <MemoryRouter initialEntries={[path]}>
        <App />
      </MemoryRouter>
    </AppProviders>,
  );
  return client;
}

async function fillLogin() {
  const user = userEvent.setup();
  await user.type(await screen.findByLabelText('Email address'), 'ada@example.com');
  await user.type(screen.getByLabelText('Password'), 'Secret123!');
  return user;
}

describe('cookie-authenticated login flow', () => {
  it('redirects a guest from home to the login form', async () => {
    renderApp();
    expect(await screen.findByRole('heading', { name: 'Good to see you.' })).toBeInTheDocument();
    expect(screen.queryByRole('button', { name: 'Log out' })).not.toBeInTheDocument();
  });

  it('logs in with credentials and greets the user from the API response', async () => {
    let body: unknown;
    let credentials: RequestCredentials | undefined;
    server.use(
      http.post('*/api/v1/login', async ({ request }) => {
        body = await request.json();
        credentials = request.credentials;
        return HttpResponse.json(profile);
      }),
    );
    renderApp('/login');
    const user = await fillLogin();
    await user.click(screen.getByRole('button', { name: 'Log in' }));
    expect(
      await screen.findByRole('heading', { name: 'Hello, Ada Lovelace.' }),
    ).toBeInTheDocument();
    expect(body).toEqual({ email: 'ada@example.com', password: 'Secret123!' });
    expect(credentials).toBe('include');
    expect(storage.local.setItem).not.toHaveBeenCalled();
    expect(storage.session.setItem).not.toHaveBeenCalled();
  });

  it('keeps incorrect credentials on the login page and clears the error when edited', async () => {
    server.use(http.post('*/api/v1/login', () => new HttpResponse(null, { status: 401 })));
    renderApp('/login');
    const user = await fillLogin();
    await user.click(screen.getByRole('button', { name: 'Log in' }));
    expect(await screen.findByRole('alert')).toHaveTextContent(
      'The email or password is incorrect.',
    );
    expect(screen.queryByRole('heading', { name: /Hello,/ })).not.toBeInTheDocument();
    await user.type(screen.getByLabelText('Password'), 'a');
    expect(screen.queryByRole('alert')).not.toBeInTheDocument();
  });

  it('prevents another submission while a login request is in flight', async () => {
    let finish: () => void = () => undefined;
    const responseReady = new Promise<void>((resolve) => {
      finish = resolve;
    });
    let calls = 0;
    server.use(
      http.post('*/api/v1/login', async () => {
        calls++;
        await responseReady;
        return HttpResponse.json(profile);
      }),
    );
    renderApp('/login');
    const user = await fillLogin();
    await user.click(screen.getByRole('button', { name: 'Log in' }));
    try {
      const button = screen.getByRole('button', { name: 'Signing in…' });
      expect(button).toBeDisabled();
      await user.keyboard('{Enter}');
      expect(calls).toBe(1);
    } finally {
      finish();
    }
    expect(
      await screen.findByRole('heading', { name: 'Hello, Ada Lovelace.' }),
    ).toBeInTheDocument();
  });

  it.each(['invalid profile', 'invalid JSON'])(
    'rejects a successful login response with %s',
    async (kind) => {
      server.use(
        http.post('*/api/v1/login', () =>
          kind === 'invalid JSON'
            ? new HttpResponse('not-json', { headers: { 'Content-Type': 'application/json' } })
            : HttpResponse.json({ first_name: 'Ada' }),
        ),
      );
      renderApp('/login');
      const user = await fillLogin();
      await user.click(screen.getByRole('button', { name: 'Log in' }));
      expect(await screen.findByRole('alert')).toHaveTextContent(
        'We could not read your account details.',
      );
      expect(screen.queryByRole('heading', { name: /Hello,/ })).not.toBeInTheDocument();
    },
  );

  it('explains rate limiting without navigating to home', async () => {
    server.use(http.post('*/api/v1/login', () => new HttpResponse(null, { status: 429 })));
    renderApp('/login');
    const user = await fillLogin();
    await user.click(screen.getByRole('button', { name: 'Log in' }));
    expect(await screen.findByRole('alert')).toHaveTextContent('Too many login attempts.');
  });

  it('validates empty fields before making a login request', async () => {
    let calls = 0;
    server.use(
      http.post('*/api/v1/login', () => {
        calls++;
        return HttpResponse.json(profile);
      }),
    );
    renderApp('/login');
    await screen.findByLabelText('Email address');
    await userEvent.click(screen.getByRole('button', { name: 'Log in' }));
    expect(await screen.findByText('Enter your password.')).toBeInTheDocument();
    expect(screen.getByText('Enter a valid email address.')).toBeInTheDocument();
    expect(calls).toBe(0);
  });

  it('restores a signed-in user from the profile endpoint on a fresh page load', async () => {
    server.use(http.get('*/api/v1/users/me', () => HttpResponse.json(profile)));
    renderApp();
    expect(
      await screen.findByRole('heading', { name: 'Hello, Ada Lovelace.' }),
    ).toBeInTheDocument();
    expect(screen.queryByText('ada@example.com')).not.toBeInTheDocument();
  });

  it('returns to login when the server reports that the session expired', async () => {
    server.use(http.get('*/api/v1/users/me', () => HttpResponse.json(profile)));
    const client = renderApp();
    await screen.findByRole('heading', { name: 'Hello, Ada Lovelace.' });
    server.use(http.get('*/api/v1/users/me', () => new HttpResponse(null, { status: 401 })));
    await act(async () => {
      await client.invalidateQueries();
    });
    expect(await screen.findByRole('heading', { name: 'Good to see you.' })).toBeInTheDocument();
  });

  it('clears the authenticated UI after successful logout', async () => {
    let authenticated = true;
    server.use(
      http.get('*/api/v1/users/me', () =>
        authenticated ? HttpResponse.json(profile) : new HttpResponse(null, { status: 401 }),
      ),
      http.post('*/api/v1/logout', () => {
        authenticated = false;
        return new HttpResponse(null, { status: 200 });
      }),
    );
    renderApp();
    await screen.findByRole('heading', { name: 'Hello, Ada Lovelace.' });
    await userEvent.click(screen.getByRole('button', { name: 'Log out' }));
    expect(await screen.findByRole('heading', { name: 'Good to see you.' })).toBeInTheDocument();
  });

  it('keeps the user on home if logout fails', async () => {
    server.use(
      http.get('*/api/v1/users/me', () => HttpResponse.json(profile)),
      http.post('*/api/v1/logout', () => new HttpResponse(null, { status: 503 })),
    );
    renderApp();
    await screen.findByRole('heading', { name: 'Hello, Ada Lovelace.' });
    await userEvent.click(screen.getByRole('button', { name: 'Log out' }));
    expect(await screen.findByRole('alert')).toHaveTextContent(
      'We could not complete your request.',
    );
    expect(screen.getByRole('heading', { name: 'Hello, Ada Lovelace.' })).toBeInTheDocument();
  });

  it('shows a retry screen when the profile service is unavailable, then recovers', async () => {
    server.use(http.get('*/api/v1/users/me', () => new HttpResponse(null, { status: 503 })));
    renderApp();
    expect(
      await screen.findByRole('heading', { name: 'We couldn’t load your account' }),
    ).toBeInTheDocument();
    server.use(http.get('*/api/v1/users/me', () => HttpResponse.json(profile)));
    await userEvent.click(screen.getByRole('button', { name: 'Try again' }));
    await waitFor(() =>
      expect(screen.getByRole('heading', { name: 'Hello, Ada Lovelace.' })).toBeInTheDocument(),
    );
  });
});
