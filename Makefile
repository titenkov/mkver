.PHONY: clean
clean:
	rm -rf dist

.PHONY: setup
setup:
	go get -u github.com/urfave/cli

.PHONY: build
build: clean setup
	go build

.PHONY: test
test:
	go test .