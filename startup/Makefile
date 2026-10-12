export PATH := $(PATH):$(shell go env GOPATH)/bin

.DEFAULT_GOAL := help

STARTUP_DIR := $(CURDIR)
PROJECT_ROOT := $(abspath $(STARTUP_DIR)/..)
DOCKER_COMPOSE_FILE := $(STARTUP_DIR)/docker/docker-compose.dev.yaml
TERRAFORM_DIR := $(STARTUP_DIR)/aws/terraform/environments/ec2
TFVARS_FILE ?= $(TERRAFORM_DIR)/terraform.tfvars
SCHEMA_FILE ?= $(PROJECT_ROOT)/authify-service/schema/users.sql
IMAGE_TAG ?= local-$(shell date +%Y%m%d%H%M%S)
K6_SCRIPT := $(STARTUP_DIR)/load-test/login.js
K6_BASE_URL ?= http://localhost:8080
export K6_BASE_URL K6_USERS
INGRESS_NGINX_MANIFEST := https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.15.1/deploy/static/provider/kind/deploy.yaml

ENV ?= local
SERVICE ?= api-gateway

# -----------------------------------------------------------------------------
# Entry points
# -----------------------------------------------------------------------------
.PHONY: help run stop

help:
	@printf '%s\n' \
		'Run an environment:' \
		'  make -C startup run ENV=local SERVICE=api-gateway' \
		'  make -C startup run ENV=local SERVICE=authify-service' \
		'  make -C startup run ENV=docker' \
		'  make -C startup run ENV=kind' \
		'' \
		'Environment commands:' \
		'  make -C startup docker-up       Start Docker Compose in the background' \
		'  make -C startup docker-down     Stop Docker Compose' \
		'  make -C startup kind-up         Deploy to the active Kind cluster' \
		'  make -C startup kind-port-forward  Expose the Kind API at localhost:8080' \
		'  make -C startup local-load-test|docker-load-test|kind-load-test (see docs/k6-load-test.md)' \
		'  make -C startup kind-down       Remove the Authify namespace and local database volume' \
		'  make -C startup stop ENV=docker|kind' \
		'' \
		'AWS with Terraform:' \
		'  make -C startup aws-ec2-plan   Preview the AWS infrastructure changes' \
		'  make -C startup aws-ec2-up     Create the stack, publish images, and wait for readiness' \
		'  make -C startup aws-ec2-down   Delete the Terraform-managed AWS stack' \
		'' \
		'Utilities:' \
		'  make -C startup proto' \
		'  make -C startup local-redis|local-prometheus|local-grafana'

run:
ifeq ($(ENV),local)
	$(MAKE) local-run SERVICE=$(SERVICE)
else ifeq ($(ENV),docker)
	$(MAKE) docker-run
else ifeq ($(ENV),kind)
	$(MAKE) kind-run
else
$(error Unsupported ENV '$(ENV)'. Choose local, docker, or kind.)
endif

stop:
ifeq ($(ENV),docker)
	$(MAKE) docker-down
else ifeq ($(ENV),kind)
	$(MAKE) kind-down
else ifeq ($(ENV),local)
	@echo 'Local services run in the foreground. Stop them with Ctrl+C.'
else
$(error Unsupported ENV '$(ENV)'. Choose local, docker, or kind.)
endif

# -----------------------------------------------------------------------------
# Shared tools
# -----------------------------------------------------------------------------
.PHONY: proto aws-cli-install

proto:
	$(MAKE) -C $(PROJECT_ROOT)/grpc-proto generate

aws-cli-install:
	brew install awscli
	brew install awsebcli
	aws --version
	eb --version

# -----------------------------------------------------------------------------
# Local processes and tools
# -----------------------------------------------------------------------------
.PHONY: local-run local-redis local-prometheus local-grafana
.PHONY: run-local redis-server prom grafana

