# Quick Start: How to Get Connected to Your MongoDB Cluster with Go

In the first tutorial, which can best be named a quick start into MongoDB development with the Go programming language (Golang), we're going to be exploring how to establish connections between the language and the database.

When it comes to future tutorials in the series, expect content on the following:

- Database create, retrieve, update, and delete (CRUD) operations.
- A look into MongoDB aggregation queries.
- Watching change streams in MongoDB.
- Multi-Document A.C.I.D. transactions.

Go is one of the more recent of the officially supported technologies with MongoDB, and in my personal opinion, it is one of the most awesome!

Throughout these tutorials, I'll be using Visual Studio Code (VSCode) for development, and I'll be connecting to a MongoDB Atlas cluster. The assumption is that you're using Go 1.13 or newer and that it is already properly installed and configured on your computer. It is also assumed that an Atlas cluster has already been created.

> Get started with an M0 cluster on [MongoDB Atlas](https://www.mongodb.com/cloud) today. It's free forever and you'll be able to work alongside this blog series. Use promo code NRABOY200 when you sign up and you'll get an extra $200.00 credit applied to your account.

If you're using a different IDE, OS, ect., the walkthrough might be slightly different, but the code will be pretty much the same.

## Getting Started

Within your **$GOPATH**, create a new project directly titled **quickstart** and add a **main.go** file to that project.

For this particular tutorial, all code will be added to the **main.go** file.

With the project created, the next step is to install the MongoDB Go Driver through the Go Package Manager. To do this, execute the following from the command line:

```bash
$ dep ensure -add "go.mongodb.org/mongo-driver/mongo@~1.1.2"
```

Note that for this tutorial we're using `dep` to manage our packages and we're using version 1.1.2 of the MongoDB Go Driver. If you don't have `dep`, you can install it through the [official documentation](https://github.com/golang/dep).

With the driver installed, open the project's **main.go** file and add the following boilerplate code:

```golang
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

```golang
func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI("<ATLAS_URI_HERE>"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
}
```

There are a few things that are happening in the above code. First we're configuring our client to use the correct URI, but we're not yet connecting to it. Assuming nothing is malformed and no error was thrown, we can define a timeout duration that we want to use when trying to connect. The ten seconds I used might be a little too generous for your needs, but feel free to play around with the value that makes the most sense to you.

After connecting, if there isn't an error, we can defer the closing of the connection for when the `main` function exits. This will keep the connection to the database open until we're done.

So if no errors were thrown, can we be sure that we're really connected? If you're concerned, we can ping the cluster from within our application:

```golang
err = client.Ping(ctx, readpref.Primary())
if err != nil {
    log.Fatal(err)
}
```

The above code uses our connected client and the context that we had previously defined. If there is no error, the ping was a success!

We can take things a step further by listing the available databases on our MongoDB Atlas cluster. Within the `main` function, add the following:

```golang
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

```golang
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
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
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

## Conclusion

You just saw how to connect to a MongoDB Atlas cluster with the Go programming language. If you decide not to use Atlas, the code will still work, you'll just have a different connection string.

Go is a powerful technology and combined with MongoDB you can accomplish anything from web applications, to desktop applications, and everything in-between.

Stay tuned for the next tutorial in the series which focuses on modeling documents with BSON and native Go data structures.