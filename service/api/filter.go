package api

type Filter interface {
	String() string
	Equals(Filter) bool
}

type StringFilter struct {
	string
}

func (f *StringFilter) Equals(other Filter) bool {
	return f.string == other.String()
}

func NewStringFilter(s string) *StringFilter {
	return &StringFilter{s}
}

func (f *StringFilter) String() string {
	return f.string
}
