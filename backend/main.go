package main

import (
	"fmt"

	"github.com/aklantan/foodcaptain/backend/handlers"
	"github.com/aklantan/foodcaptain/backend/models"
	"github.com/aklantan/foodcaptain/backend/sql/queries"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load("../.env")
	if err != nil {
		panic("Error loading env file")
	}
	db := queries.InitDB()

	var restaurant []models.Restaurant
	result := db.Find(&restaurant)

	fmt.Println(result)

	server := gin.Default()
	server.GET("/test", handlers.TestConnection)

	server.Run(":8999")
}
