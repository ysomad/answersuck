package v1

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/answersuck/vault/internal/domain"
)

const (
	accountIdKey = "aid"
	sessionIdKey = "sid"
	audienceKey  = "audience"
	deviceKey    = "device"
)

var errDeviceNotFound = errors.New("device not found in context")

// getAccountId returns account id from context
func getAccountId(c *gin.Context) (string, error) {
	aid := c.GetString(accountIdKey)

	_, err := uuid.Parse(aid)
	if err != nil {
		return "", domain.ErrAccountContextNotFound
	}

	return aid, nil
}

// getSessionId returns session id from context
func getSessionId(c *gin.Context) string { return c.GetString(sessionIdKey) }

// getAudience returns current getAudience from context
func getAudience(c *gin.Context) string { return c.GetString(audienceKey) }

func getDevice(c *gin.Context) (domain.Device, error) {
	d, exists := c.Get(deviceKey)
	if !exists {
		return domain.Device{}, errDeviceNotFound
	}

	return d.(domain.Device), nil
}
