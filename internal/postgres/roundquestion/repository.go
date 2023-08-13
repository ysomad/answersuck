package roundquestion

import "github.com/ysomad/answersuck/internal/pkg/pgclient"

const RoundQuestionsTable = "round_questions"

type Repository struct {
	*pgclient.Client
}

func NewRepository(c *pgclient.Client) *Repository {
	return &Repository{c}
}
