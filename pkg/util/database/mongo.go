package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var mongoObj *mongo.Client

func NewMongoObj() {

	host := os.Getenv("MONGO_HOST")
	user := os.Getenv("MONGO_USER")
	// dbname := os.Getenv("MONGO_DB")

	// Replace the uri string with your MongoDB deployment's connection string.
	// uri := "mongodb+srv://" + user + ":" + user + "@" + host + "/" + dbname + "?w=majority"
	uri := "mongodb://" + user + ":" + user + "@" + host + "/admin?w=majority&ssl=false"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	// defer func() {
	// 	if err = client.Disconnect(ctx); err != nil {
	// 		log.Println(err)
	// 	}
	// }()

	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Println(err)
	}

	//create connlection

	mongoObj = client
	fmt.Println("Successfully connected and pinged.")
}

func GetMongoObj() *mongo.Client {
	return mongoObj
}

func GetMongoCtx() context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return ctx
}
