package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Podcast represents the schema for the "Podcasts" collection
type Podcast struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Title  string             `bson:"title,omitempty"`
	Author string             `bson:"author,omitempty"`
	Tags   []string           `bson:"tags,omitempty"`
}

// Episode represents the schema for the "Episodes" collection
type Episode struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Podcast     primitive.ObjectID `bson:"podcast,omitempty"`
	Title       string             `bson:"title,omitempty"`
	Description string             `bson:"description,omitempty"`
	Duration    int32              `bson:"duration,omitempty"`
}

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("ATLAS_URI")))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	database := client.Database("quickstart")
	podcastsCollection := database.Collection("podcasts")
	episodesCollection := database.Collection("episodes")

	var episodes []Episode
	cursor, err := episodesCollection.Find(ctx, bson.M{"duration": bson.D{{"$gt", 25}}})
	if err != nil {
		panic(err)
	}
	if err = cursor.All(ctx, &episodes); err != nil {
		panic(err)
	}
	fmt.Println(episodes)

	podcast := Podcast{
		Title:  "The Polyglot Developer",
		Author: "Nic Raboy",
		Tags:   []string{"development", "programming", "coding"},
	}
	insertResult, err := podcastsCollection.InsertOne(ctx, podcast)
	if err != nil {
		panic(err)
	}
	fmt.Println(insertResult.InsertedID)
}
