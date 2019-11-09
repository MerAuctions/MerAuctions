package server

import (
	"net/http"
	"time"

	"github.com/MerAuctions/MerAuctions/models"
	"github.com/gin-gonic/gin"
)

func hello(c *gin.Context) {
	c.String(200, "Hello World")
}

func getAllAuctions(c *gin.Context) {
	c.String(200, "Auctions1")
}

func getAuctionsById(c *gin.Context) {
	// ⬇️for debugging
	bid := models.Bid{"bidid", "auctionid", "userid", 233, models.TimeStamp(time.Now().Unix() * 1000)}
	top_5_bids := models.BidList{
		bid,
		bid,
		bid,
		bid,
		bid,
	}
	c.HTML(http.StatusOK, "auction/index.tmpl", gin.H{
		"page_title":  "Test title",
		"name":        "マウジー　ニット",
		"description": "色はくすみピンク。丈が少し短めです。首元まであるのであったかいと思います。フリーサイズで新品未使用です。タグ付きです。元値6980円＋税。写真より薄めのピンクで普段使いしやすいです",
		"image_link":  "https://static.mercdn.net/item/detail/orig/photos/m67918189515_1.jpg?1573189516",
		"end_time":    1573292410953,
		"top_5_bids":  top_5_bids,
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
