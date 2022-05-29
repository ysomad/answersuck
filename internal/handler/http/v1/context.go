package v1

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/answersuck/vault/internal/domain/session"
)

const (
	accountIdKey = "accountId"
	nicknameKey  = "nickname"
	sessionIdKey = "sessionId"
	audienceKey  = "audience"
	deviceKey    = "device"
)

var (
	errAccountIdNotFound = errors.New("account id not found in context")
	errNicknameNotFound  = errors.New("account nickname not found in context")
	errSessionIdNotFound = errors.New("session id not found in context")
)

// getAccountId returns account id from context
func getAccountId(c *gin.Context) (string, error) {
	accountId := c.GetString(accountIdKey)

	_, err := uuid.Parse(accountId)
	if err != nil {
		return "", errAccountIdNotFound
	}

	return accountId, nil
}

func getNickname(c *gin.Context) (string, error) {
	n := c.GetString(nicknameKey)
	if n == "" {
		return "", errNicknameNotFound
	}

	return n, nil
}

// getSessionId returns session id from context
func getSessionId(c *gin.Context) (string, error) {
	s := c.GetString(sessionIdKey)
	if s == "" {
		return "", errSessionIdNotFound
	}

	return s, nil
}

func getDevice(c *gin.Context) (session.Device, error) {
	v, exists := c.Get(deviceKey)
	if !exists {
		return session.Device{}, errSessionIdNotFound
	}

	d, ok := v.(session.Device)
	if !ok {
		return session.Device{}, errSessionIdNotFound
	}

	return d, nil
}
