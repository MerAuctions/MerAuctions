# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
BINARY_NAME=merauctions

all: test build
build:
		CGO_ENABLED=0 $(GOBUILD) -o $(BINARY_NAME) -v cmd/auctions/merauctions.go
test: 
		$(GOTEST) -v ./...
clean: 
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
		rm -f $(BINARY_UNIX)
run:
		$(GORUN) cmd/auctions/merauctions.go
deps:
		$(GOMOD) download