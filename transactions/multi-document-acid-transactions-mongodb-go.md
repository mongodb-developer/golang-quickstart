# Multi-Document ACID Transactions in MongoDB with Go

The past few months have been an adventure when it comes to getting started with MongoDB using the Go programming language (Golang). We've explored everything from create, retrieve, update, and delete (CRUD) operations, to data modeling, and to change streams. To bring this series to a solid finish, we're going to take a look at a popular requirement that a lot of organizations need, and that requirement is transactions.

So why would you want transactions?

There are some situations where you might need atomicity of reads and writes to multiple documents within a single collection or multiple collections. This isn't always a necessity, but in some cases it might be.

Take the following for example.

Let's say you want to create documents in one collection that depend on documents in another collection existing. Or let's say you have schema validation rules in place on your collection. In the scenario that you're trying to create documents and the related document doesn't exist or your schema validation rules fail, you don't want the operation to proceed. Instead, you'd probably want to rollback to before it happened.

There are other reasons that you might use transactions, but you can use your imagination for those.

In this tutorial, we're going to look at what it takes to use transactions with Golang and MongoDB. Our example will rely more on schema validation rules passing, but it isn't a limitation.

## Understanding the Data Model and Applying Schema Validation

Since we've continued the same theme throughout the series, I think it'd be a good idea to have a refresher on the data model that we'll be using for this example.

In the past few tutorials we've explored working with potential podcast data in various collections. For example, our Go data model looks something like this:

```go
type Episode struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Podcast     primitive.ObjectID `bson:"podcast,omitempty"`
	Title       string             `bson:"title,omitempty"`
	Description string             `bson:"description,omitempty"`
	Duration    int32              `bson:"duration,omitempty"`
}
```

