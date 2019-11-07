# Quick Start: Creating MongoDB Documents with Go

In a [previous tutorial](https://), I demonstrated how to connect to a MongoDB cluster with the Go programming language (Golang) and verify the connection through pinging the cluster and outputting the available databases. In this tutorial, we're going to continue with the getting started material, but this time connecting to a particular collection to create new documents. This will continue to use the Go programming language.

## Tools and Versions for the Tutorial Series

I wanted to take a moment to reiterate the tools and versions that I'm using within this tutorial series:

- MongoDB Atlas with an M0 free cluster
- Visual Studio Code (VSCode)
- MongoDB Go Driver 1.1.2
- Go 1.13

If you're using different software or more recent versions, don't worry, the code should be fine, the steps just might be a little different.

> You can get started with an M0 cluster on [MongoDB Atlas](https://www.mongodb.com/cloud) for free. If sign up using the promotional code NRABOY200, you'll receive premium credit applied to your account.

The assumption is that you've configured everything prior to using this tutorial in the series.

## Understanding the Data

As a refresher, MongoDB stores data in JSON documents, which are actually Binary JSON (BSON) objects stored on disk. We won't get into the nitty gritty of how MongoDB works with JSON and BSON in this particular series, but we will familiarize ourselves with some of the data we'll be working with going forward.

Take the following MongoDB documents for example:

```json
{
    "_id": ObjectId("5d9e0173c1305d2a54eb431a"),
    "title": "The Polyglot Developer Podcast",
    "author": "Nic Raboy"
}
```

The above document might represent a podcast show that has any number of episodes. Any document that represents a show might appear in a `podcasts` collection. There will also be a document that looks like the following:

```json
{
    "_id": ObjectId("5d9f4701e9948c0f65c9165d"),
    "podcast": ObjectId("5d9e0173c1305d2a54eb431a"),
    "title": "GraphQL for API Development",
    "description": "Learn about GraphQL development in this episode of the podcast.",
    "duration": 25
}
```

The above document might represent a podcast episode. Any document that represents an episode might appear in an `episodes` collection. Neither of these two documents are particularly complex, but we'll see different variations of them as we progress through the series.

## Getting a Handle to a Specific Collection

Before data can be created or queried, a handle to a collection must be defined. It doesn't matter if the database or collection already exists on the cluster as it will be created automatically when the first document is inserted if it does not.

Since we will be using two different collections, the following can be done in Go to establish the collection handles:

```golang
quickstartDatabase := client.Database("quickstart")
podcastsCollection := quickstartDatabase.Collection("podcasts")
episodesCollection := quickstartDatabase.Collection("episodes")
```

The above code uses a `client` that is already connected to our cluster, and establishes a handle for our desired database. In this case, the database is `quickstart`. Again, if it doesn't already exist, it is fine. We are also establishing handles to two different collections, both of which don't need to exist. The `client` variable was configured in the [previous tutorial](https://) in the series.

Looking at our code thus far, we might have something that looks like the following:

```golang
package main

import (
	"context"
	"log"
	"time"
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

    quickstartDatabase := client.Database("quickstart")
    podcastsCollection := quickstartDatabase.Collection("podcasts")
    episodesCollection := quickstartDatabase.Collection("episodes")
}
```

Between this tutorial and the [first tutorial](https://) in the series, the cluster ping logic and listing of database logic was removed. Instead, we're just connecting to the cluster and creating handles to our collections in a particular database.

## Creating One or Many BSON Documents in a Single Request

Now that we have a `podcastsCollection` and an `episodesCollection` variable, we can proceed to create data and insert it into either of the collections.

For this example, I won't be using a pre-defined schema. In a future tutorial, we'll see how to map documents to native Go data structures, but for now, we're going to look at other options.

Take the following command for example:

```golang
podcastResult, err := podcastsCollection.InsertOne(ctx, bson.D{
    {Key: "title", Value: "The Polyglot Developer Podcast"},
    {Key: "author", Value: "Nic Raboy"},
})
```

The above command will insert a single document into the `podcasts` collection. While the above example is rather flat, it could be adjusted to be more complex. Take the following example:

```golang
podcastResult, err := podcastsCollection.InsertOne(ctx, bson.D{
    {Key: "title", Value: "The Polyglot Developer Podcast"},
    {Key: "author", Value: "Nic Raboy"},
    {Key: "tags", Value: bson.A{"development", "programming", "coding"}},
})
```

The above example adds a `tags` field to the document which is an array. So far we've seen `bson.D` which is a document and `bson.A` which is an array. There are other options which can be found in the documentation for the MongoDB Go driver.

If we wanted to, usage of the `bson.D` and `bson.A` data structures could be drastically simplified. Take the following simplification:

```golang
podcastResult, err := podcastsCollection.InsertOne(ctx, bson.D{
    {"title", "The Polyglot Developer Podcast"},
    {"author", "Nic Raboy"},
    {"tags", bson.A{"development", "programming", "coding"}},
})
```

Notice that in the above example, the `Key` and `Value` properties were removed. It is up to you to decide how you want to use each of the data structures that the MongoDB Go driver offers.

The `InsertOne` function returns both an `InsertOneResult` and an error. If there was no error, the `InsertOneResult`, as shown through a `podcastResult` variable in this example, has an `InsertedID` field. This is helpful if you need to reference the newly created document in future operations. We'll see more on this shortly.

In the previous few examples, the `InsertOne` function was used, which only creates a single document. If you wanted to create multiple documents, you could make use of the `InsertMany` function like follows:

```golang
episodeResult, err := episodesCollection.InsertMany(ctx, []interface{}{
    bson.D{
        {"podcast", podcastResult.InsertedID},
        {"title", "GraphQL for API Development"},
        {"description", "Learn about GraphQL from the co-creator of GraphQL, Lee Byron."},
        {"duration", 25},
    },
    bson.D{
        {"podcast", podcastResult.InsertedID},
        {"title", "Progressive Web Application Development"},
        {"description", "Learn about PWA development with Tara Manicsic."},
        {"duration", 32},
    },
})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Inserted %v documents into episode collection!\n", len(episodeResult.InsertedIDs))
```

In the above example you'll notice that we are using a slice of `interface{}` which represents each of the documents that we wish to insert. For each of the documents, the same `bson.D` rules are applied, as seen previously. Also notice that the `InsertedID` from the previous insert operation was used to reference the parent podcast for each episode inserted.

Rather than returning an `InsertOneResult`, the `InsertMany` function returns an `InsertManyResult`. However, this behaves in a similar fashion, with the exception that now we have access to `InsertedIDs` which is an `[]interface{}`. This slice will contain the ids to each of the inserted episodes for this particular example.

## Conclusion

You just saw how to insert documents into a MongoDB cluster using the Go programming language. In this particular part of the series, primitive data structures were used to save us from having to pre-define a schema. In a future tutorial, we're going to see how to make our projects more manageable by defining how MongoDB documents can be mapped to Go data structures.

With data in the database, the next part of the getting started series will include how to query for the data using Go.