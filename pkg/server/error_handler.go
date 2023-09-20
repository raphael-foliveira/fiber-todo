package server

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

// errorHandler is the default error handler for the application
func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	err = c.Status(code).JSON(fiber.Map{
		"error":       err.Error(),
		"status_code": code,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "internal server error",
		})
	}
	return nil
}
