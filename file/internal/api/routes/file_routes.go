package routes

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"

	handler "github.com/AnuragProg/printit-microservices-monorepo/file/internal/api/handlers"
	auth "github.com/AnuragProg/printit-microservices-monorepo/file/proto_gen/authentication"
)


type FileRoute struct{
	MongoClient *mongo.Client
	Router *fiber.Router
	AuthGrpcClient *auth.AuthenticationClient
}

func (fr *FileRoute)SetupRoutes(){
	router := fr.Router
	(*router).Get("/foo", handler.Foo)
}
