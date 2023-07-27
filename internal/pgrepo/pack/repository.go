package pack

import "github.com/ysomad/answersuck/internal/pkg/pgclient"

const (
	PacksTable    = "packs"
	packTagsTable = "pack_tags"
)

type repository struct {
	*pgclient.Client
}

func NewRepository(c *pgclient.Client) *repository {
	return &repository{c}
}
