package media

import "github.com/ysomad/answersuck/internal/pkg/pgclient"

const mediaTable = "media"

type Repository struct {
	*pgclient.Client
}

func NewRepository(c *pgclient.Client) *Repository {
	return &Repository{c}
}
