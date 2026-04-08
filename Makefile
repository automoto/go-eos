.PHONY: lint lint-c format-c vulncheck

lint:
	golangci-lint run

lint-c:
	./scripts/cgo-clang-format.sh --check

format-c:
	./scripts/cgo-clang-format.sh --fix

vulncheck:
	govulncheck ./...