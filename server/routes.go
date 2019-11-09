package server

import (
	"github.com/gin-gonic/gin"
)

func setupRoutes(router *gin.Engine) {

	router.GET("/hello", hello)

	router.GET("/", getAllAuctions)
	router.GET("/auctions/:auction_id", getAuctionsById)
	router.GET("/auctions/:auction_id/bids", getBidsAuctionsById)
	router.POST("/users", addNewUser)
	router.POST("/auctions/:auction_id/users/:user_id/bids", addBidAuctionIdByUserId)
	router.GET("/auctions/:auction_id/result", getResultByAuctionId)

}
