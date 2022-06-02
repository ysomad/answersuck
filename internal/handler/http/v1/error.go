package v1

import (
	"github.com/gofiber/fiber/v2"
)

type detailedError[D string | map[string]string] struct {
	Error  string `json:"error"`
	Detail D      `json:"detail"`
}

func errorResp[D string | map[string]string](c *fiber.Ctx, status int,
	err error, detail D) error {

	return c.Status(status).JSON(detailedError[D]{
		Error:  err.Error(),
		Detail: detail,
	})
}
