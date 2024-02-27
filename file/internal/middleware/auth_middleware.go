package middleware

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"

	consts "github.com/AnuragProg/printit-microservices-monorepo/file/internal/constants"
	auth "github.com/AnuragProg/printit-microservices-monorepo/file/proto_gen/authentication"
)


func GetAuthMiddleware(auth_grpc_client *auth.AuthenticationClient) fiber.Handler {
	return func (c *fiber.Ctx) error {

		// parse token
		auth_header := c.Get("authorization")
		auth_header_slice := strings.Split(auth_header, " ")
		if len(auth_header_slice) != 2{
			return fiber.NewError(fiber.StatusBadRequest, "please provide auth token")
		}
		token := auth_header_slice[1]
		log.Info("Received token = ", token)

		// verifying token
		user, err := (*auth_grpc_client).VerifyToken(context.Background(), &auth.Token{ Token: token } )
		if err != nil{
			log.Error(err.Error())
			return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
		}

		// pass user info forward for handlers
		c.Locals(consts.USER_LOCAL, user)

		return c.Next()
	}
}





