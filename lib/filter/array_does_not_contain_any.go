package filter

type ArrayDoesNotContainAnyFilter[E any] struct {
	Values []E
}

func (*ArrayDoesNotContainAnyFilter[E]) isFilter([]E) {}

func ArrayDoesNotContainAny[E any](values []E) *ArrayDoesNotContainAnyFilter[E] {
	return &ArrayDoesNotContainAnyFilter[E]{Values: values}
}
