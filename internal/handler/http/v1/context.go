package v1

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/answersuck/vault/internal/domain"
)

const (
	accountIdKey     = "id"
	sessionIdKey     = "sid"
	audienceKey      = "audience"
	lastIdKey        = "last_id"
	lastCreatedAtKey = "last_created_at"
	limitKey         = "limit"
)

var (
	ErrInvalidLastCreatedAt = errors.New("last created at must be a valid RFC3339 timestamp")
)

// GetAccountId returns account id from context
func GetAccountId(c *gin.Context) (string, error) {
	aid := c.GetString(accountIdKey)

	_, err := uuid.Parse(aid)
	if err != nil {
		return "", domain.ErrAccountContextNotFound
	}

	return aid, nil
}

// GetSessionId returns session id from context
func GetSessionId(c *gin.Context) string { return c.GetString(sessionIdKey) }

// GetAudience returns current GetAudience from context
func GetAudience(c *gin.Context) string { return c.GetString(audienceKey) }
