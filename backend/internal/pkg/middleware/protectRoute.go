package middleware

import (
	"context"
	"errors"
	"time"

	"github.com/billzayy/chat-golang/internal/db"
	"github.com/billzayy/chat-golang/internal/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ProtectRoute(userId string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	client, err := db.ConnectMongo(ctx)

	if err != nil {
		return false, errors.New("error on connecting database ")
	}

	defer client.Disconnect(ctx)

	// Convert string into ObjectId
	filter, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		return false, err
	}

	userCollection := client.Database("chat-app").Collection("users")

	var result types.ResponseUSer

	// Get user from Mongo by UserId on token
	err = userCollection.FindOne(ctx, bson.M{"_id": filter}).Decode(&result)

	if err != nil {
		return false, err
	}

	return true, nil
}
