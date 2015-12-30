GO_SRC_TAR=go1.5.2.src.tar.gz
GO_SRC_URL=https://storage.googleapis.com/golang

PYTHON_SRC_TAR=Python-3.5.1.tgz
PYTHON_SRC_URL=https://www.python.org/ftp/python/3.5.1

.PHONY: deps
deps: dataset
	pip install -r requirements.txt

.PHONY: test
test:
	nosetests --nocapture -d

.PHONY: clean
clean:
	rm -rf dataset/$(GO_SRC_TAR) dataset/$(PYTHON_SRC_TAR) dataset/go dataset/notgo

.PHONY: dataset
dataset: dataset/go dataset/notgo

dataset/go: dataset/$(GO_SRC_TAR)
	mkdir -p $@/tree
	tar xvf $< -C $@/tree '*.go'
	find $@/tree -type f -exec mv {} $@ \;
	rm -rf $@/tree

dataset/notgo: dataset/$(PYTHON_SRC_TAR)
	mkdir -p $@/tree
	tar xvf $< -C $@/tree '*.py' '*.rst' '*.c' '*.h' '*.txt'
	find $@/tree -type f -exec mv {} $@ \;
	rm -rf $@/tree

dataset/$(GO_SRC_TAR):
	curl -f -o $@ $(GO_SRC_URL)/$(GO_SRC_TAR)

dataset/$(PYTHON_SRC_TAR):
	curl -f -o $@ $(PYTHON_SRC_URL)/$(PYTHON_SRC_TAR)
