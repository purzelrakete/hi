.PHONY: build render test clean deps

build: deps
	go build -v

render: build
	percolate && open material.png

deps:
	go get -v

test: build
	go test -v

clean:
	go clean
	rm material.png