The fields in the data structure are mapped to MongoDB document fields through the BSON annotations. You can learn more about using these annotations in the [previous tutorial](https://www.mongodb.com/blog/post/quick-start-golang--mongodb--modeling-documents-with-go-data-structures) I wrote on the subject.

While we had other collections, we're going to focus strictly on the `episodes` collection for this example.

Rather than coming up with complicated code for this example to demonstrate operations that fail or should be rolled back, we're going to go with schema validation to force fail some operations. Let's assume that no episode should be less than two minutes in duration, otherwise it is not valid. Rather than implementing this, we can use features baked into MongoDB.

Take the following schema validation logic:

```json
{
    "$jsonSchema": {
        "additionalProperties": true,
        "properties": {
            "duration": {
                "bsonType": "int",
                "minimum": 2
            }
        }
    }
}
```

The above logic would be applied using the MongoDB CLI or with Compass, but we're essentially saying that our schema for the `episodes` collection can contain any fields in a document, but the `duration` field must be an integer and it must be at least two. Could our schema validation be more complex? Absolutely, but we're all about simplicity in this example. If you want to learn more about schema validation, check out [this awesome tutorial](https://www.mongodb.com/blog/post/json-schema-validation--locking-down-your-model-the-smart-way) on the subject.

Now that we know the schema and what will cause a failure, we can start implementing some transaction code that will commit or roll back changes.

## Starting and Committing Transactions

Before we dive into starting a session for our operations and committing transactions, let's establish a base point in our project. Let's assume that your project has the following boilerplate MongoDB with Go code:

```go
package main

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Episode represents the schema for the "Episodes" collection
type Episode struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Podcast     primitive.ObjectID `bson:"podcast,omitempty"`
	Title       string             `bson:"title,omitempty"`
	Description string             `bson:"description,omitempty"`
	Duration    int32              `bson:"duration,omitempty"`
}

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

Now let's also assume that you've correctly included the MongoDB Go driver as seen in a previous tutorial titled, [How to Get Connected to Your MongoDB Cluster with Go](https://www.mongodb.com/blog/post/quick-start-golang--mongodb--starting-and-setup).

The goal here will be to try to insert a document that complies with our schema validation as well as a document that doesn't, so that we have a commit that doesn't happen.

```go
// ...

func main() {
    // ...

	session, err := client.StartSession()
	if err != nil {
		panic(err)
	}
	defer session.EndSession(context.Background())

	err = mongo.WithSession(context.Background(), session, func(sessionContext mongo.SessionContext) error {
		if err = session.StartTransaction(); err != nil {
			return err
		}
		result, err := episodesCollection.InsertOne(
			sessionContext,
			Episode{
				Title:    "A Transaction Episode for the Ages",
				Duration: 15,
			},
		)
		if err != nil {
			return err
		}
        fmt.Println(result.InsertedID)
		result, err = episodesCollection.InsertOne(
			sessionContext,
			Episode{
				Title:    "Transactions for All",
				Duration: 1,
			},
		)
		if err != nil {
			return err
		}
		if err = session.CommitTransaction(sessionContext); err != nil {
			return err
        }
        fmt.Println(result.InsertedID)
		return nil
	})
	if err != nil {
		if abortErr := session.AbortTransaction(context.Background()); abortErr != nil {
			panic(abortErr)
		}
		panic(err)
	}
}
```

In the above code we start by starting a session which will encapsulate everything we want to do with atomicity. After, we start a transaction which we'll use to commit everything in the session.

A `Session` represents a MongoDB logical session and can be used to enable casual consistency for a group of operations or to execute operations in an ACID transaction. More information on how they work in Go can be found in the [documentation](https://godoc.org/go.mongodb.org/mongo-driver/mongo#Session).

Inside the session we are doing two `InsertOne` operations. The first would succeed because it doesn't violate any of our schema validation rules. It will even print out an object id when it's done. However, the second operation will fail because it is less than two minutes. The `CommitTransaction` won't ever succeed because of the error that the second operation created. When the `WithSession` function returns the error that we created, the transaction is aborted using the `AbortTransaction` function. For this reason, neither of the `InsertOne` operations will show up in the database.

## Using a Convenient Transactions API

Starting and committing transactions from within a logical session isn't the only way to work with ACID transactions using Golang and MongoDB. Instead, we can use what might be thought of as a more convenient transactions API.

Take the following adjustments to our code:

```go
// ...

func main() {
	// ...

	session, err := client.StartSession()
	if err != nil {
		panic(err)
	}
	defer session.EndSession(context.Background())

	_, err = session.WithTransaction(context.Background(), func(sessionContext mongo.SessionContext) (interface{}, error) {
		result, err := episodesCollection.InsertOne(
			sessionContext,
			Episode{
				Title:    "A Transaction Episode for the Ages",
				Duration: 15,
			},
		)
		if err != nil {
			return nil, err
		}
		result, err = episodesCollection.InsertOne(
			sessionContext,
			Episode{
				Title:    "Transactions for All",
				Duration: 2,
			},
		)
		if err != nil {
			return nil, err
		}
		return result, err
	})
	if err != nil {
		panic(err)
	}
}
```

Instead of using `WithSession`, we are now using `WithTransaction`, which handles starting a transaction, executing some application code, and then committing or aborting the transaction based on the success of that application code. Not only that, but retries can happen for specific errors if certain operations fail.

## Conclusion

You just saw how to use transactions with the MongoDB Go driver. While in this example we used schema validation to determine if a commit operation succeeds or fails, you could easily apply your own application logic within the scope of the session.

If you want to catch up on other tutorials in the getting started with Golang series, you can find some below:

- [How to Get Connected to Your MongoDB Cluster with Go](https://www.mongodb.com/blog/post/quick-start-golang--mongodb--starting-and-setup)
- [Creating MongoDB Documents with Go](https://www.mongodb.com/blog/post/quick-start-golang--mongodb--how-to-create-documents)
- [Retrieving and Querying MongoDB Documents with Go](https://www.mongodb.com/blog/post/quick-start-golang--mongodb--how-to-read-documents)
- [Updating MongoDB Documents with Go](https://www.mongodb.com/blog/post/quick-start-golang--mongodb--how-to-update-documents)
- [Deleting MongoDB Documents with Go](https://www.mongodb.com/blog/post/quick-start-golang--mongodb--how-to-delete-documents)
- [Modeling MongoDB Documents with Native Go Data Structures](https://www.mongodb.com/blog/post/quick-start-golang--mongodb--modeling-documents-with-go-data-structures)
- [Performing Complex MongoDB Data Aggregation Queries with Go](https://www.mongodb.com/blog/post/quick-start-golang--mongodb--data-aggregation-pipeline)
- [Reacting to Database Changes with MongoDB Change Streams and Go](https://)

Since transactions brings this tutorial series to a close, make sure you keep a lookout for more tutorials which focus on more niche and interesting topics that apply everything that was taught while getting started.