.PHONY: lint vulncheck

lint:
	golangci-lint run

vulncheck:
	govulncheck ./...