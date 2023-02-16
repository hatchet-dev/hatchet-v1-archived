package queryutils

import (
	"fmt"

	"github.com/hatchet-dev/hatchet/internal/repository"
	"gorm.io/gorm"
)

func ComputeQuery(opts ...repository.QueryOption) repository.Query {
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

	return q
}

func ApplyOpts(db *gorm.DB, q repository.Query) *gorm.DB {
	// apply query options to the subquery as well
	return db.
		Offset(q.Offset).
		Limit(q.Limit).
		Order(fmt.Sprintf("%s %s", q.SortBy, q.Order))
}
