package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/gofiber/fiber/v2"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	route "github.com/AnuragProg/printit-microservices-monorepo/file/internal/api/routes"
	consts "github.com/AnuragProg/printit-microservices-monorepo/file/internal/constants"
	utils "github.com/AnuragProg/printit-microservices-monorepo/file/pkg/utils"
	auth "github.com/AnuragProg/printit-microservices-monorepo/file/proto_gen/authentication"
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

	// Setup rest app
	restApp := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024, // 10 MB
	})

	// connect to mongo database
	mongoCtx, mongoCtxCancel := context.WithCancel(context.Background())
	defer mongoCtxCancel()
	mongoClient, err := mongo.Connect(mongoCtx, options.Client().ApplyURI(MONGO_URI))
	if err != nil{
		panic(err.Error())
	}
	mongoDB := mongoClient.Database(consts.FILE_METADATA_DB)
	mongoDB.CreateCollection(context.Background(), consts.FILE_METADATA_COL)
	log.Println("Connected to Mongo...")


	// connect to minio client
	minioClient, err := minio.New(MINIO_URI, &minio.Options{
		Creds: credentials.NewStaticV4(MINIO_SERVER_ACCESS_KEY, MINIO_SERVER_SECRET_KEY, ""),
		Transport: &http.Transport{
			MaxIdleConns: 100,
			IdleConnTimeout: 60*time.Second,
		},
	})
	if err != nil{
		panic(err.Error())
	}
	minioClient.MakeBucket(context.Background(), consts.FILE_BUCKET, minio.MakeBucketOptions{})
	log.Println("Connected to Minio...")

	// connect to grpc servers
	authGrpcConn, err := grpc.Dial(AUTH_GRPC_URI, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil{
		panic(err.Error())
	}
	defer authGrpcConn.Close()
	log.Println("Connected to Auth GRPC server...")

	// create grpc client
	authGrpcClient := auth.NewAuthenticationClient(authGrpcConn)

	// setup top level routes
	fileRouter := restApp.Group("/file")
	fileRoute := route.FileRoute{
		Router: &fileRouter,

		MinioClient: minioClient,

		MongoDB: mongoClient.Database(consts.FILE_METADATA_DB),
		MongoClient: mongoClient,

		AuthGrpcClient: &authGrpcClient,
	}
	fileRoute.SetupRoutes()


	// start rest server
	log.Printf("Listening on :%v\n", REST_PORT)
	restApp.Listen(fmt.Sprintf(":%v", REST_PORT))
}
