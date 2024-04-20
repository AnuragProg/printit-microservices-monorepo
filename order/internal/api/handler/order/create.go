package order

import (
	"context"

	"github.com/AnuragProg/printit-microservices-monorepo/internal/client"
	consts "github.com/AnuragProg/printit-microservices-monorepo/internal/constant"
	"github.com/AnuragProg/printit-microservices-monorepo/internal/data"

	auth "github.com/AnuragProg/printit-microservices-monorepo/proto_gen/authentication"
	file "github.com/AnuragProg/printit-microservices-monorepo/proto_gen/file"
	price "github.com/AnuragProg/printit-microservices-monorepo/proto_gen/price"
	shop "github.com/AnuragProg/printit-microservices-monorepo/proto_gen/shop"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"go.mongodb.org/mongo-driver/mongo"
)

type CreateOrderRequest struct {
	FileId		string `json:"file_id"`
	PriceId	 	string `json:"price_id"`
}


func GetCreateOrderHandler(
	orderCol *mongo.Collection,

	orderEventEmitter *client.OrderEventEmitter,

	fileGrpcClient *file.FileClient,
	shopGrpcClient *shop.ShopClient,
	priceGrpcClient *price.PriceClient,
) fiber.Handler{
	return func(c *fiber.Ctx) error {

		shopId := c.Params("shopId")

		userInfo, ok := c.Locals(consts.USER_INFO_LOCAL).(*auth.User)
		if !ok {
			return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
		}

		var createOrderRequest CreateOrderRequest
		if err := c.BodyParser(&createOrderRequest); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
		}

		// TODO: here authentication layer needs to be introduced so that only the user himself and the shopkeeper can see the file or fetch metadata
		fileInfo, err := (*fileGrpcClient).GetFileMetadataById(context.Background(), &file.FileId{ XId: createOrderRequest.FileId })
		if err != nil{
			return fiber.NewError(fiber.StatusBadRequest, "invalid file id")
		}
		shopInfo, err := (*shopGrpcClient).GetShopById(context.Background(), &shop.ShopId{ XId: shopId })
		if err != nil{
			return fiber.NewError(fiber.StatusBadRequest, "invalid shop id")
		}
		priceInfo, err := (*priceGrpcClient).GetPriceInfoByPriceIdAndShopId(
			context.Background(),
			&price.PriceIdAndShopId{
				PriceId: createOrderRequest.PriceId,
				ShopId: shopId,
			},
		)
		if err != nil{
			return fiber.NewError(fiber.StatusBadRequest, "invalid price id")
		}

		// store order on database
		orderStatus := data.ORDER_PLACED
		order, err := data.CreateOrder(
			fileInfo.GetId(),
			shopInfo.GetXId(),
			priceInfo.GetXId(),
			userInfo.GetXId(),
			string(orderStatus),
		)
		if err != nil{
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		if _, err := orderCol.InsertOne(context.Background(), order); err != nil{
			log.Error(err.Error())
			return fiber.ErrInternalServerError
		}

		// emit order event on kafka
		orderEvent := client.OrderEvent{
			ShopId: order.ShopId,
			Status: orderStatus,
			UpdatedOnOrBeforeEpochMS: order.UpdatedAt.UnixMilli(),
		}
		if err := orderEventEmitter.EmitOrderEvent(&orderEvent); err != nil {
			log.Error(err.Error()) // will just show the error and not halt the order process as such
		}

		// respond to the user
		c.JSON(struct{
			Message	string `json:"message"`
			Order		data.Order `json:"order"`
		}{
			Message: "order placed successfully",
			Order: *order,
		})
		return nil
	}
}

