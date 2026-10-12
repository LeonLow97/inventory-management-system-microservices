# Run with Docker Compose

Docker Compose starts the API Gateway, Authify Service, PostgreSQL, Redis, and the database schema initialization job.

## 1. Check Redis

The Compose stack publishes Redis on port `6379`. Stop the local Homebrew Redis service first if it is running.

```sh
brew services stop redis
```

## 2. Start the stack

Run this from the repository root:

```sh
make -C startup docker-up
```

Compose creates the `users` table during startup. Check the API Gateway readiness endpoint:

```sh
curl -i http://localhost:8080/health/readiness
```

## 3. Stop the stack

```sh
make -C startup docker-down
```

The PostgreSQL data is stored in a Docker volume and remains available when the stack is stopped.
