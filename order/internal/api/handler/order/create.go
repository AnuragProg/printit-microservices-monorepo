package order

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)



func GetCreateOrderHandler(orderCol *mongo.Collection) fiber.Handler{
	return func(c *fiber.Ctx) error {
		return nil
	}
}

