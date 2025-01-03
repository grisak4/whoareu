package getmessages

import (
	"log"
	"strconv"
	mongodbmodels "whoareu/models/mongodb_models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetMessages(c *gin.Context, mdb *mongo.Database) {
	collection := mdb.Collection("messages")

	chatId, _ := strconv.Atoi(c.Param("chat_id"))

	filter := bson.D{{Key: "chat_id", Value: chatId}}

	cursor, err := collection.Find(c, filter)
	if err != nil {
		log.Printf("[ERROR] %s", err.Error())
		c.JSON(500, gin.H{
			"message": "found error",
			"error":   err,
		})
	}
	defer cursor.Close(c)

	var results []mongodbmodels.Message
	if err := cursor.All(c, &results); err != nil {
		log.Printf("[ERROR] %s", err.Error())
		c.JSON(500, gin.H{
			"message": "found error",
			"error":   err,
		})
	}

	c.JSON(200, gin.H{
		"message": "success",
		"data":    results,
	})
}
