package filter

type AndFilter[E any] struct {
	Filters []Filter[E]
}

func (*AndFilter[E]) isFilter(E) {}

func And[E any](filters ...Filter[E]) *AndFilter[E] {
	return &AndFilter[E]{Filters: filters}
}
