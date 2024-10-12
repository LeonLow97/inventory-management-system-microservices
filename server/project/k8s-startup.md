# Kubernetes StartUp

## Kubernetes (Local Startup)

- Tool Requirements: Docker
- Push images of microservices to local Docker Registry to enable containers in Pods to pull the image.
- Run `sh k8s.sh` to automate this process, ensure that port 5000 is not used by running `lsof -i :5000` because local Docker registry container runs on port 5000.

## Useful `kubectl` commands

```
# force restart of Deployment if Pods are not reflecting changes
kubectl rollout restart deployment <deployment_name>

# force restart of StatefulSet
kubectl rollout restart statefulset <statefulset_name>

# to test connectivity with pod via health check endpoint
kubectl exec -it <pod-name> -- wget -qO- http://localhost:80/healthcheck
```
