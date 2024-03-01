package client

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	consts "github.com/AnuragProg/printit-microservices-monorepo/file/internal/constants"
)


func GetMongoClientAndDB(mongoURI string) (*mongo.Client, *mongo.Database, error){

	// create mongo client
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil{
		return nil, nil, err
	}

	// Check whether mongo is connected
	if err := mongoClient.Ping(context.Background(), nil); err!=nil{
		return nil, nil, err
	}

	// connect and create corresponding database and collections
	mongoDB := mongoClient.Database(consts.FILE_METADATA_DB)
	mongoDB.CreateCollection(context.Background(), consts.FILE_METADATA_COL)
	log.Println("Connected to Mongo...")

	return mongoClient, mongoDB, nil
}
