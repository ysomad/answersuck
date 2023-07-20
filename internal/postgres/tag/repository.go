package tag

import "github.com/ysomad/answersuck/internal/pkg/pgclient"

const tagTable = "tag"

type repository struct {
	*pgclient.Client
}

func NewRepository(c *pgclient.Client) *repository {
	return &repository{c}
}
