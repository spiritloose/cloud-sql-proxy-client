export PATH := $(GOPATH)/bin:$(PATH)

xargs := $(shell which gxargs xargs | head -n1)

all: fmt imports lint vet build

build:
	go build ./...

fmt:
	find . -name '*.go' | grep -vE '/vendor/' | $(xargs) gofmt -l | $(xargs) -r false

imports:
	find . -name '*.go' | grep -vE '/vendor/' | $(xargs) goimports -l | $(xargs) -r false

lint:
	golint ./... | grep -vE '^vendor/' | $(xargs) -r false

vet:
	go vet ./...

test:
	go test -v ./...

install:
	go install ./...

clean:
	go clean

.PHONY: all build fmt imports lint vet test install clean
