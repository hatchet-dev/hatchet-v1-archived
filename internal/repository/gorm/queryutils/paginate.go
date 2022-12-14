package queryutils

import (
	"fmt"
	"math"

	"github.com/hatchet-dev/hatchet/internal/repository"
	"gorm.io/gorm"
)

func Paginate(opts []repository.QueryOption, db *gorm.DB, pagination *repository.PaginatedResult) func(db *gorm.DB) *gorm.DB {
	q := repository.Query{
		Limit:  0,
		Offset: 0,
		Order:  repository.OrderDesc,
		SortBy: "updated_at",
	}

	for _, opt := range opts {
		opt.Apply(&q)
	}

	if q.Limit == 0 {
		// default to returing 50 results per page
		q.Limit = 50
	}

	var totalRows int64

	db.Count(&totalRows)

	pagination.NumPages = int64(math.Ceil(float64(totalRows) / float64(q.Limit)))

	if q.Offset > 0 {
		pagination.CurrentPage = int64(q.Offset / q.Limit)

		if pagination.CurrentPage < pagination.NumPages {
			pagination.NextPage = pagination.CurrentPage + 1
		} else {
			pagination.NextPage = pagination.CurrentPage
		}
	}

	return func(db *gorm.DB) *gorm.DB {
		return db.
			Offset(q.Offset).
			Limit(q.Limit).
			Order(fmt.Sprintf("%s %s", q.SortBy, q.Order))
	}
}
