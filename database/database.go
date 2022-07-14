package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetDatabase() (collection *mongo.Collection, ctx context.Context) {
	context, _ := context.WithCancel(context.Background()) //context가 cancel 혹은 timeout으로 종료되면 context의 done이 호출
	// Set client options
	clientOptions := options.Client().ApplyURI(`mongodb://leechanhui:qwer1234@localhost:20000/?connect=direct`)
	clientOptions.SetAuth(options.Credential{
		Username: "leechanhui",
		Password: "qwer1234",
	})
	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
		return
	}

	// Check the connection
	err = client.Ping(ctx, nil)

	if err != nil {
		log.Fatal(err)
		return
	}

	collection = client.Database("board").Collection("people")

	return collection, context
}
