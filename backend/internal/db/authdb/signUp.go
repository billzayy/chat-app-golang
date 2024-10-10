package authdb

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/billzayy/chat-golang/internal/db"
	"github.com/billzayy/chat-golang/internal/pkg"
	"github.com/billzayy/chat-golang/internal/pkg/middleware"
	"github.com/billzayy/chat-golang/internal/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SignUp(input types.User, w http.ResponseWriter) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	input.CreatedAt = time.Now()
	input.UpdatedAt = time.Now()

	boyProfilePic := fmt.Sprintf(`https://avatar.iran.liara.run/public/girl?username=%s`, input.FullName)
	girlProfilePic := fmt.Sprintf(`https://avatar.iran.liara.run/public/girl?username=%s`, input.FullName)

	defer cancel()

	client, err := db.ConnectMongo(ctx)

	if err != nil {
		return "", errors.New("error on connecting database ")
	}

	// Hash Password
	hashedPassword, err := pkg.HashPassword(input.Password)

	if err != nil {
		return "", errors.New("error on hash password")
	}

	input.Password = hashedPassword

	if input.Gender == "male" {
		input.ProfilePic = boyProfilePic
	} else {
		input.ProfilePic = girlProfilePic
	}

	userCollection := client.Database("chat-app").Collection("users")

	existedUser := checkUser(ctx, userCollection, input)

	if !existedUser {
		// Create a unique index on the userName field
		_, err = userCollection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
			Keys:    bson.D{{Key: "userName", Value: 1}}, // Index specification
			Options: options.Index().SetUnique(true),     // Set the index to be unique
		})

		if err != nil {
			return "", errors.New("error on create document")
		}

		// Insert Value into MongoDB
		_, err := userCollection.InsertOne(context.TODO(), input)

		if err != nil {
			return "", errors.New("error on create account")
		}

		// Generate Token JWT and Save it into Cookie
		_, err = middleware.GenerateTokenAndSetCookie(input.Id, w)

		if err != nil {
			return "", errors.New("error on generate token and set cookie")
		}
		defer client.Disconnect(ctx)

		return "Sign Up Successful", nil
	} else {
		return "User Existed", nil
	}
}

func checkUser(ctx context.Context, userCollection *mongo.Collection, input types.User) bool {
	var existedUser bool
	var findOne types.User

	err := userCollection.FindOne(ctx, bson.M{"userName": input.UserName}).Decode(&findOne)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			existedUser = false
		} else {
			log.Fatal(err)
		}
	} else {
		existedUser = true
	}

	return existedUser
}
