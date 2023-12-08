package filter

type Filter[E any] interface {
	isFilter(E)
}
