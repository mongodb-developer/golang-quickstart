# Creating, Reading, Updating, and Deleting Documents in MongoDB with Golang

Interacting with a database typically involves creating and inserting data into the database, reading data from the database, updating data that already exists in the database, and removing data from the database. These operations are known as CRUD and they are essential to pretty much every application.

In this tutorial we're going to be exploring how to establish a connection to a MongoDB cluster, as well as create, read, update, and delete documents, all using the Go programming language (Golang).

## The Requirements

There are a few requirements that must be met to be successful with this tutorial:

- A MongoDB Atlas cluster
- Go 1.13+

It is important that the MongoDB Atlas cluster is properly configured to allow connections from your locally running Go application or wherever you plan to host your application. By default, all external connections are denied.

<div class="callout">

Get started with an M0 cluster on [MongoDB Atlas](https://www.mongodb.com/cloud) today. It's free forever and you'll be able to work alongside this blog series. Use promo code NICRABOY200 when you sign up and you'll get an extra $200.00 credit applied to your account.

</div>

It's also important that the Go driver has been properly configured for development on your computer, and this includes having [dep](https://github.com/golang/dep) installed to be used as the dependency management tool.

## Connecting to a MongoDB Cluster with Golang

Within your **$GOPATH**, create a new project directory titled **quickstart** and add a **main.go** file to that project.

For this particular tutorial, all code will be added to the **main.go** file. We can start the **main.go** file with the following boilerplate code, necessary for our dependency manager:

```go
package main

func main() { }
```

The next step is to install the MongoDB Go Driver through the Go Package Manager. To do this, execute the following from the command line:

```bash
$ dep init
$ dep ensure -add "go.mongodb.org/mongo-driver/mongo@~1.3.0"
```

Note that for this tutorial we're using `dep` to manage our packages and we're using version 1.3.0 of the MongoDB Go Driver. If you don't have `dep`, you can install it through the [official documentation](https://github.com/golang/dep).

With the driver installed, open the project's **main.go** file and add the following imports to the code:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() { }
```

The above code represents the imported modules within the MongoDB Go Driver that will be used throughout the tutorial series. Most logic, at least for now, will exist in the `main` function.

Inside the `main` function, let's establish a connection to our MongoDB Atlas cluster:

```go
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
}
```

There are a few things that are happening in the above code. First we're configuring our client to use the correct URI, but we're not yet connecting to it. Assuming nothing is malformed and no error was thrown, we can define a timeout duration that we want to use when trying to connect. The ten seconds I used might be a little too generous for your needs, but feel free to play around with the value that makes the most sense to you.

In regards to the Atlas URI, you can use any of the driver URIs from the Atlas dashboard. They'll look something like this:

```
mongodb+srv://<username>:<password>@cluster0-zzart.mongodb.net/test?retryWrites=true&w=majority
```

Just remember, to use the information that Atlas provides for your particular cluster.

After connecting, if there isn't an error, we can defer the closing of the connection for when the `main` function exits. This will keep the connection to the database open until we're done.

So if no errors were thrown, can we be sure that we're really connected? If you're concerned, we can ping the cluster from within our application:

```go
err = client.Ping(ctx, readpref.Primary())
if err != nil {
    log.Fatal(err)
}
```

The above code uses our connected client and the context that we had previously defined. If there is no error, the ping was a success!

We can take things a step further by listing the available databases on our MongoDB Atlas cluster. Within the `main` function, add the following:

```go
databases, err := client.ListDatabaseNames(ctx, bson.M{})
if err != nil {
    log.Fatal(err)
}
fmt.Println(databases)
```

The above code will return a `[]string` with each of the database names. Since we don't plan to filter any of the databases in the list, the `filter` argument on the `ListDatabaseNames` function can be `bson.M{}`.

The result of the above code might be something like this:

```
[quickstart video admin local]
```

Of course, your actual databases will likely be different than mine. It is not a requirement to have specific databases or collections at this stage of the tutorial.

To bring everything together, take a look at our project thus far:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)
}
```

