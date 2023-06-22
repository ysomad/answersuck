package player

import "github.com/ysomad/answersuck/internal/pkg/pgclient"

type postgres struct {
	*pgclient.Client
}

func NewPostgres(p *pgclient.Client) *postgres {
	return &postgres{p}
}
