package modb

import (
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongosh *mongo.Client
var db *mongo.Database

func Init(uri string) error {
	clientOptions := options.Client().ApplyURI(uri)

	var err error
	mongosh, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}
	s := strings.Split(uri, "/")
	db = mongosh.Database(s[len(s)-1])
	err = mongosh.Ping(context.TODO(), nil)
	if err != nil {
		return err
	}
	return nil
}
func Disconnect() error {
	return mongosh.Disconnect(context.TODO())
}
