package filter

type NotInFilter[E any] struct {
	Values []E
}

func (*NotInFilter[E]) isFilter(E) {}

func NotIn[E any](val ...E) *NotInFilter[E] {
	return &NotInFilter[E]{Values: val}
}
