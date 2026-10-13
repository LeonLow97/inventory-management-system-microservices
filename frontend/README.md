# Authify frontend

React + TypeScript + Vite, styled with Material UI. The frontend signs in through the API Gateway and displays the authenticated user's first and last names from PostgreSQL.

## Run locally

Use Node.js 22.12+ and npm. First follow the [localhost setup](../docs/setup/localhost.md) to start PostgreSQL and Redis, then run the backend services in separate terminals:

```sh
make -C startup run ENV=local SERVICE=authify-service
make -C startup run ENV=local SERVICE=api-gateway
```

From the repository root:

```sh
npm --prefix frontend ci
npm --prefix frontend run dev
```

Open **http://localhost:3000**. Sign in with an existing database account. To create a throwaway account for development:

```sh
curl -i -X POST http://127.0.0.1:8080/api/v1/signup \
  -H 'Content-Type: application/json' \
  -d '{"first_name":"Ada","last_name":"Lovelace","email":"ada@example.com","password":"LoadTest1!"}'
```

The development server proxies `/api` to `http://127.0.0.1:8080`. It rewrites the gateway's cookie domain to a host-only cookie on the frontend origin, so cookie authentication works with the gateway's local, Docker, and Kind configurations. Keep the browser on one hostname while using the app.

For Docker, start `make -C startup docker-up`. For Kind, start `make -C startup kind-up` and keep `make -C startup kind-port-forward` running. The frontend commands above work with either environment when its gateway is reachable on port 8080. Each environment has its own database accounts.

To change the proxy target, copy `.env.example` to `.env.local`, set `API_PROXY_TARGET`, and restart Vite. `VITE_API_BASE_URL` is a public browser setting; keep its default `/api/v1` when using the proxy. Never put secrets in `VITE_*` variables.

## Source structure

```text
src/
  app/                 Application composition, routes, providers, query client, theme
  features/auth/
    api/               Login, current-user, and logout requests
    components/        Login form, Google mark, welcome artwork, protected route
    hooks/             Authentication queries and mutations
    model/             Validation schemas and domain types
  pages/               Login and home screen composition
  shared/
    api/               Credentialed HTTP client and safe request errors
    ui/                Brand and loading/error states
  test/                Test setup and HTTP fixtures
```

Keep feature-specific behavior inside its feature. Shared code should remain independent of pages and features. Use direct imports instead of broad barrel exports. TypeScript uses strict checking; ESLint checks typed code and React hooks, and Prettier defines formatting.

## Authentication behavior

- `POST /api/v1/login` returns database profile details and sets the HTTP-only authentication cookie.
- `GET /api/v1/users/me` validates the cookie and reads the authenticated profile from PostgreSQL. It restores the user after refresh and rejects inactive/deleted users.
- `/` is protected. Guests are redirected to `/login`; signed-in users visiting `/login` are redirected home.
- Authentication is rechecked on window focus and every minute while signed in. A `401` clears the authenticated UI and redirects to login. Connectivity errors show a retry screen.
- `POST /api/v1/logout` expires the cookie. The frontend clears its in-memory profile after the server confirms logout.
- The home page shows a simple greeting with the user's first and last name. The Google sign-in button is a visual placeholder and displays a message because Google OAuth is not configured yet.

Passwords and tokens are never persisted in browser storage. TanStack Query holds profile data in memory, and API responses are validated before use. The home greeting reads `first_name` and `last_name` from the backend response.

## Checks

```sh
npm --prefix frontend run typecheck
npm --prefix frontend run lint
npm --prefix frontend test
npm --prefix frontend run format:check
npm --prefix frontend run build
```

The integration tests use MSW to exercise successful login, guest redirects, invalid credentials, rate limits, profile restoration, expired sessions, logout, and connection failures. Backend tests cover the authenticated identity path and database profile query. Mocked frontend tests do not require a running database.

## Production

`npm run build` writes `dist/`. Serve those static files with a fallback to `index.html` for client routes. Route `/api` to the API Gateway on the same HTTPS origin; Vite's development proxy is not included in the production bundle. Configure the gateway's cookie domain for that hostname (or a host-only cookie) and enable secure cookies.

`npm run preview` is only a local build preview. It also proxies `/api` using `API_PROXY_TARGET` so you can inspect the built application locally.
