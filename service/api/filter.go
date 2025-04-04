package api

type Filter string

func NewFilter(s string) *Filter {
	filter := Filter(s)
	return &filter
}

func FilterEquals(this *Filter, that *Filter) bool {
	if this == nil || that == nil {
		return false
	}
	return *this == *that
}

func (f *Filter) String() string {
	return string(*f)
}
