package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"sync"

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
	sessionList    map[string]*Session
}

type Session struct {
	ID      string
	Clients map[*websocket.Conn]bool
	Mutex   sync.Mutex
}

func (c *apiConfig) webSocketHandler(g *gin.Context) {
	sessionID := g.Query("sessionID")
	fmt.Println(sessionID)

	if sessionID == "" {
		fmt.Println("No sessionID provided, generating a new one")
		// generate...
	} else {
		fmt.Println("Using sessionID from query:", sessionID)
	}

	if sessionID == "" {
		// rand.Seed(time.Now().UnixNano()) // Seed with current time for randomness
		min := int64(100000000000) // Smallest 12-digit number
		max := int64(999999999999) // Largest 12-digit number
		num := rand.Int63n(max-min+1) + min
		sessionID = "sessno" + strconv.FormatInt(num, 10)
	}
	fmt.Println(sessionID)

	if _, exists := c.sessionList[sessionID]; !exists {
		newSession := &Session{
			ID:      sessionID,
			Clients: make(map[*websocket.Conn]bool),
		}
		c.sessionList[sessionID] = newSession

	}

	ws, err := c.socketUpgrader.Upgrade(g.Writer, g.Request, nil)
	if err != nil {
		fmt.Println("Cannot upgrade connection")
		return
	}
	defer ws.Close()

	session := c.sessionList[sessionID]
	session.Mutex.Lock()
	session.Clients[ws] = true
	session.Mutex.Unlock()

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("read error:", err)
			delete(session.Clients, ws)
			break
		}

		for client := range session.Clients {
			if string(msg) == "Walls" {
				if err := client.WriteMessage(websocket.TextMessage, []byte("Sausages")); err != nil {
					fmt.Println("broadcast error:", err)
					client.Close()
					delete(session.Clients, client)
				}
				break

			}
			if err := client.WriteMessage(websocket.TextMessage, msg); err != nil {
				fmt.Println("broadcast error:", err)
				client.Close()
				delete(session.Clients, client)
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
	var sessions = make(map[string]*Session)

	apiCfg := &apiConfig{
		db_query:       db,
		socketUpgrader: upgrader,
		sessionList:    sessions,
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
