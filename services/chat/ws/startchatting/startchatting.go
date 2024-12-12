package startchatting

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
	mongodbmodels "whoareu/models/mongodb_models"

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

func ConnectChat(c *gin.Context, mdb *mongo.Client) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Printf("[ERROR | WS] %s", err.Error())
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

	mdbCol := mdb.Database("test").Collection("messages")

	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("[ERROR | WS] %s", err.Error())
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		msgModel := mongodbmodels.MongoDB_Message{
			UserID:  uint(userID),
			ChatID:  uint(chatID),
			Content: string(msg),
		}

		_, err = mdbCol.InsertOne(ctx, msgModel)
		if err != nil {
			log.Printf("[ERROR | MONGO] %v", err)
		}

		go func() {
			chat.Mutex.Lock()
			for id, client := range chat.Clients {
				if id != userID {
					_ = client.WriteMessage(messageType, msg)
				}
			}
			chat.Mutex.Unlock()
		}()
	}
}
