package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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
	episodesCollection := client.Database("quickstart").Collection("episodes")
	var podcast bson.D
	err = podcastsCollection.FindOne(ctx, bson.M{}).Decode(&podcast)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(podcast)
	cursor, err := episodesCollection.Find(ctx, bson.M{"duration": 25})
	if err != nil {
		log.Fatal(err)
	}
	for cursor.Next(ctx) {
		var episode bson.D
		cursor.Decode(&episode)
		fmt.Println(episode)
	}
	cursor.Close(ctx)
	opts := options.Find()
	opts.SetSort(bson.D{{"duration", -1}})
	cursor, err = episodesCollection.Find(ctx, bson.D{{"duration", bson.D{{"$gt", 24}}}}, opts)
	if err != nil {
		log.Fatal(err)
	}
	for cursor.Next(ctx) {
		var episode bson.D
		cursor.Decode(&episode)
		fmt.Println(episode)
	}
	cursor.Close(ctx)
}
