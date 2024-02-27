package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
	"go.mongodb.org/mongo-driver/mongo"

	handler "github.com/AnuragProg/printit-microservices-monorepo/file/internal/api/handlers"
	mid "github.com/AnuragProg/printit-microservices-monorepo/file/internal/middleware"
	auth "github.com/AnuragProg/printit-microservices-monorepo/file/proto_gen/authentication"
)


type FileRoute struct{
	Router *fiber.Router
	MinioClient *minio.Client
	MongoFileMetadataCol *mongo.Collection
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
		handler.GetUploadFileHandler(
			fr.MinioClient,
			fr.MongoFileMetadataCol,
		),
	)
}
