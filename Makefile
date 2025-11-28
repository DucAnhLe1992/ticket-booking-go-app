.PHONY: all build docker

all: build

build:
	go build ./cmd/...

docker-build-payments:
	docker build -t ghcr.io/ducanhle1992/ticket-payments:dev -f ./cmd/payments/Dockerfile .

docker-build-orders:
	docker build -t ghcr.io/ducanhle1992/ticket-orders:dev -f ./cmd/orders/Dockerfile .

docker-build-all: docker-build-payments docker-build-orders
