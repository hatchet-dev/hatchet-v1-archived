package types

type PaginationRequest struct {
	// The page to query for
	// in: query
	// example: 2
	Page int64 `schema:"page" json:"page"`
}

// swagger:model
type PaginationResponse struct {
	// the total number of pages for listing
	// example: 10
	NumPages int64 `json:"num_pages" form:"required"`

	// the current page
	// example: 2
	CurrentPage int64 `json:"current_page" form:"required"`

	// the next page
	// example: 3
	NextPage int64 `json:"next_page" form:"required"`
}
