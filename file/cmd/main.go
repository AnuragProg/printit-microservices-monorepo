package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/gofiber/fiber/v2"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	route "github.com/AnuragProg/printit-microservices-monorepo/file/internal/api/routes"
	utils "github.com/AnuragProg/printit-microservices-monorepo/file/pkg/utils"
	auth "github.com/AnuragProg/printit-microservices-monorepo/file/proto_gen/authentication"
)

var (
	MONGO_URI = utils.GetenvOrDefault(os.Getenv("MONGO_URI"), "mongodb://localhost:27017/printit")
	AUTH_GRPC_URI = utils.GetenvOrDefault(os.Getenv("AUTH_GRPC_URI"), "localhost:50051")
	REST_PORT = utils.GetenvOrDefault(os.Getenv("REST_PORT"), "3001")
)

func main(){

	// Setup rest app
	restApp := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024,
	})

	// connect to mongo database
	mongo_ctx, mongo_ctx_cancel := context.WithCancel(context.Background())
	defer mongo_ctx_cancel()
	mongo_client, err := mongo.Connect(mongo_ctx, options.Client().ApplyURI(MONGO_URI))
	if err != nil{
		panic(err.Error())
	}
	log.Println("Connected to Mongo...")

	// connect to grpc servers
	auth_grpc_conn, err := grpc.Dial(AUTH_GRPC_URI, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil{
		panic(err.Error())
	}
	defer auth_grpc_conn.Close()
	log.Println("Connected to Auth GRPC server...")

	// create grpc clients
	auth_grpc_client := auth.NewAuthenticationClient(auth_grpc_conn)


	// setup top level routes
	fileRouter := restApp.Group("/file")
	fileRoute := route.FileRoute{
		MongoClient: mongo_client,
		Router: &fileRouter,
		AuthGrpcClient: &auth_grpc_client,
	}
	fileRoute.SetupRoutes()


	// start rest server
	log.Printf("Listening on :%v\n", REST_PORT)
	restApp.Listen(fmt.Sprintf(":%v", REST_PORT))
}
