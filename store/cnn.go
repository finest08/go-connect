package store

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store struct {
	locaCollMsg  *mongo.Collection
	locaCollUser *mongo.Collection
}


func Connect() Store {
	clientOptions := options.Client().ApplyURI()
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database("data")

	return Store{
		locaCollUser: db.Collection("users"),
		locaCollMsg:  db.Collection("messages"),
	}
}
