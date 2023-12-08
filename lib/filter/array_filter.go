package filter

type ArrayFilter[E any] interface {
	isFilter([]E)
}
