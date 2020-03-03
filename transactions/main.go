package main

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Podcast struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Title  string             `bson:"title,omitempty"`
	Author string             `bson:"author,omitempty"`
	Tags   []string           `bson:"tags,omitempty"`
}

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("ATLAS_URI")))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(context.TODO())

	database := client.Database("quickstart")
	podcastsCollection := database.Collection("podcasts")

	session, err := client.StartSession()
	if err != nil {
		panic(err)
	}
	if err = session.StartTransaction(); err != nil {
		panic(err)
	}
	err = mongo.WithSession(context.Background(), session, func(sessionContext mongo.SessionContext) error {
		result, err := podcastsCollection.InsertOne(
			sessionContext,
			Podcast{
				Title:  "Transactions for All",
				Author: "Nic Raboy",
			},
		)
		if err != nil {
			panic(err)
		}
		//panic(result.InsertedID)
		if err = session.CommitTransaction(sessionContext); err != nil {
			panic(err)
		}
		fmt.Println(result.InsertedID)
		return nil
	})
	if err != nil {
		panic(err)
	}
	session.EndSession(context.Background())
}
