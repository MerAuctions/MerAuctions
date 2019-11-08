# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
BINARY_NAME=merauction

all: test build
build:
		GO111MODULE=on CGO_ENABLED=0 $(GOBUILD) -o $(BINARY_NAME) cmd/auctions/merauctions.go
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

docker: docker-push

docker-build:
		docker build -t gcr.io/kouzoh-p-harsh/merauctions:v1 .

docker-push: docker-build
		docker push gcr.io/kouzoh-p-harsh/merauctions:v1

kubernetes-delete:
	kubectl delete -f kubernetes