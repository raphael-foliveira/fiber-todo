package common

import "github.com/gofiber/fiber/v2"

func StatusCheck(c *fiber.Ctx) error {
	return c.Send([]byte("ok"))
}
