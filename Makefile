.PHONY: lint lint-c format-c vulncheck build-webapi test-webapi

build:
	go build ./...

build-webapi:
	CGO_ENABLED=0 go build ./webapi/...

test-webapi:
	CGO_ENABLED=0 go test -race ./webapi/...

lint:
	golangci-lint run

lint-c:
	./scripts/cgo-clang-format.sh --check

format-c:
	./scripts/cgo-clang-format.sh --fix

vulncheck:
	govulncheck ./...