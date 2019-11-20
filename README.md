# merAuc

[![CircleCI](https://circleci.com/gh/MerAuctions/MerAuctions/tree/master.svg?style=svg)](https://circleci.com/gh/MerAuctions/MerAuctions/tree/master)

### Code Chrysalis X Mercari Legacy Project by [Sahil](https://github.com/sahil505), [Aniket](https://github.com/aniket1743), [Liu](https://github.com/Rocuku) & [Shashank](https://github.com/shasjakka0390)

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

- If you want to populate initial data you can run `cd server/seed-data/ && go run seeddb.go`

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

## Essential Features

- Sellers can start new auctions
- Tags and Description auto generation
  - [CAP-THAT-PIC](https://github.com/CoolDogee/cap-that-pic)
- Search based on tags
- Rewards system for bidding

## Engineering Challenges

- Get public URL for images uploaded by users
- Connection with remote CAP-THAT-PIC server
- Poor MongoDB support for GoLang

## Future Goals

- Personalized Dashboard : Recommendations based on user’s interests
- User Profile
- Trending Auctions
- Google Signup
- Filters (Base Price, End Time, etc.)

## StyleGuide for Codebase

Please refer to the StyleGuide [here](https://docs.google.com/document/d/1IYPQ_jkVBBcz7EZHN6GanW2wtXnp1Lvu0Xkz9rLg6hk/edit?usp=sharing)
