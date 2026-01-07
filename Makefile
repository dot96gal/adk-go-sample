.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: govulncheck
govulncheck:
	govulncheck ./...

.PHONY: test
test:
	go test -race -v ./...

.PHONY: build
build:
	go build -o build/app ./...

.PHONY: dev
dev:
	dotenvx run -- go run ./...

.PHONY: dev-web
dev-web:
	dotenvx run -- go run ./... web api webui

.PHONY: run
run:
	dotenvx run -- ./build/app

