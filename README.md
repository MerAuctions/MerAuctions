# MerAuctions
 
[![CircleCI](https://circleci.com/gh/MerAuctions/MerAuctions/tree/master.svg?style=svg)](https://circleci.com/gh/MerAuctions/MerAuctions/tree/master)

### Mercari Greenfield Project by [Harsh](https://github.com/sipian), [Deepak](https://github.com/deepakbhatt329), [Vamshi](https://github.com/vamshiteja) & [Bowei](https://github.com/b-z)

---

## What is it?

MerAuctions is an online live auction portal where users can register for multiple auction and submit bids to buy a product. When the auction ends, the bidder who chose the highest value is declared the winner.

## How do run locally?

### Install MongoDB
```
brew tap mongodb/brew
brew install mongodb-community@4.2
mongod --config /usr/local/etc/mongod.conf --fork
ps -ef | grep mongod #check if mongodb is working
```
### Start the App
```
go run cmd/auctions/merauctions.go --mongodb-url=mongodb://localhost:27017/testing --database=testing
```
* If you want to populate initial data you can run `cd server/seed-data/ && go run seeddb.go`

## Technologies Used
### Front-end
  - Materialize CSS
  - jQuery
  - Go templates

### Back-end
  - ginkgo
  - JWT
  
### Database
  - MongoDB


### Deployment
  - Kubernetes
  - CircleCI
  - Docker
  - GCP
  
