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
	// Driver is a driver that adds sql comments (see https://google.github.io/sqlcommenter/).
	Driver struct {
		dialect.Driver // underlying driver.
		commenter
	}
	// CommentTx is a transaction implementation that adds sql comments.
	CommentTx struct {
		dialect.Tx                 // underlying transaction.
		ctx        context.Context // underlying transaction context.
		commenter
	}
)

func NewDriver(drv dialect.Driver, options ...Option) dialect.Driver {
	defaultCommenters := []Tagger{NewDriverVersionTagger()}
	opts := buildOptions(append(options, WithTagger(defaultCommenters...)))
	return &Driver{drv, commenter{opts}}
}

func (c commenter) getTags(ctx context.Context) SQLCommentTags {
	cmts := make(SQLCommentTags)
	cmts.Append(c.globalComments)
	for _, h := range c.commenters {
		cmts.Append(h.Tag(ctx))
	}
	return cmts
}

func (c commenter) withComment(ctx context.Context, query string) string {
	return fmt.Sprintf("%s /*%s*/", query, c.getTags(ctx).Marshal())
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
	return &CommentTx{tx, ctx, d.commenter}, nil
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
	return &CommentTx{tx, ctx, d.commenter}, nil
}

// Exec adds sql comments and calls the underlying transaction Exec method.
func (d *CommentTx) Exec(ctx context.Context, query string, args, v interface{}) error {
	return d.Tx.Exec(ctx, d.withComment(ctx, query), args, v)
}

// Query adds sql comments and calls the underlying transaction Query method.
func (d *CommentTx) Query(ctx context.Context, query string, args, v interface{}) error {
	return d.Tx.Query(ctx, d.withComment(ctx, query), args, v)
}

// Commit calls the underlying Tx to commit.
func (d *CommentTx) Commit() error {
	return d.Tx.Commit()
}

// Rollback calls the underlying Tx to rollback.
func (d *CommentTx) Rollback() error {
	return d.Tx.Rollback()
}
