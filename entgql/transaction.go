// Copyright 2019-present Facebook
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package entgql

import (
	"context"
	"database/sql/driver"
	"errors"
	"slices"
	"sync"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// TxOpener represents types than can open transactions.
type TxOpener interface {
	OpenTx(ctx context.Context) (context.Context, driver.Tx, error)
}

// The TxOpenerFunc type is an adapter to allow the use of
// ordinary functions as tx openers.
type TxOpenerFunc func(ctx context.Context) (context.Context, driver.Tx, error)

// OpenTx returns f(ctx).
func (f TxOpenerFunc) OpenTx(ctx context.Context) (context.Context, driver.Tx, error) {
	return f(ctx)
}

type (
	// Transactioner for graphql mutations.
	Transactioner struct {
		TxOpener
		SkipTxFunc
	}
	// SkipTxFunc allows skipping operations from
	// running under a transaction.
	SkipTxFunc func(*ast.OperationDefinition) bool
)

// SkipOperations skips the given operation names from running
// under a transaction.
func SkipOperations(names ...string) SkipTxFunc {
	return func(op *ast.OperationDefinition) bool {
		return slices.Contains(names, op.Name)
	}
}

// SkipIfHasFields skips the operation has a mutation field
// with the given names.
func SkipIfHasFields(names ...string) SkipTxFunc {
	return func(op *ast.OperationDefinition) bool {
		return slices.ContainsFunc(op.SelectionSet, func(s ast.Selection) bool {
			f, ok := s.(*ast.Field)
			return ok && slices.Contains(names, f.Name)
		})
	}
}

var _ interface {
	graphql.HandlerExtension
	graphql.OperationContextMutator
	graphql.ResponseInterceptor
} = Transactioner{}

// ExtensionName returns the extension name.
func (Transactioner) ExtensionName() string {
	return "EntGQLTransactioner"
}

// Validate is called when adding an extension to the server, it allows validation against the servers schema.
func (t Transactioner) Validate(graphql.ExecutableSchema) error {
	if t.TxOpener == nil {
		return errors.New("entgql: tx opener is nil")
	}
	return nil
}

// MutateOperationContext serializes field resolvers during mutations.
func (t Transactioner) MutateOperationContext(_ context.Context, oc *graphql.OperationContext) *gqlerror.Error {
	if !t.skipTx(oc.Operation) {
		previous := oc.ResolverMiddleware
		var mu sync.Mutex
		oc.ResolverMiddleware = func(ctx context.Context, next graphql.Resolver) (interface{}, error) {
			mu.Lock()
			defer mu.Unlock()
			return previous(ctx, next)
		}
	}
	return nil
}

// InterceptResponse runs graphql mutations under a transaction.
func (t Transactioner) InterceptResponse(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
	if t.skipTx(graphql.GetOperationContext(ctx).Operation) {
		return next(ctx)
	}
	txCtx, tx, err := t.OpenTx(ctx)
	if err != nil {
		return graphql.ErrorResponse(ctx,
			"cannot create transaction: %s", err.Error(),
		)
	}
	ctx = txCtx

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			panic(r)
		}
	}()
	rsp := next(ctx)
	if len(rsp.Errors) > 0 {
		_ = tx.Rollback()
		return &graphql.Response{
			Errors: rsp.Errors,
		}
	}
	if err := tx.Commit(); err != nil {
		return graphql.ErrorResponse(ctx,
			"cannot commit transaction: %s", err.Error(),
		)
	}
	return rsp
}

func (t Transactioner) skipTx(op *ast.OperationDefinition) bool {
	return op == nil || op.Operation != ast.Mutation || (t.SkipTxFunc != nil && t.SkipTxFunc(op))
}
