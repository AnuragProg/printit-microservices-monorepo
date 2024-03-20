package route

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"

	mid "github.com/AnuragProg/printit-microservices-monorepo/internal/middleware"
	"github.com/AnuragProg/printit-microservices-monorepo/internal/api/handler/order"
	auth "github.com/AnuragProg/printit-microservices-monorepo/proto_gen/authentication"
	"github.com/AnuragProg/printit-microservices-monorepo/proto_gen/file"
	"github.com/AnuragProg/printit-microservices-monorepo/proto_gen/shop"
	"github.com/AnuragProg/printit-microservices-monorepo/proto_gen/price"
)


type OrderRoute struct {
	Router *fiber.Router
	OrderCol *mongo.Collection
	AuthGrpcClient *auth.AuthenticationClient
	FileGrpcClient *file.FileClient
	ShopGrpcClient *shop.ShopClient
	PriceGrpcClient *price.PriceClient
}


func New(
	router fiber.Router,
	orderCol *mongo.Collection,
	authGrpcClient *auth.AuthenticationClient,
	fileGrpcClient *file.FileClient,
	shopGrpcClient *shop.ShopClient,
	priceGrpcClient *price.PriceClient,
) *OrderRoute {
	orderRoute := OrderRoute{
		Router: &router,
		OrderCol: orderCol,
		AuthGrpcClient: authGrpcClient,
		FileGrpcClient: fileGrpcClient,
		ShopGrpcClient: shopGrpcClient,
		PriceGrpcClient: priceGrpcClient,
	}
	orderRoute.SetupRoutes()
	return &orderRoute
}

func (or *OrderRoute)SetupRoutes() {

	// create order (customer) POST
	(*or.Router).Post(
		"/shop/:shopId/orders",
		mid.GetAuthMiddleware(or.AuthGrpcClient, auth.UserType_CUSTOMER),
		order.GetCreateOrderHandler(
			or.OrderCol,
			or.FileGrpcClient,
			or.ShopGrpcClient,
			or.PriceGrpcClient,
		),
	)

	// update order status PATCH
		// cancel order (customer) cancelled
		// accept order (shopkeeper) accepted -> processing
		// reject order (shopkeeper) rejected
		// complete order (shopkeeper) completed
	(*or.Router).Patch(
		"/shop/:shopId/orders/:orderId",
		order.GetOrderActionHandler(or.OrderCol),
	)

	// list my orders (customer) GET
	(*or.Router).Get(
		"/",
		mid.GetAuthMiddleware(or.AuthGrpcClient, auth.UserType_CUSTOMER),
		order.GetListCustomerOrdersHandler(or.OrderCol),
	)

	// list my orders (shopkeeper) GET
	(*or.Router).Get(
		"/shop/:shopId/orders",
		mid.GetAuthMiddleware(or.AuthGrpcClient, auth.UserType_SHOPKEEPER),
		mid.GetShopOwnershipMiddleware(or.ShopGrpcClient),
		order.GetListShopkeeperOrdersHandler(or.OrderCol),
	)
}

