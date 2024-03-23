package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc/credentials/insecure"

	grpc "google.golang.org/grpc"
	auth "github.com/AnuragProg/printit-microservices-monorepo/proto_gen/authentication"
	file "github.com/AnuragProg/printit-microservices-monorepo/proto_gen/file"
	shop "github.com/AnuragProg/printit-microservices-monorepo/proto_gen/shop"
	price "github.com/AnuragProg/printit-microservices-monorepo/proto_gen/price"
	pb "github.com/AnuragProg/printit-microservices-monorepo/proto_gen/order"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"

	"github.com/AnuragProg/printit-microservices-monorepo/internal/api/route"
	"github.com/AnuragProg/printit-microservices-monorepo/internal/client"
	"github.com/AnuragProg/printit-microservices-monorepo/internal/service"
	"github.com/AnuragProg/printit-microservices-monorepo/pkg/util"
)


var (
	REST_PORT = util.GetenvOrDefault("REST_PORT", "3004")
	GRPC_PORT = util.GetenvOrDefault("GRPC_PORT", "50055")

	MONGO_URI = util.GetenvOrDefault("MONGO_URI", "mongodb://localhost:27017")
	AUTH_GRPC_URI = util.GetenvOrDefault("AUTH_GRPC_URI", "localhost:50051")
	FILE_GRPC_URI = util.GetenvOrDefault("FILE_GRPC_URI", "localhost:50052")
	SHOP_GRPC_URI = util.GetenvOrDefault("SHOP_GRPC_URI", "localhost:50053")
	PRICE_GRPC_URI = util.GetenvOrDefault("PRICE_GRPC_URI", "localhost:50054")

	KAFKA_BROKER = util.GetenvOrDefault("KAFKA_BROKER", "localhost:9092")
)



func main(){

	mongoDB, err := client.GetMongoDB(MONGO_URI)
	if err != nil{
		panic(err.Error())
	}

	grpcApp := grpc.NewServer()
	pb.RegisterOrderServer(grpcApp, service.NewOrderService())
	go func(){
		grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%v", GRPC_PORT))
		if err != nil{
			log.Error(err.Error())
			panic(err.Error())
		}
		log.Info("(GRPC) Listening on :", GRPC_PORT)
		if err := grpcApp.Serve(grpcListener); err!=nil{
			log.Error(err.Error())
			panic(err.Error())
		}
	}()
	defer grpcApp.Stop()


	/* connect to grpc servers */

	// auth grpc
	authGrpcConn, err := grpc.Dial(AUTH_GRPC_URI, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil{
		panic(err.Error())
	}
	log.Info("Connected to Auth GRPC server...")
	defer authGrpcConn.Close()
	authGrpcClient := auth.NewAuthenticationClient(authGrpcConn)
	if _, err := authGrpcClient.HealthCheck(context.Background(), &auth.Empty{}); err != nil{
		panic(err.Error())
	}

	// file grpc
	fileGrpcConn, err := grpc.Dial(FILE_GRPC_URI, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil{
		panic(err.Error())
	}
	log.Info("Connected to File GRPC server...")
	defer fileGrpcConn.Close()
	fileGrpcClient := file.NewFileClient(fileGrpcConn)
	if _, err := fileGrpcClient.HealthCheck(context.Background(), &file.Empty{}); err != nil{
		panic(err.Error())
	}

	// shop grpc
	shopGrpcConn, err := grpc.Dial(SHOP_GRPC_URI, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil{
		panic(err.Error())
	}
	log.Info("Connected to Shop GRPC server...")
	defer shopGrpcConn.Close()
	shopGrpcClient := shop.NewShopClient(shopGrpcConn)
	if _, err := shopGrpcClient.HealthCheck(context.Background(), &shop.Empty{}); err != nil{
		panic(err.Error())
	}

	// price grpc
	priceGrpcConn, err := grpc.Dial(PRICE_GRPC_URI, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil{
		panic(err.Error())
	}
	log.Info("Connected to Price GRPC server...")
	defer priceGrpcConn.Close()
	priceGrpcClient := price.NewPriceClient(priceGrpcConn)
	if _, err := priceGrpcClient.HealthCheck(context.Background(), &price.Empty{}); err != nil{
		panic(err.Error())
	}

	// kafka client
	orderEventEmitter, err := client.NewOrderEventEmitter([]string{KAFKA_BROKER})
	if err != nil{
		panic(err.Error())
	}

	restApp := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024, // 10 MB
	})
	defer restApp.ShutdownWithTimeout(10*time.Second)

	_ = route.NewOrderRoute(
		restApp.Group("/order"),
		orderEventEmitter,
		mongoDB.OrderCol,
		&authGrpcClient,
		&fileGrpcClient,
		&shopGrpcClient,
		&priceGrpcClient,
	)
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
