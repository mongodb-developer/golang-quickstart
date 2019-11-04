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
	id, _ := primitive.ObjectIDFromHex("5d9e0173c1305d2a54eb431a")
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
	fmt.Printf("Updated %v Documents!\n", result.MatchedCount)
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
	fmt.Printf("Updated %v Documents!\n", result.MatchedCount)
}
