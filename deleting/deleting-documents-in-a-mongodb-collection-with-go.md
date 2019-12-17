# Quick Start: Deleting Documents in a MongoDB Collection with Go

In a [previous tutorial](https://), I wrote about updating documents within a collection using the Go programming language and the `UpdateOne`, `UpdateMany`, and `ReplaceOne` functions that exist in the MongoDB Go driver. This was part of our getting started exploration of CRUD with Go and MongoDB.

In this tutorial, we're going to explore the final part of CRUD, which is the deleting of documents or even collections.

## The Requirements

There are a few requirements that should be met prior to starting this tutorial if you want maximum success:

- Go 1.13+
- MongoDB Go Driver 1.1.2+
- MongoDB Atlas with an M0 cluster or better

In addition to having these requirements met, each must be properly configured.

> You can get started with a forever free M0 cluster on [MongoDB Atlas](https://www.mongodb.com/cloud). In addition, you can apply the promotional code NRABOY200 to have premium credit added to your account.

If you need help connecting your Go application to MongoDB Atlas, I encourage you to check out my tutorial [Quick Start: How to Get Connected to Your MongoDB Cluster with Go](https://www.mongodb.com/blog/post/quick-start-golang--mongodb--starting-and-setup) as it won't be explored in this particular tutorial.

## Revisiting the Data Model for the Tutorial Series

Before we jump right into the removal of documents, it probably makes sense to revisit the data model we're going to be using to avoid confusion. If you're been keeping up with the previous tutorials in the series, you'll remember we are working with a `podcasts` collection and an `episodes` collection.

The `podcasts` collection has documents that look like this:

```json
{
    "_id": ObjectId("5d9e0173c1305d2a54eb431a"),
    "title": "The Polyglot Developer Podcast",
    "author": "Nic Raboy"
}
```

The `episodes` collection has documents that look like the following:

```json
{
    "_id": ObjectId("5d9f4701e9948c0f65c9165d"),
    "podcast": ObjectId("5d9e0173c1305d2a54eb431a"),
    "title": "Episode #1",
    "description": "This is a description for the first episode.",
    "duration": 25
}
```

When we start deleting documents or dropping collections, we'll be referencing fields from the above schemas and collection names.

## Deleting a Single Document from a MongoDB Collection

Let's say that we want to delete a single document from one of our collections. We can make use of the `DeleteOne` function and provide a filter for the document that should be deleted. Take the following for example:

```go
database := client.Database("quickstart")
podcastsCollection := database.Collection("podcasts")
result, err := podcastsCollection.DeleteOne(ctx, bson.M{"title": "The Polyglot Developer Podcast"})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("DeleteOne removed %v document(s)\n", result.DeletedCount)
```

The above code obtains a handle to the database and collection that we wish to use from the client connection. Using the `DeleteOne` function, we can provide an application context as well as a filter. The filter in this example is around a podcast with a particular title.

If there was an error, we catch it and terminate the application, otherwise we use the result to print how many documents were deleted. Because we are using the `DeleteOne` function, we can only ever have a `DeletedCount` of zero or one.

Establishing a connection to the cluster and defining an application context can be seen in a [previous tutorial](https://www.mongodb.com/blog/post/quick-start-golang--mongodb--starting-and-setup) of this getting started series.

## Deleting Many Documents from a MongoDB Collection

There will often be scenarios where you need to remove more than one document from a collection in a single operation. For these tasks we can make use of the `DeleteMany` function on a collection, which behaves similar to the `InsertMany` and `UpdateMany` operations that we have already seen.

To see this in action, take a look at the following code:

```go
database := client.Database("quickstart")
episodesCollection := database.Collection("episodes")
result, err = episodesCollection.DeleteMany(ctx, bson.M{"duration": 25})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("DeleteMany removed %v document(s)\n", result.DeletedCount)
```

In the above example, pretty much everything is the same. The exception is that we're using a different collection and we're using the `DeleteMany` function. In this example, any document that has a field `duration` that equals `25` will be deleted.

In both the `DeleteOne` and `DeleteMany` examples both range and equality filters can be used. More information on the available operators can be found in the [documentation](https://docs.mongodb.com/manual/reference/operator/query/).

## Dropping a MongoDB Collection and All Documents within the Collection

Removing a single document or many documents isn't the only option. Entire collections can be dropped which would remove all documents in the collection without using a filter. An example of this can be seen below:

```go
if err = podcastsCollection.Drop(ctx); err != nil {
    log.Fatal(err)
}
```

Dropping an entire collection will return an error if something has failed. If the collection doesn't exist, the driver will mask the server error and in this case return a `nil` error.

## The Final Code

To see this tutorial come together, you can take a look at a working example of the code below:

```go
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

	database := client.Database("quickstart")
	podcastsCollection := database.Collection("podcasts")
	episodesCollection := database.Collection("episodes")

	result, err := podcastsCollection.DeleteOne(ctx, bson.M{"title": "The Polyglot Developer Podcast"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("DeleteOne removed %v document(s)\n", result.DeletedCount)

	result, err = episodesCollection.DeleteMany(ctx, bson.M{"duration": 25})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("DeleteMany removed %v document(s)\n", result.DeletedCount)

	if err = podcastsCollection.Drop(ctx); err != nil {
		log.Fatal(err)
	}

	if err = episodesCollection.Drop(ctx); err != nil {
		log.Fatal(err)
	}
}
```

If you try to run the code, don't forget to replace the MongoDB cluster information with that of your MongoDB Atlas cluster.

## Conclusion

You just saw how to delete documents as well as drop collections from MongoDB using the Go programming language. In previous tutorials of the series we explored create, retrieve, and update, all of which are CRUD operations, where the final part is around delete.

In future tutorials in the getting started with MongoDB and Go series, we'll explore binding schemas to native Go data structures, aggregation queries, change streams, and transactions.