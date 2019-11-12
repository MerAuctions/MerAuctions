# MerAuctions

[![CircleCI](https://circleci.com/gh/MerAuctions/MerAuctions.svg?style=svg)](https://circleci.com/gh/MerAuctions/MerAuctions)

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
make run
```
* If you want to populate initial data you can run cd server/seed-data/
## Minimal Viable Product (MVP) [Using a User Story]


## Essential Features


## Technologies Used

## Challenges

## Tasks & Assignment

- Sahil
  - Build the User Interface structure so that a user can upload the image from local machine or provide an image URL.
  - Try to connect the application to Instagram/Facebook/Twitter.
- Aniket
  - Make endpoints and write test cases to fetch data from MusixMatch API.
- Liu
  - Implement the algorithm to generate caption by using the lyrics obtained from MusixMatch API.
- Shashank
  - Make endpoints and write test cases to fetch data from Azure API.