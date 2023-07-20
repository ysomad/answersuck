package pgsearch

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

// conf is a postgres text search configuration.
type conf string

const (
	ConfRus conf = "russian"
	ConfEng conf = "english"
)

type order string

const (
	OrderDESC    order = "DESC"
	OrderASC     order = "ASC"
	OrderDefault order = OrderDESC
)

// Where attaches where clause to builder with query for fulltext search.
func Where(b sq.SelectBuilder, query string, c conf) sq.SelectBuilder {
	return b.Where(sq.Expr("ts @@ phraseto_tsquery('"+string(c)+"', ?)", query))
}

// OrderBy attaches order by clause to builder with order by for fulltext search.
func OrderBy(b sq.SelectBuilder, query string, c conf, o order) sq.SelectBuilder {
	return b.OrderBy(fmt.Sprintf("ts_rank(ts, to_tsquery('%s', '%s')) %s", c, query, o))
}
