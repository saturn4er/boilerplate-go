package pagination

import (
	"math"
)

type Pagination struct {
	Page    int
	PerPage int
}

type PaginatedResult[T any] struct {
	Items          []*T
	PaginationData Data
}

func NewPaginatedResult[T any](items []*T, pagination *Pagination, totalCount int) *PaginatedResult[T] {
	return &PaginatedResult[T]{
		Items:          items,
		PaginationData: NewPaginationData(pagination, totalCount),
	}
}

type Data struct {
	Page       int
	PerPage    int
	TotalPages int
	TotalCount int
}

func NewPaginationData(pagination *Pagination, count int) Data {
	if pagination == nil || pagination.PerPage == 0 {
		return Data{
			Page:       1,
			PerPage:    count,
			TotalPages: 1,
			TotalCount: count,
		}
	}

	return Data{
		Page:       pagination.Page,
		PerPage:    pagination.PerPage,
		TotalPages: int(math.Ceil(float64(count) / float64(pagination.PerPage))),
		TotalCount: count,
	}
}
