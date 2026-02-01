package database

import (
	"context"
	"database/sql"
)

type DatabaseQuery interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...any) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...any) error
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
}
