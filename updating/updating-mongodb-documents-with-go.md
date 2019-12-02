# Quick Start: Updating MongoDB Documents with Go

In a [previous tutorial](https://www.mongodb.com/blog/post/quick-start-golang--mongodb--how-to-read-documents), I demonstrated how to retrieve documents from a MongoDB collection with Go. This was part of a getting started series which focused on the Go programming language (Golang) and MongoDB. Rather than creating or querying for documents, we're going to push forward in our create, retrieve, update, delete (CRUD) demonstration and focus on updating documents within a collection.

## Tools and Versions for the Tutorial Series

To get the best results with this tutorial series, it will benefit you to use the same tools and versions of those tools that I'm using. Just to reiterate, I'm using the following:

- Go 1.13
- MongoDB Go Driver 1.1.2
- Visual Studio Code (VS Code)
- MongoDB Atlas with an M0 cluster (FREE)

While you may be able to use different versions than I'm using and find success, my recommendation is to try to match what I'm using as best as possible.

> It is FREE to use an M0 cluster on Atlas, but if you'd like some premium credit applied to your account for a more powerful cluster, you can use promotional code NRABOY200 within Atlas.

If you need help connecting your application to MongoDB Atlas, check out my tutorial [Quick Start: How to Get Connected to Your MongoDB Cluster with Go](https://www.mongodb.com/blog/post/quick-start-golang--mongodb--starting-and-setup), which goes through the basics.

## Updating Data within a Collection

If you've been keeping up with every tutorial in the series for getting started with Go, you'll remember that we're working with a `podcasts` collection and an `episodes` collection. As a quick refresher, documents in the `podcasts` collection might look like this:

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

```golang
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

```golang
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

## Replacing Documents in a Collection

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

## Conclusion

Update is an important operator when thinking about the CRUD space. It would not be very efficient for developers to have to retrieve the data they wish to change, make the change followed by a create operation, then delete the old document. Hence why being able to update is so great.

The final stage of the CRUD operation series is delete, something that we'll be focusing on in the next part of the getting started with Go and MongoDB series.