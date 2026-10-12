# Run on Kubernetes with Kind

This setup deploys the services to a local Kind cluster and exposes the API at `http://localhost:8080`.

## 1. Create or select a Kind cluster

Install `kind`, `kubectl`, and Docker or a compatible container runtime. See the [Kind Quick Start](https://kind.sigs.k8s.io/docs/user/quick-start) or the [Kind GitHub project](https://github.com/kubernetes-sigs/kind).

If you do not have a cluster yet, create the default cluster:

```sh
kind create cluster --name kind
```

Select and verify its kubectl context:

```sh
kubectl config use-context kind-kind
kubectl config current-context
```

The context should be `kind-kind`.

## 2. Deploy the application

Run from the repository root:

```sh
make -C startup kind-up
```

This builds both service images, starts the local registry, configures Kind nodes to pull from it, installs the ingress controller, Redis, and PostgreSQL, initializes the users table, then deploys the services. The registry uses the [`registry:2` image](https://hub.docker.com/_/registry) and is available on localhost port `5050`.

## 3. Forward traffic to the API

Run this in a separate terminal:

```sh
make -C startup kind-port-forward
```

Check that the gateway is ready:

```sh
curl -i http://localhost:8080/health/readiness
```

You can also run `make -C startup run ENV=kind` to deploy the stack and start port forwarding in one command.

## 4. Stop the application

Stop port forwarding with Ctrl+C. Remove the Authify namespace and its PostgreSQL volume with:

```sh
make -C startup kind-down
```

The Kind cluster and ingress controller remain. Deleting the Kind cluster also deletes its local PostgreSQL data.
