package repository

import "github.com/hatchet-dev/hatchet/api/v1/types"

type PaginatedResult struct {
	NumPages    int64
	CurrentPage int64
	NextPage    int64
}

// there is no difference between pagination result type in the `types` package
// vs this package at the moment, but we implement this to enforce a strict separation
// between the repository types and the API types
func (p *PaginatedResult) ToAPIType() *types.PaginationResponse {
	return &types.PaginationResponse{
		NumPages:    p.NumPages,
		NextPage:    p.NextPage,
		CurrentPage: p.CurrentPage,
	}
}

type Ordering string

const (
	OrderAsc  Ordering = "asc"
	OrderDesc Ordering = "desc"
)

type Query struct {
	Limit  int
	Offset int
	SortBy string
	Order  Ordering
}

type QueryOption interface {
	Apply(*Query)
}

func WithPage(paginationRequest *types.PaginationRequest) QueryOption {
	var page int

	if paginationRequest == nil {
		page = 1
	} else {
		page = int(paginationRequest.Page)
	}

	return withPage(page)
}

type withPage int

func (w withPage) Apply(q *Query) {
	q.Limit = 50
	q.Offset = 50 * int(w-1)
}

func WithLimit(limit uint) QueryOption {
	return withLimit(limit)
}

type withLimit int

func (w withLimit) Apply(q *Query) {
	q.Limit = int(w)
}

func WithOffset(offset int64) QueryOption {
	return withOffset(offset)
}

type withOffset int

func (w withOffset) Apply(q *Query) {
	q.Offset = int(w)
}

func WithOrder(order Ordering) QueryOption {
	return withOrder(order)
}

type withOrder Ordering

func (w withOrder) Apply(q *Query) {
	q.Order = Ordering(w)
}

func WithSortBy(sortBy string) QueryOption {
	return withSortBy(sortBy)
}

type withSortBy string

func (w withSortBy) Apply(q *Query) {
	q.SortBy = string(w)
}
