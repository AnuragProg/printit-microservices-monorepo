package order

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)


func GetListCustomerOrdersHandler(orderCol *mongo.Collection) fiber.Handler{
	return func (c *fiber.Ctx) error {
		c.JSON(struct{
			Message string `json:"message"`
		}{
			Message: "its alright",
		})
		return nil
	}
}

func GetListShopkeeperOrdersHandler(orderCol *mongo.Collection) fiber.Handler{
	return func (c *fiber.Ctx) error {
		return nil
	}
}
