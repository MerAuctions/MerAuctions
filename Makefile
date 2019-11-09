# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
BINARY_NAME=merauction
CLUSTER_NAME=merauction
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
		docker build -t gcr.io/kouzoh-p-harsh/merauctions:v0.1 .

docker-push: docker-build
		docker push gcr.io/kouzoh-p-harsh/merauctions:v0.1

cluster-create:
	gcloud container clusters create merauction --num-nodes=1 --machine-type=g1-small

kubernetes-build:
	gcloud container clusters get-credentials $(CLUSTER_NAME)
	kubectl apply -f kubernetes
	# kubectl exec mongo -c mongo -- mongo --eval 'db.getSiblingDB("admin").createUser({user:"main_admin",pwd:"'"$(DB_PASSWORD)"'",roles:[{role:"root",db:"admin"}]});'
