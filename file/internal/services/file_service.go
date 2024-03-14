package file_service

import (
	"time"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/AnuragProg/printit-microservices-monorepo/file/internal/data"
	pb "github.com/AnuragProg/printit-microservices-monorepo/file/proto_gen/file"
)

type FileService struct {
	pb.UnimplementedFileServer
	mongoFileMetadataCol *mongo.Collection
}

func NewFileService(mongoFileMetadataCol *mongo.Collection) *FileService {
	return &FileService{
		mongoFileMetadataCol: mongoFileMetadataCol,
	}
}

func (fs *FileService) HealthCheck(context.Context, *pb.Empty) (*pb.Empty, error){
	return &pb.Empty{}, nil
}

func (fs *FileService) GetFileMetadataById(ctx context.Context, fileId *pb.FileId) (*pb.FileMetadata, error) {
	id, err := primitive.ObjectIDFromHex(fileId.GetXId())
	if err != nil{
		return nil, status.Errorf(codes.InvalidArgument, "invalid id")
	}
	metadataRes := fs.mongoFileMetadataCol.FindOne(context.Background(), bson.M{"_id": id})
	if err := metadataRes.Err(); err != nil{
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	metadata := data.FileMetadata{}
	if err := metadataRes.Decode(&metadata); err != nil{
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	metadataResponse := &pb.FileMetadata{
		Id: metadata.Id.Hex(),
		UserId: metadata.UserId,
		FileName: metadata.FileName,
		BucketName: metadata.BucketName,
		Size: metadata.Size,
		ContentType: metadata.ContentType,
		CreatedAt: metadata.CreatedAt.Format(time.RFC3339),
		UpdatedAt: metadata.UpdatedAt.Format(time.RFC3339),
	}
	return metadataResponse, nil
}
