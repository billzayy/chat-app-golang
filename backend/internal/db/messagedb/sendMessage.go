package messagedb

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/billzayy/chat-golang/internal/db"
	"github.com/billzayy/chat-golang/internal/pkg"
	"github.com/billzayy/chat-golang/internal/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func SendMessage(cookieId string, pathVariableId string, message types.InputMessage) (types.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	client, err := db.ConnectMongo(ctx)

	if err != nil {
		return types.Message{}, errors.New("error on connecting database ")
	}

	defer client.Disconnect(ctx)

	// Convert string into ObjectID (Mongo _id Type)
	senderId, err := primitive.ObjectIDFromHex(cookieId)

	if err != nil {
		log.Println(err, "senderId")
		return types.Message{}, err
	}

	receiverId, err := primitive.ObjectIDFromHex(pathVariableId)

	if err != nil {
		log.Println(err, "receiverId")
		return types.Message{}, err
	}

	// Connect to Mongo
	conversationCollection := client.Database("chat-app").Collection("conversation")

	conversation, err := checkConversation(conversationCollection, senderId, receiverId)

	if err != nil {
		return types.Message{}, err
	}

	if !conversation {
		_, err := conversationCollection.InsertOne(ctx, types.Conversation{ // Create Conversation on Mongo
			Participants: []primitive.ObjectID{senderId, receiverId},
			Messages:     []primitive.ObjectID{},
		})

		if err != nil {
			return types.Message{}, errors.New("error on create conversation")
		}
	}

	newMessage := types.Message{
		SenderID:   senderId,
		ReceiverID: receiverId,
		Message:    message.Message,
		CreatedAt:  time.Now(),
	}

	messageId, err := createMessage(client, newMessage)

	if err != nil {
		return types.Message{}, errors.New("error on create message")
	}

	filter := bson.M{
		"participants": bson.M{
			"$all": []primitive.ObjectID{senderId, receiverId}}}

	update := bson.M{
		"$push": bson.M{"messages": messageId.Id}, // Example update operation
	}

	_, err = conversationCollection.UpdateOne(ctx, filter, update)

	if err != nil {
		return types.Message{}, errors.New("error on update conversation")
	}

	data, err := json.Marshal(messageId)

	if err != nil {
		return types.Message{}, errors.New("error on marshal messageId")
	}

	pkg.SendMessageToClient(receiverId, data)

	return messageId, nil
}

func checkConversation(collection *mongo.Collection, senderId primitive.ObjectID, receiverId primitive.ObjectID) (bool, error) {
	filter := bson.M{
		"participants": bson.M{
			"$all": []primitive.ObjectID{senderId, receiverId}}}

	cur, err := collection.Find(context.TODO(), filter)

	if err != nil {
		return false, errors.New("conversation not found")
	}

	var results []types.Conversation

	for cur.Next(context.TODO()) {
		var elem types.Conversation

		err := cur.Decode(&elem)

		if err != nil {
			return false, err
		}
		results = append(results, elem)
	}

	return results != nil, nil
}

func createMessage(client *mongo.Client, newMessage types.Message) (types.Message, error) {
	messageCollection := client.Database("chat-app").Collection("message")

	mongoResult, err := messageCollection.InsertOne(context.TODO(), newMessage)

	if err != nil {
		return types.Message{}, err
	}

	currentIdCreated := mongoResult.InsertedID.(primitive.ObjectID)

	var currentMessage types.Message

	err = messageCollection.FindOne(context.TODO(), bson.M{"_id": currentIdCreated}).Decode(&currentMessage)

	if err != nil {
		return types.Message{}, err
	}

	return currentMessage, nil
}
