GOPATH=$(shell go env GOPATH)
GOLANGCI_LINT=$(GOPATH)/bin/golangci-lint

.PHONY: clean
clean:
	rm -rf dist

.PHONY: setup
setup:
	go get

.PHONY: build
build: clean setup
	go build

.PHONY: test
test:
	@echo "==> Running tests"
	go test .

.PHONY: lint
lint: $(GOLANGCI_LINT)
	@echo "==> Linting codebase"
	@$(GOLANGCI_LINT) run