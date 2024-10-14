package api

type Filter interface {
	String() string
}

type StringFilter struct {
	string
}

func NewStringFilter(s string) *StringFilter {
	return &StringFilter{s}
}

func (f *StringFilter) String() string {
	return f.string
}
