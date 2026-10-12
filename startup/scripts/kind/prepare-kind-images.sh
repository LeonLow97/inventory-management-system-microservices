#!/usr/bin/env bash

set -euo pipefail

REGISTRY_NAME=registry
REGISTRY_PORT=5050
REGISTRY_NODE_PORT=5000
PODMAN_BIN="$(command -v podman || true)"

CURRENT_CONTEXT="$(kubectl config current-context)"
case "$CURRENT_CONTEXT" in
	kind-*) KIND_CLUSTER_NAME="${CURRENT_CONTEXT#kind-}" ;;
	*)
		echo "Current context '$CURRENT_CONTEXT' is not a Kind context." >&2
		exit 1
		;;
esac

KIND_NODES="$(kind get nodes --name "$KIND_CLUSTER_NAME")"
if [[ -z "$KIND_NODES" ]]; then
	echo "No Kind nodes found for cluster '$KIND_CLUSTER_NAME'." >&2
	exit 1
fi

while IFS= read -r node; do
	[[ -n "$node" ]] || continue
	node_status="$(docker inspect -f '{{.State.Status}}' "$node" 2>/dev/null || echo unknown)"
	if [[ "$node_status" != "running" ]]; then
		echo "Kind node '$node' is not running (status: $node_status). Start or recreate the Kind cluster before applying the stack." >&2
		exit 1
	fi
done <<< "$KIND_NODES"

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"

if ! docker network inspect kind >/dev/null 2>&1; then
	echo "The Kind Docker network was not found. Start the Kind cluster first." >&2
	exit 1
fi

if docker inspect "$REGISTRY_NAME" >/dev/null 2>&1; then
	if [[ "$(docker inspect -f '{{.State.Running}}' "$REGISTRY_NAME")" != "true" ]]; then
		docker start "$REGISTRY_NAME" >/dev/null
	fi
else
	docker run -d \
		--restart=always \
		-p "127.0.0.1:${REGISTRY_PORT}:${REGISTRY_NODE_PORT}" \
		--name "$REGISTRY_NAME" \
		registry:2
fi

if [[ "$(docker inspect -f '{{json .NetworkSettings.Networks.kind}}' "$REGISTRY_NAME")" == "null" ]]; then
	docker network connect kind "$REGISTRY_NAME"
fi

if ! docker port "$REGISTRY_NAME" "$REGISTRY_NODE_PORT/tcp" | grep -q ":${REGISTRY_PORT}$"; then
	echo "Container '$REGISTRY_NAME' must publish container port ${REGISTRY_NODE_PORT} to host port ${REGISTRY_PORT}." >&2
	exit 1
fi

curl --retry 10 --retry-connrefused --retry-delay 1 \
	--silent --show-error --fail "http://localhost:${REGISTRY_PORT}/v2/" >/dev/null

push_image() {
	local image="$1"
	local docker_image_id
	local podman_image_id

	docker_image_id="$(docker image inspect -f '{{.Id}}' "$image")"
	docker_image_id="${docker_image_id#sha256:}"

	if [[ -n "$PODMAN_BIN" ]]; then
		for attempt in {1..5}; do
			podman_image_id="$("$PODMAN_BIN" image inspect -f '{{.Id}}' "$image" 2>/dev/null || true)"
			podman_image_id="${podman_image_id#sha256:}"
			if [[ "$podman_image_id" == "$docker_image_id" ]]; then
				echo "Pushing $image to the local HTTP registry with Podman..."
				"$PODMAN_BIN" push --tls-verify=false "$image"
				return
			fi
			sleep 1
		done
	fi

	echo "No matching Podman image found for $image. Trying Docker push..." >&2
	docker push "$image"
}

REGISTRY_CONFIG_DIR="/etc/containerd/certs.d/localhost:${REGISTRY_PORT}"
while IFS= read -r node; do
	[[ -n "$node" ]] || continue
	docker exec "$node" mkdir -p "$REGISTRY_CONFIG_DIR"
	docker exec -i "$node" sh -c "cat > '$REGISTRY_CONFIG_DIR/hosts.toml'" <<EOF
[host."http://${REGISTRY_NAME}:${REGISTRY_NODE_PORT}"]
  capabilities = ["pull", "resolve", "push"]
EOF
done <<< "$KIND_NODES"

echo "Building and pushing API Gateway image..."
docker build -f "$PROJECT_ROOT/api-gateway/Dockerfile" \
	-t "localhost:${REGISTRY_PORT}/authify-api-gateway:latest" "$PROJECT_ROOT"
push_image "localhost:${REGISTRY_PORT}/authify-api-gateway:latest"

echo "Building and pushing Authify Service image..."
docker build -f "$PROJECT_ROOT/authify-service/Dockerfile" \
	-t "localhost:${REGISTRY_PORT}/authify-service:latest" "$PROJECT_ROOT"
push_image "localhost:${REGISTRY_PORT}/authify-service:latest"
