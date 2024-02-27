package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"

	auth "github.com/AnuragProg/printit-microservices-monorepo/file/proto_gen/authentication"
	route "github.com/AnuragProg/printit-microservices-monorepo/file/internal/api/routes"
	client "github.com/AnuragProg/printit-microservices-monorepo/file/internal/client"
	consts "github.com/AnuragProg/printit-microservices-monorepo/file/internal/constants"
	utils "github.com/AnuragProg/printit-microservices-monorepo/file/pkg/utils"
)

var (
	MONGO_URI = utils.GetenvOrDefault(os.Getenv("MONGO_URI"), "mongodb://localhost:27017/printit")
	AUTH_GRPC_URI = utils.GetenvOrDefault(os.Getenv("AUTH_GRPC_URI"), "localhost:50051")
	REST_PORT = utils.GetenvOrDefault(os.Getenv("REST_PORT"), "3001")

	MINIO_URI = utils.GetenvOrDefault(os.Getenv("MINIO_URI"), "localhost:9000")
	MINIO_SERVER_ACCESS_KEY = utils.GetenvOrDefault(os.Getenv("MINIO_SERVER_ACCESS_KEY"), "minio-access-key")
	MINIO_SERVER_SECRET_KEY = utils.GetenvOrDefault(os.Getenv("MINIO_SERVER_SECRET_KEY"), "minio-secret-key")
)

func main(){

	// connect to mongo database
	mongoClient, mongoDB, err := client.GetMongoClientAndDB(MONGO_URI)
	if err != nil{
		panic(err.Error())
	}
	defer mongoClient.Disconnect(context.Background())
	mongoFileMetadataCol := mongoDB.Collection(consts.FILE_METADATA_COL)

	// connect to minio client
	minioClient, err := client.GetMinioClient(MINIO_URI, MINIO_SERVER_ACCESS_KEY, MINIO_SERVER_SECRET_KEY)
	if err != nil{
		panic(err.Error())
	}

	// connect to grpc servers
	authGrpcConn, err := client.GetAuthGrpcConnAndClient(AUTH_GRPC_URI)
	if err != nil{
		panic(err.Error())
	}
	defer authGrpcConn.Close()
	authGrpcClient := auth.NewAuthenticationClient(authGrpcConn)

	// Setup rest app
	restApp := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024, // 10 MB
	})

	// setup top level routes
	fileRouter := restApp.Group("/file")
	fileRoute := route.FileRoute{
		Router: &fileRouter,
		MinioClient: minioClient,
		MongoFileMetadataCol: mongoFileMetadataCol,
		AuthGrpcClient: &authGrpcClient,
	}
	fileRoute.SetupRoutes()

	// start rest server
	log.Info("Listening on :", REST_PORT)
	restApp.Listen(fmt.Sprintf(":%v", REST_PORT))
}
