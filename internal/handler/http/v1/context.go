package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/answersuck/vault/internal/domain/account"
	"github.com/answersuck/vault/internal/domain/session"
)

const (
	accountIdKey = "accountId"
	sessionIdKey = "sessionId"
	audienceKey  = "audience"
	deviceKey    = "device"
)

// getAccountId returns account id from context
func getAccountId(c *gin.Context) (string, error) {
	accountId := c.GetString(accountIdKey)

	_, err := uuid.Parse(accountId)
	if err != nil {
		return "", account.ErrContextNotFound
	}

	return accountId, nil
}

// getSessionId returns session id from context
func getSessionId(c *gin.Context) string { return c.GetString(sessionIdKey) }

func getDevice(c *gin.Context) (session.Device, error) {
	v, exists := c.Get(deviceKey)
	if !exists {
		return session.Device{}, session.ErrDeviceContextNotFound
	}

	d, ok := v.(session.Device)
	if !ok {
		return session.Device{}, session.ErrDeviceContextNotFound
	}

	return d, nil
}
