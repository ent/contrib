package sqlcommenter

import (
	"context"
	"database/sql"
	"fmt"

	"entgo.io/ent/dialect"
)

type (
	commenter struct {
		options
	}
	// Driver is a driver that adds sql comment (see https://google.github.io/sqlcommenter/).
	Driver struct {
		dialect.Driver // underlying driver.
		commenter
	}
	// Tx is a transaction implementation that adds sql comment.
	Tx struct {
		dialect.Tx                 // underlying transaction.
		ctx        context.Context // underlying transaction context.
		commenter
	}
)

// NewDriver decorates the given driver and add sql comment to every query.
func NewDriver(drv dialect.Driver, options ...Option) dialect.Driver {
	defaultCommenters := []Tagger{NewDriverVersionTagger()}
	opts := buildOptions(append(options, WithTagger(defaultCommenters...)))
	return &Driver{drv, commenter{opts}}
}

func (c commenter) withComment(ctx context.Context, query string) string {
	tags := make(Tags)
	for _, h := range c.taggers {
		tags.Merge(h.Tag(ctx))
	}
	return fmt.Sprintf("%s /*%s*/", query, tags.Marshal())
}

// Exec adds sql comment to the original query and calls the underlying driver Exec method.
func (d *Driver) Exec(ctx context.Context, query string, args, v interface{}) error {
	return d.Driver.Exec(ctx, d.withComment(ctx, query), args, v)
}

// Query adds sql comment to the original query and calls the underlying driver Query method.
func (d *Driver) Query(ctx context.Context, query string, args, v interface{}) error {
	return d.Driver.Query(ctx, d.withComment(ctx, query), args, v)
}

// Tx wraps the underlying Tx command with a commenter.
func (d *Driver) Tx(ctx context.Context) (dialect.Tx, error) {
	tx, err := d.Driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	return &Tx{tx, ctx, d.commenter}, nil
}

// BeginTx wraps the underlying transaction with commenter and calls the underlying driver BeginTx command if it's supported.
func (d *Driver) BeginTx(ctx context.Context, opts *sql.TxOptions) (dialect.Tx, error) {
	drv, ok := d.Driver.(interface {
		BeginTx(context.Context, *sql.TxOptions) (dialect.Tx, error)
	})
	if !ok {
		return nil, fmt.Errorf("Driver.BeginTx is not supported")
	}
	tx, err := drv.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &Tx{tx, ctx, d.commenter}, nil
}

// Exec adds sql comment and calls the underlying transaction Exec method.
func (d *Tx) Exec(ctx context.Context, query string, args, v interface{}) error {
	return d.Tx.Exec(ctx, d.withComment(ctx, query), args, v)
}

// Query adds sql comment and calls the underlying transaction Query method.
func (d *Tx) Query(ctx context.Context, query string, args, v interface{}) error {
	return d.Tx.Query(ctx, d.withComment(ctx, query), args, v)
}

// Commit calls the underlying Tx to commit.
func (d *Tx) Commit() error {
	return d.Tx.Commit()
}

// Rollback calls the underlying Tx to rollback.
func (d *Tx) Rollback() error {
	return d.Tx.Rollback()
}
