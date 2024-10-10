package db

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongo(ctx context.Context) (*mongo.Client, error) {
	err := godotenv.Load("./internal/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB_URI")))

	if err != nil {
		return client, err
	}

	return client, err
}
