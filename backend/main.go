package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/aklantan/foodcaptain/backend/handlers"
	"github.com/aklantan/foodcaptain/backend/models"
	"github.com/aklantan/foodcaptain/backend/sql/queries"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type apiConfig struct {
	db_query       *gorm.DB
	tokenSecret    string
	socketUpgrader websocket.Upgrader
	clientList     map[*websocket.Conn]bool
}

func (c *apiConfig) webSocketHandler(g *gin.Context) {
	ws, err := c.socketUpgrader.Upgrade(g.Writer, g.Request, nil)
	if err != nil {
		fmt.Println("Cannot upgrade connection")
		return
	}
	defer ws.Close()

	c.clientList[ws] = true

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("read error:", err)
			delete(c.clientList, ws)
			break
		}

		for client := range c.clientList {
			if string(msg) == "Walls" {
				if err := client.WriteMessage(websocket.TextMessage, []byte("Sausages")); err != nil {
					fmt.Println("broadcast error:", err)
					client.Close()
					delete(c.clientList, client)
				}
				break

			}
			if err := client.WriteMessage(websocket.TextMessage, msg); err != nil {
				fmt.Println("broadcast error:", err)
				client.Close()
				delete(c.clientList, client)
			}
		}
	}
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

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	err := godotenv.Load("../.env")
	if err != nil {
		panic("Error loading env file")
	}
	db := queries.InitDB()
	var clients = make(map[*websocket.Conn]bool)

	apiCfg := &apiConfig{
		db_query:       db,
		socketUpgrader: upgrader,
		clientList:     clients,
	}

	server := gin.Default()
	server.Static("/static", "../frontend")
	server.LoadHTMLFiles("../frontend/index.html")
	server.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	server.GET("/test", handlers.TestConnection)
	server.GET("/restaurants", apiCfg.returnRestaurants)
	server.GET("/ws", apiCfg.webSocketHandler)

	server.Run("127.0.0.1:8999")
}
