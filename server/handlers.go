package server

import "github.com/gin-gonic/gin"

func hello(c *gin.Context) {
	c.String(200, "Hello World")
}

func getAllAuctions(c *gin.Context) {
	c.String(200, "Auctions2")
}

func getAuctionsById(c *gin.Context) {

}

func getBidsAuctionsById(c *gin.Context) {

}

func addNewUser(c *gin.Context) {

}

func bidAuctionIdByUserId(c *gin.Context) {

}

func getResultByAuctionId(c *gin.Context) {

}
