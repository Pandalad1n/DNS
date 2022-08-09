.PHONY: start
start:
	docker-compose up

.PHONY: start-web
start-web:
	docker-compose up dns

.PHONY: build
build:
	docker-compose build

.PHONY: test
test:
	docker run \
		-it \
		--rm \
		-w /app \
		-v ${PWD}:/app \
		golang:1.17 go test ./... -race -timeout 2m

.PHONY: lint
lint:
	docker run \
		--rm \
		-ti \
		-w /app \
		-v $(PWD):/app \
		golangci/golangci-lint:v1.45-alpine golangci-lint run