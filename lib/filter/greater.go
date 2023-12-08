package filter

type GreaterFilter[E any] struct {
	Value E
}

func (*GreaterFilter[E]) isFilter(_ E) {}

func Greater[E any](val E) *GreaterFilter[E] {
	return &GreaterFilter[E]{Value: val}
}

type GreaterOrEqualsFilter[E any] struct {
	Value E
}

func (*GreaterOrEqualsFilter[E]) isFilter(_ E) {}

func GreaterOrEquals[E any](val E) *GreaterOrEqualsFilter[E] {
	return &GreaterOrEqualsFilter[E]{Value: val}
}
