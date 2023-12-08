package filter

type ArrayContainsFilter[E any] struct {
	Values []E
}

func (*ArrayContainsFilter[E]) isFilter([]E) {}

func ArrayContains[E any](values []E) *ArrayContainsFilter[E] {
	return &ArrayContainsFilter[E]{Values: values}
}
