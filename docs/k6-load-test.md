# k6 Login Load Test

Run the k6 login scenario against the API Gateway on localhost, in Docker Compose, or in Kind. The scenario starts 1,000 virtual users together: each of five unique accounts sends 200 login requests concurrently, with one request from each virtual user.

The Make targets use `http://localhost:8080` by default. Install k6 first; on macOS, run `brew install k6`.

## Create five test accounts

First start the environment you plan to test and make its API Gateway reachable at `K6_BASE_URL` (default `http://localhost:8080`). For Kind, keep `make -C startup kind-port-forward` running in another terminal. Create these throwaway accounts in that environment. Each environment has its own database, so create the accounts separately for local, Docker, and Kind if you plan to test all three. If an account already exists, keep using its existing password.

```sh
for i in 1 2 3 4 5; do
  curl -i -X POST "${K6_BASE_URL:-http://localhost:8080}/api/v1/signup" \
    -H 'Content-Type: application/json' \
    -d "{\"first_name\":\"Load\",\"last_name\":\"Test\",\"email\":\"load-test-${i}@example.com\",\"password\":\"LoadTest1!\"}"
done
```

Set the five credentials in your shell. These sample accounts all use the sample password above; replace them if you created accounts with different credentials. Keep real or reused credentials out of committed files.

```sh
export K6_USERS='[
  {"email":"load-test-1@example.com","password":"LoadTest1!"},
  {"email":"load-test-2@example.com","password":"LoadTest1!"},
  {"email":"load-test-3@example.com","password":"LoadTest1!"},
  {"email":"load-test-4@example.com","password":"LoadTest1!"},
  {"email":"load-test-5@example.com","password":"LoadTest1!"}
]'
```

The script requires exactly five unique email addresses, each with a password.

## Local services

Follow the [localhost setup guide](setup/localhost.md) and keep PostgreSQL, Redis, Authify Service, and the API Gateway running. Then run:

```sh
make -C startup local-load-test
```

## Docker Compose

Start the stack and run the scenario:

```sh
make -C startup docker-up
make -C startup docker-load-test
```

Stop the stack when finished with `make -C startup docker-down`.

## Kubernetes with Kind

Deploy the stack, then keep the port-forward running in a separate terminal:

```sh
make -C startup kind-up
make -C startup kind-port-forward
```

In another terminal, set `K6_USERS` to the same five-account JSON shown above and run:

```sh
make -C startup kind-load-test
```

Remove the Authify resources when finished with `make -C startup kind-down`. The Kind cluster and ingress controller remain available.

## What this scenario measures

Login burst limiting is configured for one request per second per account and client IP. Because the 200 requests for each account arrive together, most requests for that account should receive `429 Too Many Requests`; only admitted requests reach Authify Service. This scenario measures the ingress, API Gateway, and Redis rate limiter under a concurrent burst. It does not represent 1,000 successful authentications or a direct Authify Service throughput test.

To use a different gateway address, set `K6_BASE_URL` when running a target:

```sh
K6_BASE_URL='http://127.0.0.1:8080' make -C startup local-load-test
```
