package routes

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"

	mid "github.com/AnuragProg/printit-microservices-monorepo/file/internal/middleware"
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

	(*router).Post(
		"/upload",
		mid.GetAuthMiddleware(fr.AuthGrpcClient),
		mid.GetSingleFileMiddleware(),
		mid.GetFileContentTypeCheckerMiddleware(
			map[string]interface{}{
				"application/pdf": struct{}{},
			},
		),
		handler.UploadFile,
	)
}
