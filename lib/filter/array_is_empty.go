package filter

type ArrayIsEmptyFilter[E any] struct{}

func (*ArrayIsEmptyFilter[E]) isFilter([]E) {}

func ArrayIsEmpty[E any]() *ArrayIsEmptyFilter[E] {
	return &ArrayIsEmptyFilter[E]{}
}
