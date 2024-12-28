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
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("[ERROR | WS] %s\n", err.Error())
			return
		}

		var messageData struct {
			ChatID         int    `json:"chat_id"`
			MessageContent string `json:"message_content"`
			UserID         int    `json:"user_id"`
			SenderName     string `json:"sender_name"`
		}

		if err := json.Unmarshal(msg, &messageData); err != nil {
			log.Printf("[ERROR | WS] Invalid JSON format: %s\n", err)
			continue
		}

		// Save message to MongoDB
		msgModel := mongodbmodels.Message{
			ID:         incrementids.IncrementID(mdbCol),
			UserID:     uint(messageData.UserID),
			ChatID:     uint(chatID),
			SenderName: messageData.SenderName,
			Content:    messageData.MessageContent,
			CreatedAt:  time.Now(),
		}

		_, err = mdbCol.InsertOne(context.Background(), msgModel)
		if err != nil {
			log.Printf("[ERROR | MONGO] %v\n", err)
		}

		// Add time_sent and message_id for client
		messageDataWithID := map[string]interface{}{
			"message_id":      msgModel.ID,
			"chat_id":         messageData.ChatID,
			"message_content": messageData.MessageContent,
			"user_id":         messageData.UserID,
			"sender_name":     messageData.SenderName,
			"time_sent":       msgModel.CreatedAt,
		}

		response, _ := json.Marshal(messageDataWithID)

		// Broadcast message to other clients
		go func() {
			chat.Mutex.Lock()
			defer chat.Mutex.Unlock()
			for id, client := range chat.Clients {
				if id != userID {
					_ = client.WriteMessage(websocket.TextMessage, response)
				}
			}
		}()
	}
}
