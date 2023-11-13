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

// Tests for the Todo resolver since we customized the pagination.
package todo_test

import (
	"context"
	"entgo.io/contrib/entgql/internal/todo/ent/category"
	"entgo.io/contrib/entgql/internal/todo/ent/todo"
	"fmt"
	"sort"
	"strconv"
	"testing"
	"time"

	"entgo.io/contrib/entgql"
	gen "entgo.io/contrib/entgql/internal/todo"
	"entgo.io/contrib/entgql/internal/todo/ent"
	"entgo.io/contrib/entgql/internal/todo/ent/enttest"
	"entgo.io/contrib/entgql/internal/todo/ent/migrate"
	"entgo.io/ent/dialect"
	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/AlekSi/pointer"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
)

type todoTestSuite struct {
	suite.Suite
	*client.Client
	ent *ent.Client
}

const (
	queryAll = `query {
		todos {
			totalCount
			items {
				id
			}
		}
	}`
	maxTodos = 32
	idOffset = 6 << 32
)

func (s *todoTestSuite) SetupTest() {
	time.Local = time.UTC
	s.ent = enttest.Open(s.T(), dialect.SQLite,
		fmt.Sprintf("file:%s-%d?mode=memory&cache=shared&_fk=1",
			s.T().Name(), time.Now().UnixNano(),
		),
		enttest.WithMigrateOptions(migrate.WithGlobalUniqueID(true)),
	)

	srv := handler.NewDefaultServer(gen.NewSchema(s.ent))
	srv.Use(entgql.Transactioner{TxOpener: s.ent})
	s.Client = client.New(srv)

	const mutation = `mutation($priority: Int!, $text: String!, $parent: ID) {
		createTodo(input: {status: COMPLETED, priority: $priority, text: $text, parentID: $parent}) {
			id
		}
	}`
	var (
		rsp struct {
			CreateTodo struct {
				ID string
			}
		}
		root = idOffset + 1
	)
	for i := 1; i <= maxTodos; i++ {
		id := strconv.Itoa(idOffset + i)
		var parent *int
		if i != 1 {
			if i%2 != 0 {
				parent = pointer.ToInt(idOffset + i - 2)
			} else {
				parent = pointer.ToInt(root)
			}
		}
		err := s.Post(mutation, &rsp,
			client.Var("priority", i),
			client.Var("text", id),
			client.Var("parent", parent),
		)
		s.Require().NoError(err)
		s.Require().Equal(id, rsp.CreateTodo.ID)
	}
}

func TestTodo(t *testing.T) {
	suite.Run(t, &todoTestSuite{})
}

type response struct {
	Todos struct {
		TotalCount int
		Items      []struct {
			ID   string
			Text string
		}
	}
}

func (s *todoTestSuite) TestQueryEmpty() {
	{
		var rsp struct{ ClearTodos int }
		err := s.Post(`mutation { clearTodos }`, &rsp)
		s.Require().NoError(err)
		s.Require().Equal(maxTodos, rsp.ClearTodos)
	}
	var rsp response
	err := s.Post(queryAll, &rsp)
	s.Require().NoError(err)
	s.Require().Zero(rsp.Todos.TotalCount)
	s.Require().Empty(rsp.Todos.Items)
}

func (s *todoTestSuite) TestQueryAll() {
	var rsp response
	err := s.Post(queryAll, &rsp)
	s.Require().NoError(err)

	s.Require().Equal(maxTodos, rsp.Todos.TotalCount)
	s.Require().Len(rsp.Todos.Items, maxTodos)
	for i, item := range rsp.Todos.Items {
		s.Require().Equal(strconv.Itoa(idOffset+i+1), item.ID)
	}
}

func (s *todoTestSuite) TestPageForward() {
	const (
		query = `query($offset: Int, $limit: Int) {
			todos(offset: $offset, limit: $limit) {
				totalCount
				items {
					id
				}
			}
		}`
		limit = 5
	)
	var (
		rsp response
		id  = idOffset + 1
	)
	for i := 0; i < maxTodos/limit; i++ {
		err := s.Post(query, &rsp,
			client.Var("offset", i*limit),
			client.Var("limit", limit),
		)
		s.Require().NoError(err)
		s.Require().Equal(maxTodos, rsp.Todos.TotalCount)
		s.Require().Len(rsp.Todos.Items, limit)
		for _, item := range rsp.Todos.Items {
			s.Require().Equal(strconv.Itoa(id), item.ID)
			id++
		}
	}
}

func (s *todoTestSuite) TestPageBackwards() {
	const (
		query = `query($offset: Int, $limit: Int) {
			todos(offset: $offset, limit: $limit) {
				totalCount
				items {
					id
				}
			}
		}`
		limit = 7
	)
	var (
		rsp response
		id  = idOffset + maxTodos
	)
	for i := 0; i < maxTodos/limit; i++ {
		err := s.Post(query, &rsp,
			client.Var("offset", limit*(maxTodos/limit-i)),
			client.Var("limit", limit),
		)
		s.Require().NoError(err)
		s.Require().Equal(maxTodos, rsp.Todos.TotalCount)

		for i := len(rsp.Todos.Items) - 1; i >= 0; i-- {
			item := &rsp.Todos.Items[i]
			s.Require().Equal(strconv.Itoa(id), item.ID)
			id--
		}
	}
}