local-run:
ifeq ($(SERVICE),api-gateway)
	$(MAKE) -C $(PROJECT_ROOT)/api-gateway run MODE=development
else ifeq ($(SERVICE),authify-service)
	$(MAKE) -C $(PROJECT_ROOT)/authify-service run MODE=development
else
$(error Unsupported SERVICE '$(SERVICE)'. Use api-gateway or authify-service.)
endif

local-redis:
	brew install redis
	brew services start redis

local-prometheus:
	cd $(STARTUP_DIR)/observability && prometheus --config.file=prometheus.yml

local-grafana:
	brew install grafana
	brew services start grafana

# Compatibility aliases for existing commands.
run-local: local-run
redis-server: local-redis
prom: local-prometheus
grafana: local-grafana

# -----------------------------------------------------------------------------
# Docker Compose
# -----------------------------------------------------------------------------
.PHONY: docker-run docker-up docker-down docker-build run-docker

docker-run:
	docker-compose -f $(DOCKER_COMPOSE_FILE) up --build

docker-up:
	docker-compose -f $(DOCKER_COMPOSE_FILE) up --build -d

docker-down:
	docker-compose -f $(DOCKER_COMPOSE_FILE) down

# Compatibility aliases for existing commands.
docker-build: docker-up
run-docker: docker-run

# -----------------------------------------------------------------------------
# Kubernetes on Kind
# -----------------------------------------------------------------------------
.PHONY: kind-context-check kind-images kind-infrastructure kind-database kind-services
.PHONY: kind-up kind-run kind-port-forward kind-down
.PHONY: load-test local-load-test docker-load-test kind-load-test
.PHONY: k8s-context-check k8s-apply k8s-port-forward k8s-delete

kind-context-check:
	@context="$$(kubectl config current-context)"; \
	case "$$context" in \
		kind-*) echo "Using Kind context: $$context" ;; \
		*) echo "Current context '$$context' is not a Kind context. Switch kubectl to Kind first." >&2; exit 1 ;; \
	esac

kind-up:
	$(MAKE) kind-images
	$(MAKE) kind-infrastructure
	$(MAKE) kind-database
	$(MAKE) kind-services

kind-images: kind-context-check
	./scripts/kind/prepare-kind-images.sh

kind-infrastructure: kind-context-check
	kubectl apply -f ./kubernetes/kind/namespace.yml
	kubectl apply -f ./kubernetes/kind/storageclass.yml
	kubectl apply -f $(INGRESS_NGINX_MANIFEST)
	kubectl rollout status deployment/ingress-nginx-controller -n ingress-nginx --timeout=180s
	kubectl apply -f ./kubernetes/kind/redis.yml
	kubectl rollout status deployment/redis -n authify --timeout=180s

kind-database: kind-context-check
	kubectl apply -f $(PROJECT_ROOT)/authify-service/k8s/configmap.yml \
		-f $(PROJECT_ROOT)/authify-service/k8s/secret.yml \
		-f $(PROJECT_ROOT)/authify-service/k8s/headless-clusterip.yml \
		-f $(PROJECT_ROOT)/authify-service/k8s/persistent-volume.yml \
		-f $(PROJECT_ROOT)/authify-service/k8s/statefulset.yml
	kubectl rollout status statefulset/user-postgres-statefulset -n authify --timeout=180s
	kubectl create configmap user-schema -n authify \
		--from-file=users.sql=$(PROJECT_ROOT)/authify-service/schema/users.sql \
		--dry-run=client -o yaml | kubectl apply -f -
	kubectl delete job user-schema-init -n authify --ignore-not-found
	kubectl apply -f $(PROJECT_ROOT)/authify-service/k8s/schema-init-job.yml
	kubectl wait --for=condition=complete job/user-schema-init -n authify --timeout=180s

