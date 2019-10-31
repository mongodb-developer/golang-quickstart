package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx context.Context

func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI("<ATLAS_URI_HERE>"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	podcastsCollection := client.Database("quickstart").Collection("podcasts")
	episodesCollection := client.Database("quickstart").Collection("episodes")
	podcastsCollection.InsertOne(ctx, bson.D{
		{"title", "The Polyglot Developer Podcast"},
		{"author", "Nic Raboy"},
		{"tags", bson.A{"development", "programming", "coding"}},
	})
	id, _ := primitive.ObjectIDFromHex("5dbb2038e7841e4744f57a9c")
	episodesCollection.InsertMany(ctx, []interface{}{
		bson.D{
			{"podcast", id},
			{"title", "GraphQL for API Development"},
			{"description", "Learn about GraphQL from the co-creator of GraphQL, Lee Byron."},
			{"duration", 25},
		},
		bson.D{
			{"podcast", id},
			{"title", "Progressive Web Application Development"},
			{"description", "Learn about PWA development with Tara Manicsic."},
			{"duration", 32},
		},
	})
}
