
test-all:
	go test ./...

test-domain:
	go test ./internal/domain/... -v

test-application:
	go test ./internal/application/... -v

test-interface:
	go test ./internal/interfaceadapters/... -v

test-infrastructure:
	go test ./internal/infrastructure/... -v

cover:
	go test -coverprofile=coverage.out ./... -v
	go tool cover -html=coverage.out

.PHONY: test test-domain test-application test-all
