package incrementid

import (
	"log"

	mongodbmodels "whoareu/models/mongodb_models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func IncrementIDMessages(c *gin.Context, mdb *mongo.Database) uint {
	filter := bson.D{}

	findOptions := options.FindOne().SetSort(bson.D{{Key: "message_id", Value: -1}})

	var lastMessage mongodbmodels.Message

	err := mdb.Collection("messages").FindOne(c, filter, findOptions).Decode(&lastMessage)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 1
		}
		log.Printf("[ERROR] %s", err.Error())
		return 0
	}

	return lastMessage.ID + 1
}
