package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoClient() {
	url := "mongodb://192.168.56.101:27017"
	c, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	err = c.Connect(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

}
