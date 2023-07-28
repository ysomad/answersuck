package topic

import "github.com/ysomad/answersuck/internal/pkg/pgclient"

const topicsTable = "topics"

type Repository struct {
	*pgclient.Client
}

func NewRepository(c *pgclient.Client) *Repository {
	return &Repository{c}
}
