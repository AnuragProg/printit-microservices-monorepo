package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "github.com/AnuragProg/printit-microservices-monorepo/proto_gen/order"
	grpc "google.golang.org/grpc"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"

	"github.com/AnuragProg/printit-microservices-monorepo/internal/api/route"
	"github.com/AnuragProg/printit-microservices-monorepo/internal/service"
	util "github.com/AnuragProg/printit-microservices-monorepo/pkg/util"
)


var (
	REST_PORT = util.GetenvOrDefault("REST_PORT", "3004")
	GRPC_PORT = util.GetenvOrDefault("GRPC_PORT", "50054")
)



func main(){

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

	restApp := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024, // 10 MB
	})
	defer restApp.ShutdownWithTimeout(10*time.Second)

	fileRouter := restApp.Group("/order")
	fileRoute := route.FileRoute{
		Router: &fileRouter,
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
