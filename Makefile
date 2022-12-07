CORPUS_DIR = ./corpus

all: fmt lint

install:
	go install github.com/neilpa/yajsv@latest
	go install github.com/ORKL/library/olm/olm@latest
	go install github.com/google/yamlfmt/cmd/yamlfmt@latest

lint:
	yajsv -q -s schema.json corpus/*.yaml
	yamlfmt -lint corpus/*.yaml

fmt: $(CORPUS_DIR)/*.yaml
	yamlfmt $(CORPUS_DIR)/*.yaml
