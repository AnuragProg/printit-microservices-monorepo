package routes

import (
	"github.com/gofiber/fiber/v2"

	handler "github.com/AnuragProg/printit-microservices-monorepo/file/internal/api/handlers"
	auth "github.com/AnuragProg/printit-microservices-monorepo/file/internal/proto_gen/authentication"
)


type FileRoute struct{
	Router *fiber.Router
	AuthGrpcClient *auth.AuthenticationClient
}

func (fr *FileRoute)SetupRoutes(){
	router := fr.Router
	(*router).Get("/foo", handler.Foo)
}
