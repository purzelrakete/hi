.PHONY: build render test clean

build:
	go build -v

render: build
	percolate && open material.png

test:
	go build -v
	go test -v

clean:
	go clean
	rm material.png
