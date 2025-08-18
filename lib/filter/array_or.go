package filter

type ArrayOrFilter[E any] struct {
	Filters []ArrayFilter[E]
}

func (*ArrayOrFilter[E]) isFilter([]E) {}

func ArrayOr[E any](filters ...ArrayFilter[E]) *ArrayOrFilter[E] {
	return &ArrayOrFilter[E]{Filters: filters}
}
