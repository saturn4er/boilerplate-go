package filter

type LessFilter[E any] struct {
	Value E
}

func (*LessFilter[E]) isFilter(_ E) {}

func Less[E any](val E) *LessFilter[E] {
	return &LessFilter[E]{Value: val}
}

type LessOrEqualsFilter[E any] struct {
	Value E
}

func (*LessOrEqualsFilter[E]) isFilter(_ E) {}

func LessOrEquals[E any](val E) *LessOrEqualsFilter[E] {
	return &LessOrEqualsFilter[E]{Value: val}
}
