package handlers

import (

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)


func UploadFile(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil{
		return err
	}

	file := form.File["file"][0]

	log.Info(file.Header.Get("content-type"), " received")

	return c.SendString("It's ok")
}
