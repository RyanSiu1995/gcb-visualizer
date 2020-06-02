build:
	go build -ldflags="-s -w" -o gcb-visualizer

.PHONY: test
test:
	go test -count=1 ./...

get:
	go get -v -t -d ./...

format-lint:
	gofmt -l -d .

import-lint:
	goimports -l -d .

style-lint:
	golint ./...

lint: format-lint import-lint style-lint
