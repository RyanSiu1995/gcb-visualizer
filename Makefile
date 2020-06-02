build:
	go build -ldflags="-s -w" -o gcb-visualizer

.PHONY: test
test:
	go test -count=1 ./...

get:
	go get -v -t -d ./...

format-lint:
	errors=$$(gofmt -l .); if [ "$${errors}" != "" ]; then echo "Format Lint Error Files:\n$${errors}"; exit 1; fi

import-lint:
	errors=$$(goimports -l .); if [ "$${errors}" != "" ]; then echo "Import Lint Error Files:\n$${errors}"; exit 1; fi

style-lint:
	errors=$$(golint ./...); if [ "$${errors}" != "" ]; then echo "Style Lint Error Files:\n$${errors}"; exit 1; fi

lint: format-lint import-lint style-lint
