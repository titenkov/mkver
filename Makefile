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
	go test .