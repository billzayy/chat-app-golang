package userdb

import (
	"context"
	"errors"
	"time"

	"github.com/billzayy/chat-golang/internal/db"
	"github.com/billzayy/chat-golang/internal/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAllUsers(cookieId string) ([]types.ResponseUSer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	client, err := db.ConnectMongo(ctx)

	if err != nil {
		return []types.ResponseUSer{}, errors.New("error on connecting database ")
	}

	defer client.Disconnect(ctx)

	// Convert string into ObjectID (Mongo _id Type)
	userId, err := primitive.ObjectIDFromHex(cookieId)

	if err != nil {
		return []types.ResponseUSer{}, err
	}

	// Connect to Mongo
	userCollection := client.Database("chat-app").Collection("users")

	var results []types.ResponseUSer

	cur, err := userCollection.Find(context.TODO(), bson.M{"_id": bson.M{"$ne": userId}}, options.Find().SetProjection(bson.M{"password": 0, "createdAt": 0, "updatedAt": 0}))

	if err != nil {
		return []types.ResponseUSer{}, errors.New("user not found")
	}

	for cur.Next(context.TODO()) {
		var elem types.ResponseUSer

		err := cur.Decode(&elem)

		if err != nil {
			return []types.ResponseUSer{}, err
		}
		results = append(results, elem)
	}

	if err := cur.Err(); err != nil {
		return []types.ResponseUSer{}, err
	}
	cur.Close(context.TODO())

	return results, nil
}
