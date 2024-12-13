package getchats

import (
	"log"
	"strconv"
	mongodbmodels "whoareu/models/mongodb_models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAllChats(c *gin.Context, mdb *mongo.Database) {
	col := mdb.Collection("chats")

	userId, _ := strconv.Atoi(c.Param("user_id"))

	filter := bson.D{{Key: "participants_ids", Value: userId}}

	cursor, err := col.Find(c, filter)
	if err != nil {
		log.Printf("[ERROR] %s", err.Error())
		c.JSON(500, gin.H{
			"message": "found error",
			"error":   err,
		})
	}
	defer cursor.Close(c)

	var results []mongodbmodels.Chat
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
