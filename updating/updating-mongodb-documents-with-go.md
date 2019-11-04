# Quick Start: Updating MongoDB Documents with Go

In a [previous tutorial](https://), I demonstrated how to retrieve Documents from a MongoDB collection with Go. This was part of a getting started series which focused on the Go programming language (Golang) and MongoDB. Rather than creating or querying for Documents, we're going to push forward in our create, retrieve, update, delete (CRUD) demonstration and focus on updating documents within a collection.

## Tools and Versions for the Tutorial Series

To get the best results with this tutorial series, it will benefit you to use the same tools and versions of those tools that I'm using. Just to reiterate, I'm using the following:

- Go 1.13
- MongoDB Go Driver 1.1.2
- Visual Studio Code (VS Code)
- MongoDB Atlas with an M0 cluster (FREE)

While you may be able to use different versions than I'm using and find success, my recommendation is to try to match what I'm using as best as possible.

> It is FREE to use an M0 cluster on Atlas, but if you'd like some premium credit applied to your account for a more powerful cluster, you can use promotional code NRABOY200 within Atlas.

If you need help connecting your application to MongoDB Atlas, check out my tutorial [Quick Start: How to Get Connected to Your MongoDB Cluster with Go](https://), which goes through the basics.

## Updating Data within a Collection

If you've been keeping up with every tutorial in the series for getting started with Go, you'll remember that we're working with a `podcasts` collection and an `episodes` collection. As a quick refresher, Documents in the `podcasts` collection might look like this:

```json
{
    "_id": ObjectId("5d9e0173c1305d2a54eb431a"),
    "title": "The Polyglot Developer Podcast",
    "author": "Nic Raboy"
}
```

At some point in time you're going to find yourself needing to update Documents, whether they look like my example, or reflect something of your own design. With the Go Driver for MongoDB, there are numerous options towards updating Documents.

When updating a Document, you'll need filter criteria as well as what should be updated within the data that matched the filter.

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
fmt.Printf("Updated %v Documents!\n", result.MatchedCount)
```

In the above example, we connect to our `podcasts` collection, which at this point we're going to assume has some Documents in it by the previously mentioned design. If not, fill in the gap with your own collections or Documents.

When using the `UpdateOne` function on a collection, we're choosing to update only the first Document that matches our filter, in this case an exact match on Document ID. If a match was found, the `author` field on that document will be set. The filter and change criteria can leverage the full scope of the MongoDB Query Language (MQL).

The `UpdateResult`, as seen as the `result` variable, has metric information about the operation. For example, it might say how many Documents were updated.

Let's say we wanted to update more than one Document. We could make use of the `UpdateMany` function on a collection. Take the following for example:

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
fmt.Printf("Updated %v Documents!\n", result.MatchedCount)
```

Notice that the above code is rather familiar. The change being in the use of the `UpdateMany` function and in the filter. It probably doesn't make sense to update multiple Documents based on a unique ID, so we changed it to other criteria. For every match on the filter, the `author` field is updated.

In both examples, the `author` field was part of our very loosely defined schema. We can still set fields in an update operation, even if they don't already exist on a Document, as long as the filter criteria matches.

## Conclusion

Update is an important operator when thinking about the CRUD space. It would not be very efficient for developers to have to retrieve the data they wish to change, make the change followed by a create operation, then delete the old Document. Hence why being able to update is so great.

The final stage of the CRUD operation series is delete, something that we'll be focusing on in the next part of the getting started with Go and MongoDB series.