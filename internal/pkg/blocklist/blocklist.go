package blocklist

import "sort"

type blockList struct {
	values []string
	len    int
}

// New creates new blockList with provided values
func New(opt option) *blockList {
	v := opt()

	if !sort.StringsAreSorted(v) {
		sort.Strings(v)
	}

	return &blockList{
		values: v,
		len:    len(v),
	}
}

func (l blockList) Find(s string) bool {
	i := sort.Search(l.len, func(i int) bool { return l.values[i] >= s })

	if i < l.len && l.values[i] == s {
		return true
	}

	return false
}
