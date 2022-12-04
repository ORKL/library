CORPUS_DIR = ./corpus

all: validate

install:
	go install github.com/neilpa/yajsv@latest
	go install github.com/ORKL/library/olm/olm@latest
	go install github.com/google/yamlfmt/cmd/yamlfmt@latest

validate: $(CORPUS_DIR)/*.yaml
	yamlfmt $(CORPUS_DIR)
	yajsv -q -s schema.json corpus/*.yaml
