package main

import (
	"github.com/aklantan/foodcaptain/backend/handlers"
	"github.com/gin-gonic/gin"
)

func main() {

	server := gin.Default()
	server.GET("/test", handlers.TestConnection)

	server.Run(":8999")
}
