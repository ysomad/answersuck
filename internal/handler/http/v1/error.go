package v1

import (
	"errors"

	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidRequestBody = errors.New("invalid request body")
)

type errorResponse struct {
	Error  string `json:"error"`
	Detail string `json:"details"`
}

type validationErrorResponse struct {
	Error  string            `json:"error"`
	Detail map[string]string `json:"details"`
}

func abortWithError(c *gin.Context, code int, err error, detail string) {
	c.AbortWithStatusJSON(code, errorResponse{
		Error:  err.Error(),
		Detail: detail,
	})
}

func abortWithValidationError(c *gin.Context, code int, err error, errs map[string]string) {
	c.AbortWithStatusJSON(code, validationErrorResponse{
		Error:  err.Error(),
		Detail: errs,
	})
}
