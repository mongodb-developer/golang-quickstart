# Performing Complex MongoDB Data Aggregation Queries with Go

If you've been following along with my getting started series around MongoDB and Golang, you might remember the tutorial where we took a look at [finding documents](https://www.mongodb.com/blog/post/quick-start-golang--mongodb--how-to-read-documents) in a collection. In this tutorial we saw how to use the `Find` and `FindOne` functions to filter for documents within a collection. This is essentially querying for documents within a specific collection where the filter parameters are fields within the schema of that collection.

So what happens if you need to do something a little more complex like return data that isn't within the schema, do complex manipulations prior to returning a response, or looking across collections in a single command?

This is where the MongoDB aggregation framework becomes valuable.

In this tutorial we're going to look at a few MongoDB aggregation framework examples using the Go programming language, examples that can't really be done with a basic `Find` or `FindOne` operation.

## The Requirements

To be successful with this tutorial, you'll need the following requirements to be met:

- Go 1.10+ installed and configured
- MongoDB Atlas with an M0 cluster or better
- The MongoDB Go driver

It is advisable that you've completed the other tutorials in the getting started with MongoDB and Go series as it shares [information around the schema](https://www.mongodb.com/blog/post/quick-start-golang--mongodb--modeling-documents-with-go-data-structures) that we're using as well as information on [connecting to a cluster](https://www.mongodb.com/blog/post/quick-start-golang--mongodb--starting-and-setup). However, if you feel comfortable with MongoDB and Go, it isn't an absolute requirement.

> Use promotional code [NICRABOY200](https://www.mongodb.com/cloud) to receive $200 in premium credit towards your MongoDB Atlas cluster if you'd like something more powerful than the forever free M0 cluster.

Make sure that the MongoDB Atlas cluster has been properly whitelisted so that your Go application can communicate to it. For information on installing the MongoDB Go driver and connecting to a cluster, check out my [previous tutorial](https://www.mongodb.com/blog/post/quick-start-golang--mongodb--starting-and-setup) on the subject.

## Leveraging the MongoDB Aggregation Framework in Golang

We're going to be working with a simple application for this example. To get us up to speed, within the **$GOPATH**, create a new project with a **main.go** file and within that **main.go** file, include the following:

```go
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
	episodesCollection := database.Collection("episodes")
}
```

You'll recall that we created the Go data structures with BSON annotations in a previous tutorial titled, [Modeling MongoDB Documents with Native Go Data Structures](https://www.mongodb.com/blog/post/quick-start-golang--mongodb--modeling-documents-with-go-data-structures). The logic used for connecting to a cluster and setting a handle to a particular database and collection was last seen in [How to Get Connected to Your MongoDB Cluster with Go](https://www.mongodb.com/blog/post/quick-start-golang--mongodb--starting-and-setup).

If you haven't already installed the Go driver for MongoDB, it can be installed by executing the following:

```bash
dep init
dep ensure -add "go.mongodb.org/mongo-driver/mongo"
```

If you don't have the Go dependency manager (dep) installed and configured, you can learn more about it [here](https://github.com/golang/dep).

The important thing to take note of in our boilerplate code are the native Go data structures that represent the document schema for each of the collections. The goal here is to use the aggregation framework to work with the data in those collections, but add certain groupings, manipulations, ect..

Let's assume that you have the following documents in your `episodes` collection:

```json
// Document #1
{
    "_id": ObjectId("5e3b381c1c9d4400004117e7"),
    "podcast": ObjectId("5e3b37e51c9d4400004117e6"),
    "title": "Episode #1",
    "description": "The first episode",
    "duration": 25
}

// Document #2
{
    "_id": ObjectId("5e3b38511c9d4400004117e8"),
    "podcast": ObjectId("5e3b37e51c9d4400004117e6"),
    "title": "Episode #2",
    "description": "The second episode",
    "duration": 30
}
```

The first aggregation that we're going to look at will take all the episodes for a particular podcast and get the total duration of that podcast. To be clear, I don't mean the duration of a particular episode, but the duration of the podcast as a whole.

For this aggregation query, we're going to focus on the `podcast` field as well as the `duration` field of our documents. Take the following code:

```go
id, _ := primitive.ObjectIDFromHex("5e3b37e51c9d4400004117e6")

matchStage := bson.D{{"$match", bson.D{{"podcast", id}}}}
groupStage := bson.D{{"$group", bson.D{{"_id", "$podcast"}, {"total", bson.D{{"$sum", "$duration"}}}}}}

showInfoCursor, err := episodesCollection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage})
if err != nil {
    panic(err)
}
var showsWithInfo []bson.M
if err = showInfoCursor.All(ctx, &showsWithInfo); err != nil {
    panic(err)
}
fmt.Println(showsWithInfo)
```

Because the particular `podcast` is important to us, we are taking the id of the podcast in question and converting it into an object id that MongoDB and the Go driver can understand. Next we are defining stages of the aggregation pipeline, in this case a matching stage and grouping stage. In the matching stage we are matching all documents that have the `podcast` field in question. In the grouping stage we are grouping the matches by the `podcast` field because it is non-distinct, and then we are summing each of the `duration` fields into a new `total` field. The `Aggregate` operation executes our defined pipeline.

The result would look something like this:

```plaintext
[map[_id:ObjectID("5e3b37e51c9d4400004117e6") total:55]]
```

Had we altered the aggregation to include more `podcast` values, we could have ended up with several different podcast groups and their total minutes.

Let's look at another example. For this scenario, let's say we want to "join" documents from different collections similar to how you would in a relational database. Based on our document schema, we already have a `podcast` field in the `episodes` collection referencing a document in the `podcasts` collection. That document might look like this:

```json
{
    "_id": ObjectId("5e3b37e51c9d4400004117e6"),
    "name": "The Polyglot Developer Podcast",
    "author": "Nic Raboy",
    "tags": ["development", "programming", "coding"]
}
```

So what would our aggregation query look like if we wanted to include the podcast information with the episode information? We might try to do something like this:

```go
lookupStage := bson.D{{"$lookup", bson.D{{"from", "podcasts"}, {"localField", "podcast"}, {"foreignField", "_id"}, {"as", "podcast"}}}}
unwindStage := bson.D{{"$unwind", bson.D{{"path", "$podcast"}, {"preserveNullAndEmptyArrays", false}}}}

showLoadedCursor, err := episodesCollection.Aggregate(ctx, mongo.Pipeline{lookupStage, unwindStage})
if err != nil {
    panic(err)
}
var showsLoaded []bson.M
if err = showLoadedCursor.All(ctx, &showsLoaded); err != nil {
    panic(err)
}
fmt.Println(showsLoaded)
```

In the above example, we are using the [$lookup](https://docs.mongodb.com/manual/reference/operator/aggregation/lookup/) operator to join from the `podcasts` collection using the `podcast` field in our `episodes` collection and the `_id` field in our foreign `podcasts` collection. The output from the `$lookup` operation will be an array stored as `podcast`.

After the `$lookup` operation, we make use of the [$unwind](https://docs.mongodb.com/manual/reference/operator/aggregation/unwind/) operator to flatten the array that we had previously created. Think of flattening or deconstructing an array as taking an array and now outputting each element of that array as a document in the result set.

If we were to run our aggregation, we might end up with results that look like this:

```plaintext
[map[_id:ObjectID("5e3b381c1c9d4400004117e7") description:The first episode duration:25 podcast:map[_id:ObjectID("5e3b37e51c9d4400004117e6") author:Nic Raboy name:The Polygl
ot Developer Podcast tags:[development coding programming]] title:Episode #1] map[_id:ObjectID("5e3b38511c9d4400004117e8") description:The second episode duration:30 podcast
:map[_id:ObjectID("5e3b37e51c9d4400004117e6") author:Nic Raboy name:The Polyglot Developer Podcast tags:[development coding programming]] title:Episode #2]]
```

Notice that each podcast episode in the above results is now printed with the show information. This saves you from having to execute multiple `Find` operations within your Go code.

The above query that we saw is great, but I think we can do better.

In the first example that used `$lookup` and `$unwind` we were using an `[]bson.M` to work with the results. Not the end of the world, but if we wanted to access particular fields, use an autocomplete, etc., things might get a little messy. Instead, we can create a native Go data structure to represent the results of our aggregation.

```go
// PodcastEpisode represents an aggregation result-set for two collections
type PodcastEpisode struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Podcast     Podcast            `bson:"podcast,omitempty"`
	Title       string             `bson:"title,omitempty"`
	Description string             `bson:"description,omitempty"`
	Duration    int32              `bson:"duration,omitempty"`
}
```

For the most part the above data structure will look familiar. However, pay attention to the `Podcast` field. In this example it is no longer a `primitive.ObjectID`, but instead a `Podcast` type, which is a previously defined data structure that we had created.

With this new data structure available, we can change our aggregation a bit:

```go
lookupStage := bson.D{{"$lookup", bson.D{{"from", "podcasts"}, {"localField", "podcast"}, {"foreignField", "_id"}, {"as", "podcast"}}}}
unwindStage := bson.D{{"$unwind", bson.D{{"path", "$podcast"}, {"preserveNullAndEmptyArrays", false}}}}

showLoadedStructCursor, err := episodesCollection.Aggregate(ctx, mongo.Pipeline{lookupStage, unwindStage})
if err != nil {
    panic(err)
}
var showsLoadedStruct []PodcastEpisode
if err = showLoadedStructCursor.All(ctx, &showsLoadedStruct); err != nil {
    panic(err)
}
fmt.Println(showsLoadedStruct)
```

Notice that we're now using a `[]PodcastEpisode` to store the results rather than a `[]bson.M`. While we don't demonstrate it in this example, we would have access to each field within that data structure if we wanted to.

## Conclusion

You just saw a few aggregation examples within MongoDB using the Go programming language (Golang). There are quite a few operators within the aggregation framework that MongoDB offers and you can learn more about them in the [official documentation](https://docs.mongodb.com/manual/reference/operator/aggregation-pipeline/). While the examples that I demonstrated were short and with few operators, you could end up in more advanced territory depending on your needs.

To catch up on other tutorials in the getting started with Golang series, check out these:

- [How to Get Connected to Your MongoDB Cluster with Go](https://www.mongodb.com/blog/post/quick-start-golang--mongodb--starting-and-setup)
- [Creating MongoDB Documents with Go](https://www.mongodb.com/blog/post/quick-start-golang--mongodb--how-to-create-documents)
- [Retrieving and Querying MongoDB Documents with Go](https://www.mongodb.com/blog/post/quick-start-golang--mongodb--how-to-read-documents)
- [Updating MongoDB Documents with Go](https://www.mongodb.com/blog/post/quick-start-golang--mongodb--how-to-update-documents)
- [Deleting MongoDB Documents with Go](https://www.mongodb.com/blog/post/quick-start-golang--mongodb--how-to-delete-documents)
- [Modeling MongoDB Documents with Native Go Data Structures](https://www.mongodb.com/blog/post/quick-start-golang--mongodb--modeling-documents-with-go-data-structures)

To bring the series to a close, we'll be looking at change streams and transactions using the Go programming language and MongoDB.