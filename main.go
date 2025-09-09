package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	// Get all habits
	server.GET("/habits", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"habit": "Successfully fetched",
		})
	})

	server.Run(":8080")

}
