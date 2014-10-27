.PHONY: build
build: deps
	go build -v

.PHONY: deps
deps:
	go get -v

.PHONY: test
test: build
	go test -v

.PHONY: clean
clean:
	go clean
