package tag

import "github.com/ysomad/answersuck/internal/pkg/pgclient"

const tagTable = "tags"

type repository struct {
	*pgclient.Client
}

func NewRepository(c *pgclient.Client) *repository {
	return &repository{c}
}
