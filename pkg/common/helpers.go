package common

import (
	"reflect"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func Contains[T any](arr []T, str T) bool {
	for _, a := range arr {
		if reflect.DeepEqual(a, str) {
			return true
		}
	}
	return false
}

func ParseIdFromParams(c *fiber.Ctx) (int, error) {
	intId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return 0, err
	}
	return intId, nil
}
