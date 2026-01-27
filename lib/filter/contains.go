package filter

type ContainsFilter[T any] struct {
	Substring string
}

func (*ContainsFilter[T]) isFilter(_ T) {}

func Contains[T any](substring string) *ContainsFilter[T] {
	return &ContainsFilter[T]{Substring: substring}
}
