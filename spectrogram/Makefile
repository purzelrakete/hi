.PHONY: build render test clean

build: deps
	go build -v

test: build
	go test -v

render: test
	open spectrogram.png

deps:
	go get -v

clean:
	go clean
