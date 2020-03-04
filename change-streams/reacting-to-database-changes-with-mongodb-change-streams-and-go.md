# Reacting to Database Changes with MongoDB Change Streams and Go

If you've been keeping up with my getting started with Go and MongoDB tutorial series, you'll remember that we've accomplished quite a bit so far. We've had a look at everything from CRUD interaction with the database to data modeling, and more. To play catch up with everything we've done, you can have a look at the following tutorials in the series:

- [How to Get Connected to Your MongoDB Cluster with Go](https://www.mongodb.com/blog/post/quick-start-golang--mongodb--starting-and-setup)
- [Creating MongoDB Documents with Go](https://www.mongodb.com/blog/post/quick-start-golang--mongodb--how-to-create-documents)
- [Retrieving and Querying MongoDB Documents with Go](https://www.mongodb.com/blog/post/quick-start-golang--mongodb--how-to-read-documents)
- [Updating MongoDB Documents with Go](https://www.mongodb.com/blog/post/quick-start-golang--mongodb--how-to-update-documents)
- [Deleting MongoDB Documents with Go](https://www.mongodb.com/blog/post/quick-start-golang--mongodb--how-to-delete-documents)
- [Modeling MongoDB Documents with Native Go Data Structures](https://www.mongodb.com/blog/post/quick-start-golang--mongodb--modeling-documents-with-go-data-structures)
- [Performing Complex MongoDB Data Aggregation Queries with Go](https://www.mongodb.com/blog/post/quick-start-golang--mongodb--data-aggregation-pipeline)

In this tutorial we're going to explore change streams in MongoDB and how they might be useful, all with the Go programming language (Golang).

Before we take a look at the code, let's take a step back and understand what change streams are and why there's often a need for them.

Imagine this scenario, one of many possible:

You have an application that engages with internet of things (IoT) clients. Let's say that this is a geofencing application and the IoT clients are something that can trigger the geofence as they come in and out of range. Rather than having your application constantly run queries to see if the clients are in range, wouldn't it make more sense to watch in real-time and react when it happens?

With MongoDB change streams, you can create a pipeline to watch for changes on a collection level, database level, or deployment level, and write logic within your application to do something as data comes in based on your pipeline.

## Creating a Real-Time MongoDB Change Stream with Golang

While there are many possible use-cases for change streams, we're going to continue with the example that we've been using throughout the scope of this getting started series. We're going to continue working with podcast show and podcast episode data.

Let's assume we have the following code to start:

```go
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

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("ATLAS_URI")))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(context.TODO())

	database := client.Database("quickstart")
	episodesCollection := database.Collection("episodes")
}
```

The above code is a very basic connection to a MongoDB cluster, something that we explored in the [How to Get Connected to Your MongoDB Cluster with Go](https://www.mongodb.com/blog/post/quick-start-golang--mongodb--starting-and-setup), tutorial.

To watch for changes, we can do something like the following:

```go
episodesStream, err := episodesCollection.Watch(context.TODO(), mongo.Pipeline{})
if err != nil {
    panic(err)
}
```

The above code will watch for any and all changes to documents within the `episodes` collection. The result is a cursor that we can iterate over indefinitely for data as it comes in.

We can iterate over the curser and make sense of our data using the following code:

```go
episodesStream, err := episodesCollection.Watch(context.TODO(), mongo.Pipeline{})
if err != nil {
    panic(err)
}

defer episodesStream.Close(context.TODO())

for episodesStream.Next(context.TODO()) {
    var data bson.M
    if err := episodesStream.Decode(&data); err != nil {
        panic(err)
    }
    fmt.Printf("%v\n", data)
}
```

If data were to come in, it might look something like the following:

```
map[_id:map[_data:825E4EFCB9000000012B022C0100296E5A1004D960EAE47DBE4DC8AC61034AE145240146645F696400645E3B38511C9D4400004117E80004] clusterTime:{1582234809 1} documentKey:map[_id:ObjectID("5e3b38511c9d
4400004117e8")] fullDocument:map[_id:ObjectID("5e3b38511c9d4400004117e8") description:The second episode duration:30 podcast:ObjectID("5e3b37e51c9d4400004117e6") title:Episode #3] ns:map[coll:episodes 
db:quickstart] operationType:replace]
```

In the above example, I've done a `Replace` on a particular document in the collection. In addition to information about the data, I also receive the full document that includes the change. The results will vary depending on the `operationType` that takes place.

While the code that we used would work fine, it is currently a blocking operation. If we wanted to watch for changes and continue to do other things, we'd want to use a goroutine for iterating over our change stream cursor.

We could make some changes like this:

```go
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

func iterateChangeStream(routineCtx context.Context, waitGroup sync.WaitGroup, stream *mongo.ChangeStream) {
	defer stream.Close(routineCtx)
	defer waitGroup.Done()
	for stream.Next(routineCtx) {
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

	episodesStream, err := episodesCollection.Watch(context.TODO(), mongo.Pipeline{})
	if err != nil {
		panic(err)
	}
	waitGroup.Add(1)
	routineCtx, cancelFn := context.WithCancel(context.Background())
	go iterateChangeStream(routineCtx, waitGroup, episodesStream)

	waitGroup.Wait()
}
```

A few things are happening in the above code. We've moved the stream iteration into a separate function to be used in a goroutine. However, running the application would result in it terminating quite quickly because the `main` function will terminate not too longer after creating the goroutine. To resolve this, we are making use of a `WaitGroup`. In our example, the `main` function will wait until the `WaitGroup` is empty and the `WaitGroup` only becomes empty when the goroutine terminates.

Making use of the `WaitGroup` isn't an absolute requirement as there are other ways to keep the application running while watching for changes. However, given the simplicity of this example, it made sense in order to see any changes in the stream.

To keep the `iterateChangeStream` function from running indefinitely, we are creating and passing a context that can be canceled. While we don't demonstrate canceling the function, at least we know it can be done.

## Complicating the Change Stream with the Aggregation Pipeline

In the previous example, the aggregation pipeline that we used was as basic as you can get. In other words, we were looking for any and all changes that were happening to our particular collection. While this might be good in a lot of scenarios, you'll probably get more out of using a better defined aggregation pipeline.

Take the following for example:

```go
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
```

In the above example, we're still watching for changes to the `episodes` collection. However, this time we're only watching for new documents that have a `duration` field greater than 30. Any other insert or other change stream operation won't be detected.

The results of the above code, when a match is found, might look like the following:

```
map[_id:map[_data:825E4F03CF000000012B022C0100296E5A1004D960EAE47DBE4DC8AC61034AE145240146645F696400645E4F03A01C9D44000063CCBD0004] clusterTime:{1582236623 1} documentKey:map[_id:ObjectID("5e4f03a01c9d
44000063ccbd")] fullDocument:map[_id:ObjectID("5e4f03a01c9d44000063ccbd") description:a quick start into mongodb duration:35 podcast:1234 title:getting started with mongodb] ns:map[coll:episodes db:qui
ckstart] operationType:insert]
```

With change streams, you'll have access to a subset of the MongoDB aggregation pipeline and its operators. You can learn more about what's available in the [official documentation](http://docs.mongodb.com/manual/changeStreams/#modify-change-stream-output).

## Conclusion

You just saw how to use MongoDB change streams in a Golang application using the MongoDB Go driver. As previously pointed out, change streams make it very easy to react to database, collection, and deployment changes without having to constantly query the cluster. This allows you to efficiently plan out aggregation pipelines to respond to as they happen in real-time.

If you're looking to catch up on the other tutorials in the MongoDB with Go quick start series, you can find them below:

- [How to Get Connected to Your MongoDB Cluster with Go](https://www.mongodb.com/blog/post/quick-start-golang--mongodb--starting-and-setup)
- [Creating MongoDB Documents with Go](https://www.mongodb.com/blog/post/quick-start-golang--mongodb--how-to-create-documents)
- [Retrieving and Querying MongoDB Documents with Go](https://www.mongodb.com/blog/post/quick-start-golang--mongodb--how-to-read-documents)
- [Updating MongoDB Documents with Go](https://www.mongodb.com/blog/post/quick-start-golang--mongodb--how-to-update-documents)
- [Deleting MongoDB Documents with Go](https://www.mongodb.com/blog/post/quick-start-golang--mongodb--how-to-delete-documents)
- [Modeling MongoDB Documents with Native Go Data Structures](https://www.mongodb.com/blog/post/quick-start-golang--mongodb--modeling-documents-with-go-data-structures)
- [Performing Complex MongoDB Data Aggregation Queries with Go](https://www.mongodb.com/blog/post/quick-start-golang--mongodb--data-aggregation-pipeline)

To bring the series to a close, the next tutorial will focus on transactions with the MongoDB Go driver.