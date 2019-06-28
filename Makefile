clean:
	rm -rf dist

build: clean test
	go build semver.go

test: clean
	go test .

release:
	goreleaser