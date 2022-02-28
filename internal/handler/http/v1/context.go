package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/quizlyfun/quizly-backend/internal/domain"
)

// accountID returns account id from context
func accountID(c *gin.Context) (string, error) {
	aid := c.GetString("aid")

	_, err := uuid.Parse(aid)
	if err != nil {
		return "", domain.ErrAccountContextNotFound
	}

	return aid, nil
}

// sessionID return session id from context
func sessionID(c *gin.Context) (string, error) {
	sid := c.GetString("sid")

	if sid == "" {
		return "", domain.ErrSessionContextNotFound
	}

	return sid, nil
}
