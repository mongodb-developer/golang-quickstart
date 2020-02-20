package main

import (
	"context"
	"fmt"
	"os"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ChangeStream(waitGroup sync.WaitGroup, stream *mongo.ChangeStream) {
	defer stream.Close(context.TODO())
	defer waitGroup.Done()
	for stream.Next(context.TODO()) {
		var data bson.M
		if err := stream.Decode(&data); err != nil {
			panic(err)
		}
		fmt.Printf("%v\n", data)
	}
}

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("ATLAS_URI")))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(context.TODO())

	database := client.Database("quickstart")
	episodesCollection := database.Collection("episodes")

	var waitGroup sync.WaitGroup

	matchPipeline := bson.D{
		{
			"$match", bson.D{
				{"operationType", "insert"},
				{"fullDocument.duration", bson.D{
					{"$gt", 30},
				}},
			},
		},
	}

	episodesStream, err := episodesCollection.Watch(context.TODO(), mongo.Pipeline{matchPipeline})
	if err != nil {
		panic(err)
	}
	waitGroup.Add(1)
	go ChangeStream(waitGroup, episodesStream)

	waitGroup.Wait()
}