func (s *todoTestSuite) TestPaginationOrder() {
	const (
		query = `query($offset: Int, $limit: Int, $direction: OrderDirection!, $field: TodoOrderField!) {
			todos(offset: $offset, limit: $limit, orderBy: { direction: $direction, field: $field }) {
				totalCount
				items {
					id
					text
				}
			}
		}`
		limit = 5
		steps = maxTodos/limit + 1
	)
	s.Run("ForwardAscending", func() {
		var (
			rsp response
		)
		for i := 0; i < steps; i++ {
			err := s.Post(query, &rsp,
				client.Var("offset", i*limit),
				client.Var("limit", limit),
				client.Var("direction", "ASC"),
				client.Var("field", "TEXT"),
			)
			s.Require().NoError(err)
			s.Require().Equal(maxTodos, rsp.Todos.TotalCount)
			if i < steps-1 {
				s.Require().Len(rsp.Todos.Items, limit)
			} else {
				s.Require().Len(rsp.Todos.Items, maxTodos%limit)
			}
			s.Require().True(sort.SliceIsSorted(rsp.Todos.Items, func(i, j int) bool {
				return rsp.Todos.Items[i].Text < rsp.Todos.Items[j].Text
			}))
		}
	})

	s.Run("ForwardDescending", func() {
		var (
			rsp response
		)
		for i := 0; i < steps; i++ {
			err := s.Post(query, &rsp,
				client.Var("offset", i*limit),
				client.Var("limit", limit),
				client.Var("direction", "DESC"),
				client.Var("field", "CREATED_AT"),
			)
			s.Require().NoError(err)
			s.Require().Equal(maxTodos, rsp.Todos.TotalCount)
			if i < steps-1 {
				s.Require().Len(rsp.Todos.Items, limit)
			} else {
				s.Require().Len(rsp.Todos.Items, maxTodos%limit)
			}
			s.Require().True(sort.SliceIsSorted(rsp.Todos.Items, func(i, j int) bool {
				left, _ := strconv.Atoi(rsp.Todos.Items[i].ID)
				right, _ := strconv.Atoi(rsp.Todos.Items[j].ID)
				return left > right
			}))
		}
	})
}

func (s *todoTestSuite) TestPaginationFiltering() {
	const (
		query = `query($offset: Int, $limit: Int, $status: TodoStatus, $hasParent: Boolean, $hasCategory: Boolean) {
			todos(offset: $offset, limit: $limit, where: {status: $status, hasParent: $hasParent, hasCategory: $hasCategory}) {
				totalCount
				items {
					id
				}
			}
		}`
		limit = 5
		steps = maxTodos/limit + 1
	)
	s.Run("StatusInProgress", func() {
		var rsp response
		err := s.Post(query, &rsp,
			client.Var("limit", limit),
			client.Var("offset", 0),
			client.Var("status", todo.StatusInProgress),
		)
		s.NoError(err)
		s.Zero(rsp.Todos.TotalCount)
	})
	s.Run("StatusCompleted", func() {
		var rsp response
		for i := 0; i < steps; i++ {
			err := s.Post(query, &rsp,
				client.Var("offset", i*limit),
				client.Var("limit", limit),
				client.Var("status", todo.StatusCompleted),
			)
			s.Require().NoError(err)
			s.Require().Equal(maxTodos, rsp.Todos.TotalCount)
			if i < steps-1 {
				s.Require().Len(rsp.Todos.Items, limit)
			} else {
				s.Require().Len(rsp.Todos.Items, maxTodos%limit)
			}
		}
	})
	s.Run("WithoutCategory", func() {
		var rsp response
		err := s.Post(query, &rsp,
			client.Var("limit", limit),
			client.Var("status", todo.StatusCompleted),
			client.Var("hasCategory", true),
		)
		s.Require().NoError(err)
		s.Require().Equal(0, rsp.Todos.TotalCount)
	})

	s.Run("WithCategory", func() {
		ctx := context.Background()
		id := s.ent.Todo.Query().Order(ent.Asc(todo.FieldID)).FirstIDX(ctx)
		s.ent.Category.Create().SetText("Disabled").SetStatus(category.StatusDisabled).AddTodoIDs(id).SetDuration(time.Second).ExecX(ctx)

		var (
			rsp   response
			query = `query($duration: Duration) {
					todos(where:{hasCategoryWith: {duration: $duration}}) {
						totalCount
					}
				}`
		)
		err := s.Post(query, &rsp, client.Var("duration", time.Second))
		s.NoError(err)
		s.Equal(1, rsp.Todos.TotalCount)
		err = s.Post(query, &rsp, client.Var("duration", time.Second*2))
		s.NoError(err)
		s.Zero(rsp.Todos.TotalCount)
	})
}
