package createchat

import (
	"log"
	mongodbmodels "whoareu/models/mongodb_models"
	"whoareu/utils/incrementids"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateChat(c *gin.Context, mdb *mongo.Database) {
	col := mdb.Collection("chats")

	var newchat mongodbmodels.Chat

	if err := c.BindJSON(&newchat); err != nil {
		c.JSON(400, gin.H{
			"message": "incorrect data",
		})
		return
	}

	newchat.ID = incrementids.IncrementID(col)

	if _, err := col.InsertOne(c, newchat); err != nil {
		log.Printf("[ERROR] %s", err.Error())
		c.JSON(500, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(201, gin.H{
		"message": "successfully created a chat!",
	})
}
