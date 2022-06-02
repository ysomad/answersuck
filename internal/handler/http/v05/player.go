package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/answersuck/vault/internal/config"
	"github.com/answersuck/vault/internal/domain/player"
	"github.com/answersuck/vault/pkg/logging"
	"github.com/gin-gonic/gin"
)

type PlayerService interface {
	GetByNickname(ctx context.Context, nickname string) (player.Player, error)
}

type playerHandler struct {
	t   ErrorTranslator
	cfg *config.Aggregate
	log logging.Logger

	service PlayerService
}

func newPlayerHandler(r *gin.RouterGroup, d *Deps) {
	h := &playerHandler{
		t:       d.ErrorTranslator,
		cfg:     d.Config,
		log:     d.Logger,
		service: d.PlayerService,
	}

	players := r.Group("players")
	{
		players.GET(":nickname", h.getByNickname)
	}
}

const (
	paramNickname = "nickname"
)

func (h *playerHandler) getByNickname(c *gin.Context) {
	n := c.Param(paramNickname)

	p, err := h.service.GetByNickname(c.Request.Context(), n)
	if err != nil {
		h.log.Error("http - v1 - player - getByNickname - h.service.GetByNickname: %w", err)

		if errors.Is(err, player.ErrNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, p)
}
