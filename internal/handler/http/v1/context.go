package v1

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/answersuck/vault/internal/domain"
)

const (
	accountIdKey = "accountId"
	sessionIdKey = "sessionId"
	audienceKey  = "audience"
	deviceKey    = "device"
)

var errDeviceNotFound = errors.New("device not found in context")

// getAccountId returns account id from context
func getAccountId(c *gin.Context) (string, error) {
	accountId := c.GetString(accountIdKey)

	_, err := uuid.Parse(accountId)
	if err != nil {
		return "", domain.ErrAccountContextNotFound
	}

	return accountId, nil
}

// getSessionId returns session id from context
func getSessionId(c *gin.Context) string { return c.GetString(sessionIdKey) }

func getDevice(c *gin.Context) (domain.Device, error) {
	v, exists := c.Get(deviceKey)
	if !exists {
		return domain.Device{}, errDeviceNotFound
	}

	d, ok := v.(domain.Device)
	if !ok {
		return domain.Device{}, errDeviceNotFound
	}

	return d, nil
}
