package player

import "github.com/ysomad/answersuck/internal/pkg/pgclient"

const playerTable = "player"

type Repository struct {
	*pgclient.Client
}

func NewRepository(c *pgclient.Client) *Repository {
	return &Repository{c}
}
