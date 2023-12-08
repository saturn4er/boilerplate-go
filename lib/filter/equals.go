package filter

type EqualsFilter[E any] struct {
	Value E
}

func (*EqualsFilter[E]) isFilter(_ E) {}

func Equals[E any](val E) *EqualsFilter[E] {
	return &EqualsFilter[E]{Value: val}
}
