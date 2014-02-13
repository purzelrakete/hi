LIBS := github.com/purzelrakete/hi/fftw

.PHONY: all
all: build

.PHONY: build
build: deps
	go build -v $(LIBS)

.PHONY: deps
deps:
	go get -v
	go get -v github.com/golang/lint/golint

.PHONY: test
test: deps check
	go test -v $(LIBS)

.PHONY: render
render: test
	open spectrogram.png

.PHONY: check
check:
	go tool vet . 2>&1 | wc -l | { grep 0 || { go tool vet . && false; }; }
	if find . -name '*.go' | xargs golint | grep ":"; then false; else true; fi

.PHONY: clean
clean:
	go clean
	rm spectrogram.png
