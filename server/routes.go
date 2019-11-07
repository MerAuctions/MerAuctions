package server

import (
	"github.com/gin-gonic/gin"
)

func setupRoutes(router *gin.Engine) {
	// It is good practice to version your API from the start
	v1 := router.Group("/api/v0.1")

	v1.GET("/", hello)

}
