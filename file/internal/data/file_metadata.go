package data

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)


type FileMetadata struct{
	Id				primitive.ObjectID `bson:"_id" json:"id"`
	UserId		string `bson:"user_id" json:"user_id"`
	FileName		string `bson:"file_name" json:"file_name"`
	BucketName	string `bson:"bucket_name" json:"bucket_name"`
	Size			uint32 `bson:"size" json:"size"`
	ContentType string `bson:"content_type" json:"content_type"`
	CreatedAt	time.Time `bson:"created_at" json:"-"`
	UpdatedAt	time.Time `bson:"updated_at" json:"-"`
}
