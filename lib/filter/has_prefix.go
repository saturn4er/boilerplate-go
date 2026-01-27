package filter

type HasPrefixFilter[T any] struct {
	Prefix string
}

func (*HasPrefixFilter[T]) isFilter(_ T) {}

func HasPrefix[T any](prefix string) *HasPrefixFilter[T] {
	return &HasPrefixFilter[T]{Prefix: prefix}
}
