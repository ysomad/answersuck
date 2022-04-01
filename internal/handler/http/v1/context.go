package v1

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/answersuck/vault/internal/domain"
)

var (
	ErrAudienceContextNotFound = errors.New("audience not found in context")
)

// accountId returns account id from context
func accountId(c *gin.Context) (string, error) {
	aid := c.GetString("aid")

	_, err := uuid.Parse(aid)
	if err != nil {
		return "", domain.ErrAccountContextNotFound
	}

	return aid, nil
}

// sessionId returns session id from context
func sessionId(c *gin.Context) (string, error) {
	sid := c.GetString("sid")

	if sid == "" {
		return "", domain.ErrSessionContextNotFound
	}

	return sid, nil
}

// audience returns current audience from context
func audience(c *gin.Context) (string, error) {
	aud := c.GetString("aud")

	if aud == "" {
		return "", ErrAudienceContextNotFound
	}

	return aud, nil
}
