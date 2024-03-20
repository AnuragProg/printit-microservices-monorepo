package middleware

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	consts "github.com/AnuragProg/printit-microservices-monorepo/internal/constant"

	auth "github.com/AnuragProg/printit-microservices-monorepo/proto_gen/authentication"
	"github.com/AnuragProg/printit-microservices-monorepo/proto_gen/shop"
)


func GetShopOwnershipMiddleware(shopGrpcClient *shop.ShopClient) fiber.Handler{
	return func(c *fiber.Ctx) error {
		userInfo := c.Locals(consts.USER_INFO_LOCAL).(*auth.User)
		shopId := c.Params("shopId")
		shopInfo, err := (*shopGrpcClient).GetShopByShopAndShopkeeperId(
			context.Background(),
			&shop.ShopAndShopkeeperId{
				ShopId: shopId,
				ShopkeeperId: userInfo.GetXId(),
			},
		)
		if err != nil{
			res, _ := status.FromError(err)
			if res.Code() == codes.NotFound {
				return fiber.NewError(fiber.StatusNotFound, "shop not found, do you own this shop?")
			}
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		c.Locals(consts.SHOP_INFO_LOCAL, shopInfo)
		return c.Next()
	}
}
