package handlers

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	//"github.com/gofiber/fiber/v2/log"

	//"github.com/AnuragProg/printit-microservices-monorepo/file/internal/data"
)


func GetUploadFileHandler(mongoClient *mongo.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		_, err := c.MultipartForm()
		if err != nil{
			return err
		}

		//file := form.File["file"][0]

		// save metadata to minio in a transaction
		//session, err := mongoClient.StartSession()
		//if err != nil{
		//	return fiber.ErrInternalServerError
		//}
		// TODO upload file to minio
		// commit the transaction

		return c.SendString("It's ok")

	}
}
