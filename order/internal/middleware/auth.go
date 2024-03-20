package middleware

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	auth "github.com/AnuragProg/printit-microservices-monorepo/proto_gen/authentication"
	consts "github.com/AnuragProg/printit-microservices-monorepo/internal/constant"
)

func GetAuthMiddleware(authGrpcClient *auth.AuthenticationClient, userType ...auth.UserType) fiber.Handler {
	return func(c *fiber.Ctx) error {

		authHeader := c.Get("authorization")
		authHeaderSlice := strings.Split(authHeader, " ")
		if len(authHeaderSlice) != 2 {
			return fiber.NewError(fiber.StatusUnauthorized, "bearer token not found")
		}
		token := authHeaderSlice[1]
		user, err := (*authGrpcClient).VerifyToken(context.Background(), &auth.Token{ Token: token })
		if err != nil{
			res, _ := status.FromError(err)
			if res.Code() == codes.Unauthenticated {
				return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
			}
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if len(userType) > 0 && user.GetUserType() != userType[0]{
			return fiber.NewError(fiber.StatusUnauthorized, "not authorized")
		}

		c.Locals(consts.USER_INFO_LOCAL, user)
		return c.Next()
	}
}
