package incrementids

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func IncrementID(col *mongo.Collection) uint {
	opts := options.FindOne().SetSort(bson.D{{Key: "_id", Value: -1}})

	var result struct {
		ID uint `bson:"_id"`
	}

	err := col.FindOne(context.TODO(), bson.D{}, opts).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 1
		}
	}

	return result.ID + 1
}
