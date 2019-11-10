package server

import (
	"github.com/gin-gonic/gin"
)

func setupRoutes(router *gin.Engine) {
	if mode := gin.Mode(); mode == gin.TestMode {
		router.LoadHTMLGlob("./../templates/**/*")
	} else {
		router.LoadHTMLGlob("templates/**/*")
	}
	router.Static("/js", "./static/js")
	router.Static("/css", "./static/css")
	router.GET("/hello", hello)

	router.GET("/", getAllAuctions)
	router.GET("/auctions/:auction_id", getAuctionsById)
	router.GET("/auctions/:auction_id/bids", getBidsAuctionsById)
	router.POST("/users", addNewUser)
	router.POST("/auctions/:auction_id/users/:user_id/bids", addBidAuctionIdByUserId)
	router.GET("/auctions/:auction_id/result", getResultByAuctionId)

}
