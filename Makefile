.PHONY: all build test docker clean dev

all: build

build:
	go build ./cmd/...

test:
	go test -v -race -coverprofile=coverage.out ./...

test-coverage:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Docker build targets
docker-build-auth:
	docker build -t ghcr.io/ducanhle1992/ticket-auth:dev -f ./cmd/auth/Dockerfile .

docker-build-tickets:
	docker build -t ghcr.io/ducanhle1992/ticket-tickets:dev -f ./cmd/tickets/Dockerfile .

docker-build-orders:
	docker build -t ghcr.io/ducanhle1992/ticket-orders:dev -f ./cmd/orders/Dockerfile .

docker-build-payments:
	docker build -t ghcr.io/ducanhle1992/ticket-payments:dev -f ./cmd/payments/Dockerfile .

docker-build-expiration:
	docker build -t ghcr.io/ducanhle1992/ticket-expiration:dev -f ./cmd/expiration/Dockerfile .

docker-build-all: docker-build-auth docker-build-tickets docker-build-orders docker-build-payments docker-build-expiration

# Local development
dev-infra:
	docker-compose up -d

dev-infra-down:
	docker-compose down

dev-auth:
	air -c .air.toml -- ./cmd/auth

dev-tickets:
	air -c .air.toml -- ./cmd/tickets

dev-orders:
	air -c .air.toml -- ./cmd/orders

dev-payments:
	air -c .air.toml -- ./cmd/payments

dev-expiration:
	go run ./cmd/expiration

# Database migrations
migrate-up:
	@for file in migrations/*.sql; do \
		echo "Running $$file..."; \
		psql $(DATABASE_URL) < $$file; \
	done

migrate-down:
	@echo "Manual rollback required - edit migrations as needed"

# Kubernetes
k8s-apply:
	kubectl apply -f infra/k8s/

k8s-delete:
	kubectl delete -f infra/k8s/

# Clean
clean:
	rm -rf tmp coverage.out coverage.html
	go clean
