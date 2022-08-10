.PHONY: start
start:
	docker-compose build
	docker-compose up dns

.PHONY: start-telemetry
start-telemetry:
	docker-compose up

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

.PHONY: bench
bench:
	docker run \
		-it \
		--rm \
		-w /app \
		-v ${PWD}/k6:/app \
		--network dns \
		loadimpact/k6 run test.js

.PHONY: lint
lint:
	docker run \
		--rm \
		-ti \
		-w /app \
		-v $(PWD):/app \
		golangci/golangci-lint:v1.45-alpine golangci-lint run