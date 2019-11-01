# Quick Start: Retrieving and Querying MongoDB Documents with Go

In a [previous tutorial](https://), I had written about creating Documents to be inserted into MongoDB with the Go programming language. In that tutorial we explored the `Insert` and `InsertMany` functions while making use of `bson.D`, `bson.M`, and similar MongoDB data types.

This time around, we're going to explore reading Documents from MongoDB and creating queries to retrieve Documents based on certain criteria. This will all be done with Golang and the MongoDB Go Driver.

## Tools and Versions for the Tutorial Series

I wanted to take a moment to reiterate the tools and versions that I'm using within this tutorial series:

- Go 1.13
- Visual Studio Code (VS Code)
- MongoDB Atlas with an M0 free cluster
- MongoDB Go Driver 1.1.2

To get the best experience while following this tutorial, try to match the versions as best as possible. However, other versions may still work without issue.

> You can get started with an M0 cluster on [MongoDB Atlas](https://www.mongodb.com/cloud) for free. If sign up using the promotional code NRABOY200, you'll receive premium credit applied to your account.

If you need help connecting to MongoDB Atlas, installing the MongoDB Go Driver, or getting familiar with creating Documents, I encourage you to check out one of the previous tutorials in the series.

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

Reading all data from a collection consists of making the request, then looping through the results. Knowing what fields exist on each of the documents isn't too important, only knowing the collection name itself.

A simple example of this can be done through the following:

```golang
cursor, err := episodesCollection.Find(ctx, bson.M{})
if err != nil {
    log.Fatal(err)
}
for cursor.Next(ctx) {
    var episode bson.D
    cursor.Decode(&episode)
    fmt.Println(episode)
}
cursor.Close(ctx)
```

If you think back to the [first tutorial in the series](https://), the `Find` function might look similar to the `ListDatabaseNames` function. We can provide a context and some query parameters, and get our results. In this example `bson.M` represents a map of fields in no particular order. However, because we're trying to return all Documents, there aren't any fields in our query.

Assuming no error happens, the results will exist in a MongoDB Cursor. In this simple example, the Cursor is iterated and each Document within the iteration is printed. Remember, for now at least, `bson.D` represents a MongoDB BSON Document.

To get an idea of what the printed data would look like, take the following for example:

```
[{_id ObjectID("5d9f4701e9948c0f65c9165d")} {podcast ObjectID("5d9e0173c1305d2a54eb431a")} {title GraphQL for API Development} {description Learn about GraphQL from the co-creator of GraphQL, Lee Byron.} {duration 25}]
[{_id ObjectID("5dbb21faecb6c4837b0eff62")} {podcast ObjectID("5dbb2038e7841e4744f57a9c")} {title Progressive Web Application Development} {description Learn about PWA development with Tara Manicsic.} {duration 32}]
```

We'll explore marshalling and unmarshalling the response to Go data structures in a later tutorial.

## Reading a Single Document from a Collection

Retrieving all Documents from a collection doesn't always make sense. Instead, sometimes it makes sense to only return a single document. Instead of using the `Find` function, the `FindOne` function can be used.

Take the following example:

```golang
var podcast bson.D
err = podcastsCollection.FindOne(ctx, bson.M{}).Decode(&podcast)
if err != nil {
    log.Fatal(err)
}
fmt.Println(podcast)
```

In the above example, a `FindOne` is executed without any particular query filter on the data. Rather than returning a Cursor, the single result can be decoded directly into the `bson.D` object. If there is no error, the result will print.

To get an idea of what the result might look like, take the following:

```
[{_id ObjectID("5d9e0173c1305d2a54eb431a")} {title The Polyglot Developer Podcast} {author Nicolas Raboy}]
```

Because the results in the `Find` and `FindOne` use `bson.D`, the format of the results will be the same.

## Querying Documents from a Collection with a Filter

In the previous examples of `Find` and `FindOne`, we've seen the `bson.M` filter parameter, even though it wasn't used for filtering. We can make use of the full power of the MongoDB Query Language (MQL) to filter the results of our queries, simply by populating the map.

Let's say that we want to filter our results to only include podcast episodes that are exactly 25 minutes. We could do something like the following:

```golang
cursor, err := episodesCollection.Find(ctx, bson.M{"duration": 25})
if err != nil {
    log.Fatal(err)
}
for cursor.Next(ctx) {
    var episode bson.D
    cursor.Decode(&episode)
    fmt.Println(episode)
}
cursor.Close(ctx)
```

Just to reiterate, the full scope of the MongoDB Query Language can be used in our filters, something we'll see in the next example.

## Sorting Documents in a Query

Sorting is a fundamental part of working with data. Rather than sorting the data within the Go application after the query executes, we can let MongoDB do the heavy lifting.

For sorting, we can leverage the `FindOptions` part of the MongoDB Go Driver. `FindOptions` offers more than just sorting, but it is beyond the scope of this getting started example.

Let's say we want to query for all podcast episodes that are longer than 24 minutes, but we want to list them in descending order based on the duration. We could craft a query that looks like the following, while leveraging the `FindOptions` part of the driver:

```golang
opts := options.Find()
opts.SetSort(bson.D{{"duration", -1}})
cursor, err = episodesCollection.Find(ctx, bson.D{{"duration", bson.D{{"$gt", 24}}}}, opts)
if err != nil {
    log.Fatal(err)
}
for cursor.Next(ctx) {
    var episode bson.D
    cursor.Decode(&episode)
    fmt.Println(episode)
}
cursor.Close(ctx)
```

Notice that a few things have changed in the above example. First, we're defining our `FindOptions` and the field we want to sort on. Within the `Find` function we are passing those options, but we're also using `bson.D` instead of `bson.M`.

When using `bson.M`, the order of the fields does not matter, which makes it challenging for certain queries, more specifically range queries and similar. Instead we can use `bson.D` which respects the order that each field or operator uses.

## Conclusion

There are many ways to read data from a MongoDB database, whether that be by using filters, reading a single document, sorting at a database level, or something else. We saw a few of the common strategies when it comes to application development.

In the next tutorial of the series, we're going to explore how to update Documents that have already been created, using the Go programming language.