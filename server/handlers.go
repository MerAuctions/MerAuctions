package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func hello(c *gin.Context) {
	c.String(200, "Hello World")
}

func getAllAuctions(c *gin.Context) {

}

func getAuctionsById(c *gin.Context) {
	c.HTML(http.StatusOK, "auction/index.tmpl", gin.H{
		"title": "Test title",
	})
}

func getBidsAuctionsById(c *gin.Context) {

}

func addNewUser(c *gin.Context) {

}

func bidAuctionIdByUserId(c *gin.Context) {

}

func getResultByAuctionId(c *gin.Context) {

}
