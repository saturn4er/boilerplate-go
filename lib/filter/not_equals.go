package filter

type NotEqualsFilter[E any] struct {
	Value E
}

func (*NotEqualsFilter[E]) isFilter(E) {}

func NotEquals[E any](val E) *NotEqualsFilter[E] {
	return &NotEqualsFilter[E]{Value: val}
}
