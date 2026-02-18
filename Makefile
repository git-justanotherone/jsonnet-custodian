build: go-generate
	CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o bin/custodian ./cmd/custodian

test:
	go test --cover ./...

go-generate:
	go generate ./...

PHONY: build test go-generate