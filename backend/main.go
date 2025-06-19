package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func testConnection(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func main() {

	server := gin.Default()
	server.GET("/test", testConnection)

	server.Run(":8999")
}
