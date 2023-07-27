package pack

import "github.com/ysomad/answersuck/internal/pkg/pgclient"

const packTable = "packs"

type Repository struct {
	*pgclient.Client
}

func NewRepository(c *pgclient.Client) *Repository {
	return &Repository{c}
}
