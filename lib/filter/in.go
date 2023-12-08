package filter

type InFilter[E any] struct {
	Values []E
}

func (*InFilter[E]) isFilter(E) {}

func In[E any](values ...E) *InFilter[E] {
	return &InFilter[E]{Values: values}
}
