GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
BUILDFLAGS="-v"
TRIM="-trimpath"
BINARY_NAME=c3
BINARY_UNIX=$(BINARY_NAME)_unix

PREFIX=/usr/local

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) $(BUILDFLAGS) $(TRIM) ./cmd/c3

test:
	$(GOTEST) -v ./...

install:
	cp $(BINARY_NAME) $(PREFIX)/bin/

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
