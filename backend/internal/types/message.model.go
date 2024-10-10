package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	Id         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	SenderID   primitive.ObjectID `json:"senderId" bson:"senderId"`
	ReceiverID primitive.ObjectID `json:"receiverId" bson:"receiverId"`
	Message    string             `json:"message" bson:"message"`
	CreatedAt  time.Time          `json:"createdAt" bson:"createdAt"`
}
