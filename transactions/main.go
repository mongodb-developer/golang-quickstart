package main

import (
	"context"
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

	session, err := client.StartSession()
	if err != nil {
		panic(err)
	}
	defer session.EndSession(context.Background())

	// err = mongo.WithSession(context.Background(), session, func(sessionContext mongo.SessionContext) error {
	// 	if err = session.StartTransaction(); err != nil {
	// 		panic(err)
	// 	}
	// 	result, err := episodesCollection.InsertOne(
	// 		sessionContext,
	// 		Episode{
	// 			Title:    "A Transaction Episode for the Ages",
	// 			Duration: 15,
	// 		},
	// 	)
	// 	fmt.Println(result.InsertedID)
	// 	result, err = episodesCollection.InsertOne(
	// 		sessionContext,
	// 		Episode{
	// 			Title:    "Transactions for All",
	// 			Duration: 2,
	// 		},
	// 	)
	// 	if err = session.CommitTransaction(sessionContext); err != nil {
	// 		panic(err)
	// 	}
	// 	fmt.Println(result.InsertedID)
	// 	return nil
	// })
	// if err != nil {
	// 	panic(err)
	// }

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
