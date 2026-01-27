package filter

type HasSuffixFilter[T any] struct {
	Suffix string
}

func (*HasSuffixFilter[T]) isFilter(_ T) {}

func HasSuffix[T any](suffix string) *HasSuffixFilter[T] {
	return &HasSuffixFilter[T]{Suffix: suffix}
}
