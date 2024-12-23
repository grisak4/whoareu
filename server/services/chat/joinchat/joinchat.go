package joinchat

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func JoinChat(c *gin.Context, mdb *mongo.Database) {
	col := mdb.Collection("chats")

	userId, _ := strconv.Atoi(c.Param("user_id"))
	chatId, _ := strconv.Atoi(c.Param("chat_id"))

	filter := bson.M{"_id": chatId}

	update := bson.M{
		"$push": bson.M{
			"participants_ids": userId,
		},
	}

	result, err := col.UpdateOne(c, filter, update)
	if err != nil {
		log.Printf("Ошибка обновления документа: %v", err)
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Успешный ответ
	c.JSON(201, gin.H{
		"message": "successfully joined!",
		"data":    result,
	})
}
