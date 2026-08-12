package filter

type ContainsCIFilter[T any] struct {
	Substring string
}

func (*ContainsCIFilter[T]) isFilter(_ T) {}

func ContainsCI[T any](substring string) *ContainsCIFilter[T] {
	return &ContainsCIFilter[T]{Substring: substring}
}
