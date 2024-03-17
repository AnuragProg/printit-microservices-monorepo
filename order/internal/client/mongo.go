package client

import (
	"context"

	"github.com/gofiber/fiber/v2/log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	consts "github.com/AnuragProg/printit-microservices-monorepo/internal/constant"
)

type MongoDB struct {
	Client 	*mongo.Client
	DB		  	*mongo.Database
	OrderCol	*mongo.Collection
}

func GetMongoDB(mongoURI string) (*MongoDB, error) {

	// create mongo client
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil{
		return nil, err
	}

	// Check whether mongo is connected
	if err := mongoClient.Ping(context.Background(), nil); err!=nil{
		return nil, err
	}
	log.Info("Connected to Mongo...")

	// connect and create corresponding database and collections
	mongoDB := mongoClient.Database(consts.ORDER_DB)

	mongoDB.CreateCollection(context.Background(), consts.ORDER_COL)
	orderCol := mongoDB.Collection(consts.ORDER_COL)

	return &MongoDB{
		Client: mongoClient,
		DB: mongoDB,
		OrderCol: orderCol,
	}, nil
}