Not bad for 34 lines of code, considering that about half of that was just defining the imports for packages to be used within the project.

## Creating Documents in a MongoDB Collection with Golang

If you've made it this far, you were successful in connecting to a MongoDB cluster with Go. Now we're going to look at various functions in the MongoDB Go driver for creating one or more documents at a time.

### Understanding the Data Model for the Application

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

The above document might represent a podcast episode. Any document that represents an episode might appear in an `episodes` collection. Neither of these two documents are particularly complex, but we'll see different variations of them as we progress through the tutorial.

### Getting a Handle to a Specific Collection

Before data can be created or queried, a handle to a collection must be defined. It doesn't matter if the database or collection already exists on the cluster as it will be created automatically when the first document is inserted if it does not.

Since we will be using two different collections, the following can be done in Go to establish the collection handles:

```go
quickstartDatabase := client.Database("quickstart")
podcastsCollection := quickstartDatabase.Collection("podcasts")
episodesCollection := quickstartDatabase.Collection("episodes")
```

The above code uses a `client` that is already connected to our cluster, and establishes a handle for our desired database. In this case, the database is `quickstart`. Again, if it doesn't already exist, it is fine. We are also establishing handles to two different collections, both of which don't need to exist. The `client` variable was configured in the earlier in the tutorial when we were establishing a connection to the cluster.

Looking at our code thus far, we might have something that looks like the following:

```go
package main

import (
	"context"
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

    quickstartDatabase := client.Database("quickstart")
    podcastsCollection := quickstartDatabase.Collection("podcasts")
    episodesCollection := quickstartDatabase.Collection("episodes")
}
```

The cluster ping logic and the listing of database logic was removed as it doesn't serve too much of a purpose for this particular tutorial going forward. Instead, we're just connecting to the cluster and creating handles to our collections in a particular database.

### Creating One or Many BSON Documents in a Single Request

Now that we have a `podcastsCollection` and an `episodesCollection` variable, we can proceed to create data and insert it into either of the collections.

For this example, I won't be using a pre-defined schema. In a [future tutorial](https://www.mongodb.com/blog/post/quick-start-golang--mongodb--modeling-documents-with-go-data-structures), we'll see how to map documents to native Go data structures, but for now, we're going to look at other options.

Take the following command for example:

```go
podcastResult, err := podcastsCollection.InsertOne(ctx, bson.D{
    {Key: "title", Value: "The Polyglot Developer Podcast"},
    {Key: "author", Value: "Nic Raboy"},
})
```

The above command will insert a single document into the `podcasts` collection. While the above example is rather flat, it could be adjusted to be more complex. Take the following example:

```go
podcastResult, err := podcastsCollection.InsertOne(ctx, bson.D{
    {Key: "title", Value: "The Polyglot Developer Podcast"},
    {Key: "author", Value: "Nic Raboy"},
    {Key: "tags", Value: bson.A{"development", "programming", "coding"}},
})
```

The above example adds a `tags` field to the document which is an array. So far we've seen `bson.D` which is a document and `bson.A` which is an array. There are other options which can be found in the documentation for the MongoDB Go driver.

If we wanted to, usage of the `bson.D` and `bson.A` data structures could be drastically simplified. Take the following simplification:

```go
podcastResult, err := podcastsCollection.InsertOne(ctx, bson.D{
    {"title", "The Polyglot Developer Podcast"},
    {"author", "Nic Raboy"},
    {"tags", bson.A{"development", "programming", "coding"}},
})
```

Notice that in the above example, the `Key` and `Value` properties were removed. It is up to you to decide how you want to use each of the data structures that the MongoDB Go driver offers.

The `InsertOne` function returns both an `InsertOneResult` and an error. If there was no error, the `InsertOneResult`, as shown through a `podcastResult` variable in this example, has an `InsertedID` field. This is helpful if you need to reference the newly created document in future operations. We'll see more on this shortly.

In the previous few examples, the `InsertOne` function was used, which only creates a single document. If you wanted to create multiple documents, you could make use of the `InsertMany` function like follows:

