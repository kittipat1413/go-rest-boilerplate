package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"go-rest-boilerplate/config"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type dbContextKey struct{}

func FromContext(ctx context.Context) *sqlx.DB {
	return ctx.Value(dbContextKey{}).(*sqlx.DB)
}

func NewContext(ctx context.Context, db *sqlx.DB) context.Context {
	return context.WithValue(ctx, dbContextKey{}, db)
}

func MustConnect(cfg *config.Config) *sqlx.DB {
	if db, err := Connect(cfg); err != nil {
		log.Panicln(err)
		return nil
	} else {
		return db
	}
}

func Connect(cfg *config.Config) (*sqlx.DB, error) {
	if db, err := sqlx.Open("postgres", cfg.DatabaseURL()); err != nil {
		return nil, fmt.Errorf("database: %w", err)
	} else {
		return db, nil
	}
}

func NewScope(ctx context.Context, db *sqlx.DB) (Scope, error) {
	if db == nil {
		db = FromContext(ctx)
	}

	if impl, err := newScope(ctx, db); err != nil {
		return nil, err
	} else {
		return impl, nil
	}
}

// the scope, automatically COMMIT-ing or ROLLBACK depending on the error returned.
func Run[T any](ctx context.Context, action func(Scope) (T, error)) (result T, err error) {
	var scope Scope
	if scope, err = NewScope(ctx, nil); err != nil {
		return result, err
	} else {
		defer scope.End(&err)
	}

	return action(scope)
}

// Exec creates a new scope and executes a single SQL statement in it.
func Exec(ctx context.Context, sql string, args ...interface{}) error {
	_, err := Run(ctx, func(s Scope) (struct{}, error) {
		return struct{}{}, s.Exec(sql, args...)
	})
	return err
}

// Get creates a new scope and runs the given SQL query inside it and binds the result
// to the given output argument.
func Get[T any](ctx context.Context, out T, sql string, args ...interface{}) error {
	if _, err := Run(ctx, func(s Scope) (T, error) {
		return out, s.Get(out, sql, args...)
	}); err != nil {
		return err
	} else {
		return nil
	}
}

// Select creates a new scope and runs the given SQL query inside it and binds the
// resulting records to the given output argument.
func Select[T any](ctx context.Context, out T, sql string, args ...interface{}) error {
	if _, err := Run(ctx, func(s Scope) (T, error) {
		return out, s.Select(out, sql, args...)
	}); err != nil {
		return err
	} else {
		return nil
	}
}

func IsNoRows(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

func IsErrorCode(err error, errcode pq.ErrorCode) bool {
	if pgerr, ok := err.(*pq.Error); ok {
		return pgerr.Code == errcode
	}
	return false
}
