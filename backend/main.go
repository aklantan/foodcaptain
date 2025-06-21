package main

import (
	"fmt"
	"net/http"
	"strconv"

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
	limit := 0
	l := g.Query("limit")
	if l != "" {
		lint, err := strconv.Atoi(l)
		if err != nil {
			fmt.Println("Conversion failed")
			g.JSON(http.StatusBadRequest, gin.H{"message": "cannot convert limit to number"})
			return
		}
		limit = lint
	}
	if limit == 0 {
		g.JSON(http.StatusBadRequest, gin.H{"message": "please provide a limit"})
		return
	}
	var restaurants []models.Restaurant
	c.db_query.Order("RANDOM()").Limit(limit).Find(&restaurants)
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
