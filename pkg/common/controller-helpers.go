package common

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func ParseIdFromParams(c *fiber.Ctx) (int, error) {
	intId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return 0, err
	}
	return intId, nil
}
