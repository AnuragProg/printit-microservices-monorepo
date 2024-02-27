package data

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)


type FileMetadata struct{
	Id primitive.ObjectID `bson:"_id"`
	UserId string `bson:"user_id"`
	FileId string `bson:"file_id"`
	FileName string `bson:"file_name"`
	BucketName string `bson:"bucket_name"`
	Size uint32 `bson:"size"`
	ContentType string `bson:"content_type"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}
