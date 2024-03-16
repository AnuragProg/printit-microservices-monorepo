package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)


func GetSingleFileMiddleware() fiber.Handler {
	return func (c *fiber.Ctx) error {
		form, err := c.MultipartForm()
		if err != nil{
			log.Info(err.Error())
			return fiber.NewError(fiber.StatusBadRequest, "invalid request")
		}

		files := form.File["file"]
		if len(files) != 1{
			return fiber.NewError(fiber.StatusBadRequest, "exactly one file is required")
		}
		return c.Next()
	}
}


func GetFileContentTypeCheckerMiddleware(acceptableFileTypes map[string]interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		form, err := c.MultipartForm()
		if err != nil{
			return fiber.NewError(fiber.StatusBadRequest, "invalid request")
		}

		files := form.File["file"]
		for _, file := range files{
			fileType := file.Header.Get("content-type")
			if _, ok := acceptableFileTypes[fileType]; !ok{
				return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("%v type not allowed", fileType))
			}
		}
		return c.Next()
	}
}
