package postgres

import (
	"context"
	"fmt"

	"github.com/go-pg/pg/v10"
)

type DBLogger struct{}

var _ pg.QueryHook = (*DBLogger)(nil)

func (d DBLogger) AfterQuery(ctx context.Context, q *pg.QueryEvent) error {
	formatted, _ := q.FormattedQuery()
	fmt.Printf("%v\n", string(formatted))
	return nil
}

func (d DBLogger) BeforeQuery(ctx context.Context, q *pg.QueryEvent) (context.Context, error) {
	return ctx, nil
}

func New(opts *pg.Options) *pg.DB {
	return pg.Connect(opts)
}
