version=0.1.0

.PHONY: all

all:
	@echo "make <cmd>"
	@echo ""
	@echo "commands:"
	@echo "  build         - build the source code"
	@echo "  clean         - clean the build directory"
	@echo "  fmt           - format the source code with gofmt"
	@echo "  install       - install dependencies"
	@echo "  lint          - lint the source code"
	@echo "  test          - test the source code"

clean:
	@rm -rf ./build

lint:
	@go vet $(shell glide novendor)
	@go list ./... | grep -v /vendor/ | xargs -L1 golint

test:
	@go test $(shell glide novendor)

fmt:
	@go fmt $(shell glide novendor)

build: clean lint
	@go build $(shell glide novendor)

install:
	@go get -u github.com/golang/lint/golint
	@go get -u github.com/Masterminds/glide
	@glide install
