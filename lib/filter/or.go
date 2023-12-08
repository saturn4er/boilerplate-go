package filter

type OrFilter[E any] struct {
	Filters []Filter[E]
}

func (*OrFilter[E]) isFilter(E) {}

func Or[E any](filters ...Filter[E]) *OrFilter[E] {
	return &OrFilter[E]{Filters: filters}
}
