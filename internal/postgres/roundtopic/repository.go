package roundtopic

import "github.com/ysomad/answersuck/internal/pkg/pgclient"

const roundTopicsTable = "round_topics"

type Repository struct {
	*pgclient.Client
}

func NewRepository(c *pgclient.Client) *Repository {
	return &Repository{c}
}
