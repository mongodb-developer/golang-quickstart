# Quick Start: Retrieving and Querying MongoDB Documents with Go

In a [previous tutorial](https://), I had written about creating documents to be inserted into MongoDB with the Go programming language. In that tutorial we explored the `Insert` and `InsertMany` functions while making use of `bson.D`, `bson.M`, and similar MongoDB data types.

This time around, we're going to explore reading documents from MongoDB and creating queries to retrieve documents based on certain criteria. This will all be done with Golang and the MongoDB Go Driver.

## Tools and Versions for the Tutorial Series

I wanted to take a moment to reiterate the tools and versions that I'm using within this tutorial series:

- Go 1.13
- Visual Studio Code (VS Code)
- MongoDB Atlas with an M0 free cluster
- MongoDB Go Driver 1.1.2

To get the best experience while following this tutorial, try to match the versions as best as possible. However, other versions may still work without issue.

> You can get started with an M0 cluster on [MongoDB Atlas](https://www.mongodb.com/cloud) for free. If sign up using the promotional code NRABOY200, you'll receive premium credit applied to your account.

If you need help connecting to MongoDB Atlas, installing the MongoDB Go Driver, or getting familiar with creating documents, I encourage you to check out one of the previous tutorials in the series.

## The Previously Created Data

When thinking back to the data that we created in the [previous tutorial](https://), we know that we have a `podcasts` collection with data that looks something like this:

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

## Reading All Documents from a Collection

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

If you think back to the [first tutorial in the series](https://www.mongodb.com/blog/post/quick-start-golang--mongodb--starting-and-setup), the `Find` function might look similar to the `ListDatabaseNames` function. We can provide a context and some query parameters, and get our results. In this example `bson.M` represents a map of fields in no particular order. However, because we're trying to return all documents, there aren't any fields in our query.

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
for cursor.Next(ctx) {
    var episode bson.M
    if err = cursor.Decode(&episode); err != nil {
        log.Fatal(err)
    }
    fmt.Println(episode)
}
cursor.Close(ctx)
```

In both the `*mongo.Cursor.All` and `*mongo.Cursor.Find` examples, the data is loaded into `bson.M` data structures which behave as maps. We'll explore marshalling and unmarshalling the data to custom native Go data structures in a later tutorial.

## Reading a Single Document from a Collection

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

## Querying Documents from a Collection with a Filter

In the previous examples of `Find` and `FindOne`, we've seen the `bson.M` filter parameter, even though it wasn't used for filtering. We can make use of the full power of the MongoDB Query Language (MQL) to filter the results of our queries, simply by populating the map.

Let's say that we want to filter our results to only include podcast episodes that are exactly 25 minutes. We could do something like the following:

```go
cursor, err := episodesCollection.Find(ctx, bson.M{"duration": 25})
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

To get an idea of what is a valid filter, check out the [MongoDB documentation](https://docs.mongodb.com/manual/reference/operator/query/#query-selectors) on the subject.

## Sorting Documents in a Query

Sorting is a fundamental part of working with data. Rather than sorting the data within the Go application after the query executes, we can let MongoDB do the heavy lifting.

For sorting, we can leverage the `FindOptions` struct in the MongoDB Go Driver. `FindOptions` offers more than just sorting, but it is beyond the scope of this getting started example.

Let's say we want to query for all podcast episodes that are longer than 24 minutes, but we want to list them in descending order based on the duration. We could craft a query that looks like the following, while leveraging the `FindOptions` struct of the driver:

```go
opts := options.Find()
opts.SetSort(bson.D{{"duration", -1}})
cursor, err = episodesCollection.Find(ctx, bson.D{{"duration", bson.D{{"$gt", 24}}}}, opts)
if err != nil {
    log.Fatal(err)
}
var episodes []bson.M
if err = cursor.All(ctx, &episodes); err != nil {
    log.Fatal(err)
}
fmt.Println(episodes)
```

Notice that a few things have changed in the above example. First, we're defining our `FindOptions` and the field we want to sort on. Within the `Find` function we are passing those options, but we're also using `bson.D` instead of `bson.M`. When using `bson.M`, the order of the fields does not matter, which makes it challenging for certain queries, more specifically range queries and similar. Instead we can use `bson.D` which respects the order that each field or operator uses.

## Conclusion

There are many ways to read data from a MongoDB database, whether that be by using filters, reading a single document, sorting at a database level, or something else. We saw a few of the common strategies when it comes to application development.

In the next tutorial of the series, we're going to explore how to update Documents that have already been created, using the Go programming language.