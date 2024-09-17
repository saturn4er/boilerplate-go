package filter

type IsNullFilter[E any] struct{}

func (*IsNullFilter[E]) isFilter(_ E) {}

func IsNull[E any]() *IsNullFilter[E] {
	return &IsNullFilter[E]{}
}
