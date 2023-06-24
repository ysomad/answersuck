package player

type repository interface{}

type Service struct {
	repo repository
}

func NewService(r repository) *Service {
	return &Service{repo: r}
}
