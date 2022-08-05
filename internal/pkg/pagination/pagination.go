package pagination

const (
	DefaultLimit = 10
	MaxLimit     = 50
)

type CursorParams struct {
	Limit  uint64 `json:"limit"`
	LastId uint32 `json:"last_id"`
}

type List[T any] struct {
	Result  []T  `json:"result"`
	HasNext bool `json:"has_next"`
}

func NewList[T any](objList []T, limit uint64) List[T] {
	return List[T]{
		Result:  objList,
		HasNext: uint64(len(objList)) == limit+1,
	}
}
