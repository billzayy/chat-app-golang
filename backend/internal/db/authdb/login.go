package authdb

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/billzayy/chat-golang/internal/db"
	"github.com/billzayy/chat-golang/internal/pkg/middleware"
	"github.com/billzayy/chat-golang/internal/types"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func Login(input types.RequestUser, w http.ResponseWriter) (types.ResponseUSer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	client, err := db.ConnectMongo(ctx)

	defer client.Disconnect(ctx)

	if err != nil {
		return types.ResponseUSer{}, errors.New("error on connecting database ")

	}

	userCollection := client.Database("chat-app").Collection("users")

	var findOne types.User

	err = userCollection.FindOne(ctx, bson.M{"userName": input.UserName}).Decode(&findOne)

	if err != nil {
		return types.ResponseUSer{}, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(findOne.Password), []byte(input.Password))

	if err != nil {
		return types.ResponseUSer{}, errors.New("incorrect password")
	}

	_, err = middleware.GenerateTokenAndSetCookie(findOne.Id, w)

	if err != nil {
		return types.ResponseUSer{}, errors.New("error generate token or set cookie")
	}

	result := types.ResponseUSer{
		Id:         findOne.Id,
		FullName:   findOne.FullName,
		UserName:   findOne.UserName,
		ProfilePic: findOne.ProfilePic,
	}

	return result, nil
}
