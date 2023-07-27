package tag

import "github.com/ysomad/answersuck/internal/pkg/pgclient"

const TagsTable = "tags"

type repository struct {
	*pgclient.Client
}

func NewRepository(c *pgclient.Client) *repository {
	return &repository{c}
}
