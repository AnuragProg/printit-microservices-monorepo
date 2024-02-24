package handlers

import "github.com/gofiber/fiber/v2"


func Foo(c *fiber.Ctx) error {
	return c.SendString("world")
}
