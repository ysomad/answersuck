package v1

import (
	"errors"

	"github.com/answersuck/vault/internal/domain/session"
	"github.com/gofiber/fiber/v2"
)

const (
	accountIdKey = "accoutId"
	sessionIdKey = "sessionId"
	deviceKey    = "device"
)

var (
	errAccountIdNotFound = errors.New("account id not found in context")
	errSessionIdNotFound = errors.New("session id not found in context")
	errDeviceNotFound    = errors.New("device not found in context")
)

func getAccountId(c *fiber.Ctx) (string, error) {
	v := c.Locals(accountIdKey)

	aid, ok := v.(string)
	if !ok || aid == "" {
		return "", errAccountIdNotFound
	}

	return aid, nil
}

func getSessionId(c *fiber.Ctx) (string, error) {
	v := c.Locals(sessionIdKey)

	sid, ok := v.(string)
	if !ok || sid == "" {
		return "", errSessionIdNotFound
	}

	return sid, nil
}

func getDevice(c *fiber.Ctx) (session.Device, error) {
	v := c.Locals(deviceKey)

	d, ok := v.(session.Device)
	if !ok {
		return session.Device{}, errDeviceNotFound
	}

	return d, nil
}
