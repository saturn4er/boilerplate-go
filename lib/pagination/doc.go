// Package pagination provides generic pagination types for database queries.
//
// [Pagination] specifies the page number and items per page.
// [PaginatedResult] wraps a slice of items with pagination metadata
// including total pages and total count.
//
//	result := pagination.NewPaginatedResult(items, &pagination.Pagination{
//	    Page:    1,
//	    PerPage: 20,
//	}, totalCount)
//	// result.PaginationData.TotalPages == ceil(totalCount / 20)
package pagination
