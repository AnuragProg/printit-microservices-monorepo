package service


import (
	"context"

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
	return &pb.Empty{}, nil
}
