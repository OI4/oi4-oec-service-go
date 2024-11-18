package api

type Filter interface {
	String() string
}

type StringFilter struct {
	string
}

func FilterEquals(this Filter, that Filter) bool {
	if this == nil || that == nil {
		return false
	}
	return this.String() == that.String()
}

func NewStringFilter(s string) *StringFilter {
	return &StringFilter{s}
}

func (f *StringFilter) String() string {
	return f.string
}
