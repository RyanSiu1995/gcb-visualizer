build:
	go build -ldflags="-s -w"

tests:
	go test -count=1 ./...
