package filter

type ArrayIsNotEmptyFilter[E any] struct{}

func (*ArrayIsNotEmptyFilter[E]) isFilter([]E) {}

func ArrayIsNotEmpty[E any]() *ArrayIsNotEmptyFilter[E] {
	return &ArrayIsNotEmptyFilter[E]{}
}
