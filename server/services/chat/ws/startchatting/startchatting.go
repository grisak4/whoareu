package startchatting

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
	mongodbmodels "whoareu/models/mongodb_models"
	"whoareu/utils/incrementids"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/mongo"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Chat struct {
	Clients map[int]*websocket.Conn
	Mutex   sync.Mutex
}

var chats = make(map[int]*Chat)
var chatsMutex sync.Mutex

func ConnectChat(c *gin.Context, mdb *mongo.Database) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Printf("[ERROR | WS] %s\n", err.Error())
		return
	}
	defer conn.Close()

	chatID, _ := strconv.Atoi(c.Param("chat_id"))
	userID, _ := strconv.Atoi(c.Param("user_id"))

	chatsMutex.Lock()
	if chats[chatID] == nil {
		chats[chatID] = &Chat{
			Clients: make(map[int]*websocket.Conn),
		}
	}
	chat := chats[chatID]
	chatsMutex.Unlock()

	chat.Mutex.Lock()
	chat.Clients[userID] = conn
	chat.Mutex.Unlock()

	defer func() {
		chat.Mutex.Lock()
		delete(chat.Clients, userID)
		if len(chat.Clients) == 0 {
			chatsMutex.Lock()
			delete(chats, chatID)
			chatsMutex.Unlock()
		}
		chat.Mutex.Unlock()
		conn.Close()
	}()

	mdbCol := mdb.Collection("messages")

	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("[ERROR | WS] %s\n", err.Error())
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		msgModel := mongodbmodels.Message{
			ID:        incrementids.IncrementID(mdbCol),
			UserID:    uint(userID),
			ChatID:    uint(chatID),
			Content:   string(msg),
			CreatedAt: time.Now(),
		}

		_, err = mdbCol.InsertOne(ctx, msgModel)
		if err != nil {
			log.Printf("[ERROR | MONGO] %v\n", err)
		}

		go func() {
			chat.Mutex.Lock()
			defer chat.Mutex.Unlock()
			for id, client := range chat.Clients {
				message := map[string]interface{}{
					"user_id":         userID,
					"chat_id":         chatID,
					"message_content": string(msg),
					"time_sent":       time.Now(),
				}
				if id != userID {
					jsonMessage, err := json.Marshal(message)
					if err != nil {
						log.Printf("[ERROR] Could not serialize message: %v", err)
						continue
					}
					_ = client.WriteMessage(messageType, jsonMessage)
				}
			}
		}()
	}
}
