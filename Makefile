.PHONY: build
build:
	go build -ldflags="-s -w"

.PHONY: tests
tests:
	go test -coverprofile coverage.txt -count=1 ./...
