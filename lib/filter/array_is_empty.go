package filter

type ArrayIsEmptyFilter[E any] struct {
	Values []E
}

func (*ArrayIsEmptyFilter[E]) isFilter([]E) {}

func ArrayIsEmpty[E any](values []E) *ArrayIsEmptyFilter[E] {
	return &ArrayIsEmptyFilter[E]{Values: values}
}
