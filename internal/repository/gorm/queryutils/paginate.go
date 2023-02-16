package queryutils

import (
	"math"

	"github.com/hatchet-dev/hatchet/internal/repository"
	"gorm.io/gorm"
)

func Paginate(opts []repository.QueryOption, db *gorm.DB, pagination *repository.PaginatedResult) func(db *gorm.DB) *gorm.DB {
	q := ComputeQuery(opts...)

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
		return ApplyOpts(db, q)
	}
}
