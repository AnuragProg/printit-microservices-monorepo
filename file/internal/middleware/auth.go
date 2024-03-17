package middleware

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	consts "github.com/AnuragProg/printit-microservices-monorepo/internal/constant"
	auth "github.com/AnuragProg/printit-microservices-monorepo/proto_gen/authentication"
)


func GetAuthMiddleware(authGrpcClient *auth.AuthenticationClient) fiber.Handler {
	return func (c *fiber.Ctx) error {

		// parse token
		authHeader := c.Get("authorization")
		authHeaderSlice := strings.Split(authHeader, " ")
		if len(authHeaderSlice) != 2{
			return fiber.NewError(fiber.StatusBadRequest, "please provide auth token")
		}
		token := authHeaderSlice[1]
		log.Info("Received token = ", token)

		// verifying token
		user, err := (*authGrpcClient).VerifyToken(context.Background(), &auth.Token{ Token: token } )
		if err != nil{
			log.Error(err.Error())
			res, _ := status.FromError(err)
			if res.Code() == codes.Unauthenticated{
				return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
			}
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		// pass user info forward for handlers
		c.Locals(consts.USER_LOCAL, user)

		return c.Next()
	}
}





