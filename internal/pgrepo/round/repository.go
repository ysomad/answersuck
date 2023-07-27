package round

import "github.com/ysomad/answersuck/internal/pkg/pgclient"

const roundsTable = "rounds"

type repository struct {
	*pgclient.Client
}

func NewRepository(c *pgclient.Client) *repository {
	return &repository{c}
}
