# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
BINARY_NAME=merauction
CLUSTER_NAME=merauc-cluster-1
PROJECT_NAME=kouzoh-p-s-liu
#PROJECT_NAME=kouzoh-p-sahilkhokhar
ZONE=asia-northeast1-a
DB_PASSWORD = ${MONGODB_PASSWORD}

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
		docker build -t gcr.io/$(PROJECT_NAME)/merauctions:v0.3 .

docker-push: docker-build
		docker push gcr.io/$(PROJECT_NAME)/merauctions:v0.3

cluster-create:
	gcloud container clusters create $(CLUSTER_NAME) --num-nodes=2 --machine-type=g1-small

kubernetes-build:
	gcloud container clusters get-credentials $(CLUSTER_NAME)  --zone=$(ZONE)
	kubectl apply -f kubernetes

kubernetes-delete:
	gcloud container clusters delete $(CLUSTER_NAME)  --zone=$(ZONE)
