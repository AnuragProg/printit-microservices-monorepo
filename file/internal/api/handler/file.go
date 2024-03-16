package handler

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"

	"github.com/minio/minio-go/v7"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	consts "github.com/AnuragProg/printit-microservices-monorepo/internal/constant"
	data "github.com/AnuragProg/printit-microservices-monorepo/internal/data"
	auth "github.com/AnuragProg/printit-microservices-monorepo/proto_gen/authentication"
)


func GetUploadFileHandler(
	minioClient *minio.Client,
	mongoFileMetadataCol *mongo.Collection,
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
		defer buff.Close()

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
		_, err = mongoFileMetadataCol.InsertOne(context.Background(), metadata)
		if err != nil{
			log.Error(err)
			return fiber.ErrInternalServerError
		}

		// upload file
		_, err = minioClient.PutObject(context.Background(), metadata.BucketName, metadata.Id.Hex(), buff, int64(metadata.Size), minio.PutObjectOptions{
			ContentType: metadata.ContentType,
		})
		if err != nil{
			mongoFileMetadataCol.DeleteOne(context.Background(), bson.M{
				"_id": metadata.Id,
			})
			return err
		}

		// respond
		return c.JSON(struct{
			Message string `json:"message"`
			FileInfo data.FileMetadata `json:"file_info"`
		}{
			Message: "file uploaded successfully",
			FileInfo: metadata,
		})
	}
}


func GetDownloadFileHandler(
	minioClient *minio.Client,
	mongoFileMetadataCol *mongo.Collection,
) fiber.Handler{
	return func(c *fiber.Ctx) error {

		// get file id
		fileId, err := primitive.ObjectIDFromHex(c.Params("id"))
		if err != nil{
			return fiber.NewError(fiber.StatusBadRequest, "invalid id param")
		}


		// get file metadata
		metadataRes := mongoFileMetadataCol.FindOne(context.Background(), bson.M{"_id": fileId})
		metadata := data.FileMetadata{}
		if err := metadataRes.Decode(&metadata); err != nil{
			return err
		}

		// get file
		obj, err := minioClient.GetObject(context.Background(), consts.FILE_BUCKET, metadata.Id.Hex(), minio.GetObjectOptions{})
		if err != nil{
			return fiber.NewError(fiber.StatusNotFound, "file not found")
		}
		defer obj.Close()

		// set response headers
		c.Set("content-type", metadata.ContentType)
		c.Set("content-disposition", "attachment; filename="+metadata.FileName)

		// write file to response
		if _, err = io.Copy(c.Response().BodyWriter(), obj); err != nil{
			return fiber.ErrInternalServerError
		}

		return nil
	}
}
