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
	router.GET("/auctions/:id", getAuctionsById)
	router.GET("/auctions/:id/bids", getBidsAuctionsById)
	router.POST("/users", addNewUser)
	router.POST("/auctions/:id/users/:id/bids", bidAuctionIdByUserId)
	router.GET("/auctions/:id/result", getResultByAuctionId)

}
