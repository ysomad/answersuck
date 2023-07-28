package round

import "github.com/ysomad/answersuck/internal/pkg/pgclient"

const roundsTable = "rounds"

type Repository struct {
	*pgclient.Client
}

func NewRepository(c *pgclient.Client) *Repository {
	return &Repository{c}
}
