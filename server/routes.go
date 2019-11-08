package server

import (
	"github.com/gin-gonic/gin"
)

func setupRoutes(router *gin.Engine) {
	v1 := router.Group("/api/v0.1")

	v1.GET("/", hello)

	v1.GET("/auctions/:id", getAuctionsById)
	v1.GET("/auctions/:id/bids", getBidsAuctionsById)
	v1.POST("users", addNewUser)
	v1.POST("/auctions/:id/users/:id/bids", bidAuctionByUser)
	v1.GET("/auctions/:id/result", getResultByAuctionId)

}
