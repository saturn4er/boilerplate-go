package filter

type IsNotNullFilter[E any] struct {
}

func (*IsNotNullFilter[E]) isFilter(_ E) {}

func IsNotNull[E any]() *IsNotNullFilter[E] {
	return &IsNotNullFilter[E]{}
}
