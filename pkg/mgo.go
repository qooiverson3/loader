package pkg

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Mgo struct {
	Client *mongo.Client
}

func NewMgo(url string) *Mgo {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(url))
	if err != nil {
		log.Fatal(err)
	}

	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	return &Mgo{Client: client}
}
func (mgo *Mgo) Save(event bson.D) {
	col := mgo.Client.Database("nats").Collection("LAB")
	_, err := col.InsertOne(context.TODO(), event)
	if err != nil {
		log.Fatal(err)
	}
}
