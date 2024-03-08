package file_service

import (
	"context"

	pb "github.com/AnuragProg/printit-microservices-monorepo/file/proto_gen/file"
)

type FileService struct {
	pb.UnimplementedFileServer
}

func (fs *FileService) HealthCheck(context.Context, *pb.Empty) (*pb.Empty, error){
	return &pb.Empty{}, nil
}

