# Run on Localhost

This setup runs the Go services on your computer. PostgreSQL and Redis must be available locally.

## 1. Prepare PostgreSQL

Start your local PostgreSQL service. It must accept connections on `127.0.0.1:5432`.

Run this once as a PostgreSQL administrator to create the development role and database. Skip it if they already exist.

```sh
psql postgres \
  -c "CREATE ROLE authify WITH LOGIN PASSWORD 'dba872b7'" \
  -c "CREATE DATABASE authify_users_db OWNER authify"
```

Create the users table:

```sh
psql -U authify -d authify_users_db \
  -f authify-service/schema/users.sql
```

## 2. Start Redis

The Make target installs Redis with Homebrew if needed and starts it in the background.

```sh
make -C startup local-redis
```

## 3. Start the services

Open a terminal for each service. Start Authify Service first, then the API Gateway.

```sh
make -C startup run ENV=local SERVICE=authify-service
```

```sh
make -C startup run ENV=local SERVICE=api-gateway
```

Both commands run in the foreground. Stop a service with Ctrl+C.

## 4. Check the gateway

```sh
curl -i http://localhost:8080/health/readiness
```

Prometheus and Grafana are optional. Start them with `make -C startup local-prometheus` and `make -C startup local-grafana`.
