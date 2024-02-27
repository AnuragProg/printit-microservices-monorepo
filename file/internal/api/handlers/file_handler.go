package handlers

import (
	"context"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"

	"github.com/minio/minio-go/v7"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	data "github.com/AnuragProg/printit-microservices-monorepo/file/internal/data"
	consts "github.com/AnuragProg/printit-microservices-monorepo/file/internal/constants"
	auth "github.com/AnuragProg/printit-microservices-monorepo/file/proto_gen/authentication"
)


func GetUploadFileHandler(
	minioClient *minio.Client,
	mongoClient *mongo.Client,
	fileMetadataCol *mongo.Collection,
) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// get user info
		user, ok := c.Locals(consts.USER_LOCAL).(*auth.User)
		if !ok {
			return errors.New("user not passed by auth middlewares")
		}

		// parse form data
		form, err := c.MultipartForm()
		if err != nil{
			return err
		}

		file := form.File["file"][0]
		buff, err := file.Open()
		if err != nil{
			return fiber.NewError(fiber.StatusBadRequest, "unable to open file")
		}

		// Create file metadata object
		metadata := data.FileMetadata{
			Id: primitive.NewObjectID(),
			UserId: user.XId,
			FileName: file.Filename,
			BucketName: consts.FILE_BUCKET,
			Size: uint32(file.Size),
			ContentType: file.Header.Get("content-type"),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}

		// save metadata
		_, err = fileMetadataCol.InsertOne(context.Background(), metadata)
		if err != nil{
			log.Error(err)
			return fiber.ErrInternalServerError
		}

		// upload file
		_, err = minioClient.PutObject(context.Background(), metadata.BucketName, metadata.Id.Hex(), buff, int64(metadata.Size), minio.PutObjectOptions{
			ContentType: metadata.ContentType,
		})
		if err != nil{
			fileMetadataCol.DeleteOne(context.Background(), bson.M{
				"_id": metadata.Id,
			})
			return err
		}

		// respond
		return c.JSON(map[string]string{
			"message": "file uploaded successfully",
		})
	}
}
