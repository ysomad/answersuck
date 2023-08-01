package round

import "github.com/ysomad/answersuck/internal/pkg/pgclient"

const RoundsTable = "rounds"

type Repository struct {
	*pgclient.Client
}

func NewRepository(c *pgclient.Client) *Repository {
	return &Repository{c}
}
