clean:
	rm -rf dist

build: clean test
	go build mkver.go

test: clean
	go test .

release:
	goreleaser