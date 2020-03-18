package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gladmo/kafka2mongo/conf"
)

var client *mongo.Client

func init() {
	var err error

	dsn := conf.Getenv("MONGODB_DSN", "mongodb://localhost:27017")
	opt := options.Client().ApplyURI(dsn)
	client, err = mongo.Connect(context.Background(), opt)
	if err != nil {
		panic(err)
	}
}

type Document struct {
	Database   string      `json:"database"`
	Collection string      `json:"collection"`
	Doc        interface{} `json:"doc"`
}

func Store(msg *Document) (err error) {
	collection := client.Database(msg.Database).Collection(msg.Collection)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, msg.Doc)
	if err != nil {
		return
	}
	return
}
