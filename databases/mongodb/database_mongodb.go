package database_mongodb

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBInterface interface {
	GetCollection(name string) *mongo.Collection
	Disconnect(ctx context.Context) error
}

var client *mongo.Client
var mDB *mongo.Database

func InitMongoDB() {
	uri := "mongodb://localhost:27017"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("[ERROR | MONGO] %s", err.Error())
	}

	if pingErr := client.Ping(ctx, nil); pingErr != nil {
		log.Fatalf("[ERROR | MONGO] %s", err.Error())
	}

	mDB = client.Database("test")

	log.Println("Успешно подключено к MongoDB!")
}

func GetMongoDB() *mongo.Database {
	return mDB
}

func CloseMongoDB() {
	if err := client.Disconnect(context.TODO()); err != nil {
		log.Printf("[ERROR | MONGO] %s", err.Error())
		return
	}
}
