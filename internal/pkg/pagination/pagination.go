package pagination

const (
	DefaultLimit = 10
	MaxLimit     = 100
)

type Params struct {
	LastId uint32 `json:"last_id"`
	Limit  uint64 `json:"limit"`
}

type List[T any] struct {
	Result  []T  `json:"result"`
	HasNext bool `json:"has_next"`
}

func NewList[T any](objList []T, limit uint64) List[T] {
	objLen := uint64(len(objList))
	hasNext := objLen == limit+1
	if hasNext {
		objList = objList[:objLen-1]
	}
	return List[T]{
		Result:  objList,
		HasNext: hasNext,
	}
}
