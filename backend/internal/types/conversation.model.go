package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Conversation struct {
	Participants []primitive.ObjectID `json:"participants,omitempty" bson:"participants"`
	Messages     []primitive.ObjectID `json:"messages,omitempty" bson:"messages"`
}
