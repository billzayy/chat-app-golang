package messagedb

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/billzayy/chat-golang/internal/db"
	"github.com/billzayy/chat-golang/internal/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetMessage(cookieId string, pathVariableId string) ([]types.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	client, err := db.ConnectMongo(ctx)

	if err != nil {
		return []types.Message{}, errors.New("error on connecting database ")
	}

	defer client.Disconnect(ctx)

	// Convert string into ObjectID (Mongo _id Type)
	senderId, err := primitive.ObjectIDFromHex(cookieId)

	if err != nil {
		log.Println(err, "senderId")
		return []types.Message{}, err
	}

	receiverId, err := primitive.ObjectIDFromHex(pathVariableId)

	if err != nil {
		log.Println(err, "receiverId")
		return []types.Message{}, err
	}

	// Connect to Mongo
	conversationCollection := client.Database("chat-app").Collection("conversation")

	var conversations types.Conversation

	err = conversationCollection.FindOne(context.TODO(), bson.M{"participants": bson.M{"$all": []primitive.ObjectID{senderId, receiverId}}}).Decode(&conversations)

	if err != nil {
		return []types.Message{}, errors.New("conversation not found")
	}

	result, err := fetchMessage(ctx, conversations.Messages)

	if err != nil {
		return result, err
	}

	return result, nil
}

func fetchMessage(ctx context.Context, slice []primitive.ObjectID) ([]types.Message, error) {
	var result []types.Message

	client, err := db.ConnectMongo(ctx)

	if err != nil {
		return result, err
	}

	collection := client.Database("chat-app").Collection("message")

	for _, v := range slice {

		var finding types.Message

		err := collection.FindOne(context.TODO(), bson.M{"_id": v}).Decode(&finding)

		if err != nil {
			return result, err
		}

		result = append(result, finding)
	}

	defer client.Disconnect(ctx)

	return result, nil
}
