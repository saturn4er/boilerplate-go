package filter

type ArrayContainsAnyFilter[E any] struct {
	Values []E
}

func (*ArrayContainsAnyFilter[E]) isFilter([]E) {}

func ArrayContainsAny[E any](values []E) *ArrayContainsAnyFilter[E] {
	return &ArrayContainsAnyFilter[E]{Values: values}
}