kind-services: kind-context-check
	kubectl apply -f $(PROJECT_ROOT)/authify-service/k8s/clusterip.yml \
		-f $(PROJECT_ROOT)/authify-service/k8s/deployment.yml
	kubectl rollout status deployment/user-deployment -n authify --timeout=180s
	kubectl apply -f $(PROJECT_ROOT)/api-gateway/k8s/clusterip.yml \
		-f $(PROJECT_ROOT)/api-gateway/k8s/deployment.yml
	kubectl rollout status deployment/api-gateway-deployment -n authify --timeout=180s
	kubectl apply -f ./kubernetes/kind/ingress-resource.yml

kind-run: kind-up
	$(MAKE) kind-port-forward

kind-port-forward: kind-context-check
	kubectl port-forward --address 127.0.0.1 -n ingress-nginx service/ingress-nginx-controller 8080:80

load-test:
	@command -v k6 >/dev/null 2>&1 || { echo 'k6 is required. Install it with: brew install k6' >&2; exit 1; }
	@test -n "$${K6_USERS:-}" || { echo 'Set K6_USERS to a JSON array of five unique test accounts. See docs/k6-load-test.md.' >&2; exit 1; }
	@curl --fail --silent --show-error "$(K6_BASE_URL)/health/readiness" >/dev/null || { echo 'API is not reachable at $(K6_BASE_URL). Start the selected environment and, for Kind, keep port forwarding active.' >&2; exit 1; }
	k6 run "$(K6_SCRIPT)"

local-load-test docker-load-test: load-test

kind-load-test: kind-context-check
	$(MAKE) load-test

kind-down: kind-context-check
	kubectl delete namespace authify --ignore-not-found --wait=true
	kubectl delete -f $(PROJECT_ROOT)/authify-service/k8s/persistent-volume.yml --ignore-not-found
	kubectl delete -f ./kubernetes/kind/storageclass.yml --ignore-not-found

# Compatibility aliases for existing commands.
k8s-context-check: kind-context-check
k8s-apply: kind-up
k8s-port-forward: kind-port-forward
k8s-delete: kind-down

# -----------------------------------------------------------------------------
# AWS infrastructure with Terraform
# -----------------------------------------------------------------------------
.PHONY: terraform-init aws-ec2-plan aws-ec2-up aws-ec2-down

terraform-init:
	terraform -chdir=$(TERRAFORM_DIR) init

aws-ec2-plan:
	$(MAKE) terraform-init
	@test -f "$(TFVARS_FILE)" || { echo "Missing $(TFVARS_FILE). Copy terraform.tfvars.example and set the required secrets." >&2; exit 1; }
	terraform -chdir=$(TERRAFORM_DIR) plan -var-file="$(TFVARS_FILE)" -var="schema_file=$(SCHEMA_FILE)" -var="image_tag=$(IMAGE_TAG)"

aws-ec2-up:
	$(MAKE) terraform-init
	@test -f "$(TFVARS_FILE)" || { echo "Missing $(TFVARS_FILE). Copy terraform.tfvars.example and set the required secrets." >&2; exit 1; }
	terraform -chdir=$(TERRAFORM_DIR) apply -var-file="$(TFVARS_FILE)" -var="schema_file=$(SCHEMA_FILE)" -var="image_tag=$(IMAGE_TAG)"
	bash "$(STARTUP_DIR)/aws/terraform/scripts/push-images.sh" "$(IMAGE_TAG)"
	bash "$(STARTUP_DIR)/aws/terraform/scripts/wait-for-readiness.sh"
	@printf 'API is ready at %s\n' "$$(terraform -chdir=$(TERRAFORM_DIR) output -raw api_base_url)"

aws-ec2-down:
	$(MAKE) terraform-init
	@test -f "$(TFVARS_FILE)" || { echo "Missing $(TFVARS_FILE). Keep terraform.tfvars until destroy has completed." >&2; exit 1; }
	terraform -chdir=$(TERRAFORM_DIR) destroy -var-file="$(TFVARS_FILE)" -var="schema_file=$(SCHEMA_FILE)"
