package route

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"

	mid "github.com/AnuragProg/printit-microservices-monorepo/internal/middleware"
	"github.com/AnuragProg/printit-microservices-monorepo/internal/api/handler/order"
	auth "github.com/AnuragProg/printit-microservices-monorepo/proto_gen/authentication"
)


type OrderRoute struct {
	Router *fiber.Router
	OrderCol *mongo.Collection
	AuthGrpcClient *auth.AuthenticationClient
}


func New(router fiber.Router, orderCol *mongo.Collection, authGrpcClient *auth.AuthenticationClient) *OrderRoute {
	orderRoute := OrderRoute{
		Router: &router,
		OrderCol: orderCol,
		AuthGrpcClient: authGrpcClient,
	}
	orderRoute.SetupRoutes()
	return &orderRoute
}

func (or *OrderRoute)SetupRoutes() {

	// create order (customer) POST
	(*or.Router).Post(
		"/shop/:shopId/orders",
		mid.GetAuthMiddleware(or.AuthGrpcClient, auth.UserType_CUSTOMER),
		order.GetCreateOrderHandler(or.OrderCol),
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
		"/shop/:shopId/orders/",
		mid.GetAuthMiddleware(or.AuthGrpcClient, auth.UserType_SHOPKEEPER),
		order.GetListShopkeeperOrdersHandler(or.OrderCol),
	)
}

