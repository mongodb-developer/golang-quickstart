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

// PodcastEpisode represents an aggregation result-set for two collections
type PodcastEpisode struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Podcast     Podcast            `bson:"podcast,omitempty"`
	Title       string             `bson:"title,omitempty"`
	Description string             `bson:"description,omitempty"`
	Duration    int32              `bson:"duration,omitempty"`
	PublishDate int32              `bson:"publish_date,omitempty"`
}

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("ATLAS_URI")))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	database := client.Database("quickstart")
	episodesCollection := database.Collection("episodes")

	id, _ := primitive.ObjectIDFromHex("5e3b37e51c9d4400004117e6")

	showInfoCursor, err := episodesCollection.Aggregate(ctx, []bson.D{
		bson.D{{"$match", bson.D{{"podcast", id}}}},
		bson.D{{"$group", bson.D{{"_id", "$podcast"}, {"total", bson.D{{"$sum", "$duration"}}}}}},
	})
	if err != nil {
		panic(err)
	}
	var showsWithInfo []interface{}
	if err = showInfoCursor.All(ctx, &showsWithInfo); err != nil {
		panic(err)
	}
	fmt.Println(showsWithInfo)

	showLoadedCursor, err := episodesCollection.Aggregate(ctx, []bson.D{
		bson.D{{"$lookup", bson.D{{"from", "podcasts"}, {"localField", "podcast"}, {"foreignField", "_id"}, {"as", "podcast"}}}},
		bson.D{{"$unwind", bson.D{{"path", "$podcast"}, {"preserveNullAndEmptyArrays", false}}}},
	})
	if err != nil {
		panic(err)
	}
	var showsLoaded []interface{}
	if err = showLoadedCursor.All(ctx, &showsLoaded); err != nil {
		panic(err)
	}
	fmt.Println(showsLoaded)

	showLoadedStructCursor, err := episodesCollection.Aggregate(ctx, []bson.D{
		bson.D{{"$lookup", bson.D{{"from", "podcasts"}, {"localField", "podcast"}, {"foreignField", "_id"}, {"as", "podcast"}}}},
		bson.D{{"$unwind", bson.D{{"path", "$podcast"}, {"preserveNullAndEmptyArrays", false}}}},
	})
	if err != nil {
		panic(err)
	}
	var showsLoadedStruct []PodcastEpisode
	if err = showLoadedStructCursor.All(ctx, &showsLoadedStruct); err != nil {
		panic(err)
	}
	fmt.Println(showsLoadedStruct)
}
