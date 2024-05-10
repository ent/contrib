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

package entgql_test

import (
	"context"
	"errors"
	"testing"

	"entgo.io/contrib/entgql"
	"entgo.io/contrib/entgql/mocks"
	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler/testserver"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestTransaction(t *testing.T) {
	newServer := func(opener entgql.TxOpener, skip entgql.SkipTxFunc) *testserver.TestServer {
		srv := testserver.New()
		srv.AddTransport(transport.POST{})
		srv.Use(entgql.Transactioner{TxOpener: opener, SkipTxFunc: skip})
		return srv
	}
	fwdCtx := func(ctx context.Context) context.Context {
		return ctx
	}

	t.Run("Query", func(t *testing.T) {
		t.Parallel()
		var opener mocks.TxOpener
		defer opener.AssertExpectations(t)
		srv := newServer(&opener, nil)

		c := client.New(srv)
		err := c.Post(`query { name }`, &struct{ Name string }{})
		require.NoError(t, err)
	})
	t.Run("Mutation", func(t *testing.T) {
		t.Parallel()
		t.Run("OK", func(t *testing.T) {
			var tx mocks.Tx
			tx.On("Commit").
				Return(nil).
				Once()
			defer tx.AssertExpectations(t)

			var opener mocks.TxOpener
			opener.On("OpenTx", mock.Anything).
				Return(fwdCtx, &tx, nil).
				Once()
			defer opener.AssertExpectations(t)

			srv := newServer(&opener, nil)
			srv.AroundResponses(func(context.Context, graphql.ResponseHandler) *graphql.Response {
				return &graphql.Response{Data: []byte(`{"name":"test"}`)}
			})

			c := client.New(srv)
			err := c.Post(`mutation { name }`, &struct{ Name string }{})
			require.NoError(t, err)
		})

		t.Run("SkipOperation", func(t *testing.T) {
			var (
				tx     mocks.Tx
				opener mocks.TxOpener
			)
			tx.On("Commit").
				Return(nil).
				Once()
			defer tx.AssertExpectations(t)

			srv := newServer(&opener, entgql.SkipOperations("skipped"))
			srv.AroundResponses(func(context.Context, graphql.ResponseHandler) *graphql.Response {
				return &graphql.Response{Data: []byte(`{"name":"test"}`)}
			})

			c := client.New(srv)
			err := c.Post(`mutation skipped { name }`, &struct{ Name string }{})
			require.NoError(t, err)
			opener.AssertExpectations(t)

			opener.On("OpenTx", mock.Anything).
				Return(fwdCtx, &tx, nil).
				Once()
			err = c.Post(`mutation notSkipped { name }`, &struct{ Name string }{})
			require.NoError(t, err)
			opener.AssertExpectations(t)
		})

		t.Run("SkipIfHasFields", func(t *testing.T) {
			var (
				tx     mocks.Tx
				opener mocks.TxOpener
			)
			tx.On("Commit").
				Return(nil).
				Once()
			defer tx.AssertExpectations(t)
			defer opener.AssertExpectations(t)

			srv := newServer(&opener, entgql.SkipIfHasFields("name"))
			srv.AroundResponses(func(context.Context, graphql.ResponseHandler) *graphql.Response {
				return &graphql.Response{Data: []byte(`{"name":"test"}`)}
			})
			c := client.New(srv)
			err := c.Post(`mutation { name }`, &struct{ Name string }{})
			require.NoError(t, err)
			opener.AssertExpectations(t)

			opener.On("OpenTx", mock.Anything).
				Return(fwdCtx, &tx, nil).
				Once()
			srv = newServer(&opener, entgql.SkipIfHasFields("work"))
			srv.AroundResponses(func(context.Context, graphql.ResponseHandler) *graphql.Response {
				return &graphql.Response{Data: []byte(`{"name":"test"}`)}
			})
			c = client.New(srv)
			err = c.Post(`mutation { name }`, &struct{ Name string }{})
			require.NoError(t, err)
			opener.AssertExpectations(t)
		})

		t.Run("Err", func(t *testing.T) {
			t.Parallel()
			var tx mocks.Tx
			tx.On("Rollback").
				Return(nil).
				Once()
			defer tx.AssertExpectations(t)

			var opener mocks.TxOpener
			opener.On("OpenTx", mock.Anything).
				Return(fwdCtx, &tx, nil).
				Once()
			defer opener.AssertExpectations(t)

			srv := newServer(&opener, nil)
			srv.AroundResponses(func(ctx context.Context, _ graphql.ResponseHandler) *graphql.Response {
				return graphql.ErrorResponse(ctx, "bad mutation")
			})

			c := client.New(srv)
			err := c.Post(`mutation { name }`, &struct{ Name string }{})
			require.Error(t, err)
			require.Contains(t, err.Error(), "bad mutation")
		})
		t.Run("Panic", func(t *testing.T) {
			t.Parallel()
			var tx mocks.Tx
			tx.On("Rollback").
				Return(nil).
				Once()
			defer tx.AssertExpectations(t)

			var opener mocks.TxOpener
			opener.On("OpenTx", mock.Anything).
				Return(fwdCtx, &tx, nil).
				Once()
			defer opener.AssertExpectations(t)

			srv := newServer(&opener, nil)
			srv.SetRecoverFunc(func(_ context.Context, err interface{}) error {
				return err.(error)
			})
			srv.AroundResponses(func(ctx context.Context, _ graphql.ResponseHandler) *graphql.Response {
				panic(graphql.ErrorOnPath(ctx, errors.New("oh no")))
			})

			c := client.New(srv)
			err := c.Post(`mutation { name }`, &struct{ Name string }{})
			require.Error(t, err)
			require.Contains(t, err.Error(), "oh no")
		})
		t.Run("NoTx", func(t *testing.T) {
			t.Parallel()
			var opener mocks.TxOpener
			opener.On("OpenTx", mock.Anything).
				Return(nil, nil, errors.New("bad tx")).
				Once()
			defer opener.AssertExpectations(t)

			srv := newServer(&opener, nil)
			c := client.New(srv)
			err := c.Post(`mutation { name }`, &struct{ Name string }{})
			require.Error(t, err)
			require.Contains(t, err.Error(), "bad tx")
		})
	})
}
