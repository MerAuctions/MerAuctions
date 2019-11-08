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
		"page_title":  "Test title",
		"name":        "マウジー　ニット",
		"description": "色はくすみピンク。丈が少し短めです。首元まであるのであったかいと思います。フリーサイズで新品未使用です。タグ付きです。元値6980円＋税。写真より薄めのピンクで普段使いしやすいです",
		"image_link":  "https://static.mercdn.net/item/detail/orig/photos/m67918189515_1.jpg?1573189516",
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
