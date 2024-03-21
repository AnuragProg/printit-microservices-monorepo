package order

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/AnuragProg/printit-microservices-monorepo/internal/data"
	"github.com/AnuragProg/printit-microservices-monorepo/proto_gen/authentication"
	"github.com/AnuragProg/printit-microservices-monorepo/proto_gen/shop"

	consts "github.com/AnuragProg/printit-microservices-monorepo/internal/constant"
)


func GetListCustomerOrdersHandler(
	orderCol *mongo.Collection,
) fiber.Handler{
	return func (c *fiber.Ctx) error {
		userInfo, ok := c.Locals(consts.USER_INFO_LOCAL).(*authentication.User)
		if !ok {
			return fiber.NewError(fiber.StatusInternalServerError, "interval server error")
		}

		page := c.QueryInt("page", 1) - 1
		pageSize := c.QueryInt("pageSize", 10)

		if page < 0 { page = 0 }

		userOrderFilter := bson.M{
			"customer_id": userInfo.GetXId(),
		}
		userOrderOptions := options.Find()
		userOrderOptions.SetLimit(int64(pageSize))
		userOrderOptions.SetSkip(int64(page*pageSize))
		cursor, err := orderCol.Find(
			context.Background(),
			userOrderFilter,
			userOrderOptions,
		)
		if err != nil{
			log.Error("got error: " + err.Error())
			return fiber.ErrInternalServerError
		}

		var userOrders []data.Order
		if err := cursor.All(context.Background(), &userOrders); err != nil{
			log.Error("got error: " + err.Error())
			return fiber.ErrInternalServerError
		}

		c.JSON(
			struct{
				Page int `json:"page"`
				PageSize int `json:"pageSize"`
				Orders []data.Order `json:"orders"`
			}{
				Page: page+1,
				PageSize: pageSize,
				Orders: userOrders,
			},
		)

		return nil
	}
}

func GetListShopkeeperOrdersHandler(orderCol *mongo.Collection) fiber.Handler{
	return func (c *fiber.Ctx) error {
		shopId := c.Params("shopId")

		page := c.QueryInt("page", 1) - 1
		pageSize := c.QueryInt("pageSize", 10)

		if page < 0 { page = 0 }

		shopInfo := c.Locals(consts.SHOP_INFO_LOCAL).(*shop.ShopInfo)
		if shopInfo == nil {
			// possibly shop ownership hasn't been verified
			return fiber.NewError(fiber.StatusInternalServerError, "shop verification unsuccessful")
		}

		shopOrderFilter := bson.M{
			"shop_id": shopId,
		}
		shopOrderOption := options.Find()
		shopOrderOption.SetLimit(int64(pageSize))
		shopOrderOption.SetSkip(int64(page*pageSize))
		cursor, err := orderCol.Find(context.Background(), shopOrderFilter, shopOrderOption)
		if err != nil{
			log.Error("got error: " + err.Error())
			return fiber.ErrInternalServerError
		}

		var shopOrders []data.Order
		if err := cursor.All(context.Background(), &shopOrders); err != nil{
			log.Error("got error: " + err.Error())
			return fiber.ErrInternalServerError
		}

		c.JSON(
			struct{
				Page int `json:"page"`
				PageSize int `json:"pageSize"`
				Orders []data.Order `json:"orders"`
			}{
				Page: page+1,
				PageSize: pageSize,
				Orders: shopOrders,
			},
		)

		return nil
	}
}
