package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DbClient struct {
	client *mongo.Client
	users  *mongo.Collection
}

func Dbconnect() DbClient {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mongoPort := os.Getenv("MONGODB_PORT")

	mongoUri := fmt.Sprintf("mongodb://localhost:%s", mongoPort)

	clientOptions := options.Client().ApplyURI(mongoUri)

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal("⛒ Connection Failed to Database")
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal("⛒ Connection Failed to Database")
		log.Fatal(err)
	}

	fmt.Println("⛁ Connected to Database")

	users := client.Database("onboarding").Collection("users")

	return DbClient{client: client, users: users}
}

var result struct {
	Value float64
}

func (db *DbClient) CheckIfUserExist(userId int64) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)

	defer cancel()

	err := db.users.FindOne(ctx, bson.D{{"_id", userId}}).Decode(&result)

	if err == mongo.ErrNoDocuments {
		return false
	} else if err != nil {
		log.Fatal(err)

		return false
	}

	return true
}

func (db *DbClient) AddNewUser(userId, chatId int64, username string, firstName string, lastName string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	_, err := db.users.InsertOne(ctx, bson.D{{"_id", userId}, {"chatId", chatId}, {"username", username}, {"firstName", firstName}, {"lastName", lastName}})

	if err != nil {
		log.Fatal(err)
	}
}
