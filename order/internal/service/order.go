package service

import (
	"context"
	"time"

	"github.com/AnuragProg/printit-microservices-monorepo/internal/data"
	pb "github.com/AnuragProg/printit-microservices-monorepo/proto_gen/order"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


type OrderService struct {
	pb.UnimplementedOrderServer
	orderCol *mongo.Collection
}

func NewOrderService(orderCol *mongo.Collection) *OrderService{
	return &OrderService{
		orderCol: orderCol,
	}
}

func (os *OrderService) HealthCheck(context.Context, *pb.Empty) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}
func (os *OrderService) GetShopTraffic(ctx context.Context, req *pb.GetShopTrafficRequest) (*pb.GetShopTrafficResponse, error) {
	timestamp, err := time.Parse(time.RFC3339, req.GetUpdatedOnOrBefore())
	if err != nil{
		return nil, status.Error(codes.InvalidArgument, "timestamp should be in rfc3339 format")
	}
	if _, err := primitive.ObjectIDFromHex(req.GetShopId()); err != nil{ // because we are using mongo for shop as well
		return nil, status.Error(codes.InvalidArgument, "invalid shop id")
	}
	filter := bson.M{
		"shop_id": req.GetShopId(),
		"status": bson.M{ "$in": bson.A{data.ORDER_PROCESSING} }, // all the statuses that indicates traffic on the shop
		"updated_at": bson.M{"$lte": timestamp},
	}
	traffic, err := os.orderCol.CountDocuments(context.Background(), filter)
	if err != nil{
		return nil, status.Error(codes.Internal, "unable to fetch traffic from mongo")
	}

	return &pb.GetShopTrafficResponse{
		ShopId: req.GetShopId(),
		Traffic: uint32(traffic),
	}, nil
}
