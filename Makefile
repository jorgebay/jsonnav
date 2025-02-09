.PHONY: test
test:
	go test -v ./...

.PHONY: lint
lint:
	./bin/golangci-lint run ./... --verbose

.PHONY: build
build:
	go build ./...

.PHONY: install-tools
GOLANGCI_LINT_VERSION=v1.61.0
install-tools:
	curl -sSfL \
		"https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh" \
		| sh -s -- -b bin ${GOLANGCI_LINT_VERSION}