```go
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

## Retrieving and Querying MongoDB Documents with Golang

We just explored the `Insert` and `InsertMany` functions while making use of `bson.D`, `bson.M`, and similar MongoDB data types. By now you should have at least one document within your collections because we're going to explore reading documents from MongoDB and creating queries to retrieve documents based on certain criteria.

### A Reminder of the Data in the MongoDB Collections

When thinking back to the data that we created previously, we know that we have a `podcasts` collection with data that looks something like this:

```json
{
    "_id": ObjectId("5d9e0173c1305d2a54eb431a"),
    "title": "The Polyglot Developer Podcast",
    "author": "Nic Raboy"
}
```

We also have an `episodes` collection which has similarly structured data that looks like this:

```json
{
    "_id": ObjectId("5d9f4701e9948c0f65c9165d"),
    "podcast": ObjectId("5d9e0173c1305d2a54eb431a"),
    "title": "GraphQL for API Development",
    "description": "Learn about GraphQL development in this episode of the podcast.",
    "duration": 25
}
```

Knowing the fields that exist in our data will be important to us when it comes to crafting queries to return only the data that we need, rather than everything.

### Reading All Documents from a Collection

Reading all data from a collection consists of making the request, then working with the results cursor. Knowing what fields exist on each of the documents isn't too important, only knowing the collection name itself.

A simple example of this can be done through the following:

```go
cursor, err := episodesCollection.Find(ctx, bson.M{})
if err != nil {
    log.Fatal(err)
}
var episodes []bson.M
if err = cursor.All(ctx, &episodes); err != nil {
    log.Fatal(err)
}
fmt.Println(episodes)
```

If you think back to the cluster connection part of this tutorial, the `Find` function might look similar to the `ListDatabaseNames` function. We can provide a context and some query parameters, and get our results. In this example `bson.M` represents a map of fields in no particular order. However, because we're trying to return all documents, there aren't any fields in our query.

Assuming no error happens, the results will exist in a MongoDB cursor. In this simple example, all documents are decoded into a `[]bson.M` variable, the cursor is closed, and then the variable is printed.

To get an idea of what the printed data would look like, take the following for example:

```
[map[_id:ObjectID("5dc98dea14c897e9ab808161") description:This is the first episode. duration:25 podcast:ObjectID("5dc98de914c897e9ab808160") title:Episode #1] map[_id:ObjectID
("5dc98dea14c897e9ab808162") description:This is the second episode. duration:32 podcast:ObjectID("5dc98de914c897e9ab808160") title:Episode #2]]
```

If your expected result set is large, using the `*mongo.Cursor.All` function might not be the best idea. Instead, you can iterate over your cursor and have it retrieve your data in batches. To do this, our code might change to the following:

```go
cursor, err := episodesCollection.Find(ctx, bson.M{})
if err != nil {
    log.Fatal(err)
}
defer cursor.Close(ctx)
for cursor.Next(ctx) {
    var episode bson.M
    if err = cursor.Decode(&episode); err != nil {
        log.Fatal(err)
    }
    fmt.Println(episode)
}
```

In both the `*mongo.Cursor.All` and `*mongo.Cursor.Next` examples, the data is loaded into `bson.M` data structures which behave as maps. We'll explore marshalling and unmarshalling the data to custom native Go data structures in a [later tutorial](https://www.mongodb.com/blog/post/quick-start-golang--mongodb--modeling-documents-with-go-data-structures).

### Reading a Single Document from a Collection

Retrieving all documents from a collection doesn't always make sense. Instead, sometimes it makes sense to only return a single document. Instead of using the `Find` function, the `FindOne` function can be used.

Take the following example:

```go
var podcast bson.M
if err = podcastsCollection.FindOne(ctx, bson.M{}).Decode(&podcast); err != nil {
    log.Fatal(err)
}
fmt.Println(podcast)
```

In the above example, a `FindOne` is executed without any particular query filter on the data. Rather than returning a cursor, the single result can be decoded directly into the `bson.M` object. If there is no error, the result will print.

To get an idea of what the result might look like, take the following:

```
map[_id:ObjectID("5dc98c8c9e2e56363b11b375") author:Nic Raboy title:The Polyglot Developer Podcast]
```

Because the results in the `Find` and `FindOne` use `bson.M`, the format of the results will be the same.

### Querying Documents from a Collection with a Filter

In the previous examples of `Find` and `FindOne`, we've seen the `bson.M` filter parameter, even though it wasn't used for filtering. We can make use of the full power of the MongoDB Query Language (MQL) to filter the results of our queries, simply by populating the map.

Let's say that we want to filter our results to only include podcast episodes that are exactly 25 minutes. We could do something like the following:

```go
filterCursor, err := episodesCollection.Find(ctx, bson.M{"duration": 25})
if err != nil {
    log.Fatal(err)
}
var episodesFiltered []bson.M
if err = filterCursor.All(ctx, &episodesFiltered); err != nil {
    log.Fatal(err)
}
fmt.Println(episodesFiltered)
```

To get an idea of what is a valid filter, check out the [MongoDB documentation](https://docs.mongodb.com/manual/reference/operator/query/#query-selectors) on the subject.

### Sorting Documents in a Query

Sorting is a fundamental part of working with data. Rather than sorting the data within the Go application after the query executes, we can let MongoDB do the heavy lifting.

For sorting, we can leverage the `FindOptions` struct in the MongoDB Go Driver. `FindOptions` offers more than just sorting, but it is beyond the scope of this getting started example.

Let's say we want to query for all podcast episodes that are longer than 24 minutes, but we want to list them in descending order based on the duration. We could craft a query that looks like the following, while leveraging the `FindOptions` struct of the driver:

```go
opts := options.Find()
opts.SetSort(bson.D{{"duration", -1}})
sortCursor, err := episodesCollection.Find(ctx, bson.D{{"duration", bson.D{{"$gt", 24}}}}, opts)
if err != nil {
    log.Fatal(err)
}
var episodesSorted []bson.M
if err = sortCursor.All(ctx, &episodesSorted); err != nil {
    log.Fatal(err)
}
fmt.Println(episodesSorted)
```

Notice that a few things have changed in the above example. First, we're defining our `FindOptions` and the field we want to sort on. Within the `Find` function we are passing those options, but we're also using `bson.D` instead of `bson.M`. When using `bson.M`, the order of the fields does not matter, which makes it challenging for certain queries, more specifically range queries and similar. Instead we can use `bson.D` which respects the order that each field or operator uses.

## Updating MongoDB Documents with Golang

Rather than creating or querying for documents, we're going to push forward in our create, retrieve, update, delete (CRUD) demonstration and focus on updating documents within a collection.

### Update Data within a Collection

In the previous steps of this tutorial we've been working with a `podcasts` collection and an `episodes` collection. As a quick refresher, documents in the `podcasts` collection might look like this:

```json
{
    "_id": ObjectId("5d9e0173c1305d2a54eb431a"),
    "title": "The Polyglot Developer Podcast",
    "author": "Nic Raboy"
}
```

At some point in time you're going to find yourself needing to update documents, whether they look like my example, or reflect something of your own design. With the Go Driver for MongoDB, there are numerous options towards updating documents.

When updating a document, you'll need filter criteria as well as what should be updated within the data that matched the filter.

Take the following for example:

```go
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
fmt.Printf("Updated %v Documents!\n", result.ModifiedCount)
```

In the above example, we connect to our `podcasts` collection, which at this point we're going to assume has some documents in it by the previously mentioned design. If not, fill in the gap with your own collections or documents.

When using the `UpdateOne` function on a collection, we're choosing to update only one document that matches our filter, in this case an exact match on document ID. In our example, imagine the `5d9e0173c1305d2a54eb431a` value is a hexadecimal value that you've received from a RESTful API request that needs to be converted into a valid object ID. If a match was found, the `author` field on that document will be set. The filter and update criteria can leverage the full scope of the MongoDB Query Language (MQL). However, it should be noted that update operations can only use update operators as specified in the [documentation](https://docs.mongodb.com/manual/reference/operator/update/).

The `UpdateResult`, as seen as the `result` variable, has metric information about the operation. For example, it might say how many documents were updated, which in this case would be one or zero.

Let's say we wanted to update more than one document. We could make use of the `UpdateMany` function on a collection. Take the following for example:

```go
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
fmt.Printf("Updated %v Documents!\n", result.ModifiedCount)
```

Notice that the above code is rather familiar. The change being in the use of the `UpdateMany` function and in the filter. It probably doesn't make sense to update multiple documents based on a unique ID, so we changed it to other criteria. For every document matched based on the filter, the `author` field is updated.

In both examples, the `author` field was part of our very loosely defined schema. We can still set fields in an update operation, even if they don't already exist on a document, as long as the filter criteria matches.

For example, let's say you have the following document:

```json
{
    "_id": ObjectId("5d9e0173c1305d2a54eb431a"),
    "title": "The Polyglot Developer Podcast",
    "author": "Nic Raboy"
}
```

You do an `UpdateOne` or an `UpdateMany` with the `author` field as the filter. For the `$set` operation you include a `website` field that doesn't exist currently within the document. If the document matches based on the filter, the `website` field would now exist in that document.

### Replacing Documents in a Collection

While updating fields within one or more documents is useful, there are other scenarios where you might want to replace all fields in a document while maintaining the id of that document. For this, it makes sense to make use of the `ReplaceOne` function.

An example of the `ReplaceOne` function might look like the following:

```go
result, err = podcastsCollection.ReplaceOne(
    ctx,
    bson.M{"author": "Nic Raboy"},
    bson.M{
        "title":  "The Nic Raboy Show",
        "author": "Nicolas Raboy",
    },
)
fmt.Printf("Replaced %v Documents!\n", result.ModifiedCount)
```

In the above example, a filter for documents matching the `author` happens. When a match happens, the entire document is replaced with the `title` and `author` fields in the update criteria.

When working with the `ReplaceOne` function, update operators such as `$set` cannot be used since it is a complete replace rather than an update of particular fields.

## Deleting Documents in a MongoDB Collection with Golang

As of now we've explored the create, retrieve, and update aspects of CRUD using commands like `Find`, `InsertMany`, and `UpdateMany`. To finish things off, we're going to explore the final part of CRUD, which is the deleting of documents or even entire collections.

### Revisiting the Data Model for the Tutorial Series

Before we jump right into the removal of documents, it probably makes sense to revisit the data model we're going to be using to avoid confusion. If you've been keeping up so far, you'll know we are working with a `podcasts` collection and an `episodes` collection.

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

### Deleting a Single Document from a MongoDB Collection

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

Establishing a connection to the cluster and defining an application context can be seen towards the beginning of this tutorial.

### Deleting Many Documents from a MongoDB Collection

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

### Dropping a MongoDB Collection and All Documents within the Collection

Removing a single document or many documents isn't the only option. Entire collections can be dropped which would remove all documents and meta data, such as indexes, in the collection without using a filter. An example of this can be seen below:

```go
if err = podcastsCollection.Drop(ctx); err != nil {
    log.Fatal(err)
}
```

Dropping an entire collection will return an error if something has failed. If the collection doesn't exist, the driver will mask the server error and in this case return a `nil` error.

## Conclusion

You just saw how to connect to a MongoDB Atlas cluster and perform essential CRUD operations against a database and its collections, with the Go programming language (Golang).

If you're not yet using MongoDB Atlas, the M0 cluster is part of the forever free tier. However, using promotional code [NICRABOY200](https://www.mongodb.com/cloud) will get you $200 premium credit to be used towards a more powerful cluster.

In future tutorials, we're going to take our MongoDB with Go skills to the next level and focus on [mapping BSON documents to native Go data structures](https://www.mongodb.com/blog/post/quick-start-golang--mongodb--modeling-documents-with-go-data-structures), leveraging the MongoDB data aggregation pipeline, change streams, and transactions.
