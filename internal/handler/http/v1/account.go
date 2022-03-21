package v1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/quizlyfun/quizly-backend/internal/config"
	"github.com/quizlyfun/quizly-backend/internal/domain"
	"github.com/quizlyfun/quizly-backend/internal/service"

	"github.com/quizlyfun/quizly-backend/pkg/logging"
	"github.com/quizlyfun/quizly-backend/pkg/validation"
)

const (
	accountParam = "accountId"
)

type accountHandler struct {
	t       validation.ErrorTranslator
	cfg     *config.Aggregate
	log     logging.Logger
	account service.Account
}

func newAccountHandler(handler *gin.RouterGroup, d *Deps) {
	h := &accountHandler{
		t:       d.ErrorTranslator,
		cfg:     d.Config,
		log:     d.Logger,
		account: d.AccountService,
	}

	accounts := handler.Group("accounts")
	{
		authenticated := accounts.Group("", sessionMiddleware(d.Logger, &d.Config.Session, d.SessionService))
		{
			withAccountId := authenticated.Group("", accountParamMiddleware(d.Logger))
			{
				param := urlParam(accountParam)
				withAccountId.GET(param, h.get)
				withAccountId.DELETE(param, tokenMiddleware(d.Logger, d.AuthService), h.archive)
			}
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
		abortWithError(c, http.StatusBadRequest, ErrInvalidRequestBody, h.t.TranslateError(err))
		return
	}

	_, err := h.account.Create(
		c.Request.Context(),
		&domain.Account{
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
	aid, err := accountId(c)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - account - archive - accountId: %w", err))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	sid, err := sessionId(c)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - account - archive - sessionId: %w", err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if err := h.account.Delete(c.Request.Context(), aid, sid); err != nil {
		h.log.Error(fmt.Errorf("http - v1 - account - archive - h.accountService.Delete: %w", err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.SetCookie(h.cfg.Session.CookieKey, "", -1, "", "", h.cfg.Cookie.Secure, h.cfg.Cookie.HTTPOnly)
	c.Status(http.StatusNoContent)
}

func (h *accountHandler) get(c *gin.Context) {
	aid, err := accountId(c)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - account - archive - accountId: %w", err))
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
