package main

import (
	"github.com/gofiber/fiber/v2"
	route "github.com/AnuragProg/printit-microservices-monorepo/file/internal/api/routes"
)


func main(){
	app := fiber.New()

	fileRouter := app.Group("/file")
	fileRoute := route.FileRoute{
		Router: &fileRouter,
		AuthGrpcClient: nil,
	}
	fileRoute.SetupRoutes()


	app.Listen(":3000")
}
