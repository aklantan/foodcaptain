package main

import (
	"github.com/aklantan/foodcaptain/backend/handlers"
	"github.com/aklantan/foodcaptain/backend/sql/queries"
	"github.com/gin-gonic/gin"
)

func main() {

	db := queries.InitDB()

	server := gin.Default()
	server.GET("/test", handlers.TestConnection)

	server.Run(":8999")
}
