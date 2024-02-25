package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/gofiber/fiber/v2"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	route "github.com/AnuragProg/printit-microservices-monorepo/file/internal/api/routes"
	auth "github.com/AnuragProg/printit-microservices-monorepo/file/proto_gen/authentication"
)


func main(){

	// Setup rest app
	app := fiber.New()

	// connect to mongo database
	mongo_ctx, mongo_ctx_cancel := context.WithCancel(context.Background())
	defer mongo_ctx_cancel()
	mongo_client, err := mongo.Connect(mongo_ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil{
		panic(err.Error())
	}
	log.Println("Connected to Mongo...")

	// connect to grpc servers
	auth_grpc_conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil{
		panic(err.Error())
	}
	defer auth_grpc_conn.Close()
	log.Println("Connected to Auth GRPC server...")

	// create grpc clients
	auth_grpc_client := auth.NewAuthenticationClient(auth_grpc_conn)


	// setup top level routes
	fileRouter := app.Group("/file")
	fileRoute := route.FileRoute{
		MongoClient: mongo_client,
		Router: &fileRouter,
		AuthGrpcClient: &auth_grpc_client,
	}
	fileRoute.SetupRoutes()


	// start rest server
	log.Println("Listening on 3001")
	app.Listen(":3001")
}
