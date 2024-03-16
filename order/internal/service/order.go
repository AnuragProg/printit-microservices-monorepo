package service


import (
	"context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	pb "github.com/AnuragProg/printit-microservices-monorepo/proto_gen/order"
)


type OrderService struct {
	pb.UnimplementedOrderServer
}

func NewOrderService() *OrderService{
	return &OrderService{
	}
}

func (os *OrderService) HealthCheck(context.Context, *pb.Empty) (*pb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HealthCheck not implemented")
}
