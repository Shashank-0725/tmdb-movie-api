package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func Init(uri string) error {
	var err error
	Client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	return err
}