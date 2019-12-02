package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI("<ATLAS_URI_HERE>"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	podcastsCollection := client.Database("quickstart").Collection("podcasts")

	// Update a single document based on a document id hash
	id, _ := primitive.ObjectIDFromHex("5dd890a61c9d4400003f3a31")
	result, err := podcastsCollection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.D{
			{"$set", bson.D{{"author", "Nic Raboy"}}},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Updated %v Documents!\n", result.MatchedCount)

	// Update zero or more documents based on a filter criteria
	result, err = podcastsCollection.UpdateMany(
		ctx,
		bson.M{"title": "The Polyglot Developer Podcast"},
		bson.D{
			{"$set", bson.D{{"author", "Nicolas Raboy"}}},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Updated %v Documents!\n", result.ModifiedCount)

	// Update zero or more documents and add a field that may not exist to the document
	result, err = podcastsCollection.UpdateMany(
		ctx,
		bson.M{"title": "The Polyglot Developer Podcast"},
		bson.D{
			{"$set", bson.D{{"author", "Nic Raboy"}, {"website", "thepolyglotdeveloper.com"}}},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Updated %v Documents!\n", result.ModifiedCount)

	// Repace an entire single document based on a filter criteria
	result, err = podcastsCollection.ReplaceOne(
		ctx,
		bson.M{"author": "Nic Raboy"},
		bson.M{
			"title":  "The Nic Raboy Show",
			"author": "Nicolas Raboy",
		},
	)
	fmt.Printf("Replaced %v Documents!\n", result.ModifiedCount)
}
