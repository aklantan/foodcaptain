package main

import (
	"net/http"

	"github.com/aklantan/foodcaptain/backend/handlers"
	"github.com/aklantan/foodcaptain/backend/models"
	"github.com/aklantan/foodcaptain/backend/sql/queries"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type apiConfig struct {
	db_query    *gorm.DB
	tokenSecret string
}

func (c *apiConfig) returnRestaurants(g *gin.Context) {
	var restaurants []models.Restaurant
	c.db_query.Find(&restaurants)
	g.JSON(http.StatusOK, gin.H{"restaurants": restaurants})

}

func main() {

	err := godotenv.Load("../.env")
	if err != nil {
		panic("Error loading env file")
	}
	db := queries.InitDB()

	apiCfg := &apiConfig{
		db_query: db,
	}

	server := gin.Default()
	server.Static("/static", "../frontend")
	server.LoadHTMLFiles("../frontend/index.html")
	server.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	server.GET("/test", handlers.TestConnection)
	server.GET("/restaurants", apiCfg.returnRestaurants)

	server.Run("127.0.0.1:8999")
}
