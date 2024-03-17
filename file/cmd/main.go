package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"

	"google.golang.org/grpc"

	file_service "github.com/AnuragProg/printit-microservices-monorepo/internal/service"
	pb "github.com/AnuragProg/printit-microservices-monorepo/proto_gen/file"

	route "github.com/AnuragProg/printit-microservices-monorepo/internal/api/route"
	client "github.com/AnuragProg/printit-microservices-monorepo/internal/client"
	consts "github.com/AnuragProg/printit-microservices-monorepo/internal/constant"
	utils "github.com/AnuragProg/printit-microservices-monorepo/pkg/util"
	auth "github.com/AnuragProg/printit-microservices-monorepo/proto_gen/authentication"
)

var (
	MONGO_URI = utils.GetenvOrDefault("MONGO_URI", "mongodb://localhost:27017")
	AUTH_GRPC_URI = utils.GetenvOrDefault("AUTH_GRPC_URI", "localhost:50051")
	REST_PORT = utils.GetenvOrDefault("REST_PORT", "3001")
	GRPC_PORT = utils.GetenvOrDefault("GRPC_PORT", "50052")

	MINIO_URI = utils.GetenvOrDefault("MINIO_URI", "localhost:9000")
	MINIO_SERVER_ACCESS_KEY = utils.GetenvOrDefault("MINIO_SERVER_ACCESS_KEY", "minio-access-key")
	MINIO_SERVER_SECRET_KEY = utils.GetenvOrDefault("MINIO_SERVER_SECRET_KEY", "minio-secret-key")
)

func main(){

	// connect to mongo database
	mongoClient, mongoDB, err := client.GetMongoClientAndDB(MONGO_URI)
	if err != nil{
		log.Error(err.Error())
		panic(err.Error())
	}
	defer mongoClient.Disconnect(context.Background())
	mongoFileMetadataCol := mongoDB.Collection(consts.FILE_METADATA_COL)

	// connect to minio client
	minioClient, err := client.GetMinioClient(MINIO_URI, MINIO_SERVER_ACCESS_KEY, MINIO_SERVER_SECRET_KEY)
	if err != nil{
		log.Error(err.Error())
		panic(err.Error())
	}

	// setup grpc server for file service functionalitities
	grpcServer := grpc.NewServer()
	pb.RegisterFileServer(grpcServer, file_service.NewFileService(mongoFileMetadataCol))
	go func(){
		grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%v", GRPC_PORT))
		if err != nil{
			log.Error(err.Error())
			panic(err.Error())
		}
		log.Info("(GRPC) Listening on :", GRPC_PORT)
		if err := grpcServer.Serve(grpcListener); err!=nil{
			log.Error(err.Error())
			panic(err.Error())
		}
	}()
	defer grpcServer.Stop()

	// connect to grpc servers
	authGrpcConn, err := client.GetAuthGrpcConnAndClient(AUTH_GRPC_URI)
	if err != nil{
		log.Error(err.Error())
		panic(err.Error())
	}
	defer authGrpcConn.Close()
	authGrpcClient := auth.NewAuthenticationClient(authGrpcConn)
	if _, err = authGrpcClient.HealthCheck(context.Background(), &auth.Empty{}); err != nil{
		log.Error(err.Error())
		panic(err.Error())
	}

	// Setup rest app
	restApp := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024, // 10 MB
	})
	defer restApp.ShutdownWithTimeout(10*time.Second)

	// setup top level routes
	fileRouter := restApp.Group("/file")
	fileRoute := route.FileRoute{
		Router: &fileRouter,
		MinioClient: minioClient,
		MongoFileMetadataCol: mongoFileMetadataCol,
		AuthGrpcClient: &authGrpcClient,
	}
	fileRoute.SetupRoutes()

	go func(){
		// start rest server
		log.Info("(REST) Listening on :", REST_PORT)
		restApp.Listen(fmt.Sprintf(":%v", REST_PORT))
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Info("Shutting down file service")
}
