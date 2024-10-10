package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FullName   string             `json:"fullName,omitempty" bson:"fullName, omitempty"`
	UserName   string             `json:"userName,omitempty" bson:"userName, omitempty"`
	Password   string             `json:"password,omitempty" bson:"password, omitempty"`
	Gender     string             `json:"gender,omitempty" bson:"gender,omitempty" validate:"oneof=male female"`
	ProfilePic string             `json:"profilePic,omitempty" bson:"profilePic, omitempty"`
	CreatedAt  time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt  time.Time          `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

type ResponseUSer struct {
	Id         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FullName   string             `json:"fullName,omitempty" bson:"fullName, omitempty"`
	UserName   string             `json:"userName,omitempty" bson:"userName, omitempty"`
	ProfilePic string             `json:"profilePic,omitempty" bson:"profilePic, omitempty"`
}

type RequestUser struct {
	UserName string `json:"userName,omitempty" bson:"userName, omitempty"`
	Password string `json:"password,omitempty" bson:"password, omitempty"`
}
