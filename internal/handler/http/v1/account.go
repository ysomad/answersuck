package v1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Quizish/quizish-backend/internal/app"
	"github.com/Quizish/quizish-backend/internal/domain"
	"github.com/Quizish/quizish-backend/internal/service"
	"github.com/Quizish/quizish-backend/pkg/logging"
	"github.com/Quizish/quizish-backend/pkg/validation"
)

type accountHandler struct {
	validation.ErrorTranslator
	cfg     *app.Config
	log     logging.Logger
	account service.Account
}

func newAccountHandler(handler *gin.RouterGroup, d *Deps) {
	h := &accountHandler{
		d.ErrorTranslator,
		d.Config,
		d.Logger,
		d.AccountService,
	}

	accounts := handler.Group("accounts")
	{
		authenticated := accounts.Group("", sessionMiddleware(d.Logger, d.Config, d.SessionService))
		{
			secure := authenticated.Group("", tokenMiddleware(d.Logger, d.AuthService))
			{
				secure.DELETE("", h.archive)
			}

			authenticated.GET("", h.get)
		}

		accounts.POST("", h.create)
	}
}

type accountCreateRequest struct {
	Email    string `json:"email" binding:"required,email,lte=255"`
	Username string `json:"username" binding:"required,alphanum,gte=4,lte=16"`
	Password string `json:"password" binding:"required,gte=8,lte=64"`
}

func (h *accountHandler) create(c *gin.Context) {
	var r accountCreateRequest

	if err := c.ShouldBindJSON(&r); err != nil {
		abortWithValidationError(c, http.StatusBadRequest, ErrInvalidRequestBody, h.TranslateError(err))
		return
	}

	_, err := h.account.Create(
		c.Request.Context(),
		domain.Account{
			Email:    r.Email,
			Username: r.Username,
			Password: r.Password,
		},
	)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - account - create: %w", err))

		if errors.Is(err, domain.ErrAccountAlreadyExist) {
			abortWithError(c, http.StatusConflict, domain.ErrAccountAlreadyExist, "")
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusCreated)
}

func (h *accountHandler) archive(c *gin.Context) {
	aid, err := accountID(c)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - account - archive - accountID: %w", err))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	sid, err := sessionID(c)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - account - archive - sessionID: %w", err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if err := h.account.Delete(c.Request.Context(), aid, sid); err != nil {
		h.log.Error(fmt.Errorf("http - v1 - account - archive - h.accountService.Delete: %w", err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.SetCookie(h.cfg.SessionCookie, "", -1, "", "", h.cfg.CookieSecure, h.cfg.CookieHTTPOnly)
	c.Status(http.StatusNoContent)
}

func (h *accountHandler) get(c *gin.Context) {
	aid, err := accountID(c)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - account - archive - accountID: %w", err))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	acc, err := h.account.GetByID(c.Request.Context(), aid)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - account - get: %w", err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, acc)
}
