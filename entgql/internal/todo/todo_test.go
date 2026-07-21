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

package todo_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"entgo.io/ent/dialect/sql"

	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/require"

	"entgo.io/contrib/entgql"
	gen "entgo.io/contrib/entgql/internal/todo"
	"entgo.io/contrib/entgql/internal/todo/ent"
	"entgo.io/contrib/entgql/internal/todo/ent/category"
	"entgo.io/contrib/entgql/internal/todo/ent/enttest"
	"entgo.io/contrib/entgql/internal/todo/ent/migrate"
	"entgo.io/contrib/entgql/internal/todo/ent/todo"
	"entgo.io/ent/dialect"
	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/stretchr/testify/suite"
	"github.com/vektah/gqlparser/v2/gqlerror"

	_ "github.com/mattn/go-sqlite3"
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
			edges {
				node {
					id
					status
				}
				cursor
			}
			pageInfo {
				hasNextPage
				hasPreviousPage
				startCursor
				endCursor
			}
		}
	}`
	maxTodos = 32
	idOffset = 7 << 32
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
		Edges      []struct {
			Node struct {
				ID            string
				CreatedAt     string
				PriorityOrder int
				Status        todo.Status
				Text          string
				Parent        struct {
					ID string
				}
			}
			Cursor string
		}
		PageInfo struct {
			HasNextPage     bool
			HasPreviousPage bool
			StartCursor     *string
			EndCursor       *string
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
	s.Require().Empty(rsp.Todos.Edges)
	s.Require().False(rsp.Todos.PageInfo.HasNextPage)
	s.Require().False(rsp.Todos.PageInfo.HasPreviousPage)
	s.Require().Nil(rsp.Todos.PageInfo.StartCursor)
	s.Require().Nil(rsp.Todos.PageInfo.EndCursor)
}

func (s *todoTestSuite) TestQueryAll() {
	var rsp response
	err := s.Post(queryAll, &rsp)
	s.Require().NoError(err)

	s.Require().Equal(maxTodos, rsp.Todos.TotalCount)
	s.Require().Len(rsp.Todos.Edges, maxTodos)
	s.Require().False(rsp.Todos.PageInfo.HasNextPage)
	s.Require().False(rsp.Todos.PageInfo.HasPreviousPage)
	s.Require().Equal(
		rsp.Todos.Edges[0].Cursor,
		*rsp.Todos.PageInfo.StartCursor,
	)
	s.Require().Equal(
		rsp.Todos.Edges[len(rsp.Todos.Edges)-1].Cursor,
		*rsp.Todos.PageInfo.EndCursor,
	)
	for i, edge := range rsp.Todos.Edges {
		s.Require().Equal(strconv.Itoa(idOffset+i+1), edge.Node.ID)
		s.Require().EqualValues(todo.StatusCompleted, edge.Node.Status)
		s.Require().NotEmpty(edge.Cursor)
	}
}

func (s *todoTestSuite) TestPageForward() {
	const (
		query = `query($after: Cursor, $first: Int) {
			todos(after: $after, first: $first) {
				totalCount
				edges {
					node {
						id
					}
					cursor
				}
				pageInfo {
					hasNextPage
					endCursor
				}
			}
		}`
		first = 5
	)
	var (
		after interface{}
		rsp   response
		id    = idOffset + 1
	)
	for i := 0; i < maxTodos/first; i++ {
		err := s.Post(query, &rsp,
			client.Var("after", after),
			client.Var("first", first),
		)
		s.Require().NoError(err)
		s.Require().Equal(maxTodos, rsp.Todos.TotalCount)
		s.Require().Len(rsp.Todos.Edges, first)
		s.Require().True(rsp.Todos.PageInfo.HasNextPage)
		s.Require().NotEmpty(rsp.Todos.PageInfo.EndCursor)

		for _, edge := range rsp.Todos.Edges {
			s.Require().Equal(strconv.Itoa(id), edge.Node.ID)
			s.Require().NotEmpty(edge.Cursor)
			id++
		}
		after = rsp.Todos.PageInfo.EndCursor
	}

	err := s.Post(query, &rsp,
		client.Var("after", rsp.Todos.PageInfo.EndCursor),
		client.Var("first", first),
	)
	s.Require().NoError(err)
	s.Require().Equal(maxTodos, rsp.Todos.TotalCount)
	s.Require().NotEmpty(rsp.Todos.Edges)
	s.Require().Len(rsp.Todos.Edges, maxTodos%first)
	s.Require().False(rsp.Todos.PageInfo.HasNextPage)
	s.Require().NotEmpty(rsp.Todos.PageInfo.EndCursor)

	for _, edge := range rsp.Todos.Edges {
		s.Require().Equal(strconv.Itoa(id), edge.Node.ID)
		s.Require().NotEmpty(edge.Cursor)
		id++
	}

	after = rsp.Todos.PageInfo.EndCursor
	rsp = response{}
	err = s.Post(query, &rsp,
		client.Var("after", after),
		client.Var("first", first),
	)
	s.Require().NoError(err)
	s.Require().Equal(maxTodos, rsp.Todos.TotalCount)
	s.Require().Empty(rsp.Todos.Edges)
	s.Require().Empty(rsp.Todos.PageInfo.EndCursor)
	s.Require().False(rsp.Todos.PageInfo.HasNextPage)
}

func (s *todoTestSuite) TestPageBackwards() {
	const (
		query = `query($before: Cursor, $last: Int) {
			todos(before: $before, last: $last) {
				totalCount
				edges {
					node {
						id
					}
					cursor
				}
				pageInfo {
					hasPreviousPage
					startCursor
				}
			}
		}`
		last = 7
	)
	var (
		before interface{}
		rsp    response
		id     = idOffset + maxTodos
	)
	for i := 0; i < maxTodos/last; i++ {
		err := s.Post(query, &rsp,
			client.Var("before", before),
			client.Var("last", last),
		)
		s.Require().NoError(err)
		s.Require().Equal(maxTodos, rsp.Todos.TotalCount)
		s.Require().Len(rsp.Todos.Edges, last)
		s.Require().True(rsp.Todos.PageInfo.HasPreviousPage)
		s.Require().NotEmpty(rsp.Todos.PageInfo.StartCursor)

		for i := len(rsp.Todos.Edges) - 1; i >= 0; i-- {
			edge := &rsp.Todos.Edges[i]
			s.Require().Equal(strconv.Itoa(id), edge.Node.ID)
			s.Require().NotEmpty(edge.Cursor)
			id--
		}
		before = rsp.Todos.PageInfo.StartCursor
	}

	err := s.Post(query, &rsp,
		client.Var("before", rsp.Todos.PageInfo.StartCursor),
		client.Var("last", last),
	)
	s.Require().NoError(err)
	s.Require().Equal(maxTodos, rsp.Todos.TotalCount)
	s.Require().NotEmpty(rsp.Todos.Edges)
	s.Require().Len(rsp.Todos.Edges, maxTodos%last)
	s.Require().False(rsp.Todos.PageInfo.HasPreviousPage)
	s.Require().NotEmpty(rsp.Todos.PageInfo.StartCursor)

	for i := len(rsp.Todos.Edges) - 1; i >= 0; i-- {
		edge := &rsp.Todos.Edges[i]
		s.Require().Equal(strconv.Itoa(id), edge.Node.ID)
		s.Require().NotEmpty(edge.Cursor)
		id--
	}
	s.Require().Equal(idOffset, id)

	before = rsp.Todos.PageInfo.StartCursor
	rsp = response{}
	err = s.Post(query, &rsp,
		client.Var("before", before),
		client.Var("last", last),
	)
	s.Require().NoError(err)
	s.Require().Equal(maxTodos, rsp.Todos.TotalCount)
	s.Require().Empty(rsp.Todos.Edges)
	s.Require().Empty(rsp.Todos.PageInfo.StartCursor)
	s.Require().False(rsp.Todos.PageInfo.HasPreviousPage)
}

func (s *todoTestSuite) TestPaginationOrder() {
	const (
		query = `query($after: Cursor, $first: Int, $before: Cursor, $last: Int, $direction: OrderDirection!, $field: TodoOrderField!) {
			todos(after: $after, first: $first, before: $before, last: $last, orderBy: { direction: $direction, field: $field }) {
				totalCount
				edges {
					node {
						id
						createdAt
						priorityOrder
						status
						text
					}
					cursor
				}
				pageInfo {
					hasNextPage
					hasPreviousPage
					startCursor
					endCursor
				}
			}
		}`
		step  = 5
		steps = maxTodos/step + 1
	)
	s.Run("ForwardAscending", func() {
		var (
			rsp     response
			endText string
		)
		for i := 0; i < steps; i++ {
			err := s.Post(query, &rsp,
				client.Var("after", rsp.Todos.PageInfo.EndCursor),
				client.Var("first", step),
				client.Var("direction", "ASC"),
				client.Var("field", "TEXT"),
			)
			s.Require().NoError(err)
			s.Require().Equal(maxTodos, rsp.Todos.TotalCount)
			if i < steps-1 {
				s.Require().Len(rsp.Todos.Edges, step)
				s.Require().True(rsp.Todos.PageInfo.HasNextPage)
			} else {
				s.Require().Len(rsp.Todos.Edges, maxTodos%step)
				s.Require().False(rsp.Todos.PageInfo.HasNextPage)
			}
			s.Require().True(sort.SliceIsSorted(rsp.Todos.Edges, func(i, j int) bool {
				return rsp.Todos.Edges[i].Node.Text < rsp.Todos.Edges[j].Node.Text
			}))
			s.Require().NotNil(rsp.Todos.PageInfo.StartCursor)
			s.Require().Equal(*rsp.Todos.PageInfo.StartCursor, rsp.Todos.Edges[0].Cursor)
			s.Require().NotNil(rsp.Todos.PageInfo.EndCursor)
			end := rsp.Todos.Edges[len(rsp.Todos.Edges)-1]
			s.Require().Equal(*rsp.Todos.PageInfo.EndCursor, end.Cursor)
			if i > 0 {
				s.Require().Less(endText, rsp.Todos.Edges[0].Node.Text)
			}
			endText = end.Node.Text
		}
	})
	s.Run("ForwardDescending", func() {
		var (
			rsp   response
			endID int
		)
		for i := 0; i < steps; i++ {
			err := s.Post(query, &rsp,
				client.Var("after", rsp.Todos.PageInfo.EndCursor),
				client.Var("first", step),
				client.Var("direction", "DESC"),
				client.Var("field", "CREATED_AT"),
			)
			s.Require().NoError(err)
			s.Require().Equal(maxTodos, rsp.Todos.TotalCount)
			if i < steps-1 {
				s.Require().Len(rsp.Todos.Edges, step)
				s.Require().True(rsp.Todos.PageInfo.HasNextPage)
			} else {
				s.Require().Len(rsp.Todos.Edges, maxTodos%step)
				s.Require().False(rsp.Todos.PageInfo.HasNextPage)
			}
			s.Require().True(sort.SliceIsSorted(rsp.Todos.Edges, func(i, j int) bool {
				left, _ := strconv.Atoi(rsp.Todos.Edges[i].Node.ID)
				right, _ := strconv.Atoi(rsp.Todos.Edges[j].Node.ID)
				return left > right
			}))
			s.Require().NotNil(rsp.Todos.PageInfo.StartCursor)
			s.Require().Equal(*rsp.Todos.PageInfo.StartCursor, rsp.Todos.Edges[0].Cursor)
			s.Require().NotNil(rsp.Todos.PageInfo.EndCursor)
			end := rsp.Todos.Edges[len(rsp.Todos.Edges)-1]
			s.Require().Equal(*rsp.Todos.PageInfo.EndCursor, end.Cursor)
			if i > 0 {
				id, _ := strconv.Atoi(rsp.Todos.Edges[0].Node.ID)
				s.Require().Greater(endID, id)
			}
			endID, _ = strconv.Atoi(end.Node.ID)
		}
	})
	s.Run("BackwardAscending", func() {
		var (
			rsp           response
			startPriority int
		)
		for i := 0; i < steps; i++ {
			err := s.Post(query, &rsp,
				client.Var("before", rsp.Todos.PageInfo.StartCursor),
				client.Var("last", step),
				client.Var("direction", "ASC"),
				client.Var("field", "PRIORITY_ORDER"),
			)
			s.Require().NoError(err)
			s.Require().Equal(maxTodos, rsp.Todos.TotalCount)
			if i < steps-1 {
				s.Require().Len(rsp.Todos.Edges, step)
				s.Require().True(rsp.Todos.PageInfo.HasPreviousPage)
			} else {
				s.Require().Len(rsp.Todos.Edges, maxTodos%step)
				s.Require().False(rsp.Todos.PageInfo.HasPreviousPage)
			}
			s.Require().True(sort.SliceIsSorted(rsp.Todos.Edges, func(i, j int) bool {
				return rsp.Todos.Edges[i].Node.PriorityOrder < rsp.Todos.Edges[j].Node.PriorityOrder
			}))
			s.Require().NotNil(rsp.Todos.PageInfo.StartCursor)
			start := rsp.Todos.Edges[0]
			s.Require().Equal(*rsp.Todos.PageInfo.StartCursor, start.Cursor)
			s.Require().NotNil(rsp.Todos.PageInfo.EndCursor)
			end := rsp.Todos.Edges[len(rsp.Todos.Edges)-1]
			s.Require().Equal(*rsp.Todos.PageInfo.EndCursor, end.Cursor)
			if i > 0 {
				s.Require().Greater(startPriority, end.Node.PriorityOrder)
			}
			startPriority = start.Node.PriorityOrder
		}
	})
	s.Run("BackwardDescending", func() {
		var (
			rsp            response
			startCreatedAt time.Time
		)
		for i := 0; i < steps; i++ {
			err := s.Post(query, &rsp,
				client.Var("before", rsp.Todos.PageInfo.StartCursor),
				client.Var("last", step),
				client.Var("direction", "DESC"),
				client.Var("field", "CREATED_AT"),
			)
			s.Require().NoError(err)
			s.Require().Equal(maxTodos, rsp.Todos.TotalCount)
			if i < steps-1 {
				s.Require().Len(rsp.Todos.Edges, step)
				s.Require().True(rsp.Todos.PageInfo.HasPreviousPage)
			} else {
				s.Require().Len(rsp.Todos.Edges, maxTodos%step)
				s.Require().False(rsp.Todos.PageInfo.HasPreviousPage)
			}
			s.Require().True(sort.SliceIsSorted(rsp.Todos.Edges, func(i, j int) bool {
				left, _ := time.Parse(time.RFC3339, rsp.Todos.Edges[i].Node.CreatedAt)
				right, _ := time.Parse(time.RFC3339, rsp.Todos.Edges[j].Node.CreatedAt)
				return left.After(right)
			}))
			s.Require().NotNil(rsp.Todos.PageInfo.StartCursor)
			start := rsp.Todos.Edges[0]
			s.Require().Equal(*rsp.Todos.PageInfo.StartCursor, start.Cursor)
			s.Require().NotNil(rsp.Todos.PageInfo.EndCursor)
			end := rsp.Todos.Edges[len(rsp.Todos.Edges)-1]
			s.Require().Equal(*rsp.Todos.PageInfo.EndCursor, end.Cursor)
			if i > 0 {
				endCreatedAt, _ := time.Parse(time.RFC3339, end.Node.CreatedAt)
				s.Require().True(startCreatedAt.Before(endCreatedAt) || startCreatedAt.Equal(endCreatedAt))
			}
			startCreatedAt, _ = time.Parse(time.RFC3339, start.Node.CreatedAt)
		}
	})
}

func (s *todoTestSuite) TestPaginationFiltering() {
	const (
		query = `query($after: Cursor, $first: Int, $before: Cursor, $last: Int, $status: TodoStatus, $hasParent: Boolean, $hasCategory: Boolean) {
			todos(after: $after, first: $first, before: $before, last: $last, where: {status: $status, hasParent: $hasParent, hasCategory: $hasCategory}) {
				totalCount
				edges {
					node {
						id
						parent {
							id
						}
					}
					cursor
				}
				pageInfo {
					hasNextPage
					hasPreviousPage
					startCursor
					endCursor
				}
			}
		}`
		step  = 5
		steps = maxTodos/step + 1
	)
	s.Run("StatusInProgress", func() {
		var rsp response
		err := s.Post(query, &rsp,
			client.Var("first", step),
			client.Var("status", todo.StatusInProgress),
		)
		s.NoError(err)
		s.Zero(rsp.Todos.TotalCount)
	})
	s.Run("StatusCompleted", func() {
		var rsp response
		for i := 0; i < steps; i++ {
			err := s.Post(query, &rsp,
				client.Var("after", rsp.Todos.PageInfo.EndCursor),
				client.Var("first", step),
				client.Var("status", todo.StatusCompleted),
			)
			s.Require().NoError(err)
			s.Require().Equal(maxTodos, rsp.Todos.TotalCount)
			if i < steps-1 {
				s.Require().Len(rsp.Todos.Edges, step)
				s.Require().True(rsp.Todos.PageInfo.HasNextPage)
			} else {
				s.Require().Len(rsp.Todos.Edges, maxTodos%step)
				s.Require().False(rsp.Todos.PageInfo.HasNextPage)
			}
		}
	})
	s.Run("WithParent", func() {
		var rsp response
		err := s.Post(query, &rsp,
			client.Var("first", step),
			client.Var("status", todo.StatusCompleted),
			client.Var("hasParent", true),
		)
		s.Require().NoError(err)
		s.Require().Equal(maxTodos-1, rsp.Todos.TotalCount, "All todo items without the root")
	})
	s.Run("WithoutParent", func() {
		var rsp response
		err := s.Post(query, &rsp,
			client.Var("first", step),
			client.Var("status", todo.StatusCompleted),
			client.Var("hasParent", false),
		)
		s.Require().NoError(err)
		s.Require().Equal(1, rsp.Todos.TotalCount, "Only the root item")
	})
	s.Run("WithoutCategory", func() {
		var rsp response
		err := s.Post(query, &rsp,
			client.Var("first", step),
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

	s.Run("EmptyFilter", func() {
		var (
			rsp   response
			query = `query {
				todos(where:{}) {
					totalCount
				}
			}`
		)
		err := s.Post(query, &rsp)
		s.NoError(err)
		s.Equal(s.ent.Todo.Query().CountX(context.Background()), rsp.Todos.TotalCount)
	})

	s.Run("Zero first", func() {
		var (
			rsp   response
			query = `query {
				todos(first: 0) {
					totalCount
				}
			}`
		)
		err := s.Post(query, &rsp)
		s.NoError(err)
		s.Equal(s.ent.Todo.Query().CountX(context.Background()), rsp.Todos.TotalCount)
	})

	s.Run("Zero last", func() {
		var (
			rsp   response
			query = `query {
				todos(last: 0) {
					totalCount
				}
			}`
		)
		err := s.Post(query, &rsp)
		s.NoError(err)
		s.Equal(s.ent.Todo.Query().CountX(context.Background()), rsp.Todos.TotalCount)
	})
}

func (s *todoTestSuite) TestFilteringWithCustomPredicate() {
	ctx := context.Background()
	td1 := s.ent.Todo.Create().
		SetStatus(todo.StatusCompleted).
		SetText("test1").
		SetCreatedAt(time.Now().
			Add(48 * time.Hour)).
		SaveX(ctx)
	td2 := s.ent.Todo.Create().
		SetStatus(todo.StatusCompleted).
		SetText("test2").
		SetCreatedAt(time.Now().Add(-48 * time.Hour)).
		SaveX(ctx)
	td3 := s.ent.Todo.Create().
		SetStatus(todo.StatusCompleted).
		SetText("test2").
		SetCreatedAt(time.Now()).
		SaveX(ctx)
	td4 := s.ent.Todo.Create().
		SetStatus(todo.StatusCompleted).
		SetText("test3").
		SetCreatedAt(time.Now().Add(-48*time.Hour)).
		AddChildren(td1, td2, td3).
		SaveX(ctx)

	s.Run("createdToday true using interface", func() {
		var rsp struct {
			Todo struct {
				Children struct {
					TotalCount int
				}
			}
		}
		err := s.Post(`query($id: ID!, $createdToday: Boolean) {
			todo: node(id: $id) {
				... on Todo {
					children (where: {createdToday: $createdToday}) {
						totalCount
					}
				}
			}
		}`, &rsp,
			client.Var("id", td4.ID),
			client.Var("createdToday", true),
		)
		s.NoError(err)
		s.Equal(1, rsp.Todo.Children.TotalCount)
	})

	s.Run("createdToday false using interface", func() {
		var rsp struct {
			Todo struct {
				Children struct {
					TotalCount int
				}
			}
		}
		err := s.Post(`query($id: ID!, $createdToday: Boolean) {
			todo: node(id: $id) {
				... on Todo {
					children (where: {createdToday: $createdToday}) {
						totalCount
					}
				}
			}
		}`, &rsp,
			client.Var("id", td4.ID),
			client.Var("createdToday", false),
		)
		s.NoError(err)
		s.Equal(2, rsp.Todo.Children.TotalCount)
	})

	s.Run("createdToday true", func() {
		var rsp response
		err := s.Post(`query($createdToday: Boolean) {
			todos(where: {createdToday: $createdToday}) {
				totalCount
			}
		}`, &rsp,
			client.Var("createdToday", true),
		)
		s.NoError(err)
		s.Equal(maxTodos+1, rsp.Todos.TotalCount)
	})

	s.Run("createdToday false", func() {
		var rsp response
		err := s.Post(`query($createdToday: Boolean) {
			todos(where: {createdToday: $createdToday}) {
				totalCount
			}
		}`, &rsp,
			client.Var("createdToday", false),
		)
		s.NoError(err)
		s.Equal(3, rsp.Todos.TotalCount)
	})

	s.Run("not createdToday true", func() {
		var rsp response
		err := s.Post(`query($createdToday: Boolean) {
			todos(where: {not:{createdToday: $createdToday}}) {
				totalCount
			}
		}`, &rsp,
			client.Var("createdToday", true),
		)
		s.NoError(err)
		s.Equal(3, rsp.Todos.TotalCount)
	})

	s.Run("not createdToday false", func() {
		var rsp response
		err := s.Post(`query($createdToday: Boolean) {
			todos(where: {not:{createdToday: $createdToday}}) {
				totalCount
			}
		}`, &rsp,
			client.Var("createdToday", false),
		)
		s.NoError(err)
		s.Equal(maxTodos+1, rsp.Todos.TotalCount)
	})

	s.Run("or createdToday", func() {
		var rsp response
		err := s.Post(`query($createdToday1: Boolean, $createdToday2: Boolean) {
			todos(where: {or:[{createdToday: $createdToday1}, {createdToday: $createdToday2}]}) {
				totalCount
			}
		}`, &rsp,
			client.Var("createdToday1", true),
			client.Var("createdToday2", false),
		)
		s.NoError(err)
		s.Equal(maxTodos+4, rsp.Todos.TotalCount)
	})

	s.Run("and createdToday", func() {
		var rsp response
		err := s.Post(`query($createdToday1: Boolean, $createdToday2: Boolean) {
			todos(where: {and:[{createdToday: $createdToday1}, {createdToday: $createdToday2}]}) {
				totalCount
			}
		}`, &rsp,
			client.Var("createdToday1", true),
			client.Var("createdToday2", false),
		)
		s.NoError(err)
		s.Equal(0, rsp.Todos.TotalCount)
	})
}

func (s *todoTestSuite) TestNode() {
	const (
		query = `query($id: ID!) {
			todo: node(id: $id) {
				... on Todo {
					priorityOrder
				}
			}
		}`
	)
	var rsp struct{ Todo struct{ PriorityOrder int } }
	err := s.Post(query, &rsp, client.Var("id", idOffset+maxTodos))
	s.Require().NoError(err)
	err = s.Post(query, &rsp, client.Var("id", -1))
	var jerr client.RawJsonError
	s.Require().True(errors.As(err, &jerr))
	var errs gqlerror.List
	err = json.Unmarshal(jerr.RawMessage, &errs)
	s.Require().NoError(err)
	s.Require().Len(errs, 1)
	s.Require().Equal("Could not resolve to a node with the global id of '-1'", errs[0].Message)
	s.Require().Equal("NOT_FOUND", errs[0].Extensions["code"])
}

func (s *todoTestSuite) TestNodes() {
	const (
		query = `query($ids: [ID!]!) {
			todos: nodes(ids: $ids) {
				... on Todo {
					text
				}
			}
		}`
	)
	var rsp struct{ Todos []*struct{ Text string } }
	ids := []int{1, 2, 3, 3, 3, maxTodos + 1, 2, 2, maxTodos + 5}
	for i := range ids {
		ids[i] = idOffset + ids[i]
	}
	err := s.Post(query, &rsp, client.Var("ids", ids))
	s.Require().Error(err)
	s.Require().Len(rsp.Todos, len(ids))
	errmsgs := make([]string, 0, 2)
	for i, id := range ids {
		if id <= idOffset+maxTodos {
			s.Require().Equal(strconv.Itoa(id), rsp.Todos[i].Text)
		} else {
			s.Require().Nil(rsp.Todos[i])
			errmsgs = append(errmsgs,
				fmt.Sprintf("Could not resolve to a node with the global id of '%d'", id),
			)
		}
	}

	var jerr client.RawJsonError
	s.Require().True(errors.As(err, &jerr))
	var errs gqlerror.List
	err = json.Unmarshal(jerr.RawMessage, &errs)
	s.Require().NoError(err)
	s.Require().Len(errs, len(errmsgs))
	for i, err := range errs {
		s.Require().Equal(errmsgs[i], err.Message)
		s.Require().Equal("NOT_FOUND", err.Extensions["code"])
	}
}

func (s *todoTestSuite) TestNodeCollection() {
	const (
		query = `query($id: ID!) {
			todo: node(id: $id) {
				... on Todo {
					parent {
						text
						parent {
							text
						}
					}
					children {
						edges {
							node {
								text
								children {
									edges {
										node {
											text
										}
									}
								}
							}
						}
					}
				}
			}
		}`
	)
	var rsp struct {
		Todo struct {
			Parent *struct {
				Text   string
				Parent *struct {
					Text string
				}
			}
			Children struct {
				Edges []struct {
					Node struct {
						Text     string
						Children struct {
							Edges []struct {
								Node struct {
									Text string
								}
							}
						}
					}
				}
			}
		}
	}
	err := s.Post(query, &rsp, client.Var("id", idOffset+1))
	s.Require().NoError(err)
	s.Require().Nil(rsp.Todo.Parent)
	s.Require().Len(rsp.Todo.Children.Edges, maxTodos/2+1)
	s.Require().Condition(func() bool {
		for _, child := range rsp.Todo.Children.Edges {
			if child.Node.Text == strconv.Itoa(idOffset+3) {
				s.Require().Len(child.Node.Children.Edges, 1)
				s.Require().Equal(strconv.Itoa(idOffset+5), child.Node.Children.Edges[0].Node.Text)
				return true
			}
		}
		return false
	})

	err = s.Post(query, &rsp, client.Var("id", idOffset+4))
	s.Require().NoError(err)
	s.Require().NotNil(rsp.Todo.Parent)
	s.Require().Equal(strconv.Itoa(idOffset+1), rsp.Todo.Parent.Text)
	s.Require().Empty(rsp.Todo.Children.Edges)

	err = s.Post(query, &rsp, client.Var("id", strconv.Itoa(idOffset+5)))
	s.Require().NoError(err)
	s.Require().NotNil(rsp.Todo.Parent)
	s.Require().Equal(strconv.Itoa(idOffset+3), rsp.Todo.Parent.Text)
	s.Require().NotNil(rsp.Todo.Parent.Parent)
	s.Require().Equal(strconv.Itoa(idOffset+1), rsp.Todo.Parent.Parent.Text)
	s.Require().Len(rsp.Todo.Children.Edges, 1)
	s.Require().Equal(strconv.Itoa(idOffset+7), rsp.Todo.Children.Edges[0].Node.Text)
}

func (s *todoTestSuite) TestConnCollection() {
	const (
		query = `query {
			todos {
				edges {
					node {
						id
						parent {
							id
						}
						children {
							edges {
								node {
									id
								}
							}
						}
					}
				}
			}
		}`
	)
	var rsp struct {
		Todos struct {
			Edges []struct {
				Node struct {
					ID     string
					Parent *struct {
						ID string
					}
					Children struct {
						Edges []struct {
							Node struct {
								ID string
							}
						}
					}
				}
			}
		}
	}

	err := s.Post(query, &rsp)
	s.Require().NoError(err)
	s.Require().Len(rsp.Todos.Edges, maxTodos)

	for i, edge := range rsp.Todos.Edges {
		switch {
		case i == 0:
			s.Require().Nil(edge.Node.Parent)
			s.Require().Len(edge.Node.Children.Edges, maxTodos/2+1)
		case i%2 == 0:
			s.Require().NotNil(edge.Node.Parent)
			id, err := strconv.Atoi(edge.Node.Parent.ID)
			s.Require().NoError(err)
			s.Require().Equal(idOffset+i-1, id)
			if i < len(rsp.Todos.Edges)-2 {
				s.Require().Len(edge.Node.Children.Edges, 1)
			} else {
				s.Require().Empty(edge.Node.Children.Edges)
			}
		case i%2 != 0:
			s.Require().NotNil(edge.Node.Parent)
			s.Require().Equal(strconv.Itoa(idOffset+1), edge.Node.Parent.ID)
			s.Require().Empty(edge.Node.Children.Edges)
		}
	}
}

func (s *todoTestSuite) TestEnumEncoding() {
	s.Run("Encode", func() {
		const status = todo.StatusCompleted
		s.Require().Implements((*graphql.Marshaler)(nil), status)
		var b strings.Builder
		status.MarshalGQL(&b)
		str := b.String()
		const quote = `"`
		s.Require().Equal(quote, str[:1])
		s.Require().Equal(quote, str[len(str)-1:])
		str = str[1 : len(str)-1]
		s.Require().EqualValues(status, str)
	})
	s.Run("Decode", func() {
		const want = todo.StatusInProgress
		var got todo.Status
		s.Require().Implements((*graphql.Unmarshaler)(nil), &got)
		err := got.UnmarshalGQL(want.String())
		s.Require().NoError(err)
		s.Require().Equal(want, got)
	})
}

func (s *todoTestSuite) TestNodeOptions() {
	ctx := context.Background()
	td := s.ent.Todo.Create().SetText("text").SetStatus(todo.StatusInProgress).SaveX(ctx)

	nr, err := s.ent.Noder(ctx, td.ID)
	s.Require().NoError(err)
	s.Require().IsType(nr, (*ent.Todo)(nil))
	s.Require().Equal(td.ID, nr.(*ent.Todo).ID)

	nr, err = s.ent.Noder(ctx, td.ID, ent.WithFixedNodeType(todo.Table))
	s.Require().NoError(err)
	s.Require().IsType(nr, (*ent.Todo)(nil))
	s.Require().Equal(td.ID, nr.(*ent.Todo).ID)

	_, err = s.ent.Noder(ctx, td.ID, ent.WithNodeType(func(context.Context, int) (string, error) {
		return "", errors.New("bad node type")
	}))
	s.Require().EqualError(err, "bad node type")
}

func (s *todoTestSuite) TestMutationFieldCollection() {
	var rsp struct {
		CreateTodo struct {
			Text   string
			Parent struct {
				ID   string
				Text string
			}
		}
	}
	err := s.Post(`mutation ($parentID: ID!) {
		createTodo(input: { status: IN_PROGRESS, priority: 0, text: "OKE", parentID: $parentID }) {
			parent {
				id
				text
			}
			text
		}
	}`, &rsp, client.Var("parentID", strconv.Itoa(idOffset+1)))
	s.Require().NoError(err)
	s.Require().Equal("OKE", rsp.CreateTodo.Text)
	s.Require().Equal(strconv.Itoa(idOffset+1), rsp.CreateTodo.Parent.ID)
	s.Require().Equal(strconv.Itoa(idOffset+1), rsp.CreateTodo.Parent.Text)
}

func (s *todoTestSuite) TestQueryJSONFields() {
	var (
		ctx = context.Background()
		cat = s.ent.Category.Create().SetText("Disabled").SetStatus(category.StatusDisabled).SetStrings([]string{"a", "b"}).SetText("category").SaveX(ctx)
		rsp struct {
			Node struct {
				Text    string
				Strings []string
			}
		}
	)
	err := s.Post(`query node($id: ID!) {
	    node(id: $id) {
	    	... on Category {
				text
				strings
			}
		}
	}`, &rsp, client.Var("id", cat.ID))
	s.Require().NoError(err)
	s.Require().Equal(cat.Text, rsp.Node.Text)
	s.Require().Equal(cat.Strings, rsp.Node.Strings)
}

// TestInterfaceEagerLoad verifies that edges accessed through a GQL interface
// (BookmarkItem) are properly eager-loaded when queried through the interface.
func TestPageInfo(t *testing.T) {
	ctx := context.Background()
	ec := enttest.Open(
		t, dialect.SQLite,
		fmt.Sprintf("file:%s?mode=memory&cache=shared&_fk=1", t.Name()),
		enttest.WithMigrateOptions(migrate.WithGlobalUniqueID(true)),
	)
	for i := 1; i <= 5; i++ {
		ec.Todo.Create().SetText(strconv.Itoa(i)).SetStatus(todo.StatusInProgress).SaveX(ctx)
	}

	var (
		srv   = handler.NewDefaultServer(gen.NewSchema(ec))
		gqlc  = client.New(srv)
		query = `query ($after: Cursor, $first: Int, $before: Cursor, $last: Int $direction: OrderDirection!, $field: TodoOrderField!) {
			todos(after: $after, first: $first, before: $before, last: $last, orderBy: { direction: $direction, field: $field }) {
				edges {
					cursor
					node {
						text
					}
				}
				pageInfo {
					startCursor
					endCursor
					hasNextPage
					hasPreviousPage
				}
				totalCount
			}
		}`
		rsp struct {
			Todos struct {
				TotalCount int
				Edges      []struct {
					Cursor string
					Node   struct {
						Text string
					}
				}
				PageInfo struct {
					HasNextPage     bool
					HasPreviousPage bool
					StartCursor     *string
					EndCursor       *string
				}
			}
		}
		ascOrder  = []client.Option{client.Var("direction", "ASC"), client.Var("field", "TEXT")}
		descOrder = []client.Option{client.Var("direction", "DESC"), client.Var("field", "TEXT")}
		texts     = func() (s []string) {
			for _, n := range rsp.Todos.Edges {
				s = append(s, n.Node.Text)
			}
			return
		}
	)

	err := gqlc.Post(query, &rsp, ascOrder...)
	require.NoError(t, err)
	require.Equal(t, []string{"1", "2", "3", "4", "5"}, texts())
	require.Equal(t, 5, rsp.Todos.TotalCount)
	require.False(t, rsp.Todos.PageInfo.HasNextPage)
	require.False(t, rsp.Todos.PageInfo.HasPreviousPage)

	err = gqlc.Post(query, &rsp, append(ascOrder, client.Var("first", 2))...)
	require.NoError(t, err)
	require.Equal(t, []string{"1", "2"}, texts())
	require.Equal(t, 5, rsp.Todos.TotalCount)
	require.True(t, rsp.Todos.PageInfo.HasNextPage)
	require.False(t, rsp.Todos.PageInfo.HasPreviousPage)
	require.Equal(t, rsp.Todos.Edges[0].Cursor, *rsp.Todos.PageInfo.StartCursor)
	require.Equal(t, rsp.Todos.Edges[1].Cursor, *rsp.Todos.PageInfo.EndCursor)

	err = gqlc.Post(query, &rsp, append(ascOrder, client.Var("first", 2), client.Var("after", rsp.Todos.PageInfo.EndCursor))...)
	require.NoError(t, err)
	require.Equal(t, []string{"3", "4"}, texts())
	require.Equal(t, 5, rsp.Todos.TotalCount)
	require.True(t, rsp.Todos.PageInfo.HasNextPage)
	require.True(t, rsp.Todos.PageInfo.HasPreviousPage)

	err = gqlc.Post(query, &rsp, append(ascOrder, client.Var("first", 2), client.Var("after", rsp.Todos.PageInfo.EndCursor))...)
	require.NoError(t, err)
	require.Equal(t, []string{"5"}, texts())
	require.Equal(t, 5, rsp.Todos.TotalCount)
	require.False(t, rsp.Todos.PageInfo.HasNextPage)
	require.True(t, rsp.Todos.PageInfo.HasPreviousPage)

	err = gqlc.Post(query, &rsp, append(ascOrder, client.Var("last", 2), client.Var("before", rsp.Todos.PageInfo.EndCursor))...)
	require.NoError(t, err)
	require.Equal(t, []string{"3", "4"}, texts())
	require.Equal(t, 5, rsp.Todos.TotalCount)
	require.True(t, rsp.Todos.PageInfo.HasNextPage)
	require.True(t, rsp.Todos.PageInfo.HasPreviousPage)

	err = gqlc.Post(query, &rsp, append(ascOrder, client.Var("last", 2), client.Var("before", rsp.Todos.PageInfo.StartCursor))...)
	require.NoError(t, err)
	require.Equal(t, []string{"1", "2"}, texts())
	require.Equal(t, 5, rsp.Todos.TotalCount)
	require.True(t, rsp.Todos.PageInfo.HasNextPage)
	require.False(t, rsp.Todos.PageInfo.HasPreviousPage)

	err = gqlc.Post(query, &rsp, descOrder...)
	require.NoError(t, err)
	require.Equal(t, []string{"5", "4", "3", "2", "1"}, texts())
	require.Equal(t, 5, rsp.Todos.TotalCount)
	require.False(t, rsp.Todos.PageInfo.HasNextPage)
	require.False(t, rsp.Todos.PageInfo.HasPreviousPage)

	err = gqlc.Post(query, &rsp, append(descOrder, client.Var("first", 2))...)
	require.NoError(t, err)
	require.Equal(t, []string{"5", "4"}, texts())
	require.Equal(t, 5, rsp.Todos.TotalCount)
	require.True(t, rsp.Todos.PageInfo.HasNextPage)
	require.False(t, rsp.Todos.PageInfo.HasPreviousPage)

	err = gqlc.Post(query, &rsp, append(descOrder, client.Var("first", 2), client.Var("after", rsp.Todos.PageInfo.EndCursor))...)
	require.NoError(t, err)
	require.Equal(t, []string{"3", "2"}, texts())
	require.Equal(t, 5, rsp.Todos.TotalCount)
	require.True(t, rsp.Todos.PageInfo.HasNextPage)
	require.True(t, rsp.Todos.PageInfo.HasPreviousPage)

	err = gqlc.Post(query, &rsp, append(descOrder, client.Var("first", 2), client.Var("after", rsp.Todos.PageInfo.EndCursor))...)
	require.NoError(t, err)
	require.Equal(t, []string{"1"}, texts())
	require.Equal(t, 5, rsp.Todos.TotalCount)
	require.False(t, rsp.Todos.PageInfo.HasNextPage)
	require.True(t, rsp.Todos.PageInfo.HasPreviousPage)

	err = gqlc.Post(query, &rsp, append(descOrder, client.Var("last", 2), client.Var("before", rsp.Todos.PageInfo.EndCursor))...)
	require.NoError(t, err)
	require.Equal(t, []string{"3", "2"}, texts())
	require.Equal(t, 5, rsp.Todos.TotalCount)
	require.True(t, rsp.Todos.PageInfo.HasNextPage)
	require.True(t, rsp.Todos.PageInfo.HasPreviousPage)

	err = gqlc.Post(query, &rsp, append(descOrder, client.Var("last", 2), client.Var("before", rsp.Todos.PageInfo.StartCursor))...)
	require.NoError(t, err)
	require.Equal(t, []string{"5", "4"}, texts())
	require.Equal(t, 5, rsp.Todos.TotalCount)
	require.True(t, rsp.Todos.PageInfo.HasNextPage)
	require.False(t, rsp.Todos.PageInfo.HasPreviousPage)
}

type queryCount struct {
	n uint64
	dialect.Driver
}

func (q *queryCount) reset()        { atomic.StoreUint64(&q.n, 0) }
func (q *queryCount) value() uint64 { return atomic.LoadUint64(&q.n) }

func (q *queryCount) Query(ctx context.Context, query string, args, v interface{}) error {
	atomic.AddUint64(&q.n, 1)
	return q.Driver.Query(ctx, query, args, v)
}

func TestNestedConnection(t *testing.T) {
	ctx := context.Background()
	drv, err := sql.Open(dialect.SQLite, fmt.Sprintf("file:%s?mode=memory&cache=shared&_fk=1", t.Name()))
	require.NoError(t, err)
	count := &queryCount{Driver: drv}
	ec := enttest.NewClient(t,
		enttest.WithOptions(ent.Driver(count)),
		enttest.WithMigrateOptions(migrate.WithGlobalUniqueID(true)),
	)
	srv := handler.NewDefaultServer(gen.NewSchema(ec))
	gqlc := client.New(srv)

	bulkG := make([]*ent.GroupCreate, 10)
	for i := range bulkG {
		bulkG[i] = ec.Group.Create().SetName(fmt.Sprintf("group-%d", i))
	}
	groups := ec.Group.CreateBulk(bulkG...).SaveX(ctx)
	bulkU := make([]*ent.UserCreate, 10)
	for i := range bulkU {
		bulkU[i] = ec.User.Create().
			SetName(fmt.Sprintf("user-%d", i)).
			AddGroups(groups[:len(groups)-i]...).
			SetRequiredMetadata(map[string]any{})
	}
	users := ec.User.CreateBulk(bulkU...).SaveX(ctx)
	users[0].Update().AddFriends(users[1:]...).SaveX(ctx) // user 0 is friends with all
	users[1].Update().AddFriends(users[2:]...).SaveX(ctx) // user 1 is friends with all

	t.Run("After Cursor", func(t *testing.T) {
		var (
			query = `query ($id: ID!, $after: Cursor) {
				 user: node(id: $id) {
					... on User {
						id
						name
						friends(after: $after) {
							totalCount
							edges {
								cursor
								node {
									id
									name
									friends {
										totalCount
										edges {
											node {
												id
												name
											}
										}
									}
								}
							}
						}
					}
				}
			}`
			rsp struct {
				User struct {
					ID      string
					Name    string
					Friends struct {
						TotalCount int
						Edges      []struct {
							Cursor string
							Node   struct {
								ID      string
								Name    string
								Friends struct {
									TotalCount int
									Edges      []struct {
										Node struct {
											ID   string
											Name string
										}
									}
								}
							}
						}
					}
				}
			}
			after any
		)
		err = gqlc.Post(query, &rsp, client.Var("id", users[0].ID), client.Var("after", after))
		require.NoError(t, err)
		require.Equal(t, 9, rsp.User.Friends.TotalCount)
		require.Len(t, rsp.User.Friends.Edges, 9, "All users are friends with user 0")
		// First friend of user 0 is user 1.
		require.Equal(t, strconv.Itoa(users[1].ID), rsp.User.Friends.Edges[0].Node.ID)
		require.Len(t, rsp.User.Friends.Edges[0].Node.Friends.Edges, 9, "All users are friends with user 1")
		// All other users have 2 friends (user 0 and user 1).
		for _, u := range rsp.User.Friends.Edges[1:] {
			require.Len(t, u.Node.Friends.Edges, 2)
		}

		// Paginate over the friends of user 0.
		n := len(rsp.User.Friends.Edges)
		for i := 0; i < n; i++ {
			err = gqlc.Post(query, &rsp, client.Var("id", users[0].ID), client.Var("after", after))
			require.NoError(t, err)
			require.Equal(t, 9, rsp.User.Friends.TotalCount)
			require.Lenf(t, rsp.User.Friends.Edges, n-i, "There are %d friends after %v", n-i, after)
			after = rsp.User.Friends.Edges[0].Cursor
		}
	})

	t.Run("TotalCount", func(t *testing.T) {
		var (
			query = `query ($first: Int) {
				users (first: $first) {
					totalCount
					edges {
						node {
							name
							groups {
								totalCount
							}
							friends {
								edges {
									node {
										name
									}
								}
							}
						}
					}
				}
			}`
			rsp struct {
				Users struct {
					TotalCount int
					Edges      []struct {
						Node struct {
							Name   string
							Groups struct {
								TotalCount int
							}
							Friends struct {
								Edges []struct {
									Node struct {
										Name string
									}
								}
							}
						}
					}
				}
			}
		)
		count.reset()
		err = gqlc.Post(query, &rsp, client.Var("first", nil))
		require.NoError(t, err)
		// One query for loading all users, and one for getting the groups of each user.
		// The totalCount of the root query can be inferred from the length of the user edges.
		require.EqualValues(t, 3, count.value())
		require.Equal(t, 10, rsp.Users.TotalCount)
		require.Equal(t, 9, len(rsp.Users.Edges[0].Node.Friends.Edges))

		for n := 1; n <= 10; n++ {
			count.reset()
			err = gqlc.Post(query, &rsp, client.Var("first", n))
			require.NoError(t, err)
			// Two queries for getting the users and their totalCount.
			// And another one for getting the totalCount of each user.
			require.EqualValues(t, 4, count.value())
			require.Equal(t, 10, rsp.Users.TotalCount)
			for i, e := range rsp.Users.Edges {
				require.Equal(t, users[i].Name, e.Node.Name)
				// Each user i, is connected to 10-i groups.
				require.Equal(t, 10-i, e.Node.Groups.TotalCount)
			}
		}
	})

	t.Run("FirstN", func(t *testing.T) {
		var (
			query = `query ($first: Int) {
				users {
					totalCount
					edges {
						node {
							name
							groups (first: $first) {
								totalCount
								edges {
									node {
										name
									}
								}
							}
						}
					}
				}
			}`
			rsp struct {
				Users struct {
					TotalCount int
					Edges      []struct {
						Node struct {
							Name   string
							Groups struct {
								TotalCount int
								Edges      []struct {
									Node struct {
										Name string
									}
								}
							}
						}
					}
				}
			}
		)
		count.reset()
		err = gqlc.Post(query, &rsp, client.Var("first", nil))
		require.NoError(t, err)
		// One for getting all users, and one for getting all groups.
		// The totalCount is derived from len(User.Edges.Groups).
		require.EqualValues(t, 2, count.value())
		require.Equal(t, 10, rsp.Users.TotalCount)

		for n := 1; n <= 10; n++ {
			count.reset()
			err = gqlc.Post(query, &rsp, client.Var("first", n))
			require.NoError(t, err)
			// One query for getting the users (totalCount is derived), and another
			// two queries for getting the groups and the totalCount of each user.
			require.EqualValues(t, 3, count.value())
			require.Equal(t, 10, rsp.Users.TotalCount)
			for i, e := range rsp.Users.Edges {
				require.Equal(t, users[i].Name, e.Node.Name)
				require.Equal(t, 10-i, e.Node.Groups.TotalCount)
				require.Len(t, e.Node.Groups.Edges, int(math.Min(float64(n), float64(10-i))))
				for j, g := range e.Node.Groups.Edges {
					require.Equal(t, groups[j].Name, g.Node.Name)
				}
			}
		}
	})

	t.Run("Paginate", func(t *testing.T) {
		var (
			query = `query ($first: Int, $after: Cursor) {
				users (first: 1) {
					totalCount
					edges {
						node {
							name
							groups (first: $first, after: $after) {
								totalCount
								edges {
									node {
										name
										users (first: 1) {
											edges {
												node {
													name
												}
											}
										}
									}
									cursor
								}
							}
						}
					}
				}
			}`
			rsp struct {
				Users struct {
					TotalCount int
					Edges      []struct {
						Node struct {
							Name   string
							Groups struct {
								TotalCount int
								Edges      []struct {
									Node struct {
										Name  string
										Users struct {
											Edges []struct {
												Node struct {
													Name string
												}
											}
										}
									}
									Cursor string
								}
							}
						}
					}
				}
			}
			after interface{}
		)
		for i := 0; i < 10; i++ {
			count.reset()
			err = gqlc.Post(query, &rsp, client.Var("first", 1), client.Var("after", after))
			require.NoError(t, err)
			require.EqualValues(t, 5, count.value())
			require.Len(t, rsp.Users.Edges, 1)
			require.Len(t, rsp.Users.Edges[0].Node.Groups.Edges, 1)
			require.Equal(t, groups[i].Name, rsp.Users.Edges[0].Node.Groups.Edges[0].Node.Name)
			require.Len(t, rsp.Users.Edges[0].Node.Groups.Edges[0].Node.Users.Edges, 1)
			require.Equal(t, users[0].Name, rsp.Users.Edges[0].Node.Groups.Edges[0].Node.Users.Edges[0].Node.Name)
			after = rsp.Users.Edges[0].Node.Groups.Edges[0].Cursor
		}
	})

	t.Run("Nodes", func(t *testing.T) {
		var (
			query = `query ($ids: [ID!]!) {
				groups: nodes(ids: $ids) {
					... on Group {
						name
						users(last: 1) {
							totalCount
							edges {
								node {
									name
								}
							}
						}
					}
				}
			}`
			rsp struct {
				Groups []struct {
					Name  string
					Users struct {
						TotalCount int
						Edges      []struct {
							Node struct {
								Name string
							}
						}
					}
				}
			}
		)
		// One query to trigger the loading of the ent_types content.
		err = gqlc.Post(query, &rsp, client.Var("ids", []int{groups[0].ID}))
		require.NoError(t, err)
		for i := 1; i <= 10; i++ {
			ids := make([]int, 0, i)
			for _, g := range groups {
				ids = append(ids, g.ID)
			}
			count.reset()
			err = gqlc.Post(query, &rsp, client.Var("ids", ids))
			require.NoError(t, err)
			require.Len(t, rsp.Groups, 10)
			for _, g := range rsp.Groups {
				require.Len(t, g.Users.Edges, 1)
			}
			require.EqualValues(t, 3, count.value())
		}
	})

	t.Run("Node-cursor", func(t *testing.T) {
		var (
			allQuery = `query ($id: ID!) {
				group: node(id: $id) {
					... on Group {
						users(last: 2) {
							edges {
								cursor
							}
						}
					}
				}
			}`
			allRsp struct {
				Group struct {
					Users struct {
						Edges []struct {
							Cursor string
						}
					}
				}
			}
		)
		err = gqlc.Post(allQuery, &allRsp, client.Var("id", groups[0].ID))
		require.NoError(t, err)
		require.True(t, len(allRsp.Group.Users.Edges) >= 2)
		lastCursor := allRsp.Group.Users.Edges[len(allRsp.Group.Users.Edges)-1].Cursor
		secondLastCursor := allRsp.Group.Users.Edges[len(allRsp.Group.Users.Edges)-2].Cursor

		var (
			query = `query ($id: ID!, $cursor: Cursor) {
				group: node(id: $id) {
					... on Group {
						name
						users(last: 1, before: $cursor) {
							totalCount
							edges {
								cursor
								node {
									name
								}
							}
						}
					}
				}
			}`
			rsp struct {
				Group struct {
					Name  string
					Users struct {
						TotalCount int
						Edges      []struct {
							Cursor string
							Node   struct {
								Name string
							}
						}
					}
				}
			}
		)
		err = gqlc.Post(query, &rsp,
			client.Var("id", groups[0].ID),
			client.Var("cursor", lastCursor),
		)
		require.NoError(t, err)
		require.Equal(t, 1, len(rsp.Group.Users.Edges))
		require.Equal(t, secondLastCursor, rsp.Group.Users.Edges[0].Cursor)
	})
}

func TestEdgesFiltering(t *testing.T) {
	ctx := context.Background()
	drv, err := sql.Open(dialect.SQLite, fmt.Sprintf("file:%s?mode=memory&cache=shared&_fk=1", t.Name()))
	require.NoError(t, err)
	count := &queryCount{Driver: drv}
	ec := enttest.NewClient(t,
		enttest.WithOptions(ent.Driver(count)),
		enttest.WithMigrateOptions(migrate.WithGlobalUniqueID(true)),
	)
	srv := handler.NewDefaultServer(gen.NewSchema(ec))
	gqlc := client.New(srv)

	root := ec.Todo.CreateBulk(
		ec.Todo.Create().SetText("t0.1").SetStatus(todo.StatusInProgress),
		ec.Todo.Create().SetText("t0.2").SetStatus(todo.StatusInProgress),
		ec.Todo.Create().SetText("t0.3").SetStatus(todo.StatusCompleted),
	).SaveX(ctx)

	child := ec.Todo.CreateBulk(
		ec.Todo.Create().SetText("t1.1").SetParent(root[0]).SetStatus(todo.StatusInProgress),
		ec.Todo.Create().SetText("t1.2").SetParent(root[0]).SetStatus(todo.StatusCompleted),
		ec.Todo.Create().SetText("t1.3").SetParent(root[0]).SetStatus(todo.StatusCompleted),
	).SaveX(ctx)

	grandchild := ec.Todo.CreateBulk(
		ec.Todo.Create().SetText("t2.1").SetParent(child[0]).SetStatus(todo.StatusCompleted),
		ec.Todo.Create().SetText("t2.2").SetParent(child[0]).SetStatus(todo.StatusCompleted),
		ec.Todo.Create().SetText("t2.3").SetParent(child[0]).SetStatus(todo.StatusInProgress),
		ec.Todo.Create().SetText("t2.4").SetParent(child[1]).SetStatus(todo.StatusInProgress),
		ec.Todo.Create().SetText("t2.5").SetParent(child[1]).SetStatus(todo.StatusInProgress),
	).SaveX(ctx)

	ec.Todo.CreateBulk(
		ec.Todo.Create().SetText("t3.1").SetParent(grandchild[0]).SetStatus(todo.StatusInProgress),
		ec.Todo.Create().SetText("t3.2").SetParent(grandchild[0]).SetStatus(todo.StatusInProgress),
		ec.Todo.Create().SetText("t3.3").SetParent(grandchild[0]).SetStatus(todo.StatusCompleted),
		ec.Todo.Create().SetText("t3.4").SetParent(grandchild[1]).SetStatus(todo.StatusInProgress),
		ec.Todo.Create().SetText("t3.5").SetParent(grandchild[1]).SetStatus(todo.StatusInProgress),
		ec.Todo.Create().SetText("t3.6").SetParent(grandchild[1]).SetStatus(todo.StatusCompleted),
		ec.Todo.Create().SetText("t3.7").SetParent(grandchild[1]).SetStatus(todo.StatusCompleted),
	).ExecX(ctx)

	query := `query todos($id: ID!, $lv2Status: TodoStatus!) {
		todos(where:{id: $id}) {
			edges {
				node {
					children(where: {statusNEQ: COMPLETED}) {
						totalCount
						edges {
							node {
								text
								children(where: {statusNEQ: $lv2Status}) {
									totalCount
									edges {
										node {
											text
											children {
												totalCount
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}`

	var rsp struct {
		Todos struct {
			Edges []struct {
				Node struct {
					Children struct {
						TotalCount int
						Edges      []struct {
							Node struct {
								Text     string
								Children struct {
									TotalCount int
									Edges      []struct {
										Node struct {
											Text     string
											Children struct {
												TotalCount int
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	t.Run("query level 2 NEQ IN_PROGRESS", func(t *testing.T) {
		count.reset()
		err := gqlc.Post(query, &rsp, client.Var("id", root[0].ID), client.Var("lv2Status", "IN_PROGRESS"))
		require.NoError(t, err)

		require.Equal(t, 1, rsp.Todos.Edges[0].Node.Children.TotalCount)
		require.Equal(t, child[0].Text, rsp.Todos.Edges[0].Node.Children.Edges[0].Node.Text)

		n := rsp.Todos.Edges[0].Node.Children.Edges[0].Node
		require.Equal(t, 2, n.Children.TotalCount)
		require.Equal(t, grandchild[0].Text, n.Children.Edges[0].Node.Text)
		require.Equal(t, 3, n.Children.Edges[0].Node.Children.TotalCount)
		require.Equal(t, grandchild[1].Text, n.Children.Edges[1].Node.Text)
		require.Equal(t, 4, n.Children.Edges[1].Node.Children.TotalCount)

		// Top-level todos, children, grand-children and totalCount of great-children.
		require.EqualValues(t, 4, count.n)
	})

	t.Run("query level 2 NEQ COMPLETED", func(t *testing.T) {
		count.reset()
		err := gqlc.Post(query, &rsp, client.Var("id", root[0].ID), client.Var("lv2Status", "COMPLETED"))
		require.NoError(t, err)

		require.Equal(t, 1, rsp.Todos.Edges[0].Node.Children.TotalCount)
		require.Equal(t, child[0].Text, rsp.Todos.Edges[0].Node.Children.Edges[0].Node.Text)

		n := rsp.Todos.Edges[0].Node.Children.Edges[0].Node
		require.Equal(t, 1, n.Children.TotalCount)
		require.Equal(t, grandchild[2].Text, n.Children.Edges[0].Node.Text)
		require.Zero(t, n.Children.Edges[0].Node.Children.TotalCount)

		// Top-level todos, children, grand-children and totalCount of great-children.
		require.EqualValues(t, 4, count.n)
	})

	query = `query todos($id: ID!) {
		todos(where:{id: $id}) {
			edges {
				node {
					completed: children(where: {status: COMPLETED}) {
						totalCount
						edges {
							node {
								text
							}
						}
					}
					inProgress: children(where: {statusNEQ: COMPLETED}) {
						totalCount
						edges {
							node {
								text
							}
						}
					}
				}
			}
		}
	}`

	var rsp2 struct {
		Todos struct {
			Edges []struct {
				Node struct {
					Completed struct {
						TotalCount int
						Edges      []struct {
							Node struct {
								Text string
							}
						}
					}
					InProgress struct {
						TotalCount int
						Edges      []struct {
							Node struct {
								Text string
							}
						}
					}
				}
			}
		}
	}

	t.Run("filter connection with alias", func(t *testing.T) {
		count.reset()
		err := gqlc.Post(query, &rsp2, client.Var("id", root[0].ID), client.Var("lv2Status", "IN_PROGRESS"))
		require.NoError(t, err)

		require.Equal(t, 2, rsp2.Todos.Edges[0].Node.Completed.TotalCount)
		n := rsp2.Todos.Edges[0].Node.Completed
		require.Equal(t, "t1.2", n.Edges[0].Node.Text)
		require.Equal(t, "t1.3", n.Edges[1].Node.Text)

		require.Equal(t, 1, rsp2.Todos.Edges[0].Node.InProgress.TotalCount)
		n = rsp2.Todos.Edges[0].Node.InProgress
		require.Equal(t, "t1.1", n.Edges[0].Node.Text)
	})
}

func TestMutation_CreateCategory(t *testing.T) {
	ec := enttest.Open(t, dialect.SQLite,
		fmt.Sprintf("file:%s?mode=memory&cache=shared&_fk=1", t.Name()),
		enttest.WithMigrateOptions(migrate.WithGlobalUniqueID(true)),
	)
	srv := handler.NewDefaultServer(gen.NewSchema(ec))
	srv.Use(entgql.Transactioner{TxOpener: ec})
	gqlc := client.New(srv)

	// Create a category.
	var rsp struct {
		CreateCategory struct {
			ID     string
			Text   string
			Status string
			Todos  struct {
				TotalCount int
				Edges      []struct {
					Node struct {
						Text   string
						Status string
					}
				}
			}
		}
	}
	err := gqlc.Post(`
	mutation createCategory {
		createCategory(input: {
			text: "cate1"
			status: ENABLED
			createTodos: [
				{ status: IN_PROGRESS, text: "c1.t1" },
				{ status: IN_PROGRESS, text: "c1.t2" }
			]
		}) {
			id
			text
			status
			todos {
				totalCount
				edges {
					node {
						text
						status
					}
				}
			}
		}
	}
	`, &rsp)
	require.NoError(t, err)

	require.Equal(t, 2, rsp.CreateCategory.Todos.TotalCount)
	n := rsp.CreateCategory.Todos
	require.Equal(t, "c1.t1", n.Edges[0].Node.Text)
	require.Equal(t, "c1.t2", n.Edges[1].Node.Text)
}

func TestMutation_ClearChildren(t *testing.T) {
	ec := enttest.Open(t, dialect.SQLite,
		fmt.Sprintf("file:%s?mode=memory&cache=shared&_fk=1", t.Name()),
		enttest.WithMigrateOptions(migrate.WithGlobalUniqueID(true)),
	)
	srv := handler.NewDefaultServer(gen.NewSchema(ec))
	srv.Use(entgql.Transactioner{TxOpener: ec})
	gqlc := client.New(srv)

	ctx := context.Background()
	root := ec.Todo.Create().SetText("t0.1").SetStatus(todo.StatusInProgress).SaveX(ctx)
	ec.Todo.CreateBulk(
		ec.Todo.Create().SetText("t1.1").SetParent(root).SetStatus(todo.StatusInProgress),
		ec.Todo.Create().SetText("t1.2").SetParent(root).SetStatus(todo.StatusCompleted),
		ec.Todo.Create().SetText("t1.3").SetParent(root).SetStatus(todo.StatusCompleted),
	).ExecX(ctx)
	require.True(t, root.QueryChildren().ExistX(ctx))

	var rsp struct {
		UpdateTodo struct {
			ID string
		}
	}
	err := gqlc.Post(`
	mutation cleanChildren($id: ID!){
		updateTodo(id: $id, input: {clearChildren: true}) {
			id
		}
	}
	`, &rsp, client.Var("id", root.ID))
	require.NoError(t, err)
	require.Equal(t, strconv.Itoa(root.ID), rsp.UpdateTodo.ID)
	require.False(t, root.QueryChildren().ExistX(ctx))
}

func TestMutation_ClearFriend(t *testing.T) {
	ec := enttest.Open(t, dialect.SQLite,
		fmt.Sprintf("file:%s?mode=memory&cache=shared&_fk=1", t.Name()),
		enttest.WithMigrateOptions(migrate.WithGlobalUniqueID(true)),
	)
	srv := handler.NewDefaultServer(gen.NewSchema(ec))
	srv.Use(entgql.Transactioner{TxOpener: ec})
	gqlc := client.New(srv)

	ctx := context.Background()
	user := ec.User.Create().SetRequiredMetadata(map[string]any{}).SaveX(ctx)
	friend := ec.User.Create().SetRequiredMetadata(map[string]any{}).AddFriends(user).SaveX(ctx)
	friendship := user.QueryFriendships().FirstX(ctx)

	require.True(t, user.QueryFriends().ExistX(ctx))
	require.True(t, friend.QueryFriends().ExistX(ctx))

	var rsp struct {
		UpdateFriendship struct {
			ID string
		}
	}
	err := gqlc.Post(`
	mutation clearFriend($id: ID!){
		updateFriendship(id: $id, input: {clearFriend: true}) {
			id
		}
	}
	`, &rsp, client.Var("id", friendship.ID))

	require.ErrorContains(t, err, "\\\"clearFriend\\\" is not defined by type \\\"UpdateFriendshipInput\\\"")
}

func TestDescendingIDs(t *testing.T) {
	ctx := context.Background()
	ec := enttest.Open(t, dialect.SQLite,
		fmt.Sprintf("file:%s?mode=memory&_fk=1", t.Name()),
		enttest.WithMigrateOptions(migrate.WithGlobalUniqueID(true)),
	)
	srv := handler.NewDefaultServer(gen.NewSchema(ec))
	srv.Use(entgql.Transactioner{TxOpener: ec})
	ec.Category.CreateBulk(
		ec.Category.Create().SetID(1).SetText("c1").SetStatus(category.StatusEnabled),
		ec.Category.Create().SetID(2).SetText("c2").SetStatus(category.StatusEnabled),
		ec.Category.Create().SetID(3).SetText("c3").SetStatus(category.StatusEnabled),
	).SaveX(ctx)

	var (
		gqlc = client.New(srv)
		// language=GraphQL
		query = `query ($after: Cursor){
		  categories(orderBy: [{direction: DESC, field: ID}], first: 2, after: $after) {
		    edges {
		      node {
		        id
		      }
		      cursor
		    }
		  }
		}`
		rsp struct {
			Categories struct {
				Edges []struct {
					Node struct {
						ID string
					}
					Cursor string
				}
			}
		}
		after any
	)
	err := gqlc.Post(query, &rsp, client.Var("after", after))
	require.NoError(t, err)
	require.Equal(t, "3", rsp.Categories.Edges[0].Node.ID)
	require.Equal(t, "2", rsp.Categories.Edges[1].Node.ID)
	after = rsp.Categories.Edges[1].Cursor

	err = gqlc.Post(query, &rsp, client.Var("after", after))
	require.NoError(t, err)
	require.Len(t, rsp.Categories.Edges, 1)
	require.Equal(t, "1", rsp.Categories.Edges[0].Node.ID)
	after = rsp.Categories.Edges[0].Cursor

	err = gqlc.Post(query, &rsp, client.Var("after", after))
	require.NoError(t, err)
	require.Empty(t, rsp.Categories.Edges)
}

func TestMultiFieldsOrder(t *testing.T) {
	ctx := context.Background()
	ec := enttest.Open(t, dialect.SQLite,
		fmt.Sprintf("file:%s?mode=memory&_fk=1", t.Name()),
		enttest.WithMigrateOptions(migrate.WithGlobalUniqueID(true)),
	)
	srv := handler.NewDefaultServer(gen.NewSchema(ec))
	srv.Use(entgql.Transactioner{TxOpener: ec})
	cats := ec.Category.CreateBulk(
		ec.Category.Create().SetID(1).SetText("category-1").SetCount(2).SetStatus(category.StatusDisabled),
		ec.Category.Create().SetID(2).SetText("category-1").SetCount(1).SetStatus(category.StatusDisabled),
		ec.Category.Create().SetID(3).SetText("category-2").SetCount(2).SetStatus(category.StatusDisabled),
		ec.Category.Create().SetID(4).SetText("category-2").SetCount(1).SetStatus(category.StatusDisabled),
		ec.Category.Create().SetID(5).SetText("category-3").SetCount(2).SetStatus(category.StatusDisabled),
		ec.Category.Create().SetID(6).SetText("category-3").SetCount(1).SetStatus(category.StatusDisabled),
	).SaveX(ctx)
	require.Len(t, cats, 6)

	var (
		gqlc = client.New(srv)
		// language=GraphQL
		query = `query ($after: Cursor, $before: Cursor, $first: Int, $last: Int, $textDirection: OrderDirection = ASC, $countDirection: OrderDirection = ASC){
		  categories(
			after: $after,
			before: $before,
			first: $first,
			last: $last,
		  	orderBy: [{field: TEXT, direction: $textDirection}, {field: COUNT, direction: $countDirection}],
		  ) {
		    edges {
		      node {
		        id
		      }
		      cursor
		    }
		    pageInfo {
		      hasNextPage
		      hasPreviousPage
		      startCursor
		      endCursor
		    }
		  }
		}`
		rsp struct {
			Categories struct {
				Edges []struct {
					Node struct {
						ID string
					}
					Cursor string
				}
				PageInfo struct {
					HasNextPage     bool
					HasPreviousPage bool
					StartCursor     *string
					EndCursor       *string
				}
			}
		}
		ids = func() (s []string) {
			for _, n := range rsp.Categories.Edges {
				s = append(s, n.Node.ID)
			}
			return
		}
		after, before any
	)
	t.Run("CountAsc", func(t *testing.T) {
		err := gqlc.Post(query, &rsp, client.Var("after", after), client.Var("before", before), client.Var("first", 3))
		require.NoError(t, err)
		require.Len(t, rsp.Categories.Edges, 3)
		// cats[1], cats[0], cats[3].
		require.Equal(t, []string{"2", "1", "4"}, ids())
		require.True(t, rsp.Categories.PageInfo.HasNextPage)

		after = rsp.Categories.PageInfo.EndCursor
		err = gqlc.Post(query, &rsp, client.Var("after", after), client.Var("first", 3))
		require.NoError(t, err)
		require.Len(t, rsp.Categories.Edges, 3)
		// cats[2], cats[5], cats[4].
		require.Equal(t, []string{"3", "6", "5"}, ids())
		require.False(t, rsp.Categories.PageInfo.HasNextPage)

		after = rsp.Categories.PageInfo.EndCursor
		err = gqlc.Post(query, &rsp, client.Var("after", after))
		require.NoError(t, err)
		require.Empty(t, rsp.Categories.Edges)
		require.False(t, rsp.Categories.PageInfo.HasNextPage)
	})

	t.Run("CountDesc", func(t *testing.T) {
		err := gqlc.Post(query, &rsp, client.Var("countDirection", "DESC"), client.Var("first", 3))
		require.NoError(t, err)
		require.Len(t, rsp.Categories.Edges, 3)
		// cats[0], cats[1], cats[2].
		require.Equal(t, []string{"1", "2", "3"}, ids())
		require.True(t, rsp.Categories.PageInfo.HasNextPage)

		after = rsp.Categories.PageInfo.EndCursor
		err = gqlc.Post(query, &rsp, client.Var("after", after), client.Var("countDirection", "DESC"), client.Var("first", 3))
		require.NoError(t, err)
		require.Len(t, rsp.Categories.Edges, 3)
		// cats[3], cats[4], cats[5].
		require.Equal(t, []string{"4", "5", "6"}, ids())
		require.False(t, rsp.Categories.PageInfo.HasNextPage)

		after = rsp.Categories.PageInfo.EndCursor
		err = gqlc.Post(query, &rsp, client.Var("after", after), client.Var("countDirection", "DESC"))
		require.NoError(t, err)
		require.Empty(t, rsp.Categories.Edges)
		require.False(t, rsp.Categories.PageInfo.HasNextPage)
	})

	t.Run("TextCountDesc", func(t *testing.T) {
		// Page forward.
		err := gqlc.Post(query, &rsp, client.Var("textDirection", "DESC"), client.Var("countDirection", "DESC"), client.Var("first", 3))
		require.NoError(t, err)
		require.Len(t, rsp.Categories.Edges, 3)
		// cats[4], cats[5], cats[2].
		require.Equal(t, []string{"5", "6", "3"}, ids())
		require.True(t, rsp.Categories.PageInfo.HasNextPage)

		after = rsp.Categories.PageInfo.EndCursor
		err = gqlc.Post(query, &rsp, client.Var("after", after), client.Var("textDirection", "DESC"), client.Var("countDirection", "DESC"), client.Var("first", 3))
		require.NoError(t, err)
		require.Len(t, rsp.Categories.Edges, 3)
		// cats[3], cats[0], cats[1].
		require.Equal(t, []string{"4", "1", "2"}, ids())
		require.False(t, rsp.Categories.PageInfo.HasNextPage)

		after = rsp.Categories.PageInfo.EndCursor
		err = gqlc.Post(query, &rsp, client.Var("after", after), client.Var("textDirection", "DESC"), client.Var("countDirection", "DESC"))
		require.NoError(t, err)
		require.Empty(t, rsp.Categories.Edges)
		require.False(t, rsp.Categories.PageInfo.HasNextPage)

		// All categories.
		err = gqlc.Post(query, &rsp, client.Var("textDirection", "DESC"), client.Var("countDirection", "DESC"))
		require.NoError(t, err)
		require.Len(t, rsp.Categories.Edges, 6)
		// cats[4], cats[5], cats[2], cats[3], cats[0], cats[1].
		require.Equal(t, []string{"5", "6", "3", "4", "1", "2"}, ids())
		require.False(t, rsp.Categories.PageInfo.HasNextPage)

		// Page backward.
		err = gqlc.Post(query, &rsp, client.Var("textDirection", "DESC"), client.Var("countDirection", "DESC"), client.Var("last", 3))
		require.NoError(t, err)
		require.Len(t, rsp.Categories.Edges, 3)
		// cats[3], cats[0], cats[1].
		require.Equal(t, []string{"4", "1", "2"}, ids())
		require.True(t, rsp.Categories.PageInfo.HasPreviousPage)

		before = rsp.Categories.PageInfo.StartCursor
		err = gqlc.Post(query, &rsp, client.Var("before", before), client.Var("textDirection", "DESC"), client.Var("countDirection", "DESC"), client.Var("last", 3))
		require.NoError(t, err)
		require.Len(t, rsp.Categories.Edges, 3)
		// cats[4], cats[5], cats[2].
		require.Equal(t, []string{"5", "6", "3"}, ids())
		require.False(t, rsp.Categories.PageInfo.HasPreviousPage)
	})
}

type queryRecorder struct {
	queries []string
	dialect.Driver
}

func (r *queryRecorder) reset() {
	r.queries = nil
}

func (r *queryRecorder) Query(ctx context.Context, query string, args, v interface{}) error {
	r.queries = append(r.queries, query)
	return r.Driver.Query(ctx, query, args, v)
}

func TestReduceQueryComplexity(t *testing.T) {
	ctx := context.Background()
	drv, err := sql.Open(dialect.SQLite, fmt.Sprintf("file:%s?mode=memory&cache=shared&_fk=1", t.Name()))
	require.NoError(t, err)
	rec := &queryRecorder{Driver: drv}
	ec := enttest.NewClient(t,
		enttest.WithOptions(ent.Driver(rec)),
		enttest.WithMigrateOptions(migrate.WithGlobalUniqueID(true)),
	)
	var (
		// language=GraphQL
		query = `query Todo($id: ID!) {
			node(id: $id) {
				... on Todo {
					text
					children (first: 10) {
						edges {
							node {
								text
							}
						}
					}
				}
			}
		}`
		gqlc = client.New(handler.NewDefaultServer(gen.NewSchema(ec)))
	)
	t1 := ec.Todo.Create().SetText("t1").SetStatus(todo.StatusInProgress).SaveX(ctx)
	rec.reset()
	require.NoError(t, gqlc.Post(query, new(any), client.Var("id", t1.ID)))
	require.Equal(t, []string{
		// Node mapping (cached).
		"SELECT `type` FROM `ent_types` ORDER BY `id` ASC",
		// Top-level todo.
		"SELECT `todos`.`id`, `todos`.`text` FROM `todos` WHERE `todos`.`id` = ? LIMIT 2",
		// Children todos (without CTE).
		"SELECT `todos`.`id`, `todos`.`text`, `todos`.`project_todos`, `todos`.`todo_children`, `todos`.`todo_secret` FROM `todos` WHERE `todos`.`todo_children` IN (?) ORDER BY `todos`.`id` LIMIT 11",
	}, rec.queries)

	// language=GraphQL
	query = `query Todos($ids: [ID!]!) {
		todos: nodes (ids: $ids) {
			... on Todo {
				text
				children (first: 10) {
					edges {
						node {
							text
						}
					}
				}
			}
		}
	}`
	rec.reset()
	require.NoError(t, gqlc.Post(query, new(any), client.Var("ids", []int{t1.ID})))
	// A single ID is implemented by the `node` query.
	require.Equal(t, []string{
		// Top-level todo.
		"SELECT `todos`.`id`, `todos`.`text` FROM `todos` WHERE `todos`.`id` = ? LIMIT 2",
		// Children todos (without CTE).
		"SELECT `todos`.`id`, `todos`.`text`, `todos`.`project_todos`, `todos`.`todo_children`, `todos`.`todo_secret` FROM `todos` WHERE `todos`.`todo_children` IN (?) ORDER BY `todos`.`id` LIMIT 11",
	}, rec.queries)

	rec.reset()
	require.NoError(t, gqlc.Post(query, new(any), client.Var("ids", []int{t1.ID, t1.ID})))
	require.Equal(t, []string{
		// Top-level todo.
		"SELECT `todos`.`id`, `todos`.`text` FROM `todos` WHERE `todos`.`id` IN (?, ?)",
		// Children todos (with CTE).
		"WITH `src_query` AS (SELECT `todos`.`id`, `todos`.`text`, `todos`.`project_todos`, `todos`.`todo_children`, `todos`.`todo_secret` FROM `todos` WHERE `todos`.`todo_children` IN (?)), `limited_query` AS (SELECT *, (ROW_NUMBER() OVER (PARTITION BY `todo_children` ORDER BY `id` ASC)) AS `row_number` FROM `src_query`) SELECT `id`, `text`, `project_todos`, `todo_children`, `todo_secret` FROM `limited_query` AS `todos` WHERE `todos`.`row_number` <= ?",
	}, rec.queries)

	// Propagate uniqueness to one-child edges.
	// language=GraphQL
	query = `query Todo($id: ID!) {
			node(id: $id) {
				... on Todo {
					parent {
						text
						children (first: 5) {
							edges {
								node {
									text
								}
							}
						}
					}
					category {
						text
						todos (first: 10) {
							edges {
								node {
									text
								}
							}
						}
					}
				}
			}
		}`
	ec.Todo.Create().SetText("t0").SetStatus(todo.StatusInProgress).AddChildren(t1).SaveX(ctx)
	ec.Category.Create().AddTodos(t1).SetText("c0").SetStatus(category.StatusEnabled).SaveX(ctx)
	rec.reset()
	require.NoError(t, gqlc.Post(query, new(any), client.Var("id", t1.ID)))
	require.Equal(t, []string{
		// Top-level todo.
		"SELECT `todos`.`id`, `todos`.`category_id`, `todos`.`project_todos`, `todos`.`todo_children`, `todos`.`todo_secret` FROM `todos` WHERE `todos`.`id` = ? LIMIT 2",
		// Parent todo.
		"SELECT `todos`.`id`, `todos`.`text` FROM `todos` WHERE `todos`.`id` IN (?)",
		// Parent children.
		"SELECT `todos`.`id`, `todos`.`text`, `todos`.`project_todos`, `todos`.`todo_children`, `todos`.`todo_secret` FROM `todos` WHERE `todos`.`todo_children` IN (?) ORDER BY `todos`.`id` LIMIT 6",
		// Category.
		"SELECT `categories`.`id`, `categories`.`text` FROM `categories` WHERE `categories`.`id` IN (?)",
		// Category todos.
		"SELECT `todos`.`id`, `todos`.`text`, `todos`.`category_id`, `todos`.`project_todos`, `todos`.`todo_children`, `todos`.`todo_secret` FROM `todos` WHERE `todos`.`category_id` IN (?) ORDER BY `todos`.`id` LIMIT 11",
	}, rec.queries)

	// Same as above, but with multiple IDs.
	// language=GraphQL
	query = `query Todo($id: ID!) {
			nodes(ids: [$id, $id]) {
				... on Todo {
					parent {
						text
						children (first: 5) {
							edges {
								node {
									text
								}
							}
						}
					}
					category {
						text
						todos (first: 10) {
							edges {
								node {
									text
								}
							}
						}
					}
				}
			}
		}`
	rec.reset()
	require.NoError(t, gqlc.Post(query, new(any), client.Var("id", t1.ID)))
	require.Equal(t, []string{
		// Root nodes.
		"SELECT `todos`.`id`, `todos`.`category_id`, `todos`.`project_todos`, `todos`.`todo_children`, `todos`.`todo_secret` FROM `todos` WHERE `todos`.`id` IN (?, ?)",
		// Their parents (2 max).
		"SELECT `todos`.`id`, `todos`.`text` FROM `todos` WHERE `todos`.`id` IN (?)",
		// 5 children for each parent.
		"WITH `src_query` AS (SELECT `todos`.`id`, `todos`.`text`, `todos`.`project_todos`, `todos`.`todo_children`, `todos`.`todo_secret` FROM `todos` WHERE `todos`.`todo_children` IN (?)), `limited_query` AS (SELECT *, (ROW_NUMBER() OVER (PARTITION BY `todo_children` ORDER BY `id` ASC)) AS `row_number` FROM `src_query`) SELECT `id`, `text`, `project_todos`, `todo_children`, `todo_secret` FROM `limited_query` AS `todos` WHERE `todos`.`row_number` <= ?",
		// Category.
		"SELECT `categories`.`id`, `categories`.`text` FROM `categories` WHERE `categories`.`id` IN (?)",
		// 10 todos for each category.
		"WITH `src_query` AS (SELECT `todos`.`id`, `todos`.`text`, `todos`.`category_id`, `todos`.`project_todos`, `todos`.`todo_children`, `todos`.`todo_secret` FROM `todos` WHERE `todos`.`category_id` IN (?)), `limited_query` AS (SELECT *, (ROW_NUMBER() OVER (PARTITION BY `category_id` ORDER BY `id` ASC)) AS `row_number` FROM `src_query`) SELECT `id`, `text`, `category_id`, `project_todos`, `todo_children`, `todo_secret` FROM `limited_query` AS `todos` WHERE `todos`.`row_number` <= ?",
	}, rec.queries)
}

func TestFieldSelection(t *testing.T) {
	ctx := context.Background()
	drv, err := sql.Open(dialect.SQLite, fmt.Sprintf("file:%s?mode=memory&cache=shared&_fk=1", t.Name()))
	require.NoError(t, err)
	rec := &queryRecorder{Driver: drv}
	ec := enttest.NewClient(t,
		enttest.WithOptions(ent.Driver(rec)),
		enttest.WithMigrateOptions(migrate.WithGlobalUniqueID(true)),
	)
	root := ec.Todo.CreateBulk(
		ec.Todo.Create().SetText("t0.1").SetStatus(todo.StatusInProgress),
		ec.Todo.Create().SetText("t0.2").SetStatus(todo.StatusInProgress),
		ec.Todo.Create().SetText("t0.3").SetStatus(todo.StatusCompleted),
	).SaveX(ctx)
	ec.Todo.CreateBulk(
		ec.Todo.Create().SetText("t1.1").SetParent(root[0]).SetStatus(todo.StatusInProgress),
		ec.Todo.Create().SetText("t1.2").SetParent(root[0]).SetStatus(todo.StatusCompleted),
		ec.Todo.Create().SetText("t1.3").SetParent(root[0]).SetStatus(todo.StatusCompleted),
	).SaveX(ctx)
	var (
		// language=GraphQL
		query = `query {
			todos {
				edges {
					node {
						children {
							totalCount
							edges {
								node {
									text
								}
							}
						}
					}
				}
			}
		}`
		rsp struct {
			Todos struct {
				Edges []struct {
					Node struct {
						Children struct {
							TotalCount int
							Edges      []struct {
								Node struct {
									Text string
								}
							}
						}
					}
				}
			}
		}
		gqlc = client.New(handler.NewDefaultServer(gen.NewSchema(ec)))
	)
	rec.reset()
	gqlc.MustPost(query, &rsp)
	require.Equal(t, []string{
		// No fields were selected besides the "id" field.
		"SELECT `todos`.`id` FROM `todos` ORDER BY `todos`.`id`",
		// The "id" and the "text" fields were selected + all foreign keys (see, `withFKs` query field).
		"SELECT `todos`.`id`, `todos`.`text`, `todos`.`project_todos`, `todos`.`todo_children`, `todos`.`todo_secret` FROM `todos` WHERE `todos`.`todo_children` IN (?, ?, ?, ?, ?, ?) ORDER BY `todos`.`id`",
	}, rec.queries)

	ec.Category.CreateBulk(
		ec.Category.Create().AddTodos(root[0]).SetText("c0").SetStatus(category.StatusEnabled),
		ec.Category.Create().AddTodos(root[1]).SetText("c1").SetStatus(category.StatusEnabled),
		ec.Category.Create().AddTodos(root[2]).SetText("c2").SetStatus(category.StatusEnabled),
	).SaveX(ctx)
	var (
		// language=GraphQL
		query2 = `query {
			todos {
				edges {
					node {
						category {
							text
						}
					}
				}
			}
		}`
		rsp2 struct {
			Todos struct {
				Edges []struct {
					Node struct {
						Category struct {
							Text string
						}
					}
				}
			}
		}
	)
	rec.reset()
	client.New(handler.NewDefaultServer(gen.NewSchema(ec))).
		MustPost(query2, &rsp2)
	require.Equal(t, []string{
		// Also query the "category_id" field for the "category" selection.
		"SELECT `todos`.`id`, `todos`.`category_id` FROM `todos` ORDER BY `todos`.`id`",
		// Select the "text" field for the "category" selection.
		"SELECT `categories`.`id`, `categories`.`text` FROM `categories` WHERE `categories`.`id` IN (?, ?, ?, ?)",
	}, rec.queries)

	var (
		// language=GraphQL
		query3 = `query {
			todos {
				edges {
					node {
						__typename
						text
					}
				}
			}
		}`
		rsp3 struct {
			Todos struct {
				Edges []struct {
					Node struct {
						TypeName string `json:"__typename"`
						Text     string
					}
				}
			}
		}
	)
	rec.reset()
	client.New(handler.NewDefaultServer(gen.NewSchema(ec))).
		MustPost(query3, &rsp3)
	require.Equal(t, []string{
		// Ignore the __typename meta field.
		"SELECT `todos`.`id`, `todos`.`text` FROM `todos` ORDER BY `todos`.`id`",
	}, rec.queries)

	var (
		// language=GraphQL
		query4 = `query {
			todos {
				edges {
					node {
						text
						extendedField
					}
				}
			}
		}`
		rsp4 struct {
			Todos struct {
				Edges []struct {
					Node struct {
						Text          string
						ExtendedField string
					}
				}
			}
		}
	)
	rec.reset()
	client.New(handler.NewDefaultServer(gen.NewSchema(ec))).
		MustPost(query4, &rsp4)
	require.Equal(t, []string{
		// Unknown fields enforce query all columns.
		"SELECT `todos`.`id`, `todos`.`created_at`, `todos`.`status`, " +
			"`todos`.`priority`, `todos`.`text`, `todos`.`name`, `todos`.`blob`, " +
			"`todos`.`category_id`, `todos`.`init`, `todos`.`custom`, " +
			"`todos`.`customp`, `todos`.`value` FROM `todos` ORDER BY `todos`.`id`",
	}, rec.queries)

	rootO2M := ec.OneToMany.CreateBulk(
		ec.OneToMany.Create().SetName("t0.1"),
		ec.OneToMany.Create().SetName("t0.2"),
		ec.OneToMany.Create().SetName("t0.3"),
	).SaveX(ctx)
	ec.OneToMany.CreateBulk(
		ec.OneToMany.Create().SetName("t1.1").SetParent(rootO2M[0]),
		ec.OneToMany.Create().SetName("t1.2").SetParent(rootO2M[0]),
		ec.OneToMany.Create().SetName("t1.3").SetParent(rootO2M[0]),
	).SaveX(ctx)
	var (
		// language=GraphQL
		queryO2M = `query {
			oneToMany {
				edges {
					node {
						id
						name
						children {
							name
						}
					}
				}
			}
		}`
		rspO2M struct {
			OneToMany struct {
				Edges []struct {
					Node struct {
						ID       string
						Name     string
						Children []struct {
							ID   string
							Name string
						}
					}
				}
			}
		}
		gqlcO2M = client.New(handler.NewDefaultServer(gen.NewSchema(ec)))
	)
	rec.reset()
	gqlcO2M.MustPost(queryO2M, &rspO2M)
	require.Equal(t, []string{
		"SELECT `one_to_manies`.`id`, `one_to_manies`.`name` FROM `one_to_manies` ORDER BY `one_to_manies`.`id`",
		"SELECT `one_to_manies`.`id`, `one_to_manies`.`name`, `one_to_manies`.`parent_id` FROM `one_to_manies` WHERE `one_to_manies`.`parent_id` IN (?, ?, ?, ?, ?, ?)",
	}, rec.queries)

	// Test the CollectedFor annotation.
	ec.Todo.Create().SetText("test").SetName("hello").SetStatus(todo.StatusCompleted).SaveX(ctx)
	var (
		// language=GraphQL
		query5 = `query {
			todos (where: {name: "hello"}) {
				edges {
					node {
						uppercaseName
					}
				}
			}
		}`
		rsp5 struct {
			Todos struct {
				Edges []struct {
					Node struct {
						UppercaseName *string
					}
				}
			}
		}
	)
	rec.reset()
	client.New(handler.NewDefaultServer(gen.NewSchema(ec))).
		MustPost(query5, &rsp5)
	require.Equal(t, []string{
		"SELECT `todos`.`id`, `todos`.`name` FROM `todos` WHERE `todos`.`name` = ? ORDER BY `todos`.`id`",
	}, rec.queries)
	require.Len(t, rsp5.Todos.Edges, 1)
	require.NotNil(t, rsp5.Todos.Edges[0].Node.UppercaseName)
	require.Equal(t, "HELLO", *rsp5.Todos.Edges[0].Node.UppercaseName)
}

func TestCollectedForMultipleFields(t *testing.T) {
	ctx := context.Background()
	drv, err := sql.Open(dialect.SQLite, fmt.Sprintf("file:%s?mode=memory&cache=shared&_fk=1", t.Name()))
	require.NoError(t, err)
	rec := &queryRecorder{Driver: drv}
	ec := enttest.NewClient(t,
		enttest.WithOptions(ent.Driver(rec)),
		enttest.WithMigrateOptions(migrate.WithGlobalUniqueID(true)),
	)
	ec.User.Create().
		SetFirstName("John").
		SetLastName("Doe").
		SetRequiredMetadata(map[string]any{}).
		SaveX(ctx)

	var (
		// language=GraphQL
		query = `query {
			users {
				edges {
					node {
						fullName
					}
				}
			}
		}`
		rsp struct {
			Users struct {
				Edges []struct {
					Node struct {
						FullName *string
					}
				}
			}
		}
	)
	rec.reset()
	client.New(handler.NewDefaultServer(gen.NewSchema(ec))).
		MustPost(query, &rsp)
	require.Equal(t, []string{
		"SELECT `users`.`id`, `users`.`first_name`, `users`.`last_name` FROM `users` ORDER BY `users`.`id`",
	}, rec.queries)
	require.Len(t, rsp.Users.Edges, 1)
	require.NotNil(t, rsp.Users.Edges[0].Node.FullName)
	require.Equal(t, "John Doe", *rsp.Users.Edges[0].Node.FullName)
}

func TestOrderByEdgeCount(t *testing.T) {
	ctx := context.Background()
	ec := enttest.Open(
		t, dialect.SQLite,
		fmt.Sprintf("file:%s?mode=memory&cache=shared&_fk=1", t.Name()),
		enttest.WithMigrateOptions(migrate.WithGlobalUniqueID(true)),
	)
	gqlc := client.New(handler.NewDefaultServer(gen.NewSchema(ec)))
	cats := ec.Category.CreateBulk(
		ec.Category.Create().SetText("parents").SetStatus(category.StatusEnabled),
		ec.Category.Create().SetText("children").SetStatus(category.StatusEnabled),
	).SaveX(ctx)
	root := ec.Todo.CreateBulk(
		ec.Todo.Create().SetText("t0.1").SetStatus(todo.StatusPending).SetCategory(cats[0]),
		ec.Todo.Create().SetText("t0.2").SetStatus(todo.StatusInProgress).SetCategory(cats[0]),
		ec.Todo.Create().SetText("t0.3").SetStatus(todo.StatusCompleted).SetCategory(cats[0]),
	).SaveX(ctx)
	ec.Todo.CreateBulk(
		ec.Todo.Create().SetText("t1.1").SetParent(root[0]).SetStatus(todo.StatusInProgress).SetCategory(cats[1]),
		ec.Todo.Create().SetText("t1.2").SetParent(root[0]).SetStatus(todo.StatusCompleted).SetCategory(cats[1]),
		ec.Todo.Create().SetText("t1.3").SetParent(root[0]).SetStatus(todo.StatusCompleted).SetCategory(cats[1]),
		ec.Todo.Create().SetText("t2.1").SetParent(root[1]).SetStatus(todo.StatusInProgress).SetCategory(cats[1]),
		ec.Todo.Create().SetText("t2.2").SetParent(root[1]).SetStatus(todo.StatusCompleted).SetCategory(cats[1]),
		ec.Todo.Create().SetText("t3.1").SetParent(root[2]).SetStatus(todo.StatusInProgress).SetCategory(cats[1]),
	).SaveX(ctx)

	t.Run("ChildrenCount", func(t *testing.T) {
		var (
			// language=GraphQL
			query = `query TodosByChildCount($direction: OrderDirection = ASC){
				todos(
					# Filter only those with children.
					where: {hasChildren: true},
					orderBy: {field: CHILDREN_COUNT, direction: $direction},
				) {
					edges {
						node {
							id
						}
					}
				}
			}`
			rsp struct {
				Todos struct {
					Edges []struct {
						Node struct {
							ID string
						}
					}
				}
			}
		)
		gqlc.MustPost(query, &rsp, client.Var("direction", "DESC"))
		require.Len(t, rsp.Todos.Edges, 3)
		for i, r := range root {
			require.Equal(t, rsp.Todos.Edges[i].Node.ID, strconv.Itoa(r.ID))
		}
		gqlc.MustPost(query, &rsp, client.Var("direction", "ASC"))
		require.Len(t, rsp.Todos.Edges, 3)
		for i, r := range root {
			require.Equal(t, rsp.Todos.Edges[len(rsp.Todos.Edges)-i-1].Node.ID, strconv.Itoa(r.ID))
		}
	})

	t.Run("NestedEdgeCountOrdering", func(t *testing.T) {
		var (
			// language=GraphQL
			query = `query CategoryByTodosCount {
				categories(
					orderBy: {field: TODOS_COUNT, direction: DESC},
				) {
					edges {
						node {
							id
							todos(orderBy: {field: CHILDREN_COUNT, direction: DESC}) {
								edges {
									node {
										id
									}
								}
							}
						}
					}
				}
			}`
			rsp struct {
				Categories struct {
					Edges []struct {
						Node struct {
							ID    string
							Todos struct {
								Edges []struct {
									Node struct {
										ID string
									}
								}
							}
						}
					}
				}
			}
		)
		gqlc.MustPost(query, &rsp)
		require.Len(t, rsp.Categories.Edges, 2)
		childC, parentC := rsp.Categories.Edges[0].Node, rsp.Categories.Edges[1].Node
		// Second categories holds todos without children.
		require.Equal(t, childC.ID, strconv.Itoa(cats[1].ID))
		require.Len(t, childC.Todos.Edges, 6)
		// First categories holds parent todos.
		require.Equal(t, parentC.ID, strconv.Itoa(cats[0].ID))
		require.Len(t, parentC.Todos.Edges, 3)
		for i, r := range root {
			require.Equal(t, parentC.Todos.Edges[i].Node.ID, strconv.Itoa(r.ID))
		}
	})

	t.Run("EdgeFieldOrdering", func(t *testing.T) {
		var (
			// language=GraphQL
			query = `query TodosByParentStatus($direction: OrderDirection = ASC) {
				todos(
					# Filter out parent todos.
					where: {hasParent: true},
					orderBy: {field: PARENT_STATUS, direction: $direction},
				) {
					edges {
						node {
							parent {
								status
							}
						}
					}
				}
			}`
			rsp struct {
				Todos struct {
					Edges []struct {
						Node struct {
							Parent struct {
								Status todo.Status
							}
						}
					}
				}
			}
			expected = []todo.Status{
				todo.StatusCompleted, todo.StatusInProgress, todo.StatusInProgress,
				todo.StatusPending, todo.StatusPending, todo.StatusPending,
			}
		)
		gqlc.MustPost(query, &rsp, client.Var("direction", "ASC"))
		require.Len(t, rsp.Todos.Edges, 6)
		for i, p := range rsp.Todos.Edges {
			require.Equal(t, expected[i], p.Node.Parent.Status)
		}
		// Reverse the order.
		gqlc.MustPost(query, &rsp, client.Var("direction", "DESC"))
		require.Len(t, rsp.Todos.Edges, 6)
		for i, p := range rsp.Todos.Edges {
			require.Equal(t, expected[len(expected)-i-1], p.Node.Parent.Status)
		}
	})

	t.Run("ExposeOrderField", func(t *testing.T) {
		var (
			// language=GraphQL
			query = `query CategoryByTodosCount {
				categories(
					orderBy: {field: TODOS_COUNT, direction: DESC},
				) {
					edges {
						node {
							todosCount
						}
					}
				}
			}`
			rsp struct {
				Categories struct {
					Edges []struct {
						Node struct {
							TodosCount int
						}
					}
				}
			}
		)
		gqlc.MustPost(query, &rsp)
		require.Len(t, rsp.Categories.Edges, 2)
		require.Equal(t, rsp.Categories.Edges[0].Node.TodosCount, 6)
		require.Equal(t, rsp.Categories.Edges[1].Node.TodosCount, 3)
	})
}

func TestSatisfiesFragments(t *testing.T) {
	ctx := context.Background()
	ec := enttest.Open(
		t, dialect.SQLite,
		fmt.Sprintf("file:%s?mode=memory&cache=shared&_fk=1", t.Name()),
		enttest.WithMigrateOptions(migrate.WithGlobalUniqueID(true)),
	)
	gqlc := client.New(handler.NewDefaultServer(gen.NewSchema(ec)))
	cat := ec.Category.Create().SetText("cat").SetStatus(category.StatusEnabled).SaveX(ctx)
	todos := ec.Todo.CreateBulk(
		ec.Todo.Create().SetText("t1").SetStatus(todo.StatusPending).SetCategory(cat),
		ec.Todo.Create().SetText("t2").SetStatus(todo.StatusInProgress).SetCategory(cat),
		ec.Todo.Create().SetText("t3").SetStatus(todo.StatusCompleted).SetCategory(cat),
	).SaveX(ctx)
	var (
		// language=GraphQL
		query = `query CategoryTodo($id: ID!) {
		  category: node(id: $id) {
		    __typename
		    id
		    ... on Category {
			  text
			  ...CategoryTodos
		    }
		  }
		}

		fragment CategoryTodos on Category {
		  todos (orderBy: {field: TEXT}) {
		    edges {
		      node {
		        id
				...TodoFields
		      }
		    }
		  }
		}

		fragment TodoFields on Todo {
		  id
		  text
		  createdAt
		}
		`
		rsp struct {
			Category struct {
				TypeName string `json:"__typename"`
				ID, Text string
				Todos    struct {
					Edges []struct {
						Node struct {
							ID, Text, CreatedAt string
						}
					}
				}
			}
		}
	)
	gqlc.MustPost(query, &rsp, client.Var("id", cat.ID))
	require.Equal(t, strconv.Itoa(cat.ID), rsp.Category.ID)
	require.Len(t, rsp.Category.Todos.Edges, 3)
	for i := range todos {
		require.Equal(t, strconv.Itoa(todos[i].ID), rsp.Category.Todos.Edges[i].Node.ID)
		require.Equal(t, todos[i].Text, rsp.Category.Todos.Edges[i].Node.Text)
		ts, err := todos[i].CreatedAt.MarshalText()
		require.NoError(t, err)
		require.Equal(t, string(ts), rsp.Category.Todos.Edges[i].Node.CreatedAt)
	}
}

func TestSatisfiesDeeperFragments(t *testing.T) {
	ctx := context.Background()
	ec := enttest.Open(
		t, dialect.SQLite,
		fmt.Sprintf("file:%s?mode=memory&cache=shared&_fk=1", t.Name()),
		enttest.WithMigrateOptions(migrate.WithGlobalUniqueID(true)),
	)
	gqlc := client.New(handler.NewDefaultServer(gen.NewSchema(ec)))
	cat := ec.Category.Create().SetText("cat").SetStatus(category.StatusEnabled).SaveX(ctx)
	todos := ec.Todo.CreateBulk(
		ec.Todo.Create().SetText("t1").SetStatus(todo.StatusPending).SetCategory(cat),
		ec.Todo.Create().SetText("t2").SetStatus(todo.StatusInProgress).SetCategory(cat),
		ec.Todo.Create().SetText("t3").SetStatus(todo.StatusCompleted).SetCategory(cat),
	).SaveX(ctx)
	var (
		// language=GraphQL
		query = `query Node($id: ID!) {
			todo: node(id: $id) {
				__typename
				... on Todo {
					... MainFra
				}
				id
			}
		}

		fragment MainFra on Todo {
			...Child1
			id
			category {
				id
			}
		}

		fragment Child2 on Category {
			id
			text
		}

		fragment Child1 on Todo {
			text
			category {
				id
				... Child2
			}
		}`
		rsp struct {
			Todo struct {
				TypeName string `json:"__typename"`
				ID, Text string
				Category struct {
					ID, Text string
				}
			}
		}
	)

	gqlc.MustPost(query, &rsp, client.Var("id", todos[0].ID))
	require.Equal(t, "cat", cat.Text)
	require.Equal(t, "cat", rsp.Todo.Category.Text)
}

func TestRenamedType(t *testing.T) {
	ctx := context.Background()
	ec := enttest.Open(
		t, dialect.SQLite,
		fmt.Sprintf("file:%s?mode=memory&cache=shared&_fk=1", t.Name()),
		enttest.WithMigrateOptions(migrate.WithGlobalUniqueID(true)),
	)
	gqlc := client.New(handler.NewDefaultServer(gen.NewSchema(ec)))
	wr := ec.Workspace.Create().SetName("Ariga").SaveX(ctx)
	var (
		// language=GraphQL
		query = `query Node($id: ID!) {
			text: node(id: $id) {
				id
				... on Organization {
					name
				}
			}
		}`
		rsp struct {
			Text struct {
				ID, Name string
			}
		}
	)
	gqlc.MustPost(query, &rsp, client.Var("id", wr.ID))
	require.Equal(t, "Ariga", rsp.Text.Name)
}

func TestSatisfiesNodeFragments(t *testing.T) {
	ctx := context.Background()
	ec := enttest.Open(
		t, dialect.SQLite,
		fmt.Sprintf("file:%s?mode=memory&cache=shared&_fk=1", t.Name()),
		enttest.WithMigrateOptions(migrate.WithGlobalUniqueID(true)),
	)
	gqlc := client.New(handler.NewDefaultServer(gen.NewSchema(ec)))
	t1 := ec.Todo.Create().SetText("t1").SetStatus(todo.StatusPending).SaveX(ctx)
	var (
		// language=GraphQL
		query = `query Node($id: ID!) {
			todo: node(id: $id) {
				id
				...NodeFragment
			}
		}
		fragment NodeFragment on Node {
			... on Todo {
				createdAt
				status
				text
			}
		}`
		rsp struct {
			Todo struct {
				ID, Text, CreatedAt string
				Status              todo.Status
			}
		}
	)
	gqlc.MustPost(query, &rsp, client.Var("id", t1.ID))
	require.Equal(t, strconv.Itoa(t1.ID), rsp.Todo.ID)
	require.Equal(t, "t1", rsp.Todo.Text)
	require.NotEmpty(t, rsp.Todo.Status)
	require.NotEmpty(t, rsp.Todo.CreatedAt)

	g1 := ec.Group.Create().SetName("g1").SaveX(ctx)
	var (
		// language=GraphQL
		query1 = `query Node($id: ID!) {
			group: node(id: $id) {
				id
				...NamedNodeFragment
			}
		}
		fragment NamedNodeFragment on NamedNode {
			... on Group {
				name
			}
		}`
		rsp1 struct {
			Group struct {
				ID, Name string
			}
		}
	)
	gqlc.MustPost(query1, &rsp1, client.Var("id", g1.ID))
	require.Equal(t, strconv.Itoa(g1.ID), rsp1.Group.ID)
	require.Equal(t, "g1", rsp1.Group.Name)
}

func TestUnionConnectionEagerLoadFragmentOrder(t *testing.T) {
	ctx := context.Background()
	drv, err := sql.Open(dialect.SQLite, fmt.Sprintf("file:%s?mode=memory&cache=shared&_fk=1", t.Name()))
	require.NoError(t, err)
	rec := &queryRecorder{Driver: drv}
	ec := enttest.NewClient(t,
		enttest.WithOptions(ent.Driver(rec)),
		enttest.WithMigrateOptions(migrate.WithGlobalUniqueID(true)),
	)
	gqlc := client.New(handler.NewDefaultServer(gen.NewSchema(ec)))

	dataQueries := func(queries []string) []string {
		out := make([]string, 0, len(queries))
		for _, q := range queries {
			if strings.Contains(q, "ent_types") {
				continue
			}
			out = append(out, q)
		}
		return out
	}

	cat := ec.Category.Create().SetText("cat").SetStatus(category.StatusEnabled).SaveX(ctx)
	const todoCount = 5
	for i := 0; i < todoCount; i++ {
		ec.Todo.Create().
			SetText(fmt.Sprintf("t%d", i)).
			SetStatus(todo.StatusInProgress).
			SetCategory(cat).
			SaveX(ctx)
	}

	run := func(projectFirst bool) ([]string, []string) {
		var fragments string
		if projectFirst {
			fragments = `... on Project {
				todos(first: 10) {
					edges {
						node {
							text
						}
					}
				}
			}
			... on Category {
				todos(first: 10) {
					edges {
						node {
							category {
								text
							}
						}
					}
				}
			}`
		} else {
			fragments = `... on Category {
				todos(first: 10) {
					edges {
						node {
							category {
								text
							}
						}
					}
				}
			}
			... on Project {
				todos(first: 10) {
					edges {
						node {
							text
						}
					}
				}
			}`
		}
		query := fmt.Sprintf(`query Node($id: ID!) {
			category: node(id: $id) {
				%s
			}
		}`, fragments)
		var rsp struct {
			Category struct {
				Todos struct {
					Edges []struct {
						Node struct {
							Category struct {
								Text string
							}
						}
					}
				}
			}
		}
		rec.reset()
		gqlc.MustPost(query, &rsp, client.Var("id", cat.ID))
		texts := make([]string, len(rsp.Category.Todos.Edges))
		for i, e := range rsp.Category.Todos.Edges {
			texts[i] = e.Node.Category.Text
		}
		return dataQueries(rec.queries), texts
	}

	const wantQueries = 3 // category node + batched todos + batched category edges.

	queries1, texts1 := run(true)
	queries2, texts2 := run(false)
	require.Equal(t, wantQueries, len(queries1), "expected eager-load queries, got: %v", queries1)
	require.Equal(t, wantQueries, len(queries2), "expected eager-load queries, got: %v", queries2)
	require.Equal(t, queries1, queries2, "queries must not depend on fragment order")
	require.Equal(t, []string{"cat", "cat", "cat", "cat", "cat"}, texts1)
	require.Equal(t, texts1, texts2)
	require.Contains(t, queries1[1], "category_id` IN")
	require.Contains(t, queries1[2], "categories`.`id` IN")
}

// TestInterfaceConnectionPagination verifies that the polymorphic BookmarkItem
// connection exposed on Category.items is a real Relay connection: it accepts
// first/after/before/last, returns unique per-edge cursors, and reports a
// correct pageInfo/totalCount while walking through heterogeneous members.
// TestInterfaceViewConnection verifies the global `bookmarkItems` connection, which
// is backed by the BookmarkItemView SQL view. It paginates the view at the database
// level and resolves each row to its concrete node (Todo/Project) so that
// type-specific fragments work.
func TestInterfaceViewConnection(t *testing.T) {
	ctx := context.Background()
	drv, err := sql.Open(dialect.SQLite,
		fmt.Sprintf("file:%s?mode=memory&cache=shared&_fk=1", t.Name()))
	require.NoError(t, err)
	ec := enttest.NewClient(t,
		enttest.WithOptions(ent.Driver(drv)),
		// Global unique ids let the view rows resolve back to their concrete type.
		enttest.WithMigrateOptions(migrate.WithGlobalUniqueID(true)),
	)
	// The built-in migrator does not create views, so create it manually. The
	// definition mirrors the BookmarkItemView schema's entsql.ViewFor query.
	_, err = drv.DB().ExecContext(ctx,
		"CREATE VIEW bookmark_item_views AS SELECT id, 'Todo' AS kind, category_id, text, name FROM todos UNION ALL SELECT id, 'Project' AS kind, category_projects AS category_id, text, name FROM projects")
	require.NoError(t, err)
	gqlc := client.New(handler.NewDefaultServer(gen.NewSchema(ec)))

	// 3 todos + 2 projects across two categories = 5 interface members.
	cat := ec.Category.Create().SetText("cat1").SetStatus(category.StatusEnabled).SaveX(ctx)
	ec.Todo.CreateBulk(
		ec.Todo.Create().SetText("t1").SetStatus(todo.StatusInProgress).SetCategory(cat),
		ec.Todo.Create().SetText("t2").SetStatus(todo.StatusCompleted).SetCategory(cat),
		ec.Todo.Create().SetText("t3").SetStatus(todo.StatusInProgress).SetCategory(cat),
	).SaveX(ctx)
	ec.Project.CreateBulk(
		ec.Project.Create().SetText("project").SetCategory(cat),
		ec.Project.Create().SetText("project").SetCategory(cat),
	).SaveX(ctx)

	type page struct {
		BookmarkItems struct {
			TotalCount int
			PageInfo   struct {
				HasNextPage bool
				EndCursor   *string
			}
			Edges []struct {
				Cursor string
				Node   struct {
					Typename string `json:"__typename"`
					Text     string // present only for Todo
					Owner    *struct{ Text string }
				}
			}
		}
	}
	const q = `query($first: Int, $after: Cursor) {
		bookmarkItems(first: $first, after: $after) {
			totalCount
			pageInfo { hasNextPage endCursor }
			edges {
				cursor
				node {
					__typename
					... on BookmarkItem { owner { text } }
					... on Todo { text }
				}
			}
		}
	}`

	// Full listing: all 5 members, each resolved to its concrete type.
	var all page
	gqlc.MustPost(q, &all)
	require.Equal(t, 5, all.BookmarkItems.TotalCount)
	require.Len(t, all.BookmarkItems.Edges, 5)
	types := map[string]int{}
	for _, e := range all.BookmarkItems.Edges {
		types[e.Node.Typename]++
		require.NotNil(t, e.Node.Owner)
		require.Equal(t, "cat1", e.Node.Owner.Text)
	}
	require.Equal(t, 3, types["Todo"], "expected 3 todos resolved to concrete type")
	require.Equal(t, 2, types["Project"], "expected 2 projects resolved to concrete type")

	// DB-level pagination: first 2, then continue after the end cursor.
	var p1 page
	gqlc.MustPost(q, &p1, client.Var("first", 2))
	require.Len(t, p1.BookmarkItems.Edges, 2)
	require.True(t, p1.BookmarkItems.PageInfo.HasNextPage)
	require.Equal(t, 5, p1.BookmarkItems.TotalCount)
	require.NotNil(t, p1.BookmarkItems.PageInfo.EndCursor)

	var p2 page
	gqlc.MustPost(q, &p2, client.Var("first", 2), client.Var("after", *p1.BookmarkItems.PageInfo.EndCursor))
	require.Len(t, p2.BookmarkItems.Edges, 2)
	// No overlap between pages.
	require.NotEqual(t, p1.BookmarkItems.Edges[1].Cursor, p2.BookmarkItems.Edges[0].Cursor)
}

// TestInterfaceConnectionConcreteTypes verifies that nodes of the edge-based
// BookmarkItem connection are the concrete Todo/Project types, so both the shared
// interface fragment and type-specific fragments resolve.
// TestInterfaceConnectionScoping verifies that the edge-based BookmarkItem
// connection only contains the members belonging to the queried parent, and
// that an empty connection is reported correctly.
// TestInterfaceViewConnectionReverse verifies backward pagination (last/before)
// on the global, view-backed BookmarkItem connection.
func TestInterfaceViewConnectionReverse(t *testing.T) {
	ctx := context.Background()
	drv, err := sql.Open(dialect.SQLite,
		fmt.Sprintf("file:%s?mode=memory&cache=shared&_fk=1", t.Name()))
	require.NoError(t, err)
	ec := enttest.NewClient(t,
		enttest.WithOptions(ent.Driver(drv)),
		enttest.WithMigrateOptions(migrate.WithGlobalUniqueID(true)),
	)
	_, err = drv.DB().ExecContext(ctx,
		"CREATE VIEW bookmark_item_views AS SELECT id, 'Todo' AS kind, category_id, text, name FROM todos UNION ALL SELECT id, 'Project' AS kind, category_projects AS category_id, text, name FROM projects")
	require.NoError(t, err)
	gqlc := client.New(handler.NewDefaultServer(gen.NewSchema(ec)))

	cat := ec.Category.Create().SetText("cat1").SetStatus(category.StatusEnabled).SaveX(ctx)
	ec.Todo.CreateBulk(
		ec.Todo.Create().SetText("t1").SetStatus(todo.StatusInProgress).SetCategory(cat),
		ec.Todo.Create().SetText("t2").SetStatus(todo.StatusCompleted).SetCategory(cat),
		ec.Todo.Create().SetText("t3").SetStatus(todo.StatusInProgress).SetCategory(cat),
	).SaveX(ctx)
	ec.Project.CreateBulk(
		ec.Project.Create().SetText("project").SetCategory(cat),
		ec.Project.Create().SetText("project").SetCategory(cat),
	).SaveX(ctx)

	type page struct {
		BookmarkItems struct {
			TotalCount int
			PageInfo   struct {
				HasNextPage     bool
				HasPreviousPage bool
			}
			Edges []struct {
				Cursor string
			}
		}
	}
	const q = `query($first: Int, $after: Cursor, $last: Int, $before: Cursor) {
		bookmarkItems(first: $first, after: $after, last: $last, before: $before) {
			totalCount
			pageInfo { hasNextPage hasPreviousPage }
			edges { cursor }
		}
	}`

	// Full ordered listing to derive cursors.
	var all page
	gqlc.MustPost(q, &all)
	require.Len(t, all.BookmarkItems.Edges, 5)

	// last: 2 returns the final two members and reports a previous page.
	var lastRsp page
	gqlc.MustPost(q, &lastRsp, client.Var("last", 2))
	require.Len(t, lastRsp.BookmarkItems.Edges, 2)
	require.True(t, lastRsp.BookmarkItems.PageInfo.HasPreviousPage)
	require.Equal(t, all.BookmarkItems.Edges[3].Cursor, lastRsp.BookmarkItems.Edges[0].Cursor)
	require.Equal(t, all.BookmarkItems.Edges[4].Cursor, lastRsp.BookmarkItems.Edges[1].Cursor)

	// last + before the 3rd cursor returns the two members preceding it
	// (backward traversal, as per the Relay first/after vs last/before pairing).
	var beforeRsp page
	gqlc.MustPost(q, &beforeRsp,
		client.Var("last", 2), client.Var("before", all.BookmarkItems.Edges[2].Cursor))
	require.Len(t, beforeRsp.BookmarkItems.Edges, 2)
	require.Equal(t, all.BookmarkItems.Edges[0].Cursor, beforeRsp.BookmarkItems.Edges[0].Cursor)
	require.Equal(t, all.BookmarkItems.Edges[1].Cursor, beforeRsp.BookmarkItems.Edges[1].Cursor)
}

// TestInterfaceViewConnectionOrderAndFilter verifies entgql-generated orderBy
// and where support on the global, view-backed BookmarkItem connection.
func TestInterfaceViewConnectionOrderAndFilter(t *testing.T) {
	ctx := context.Background()
	drv, err := sql.Open(dialect.SQLite,
		fmt.Sprintf("file:%s?mode=memory&cache=shared&_fk=1", t.Name()))
	require.NoError(t, err)
	ec := enttest.NewClient(t,
		enttest.WithOptions(ent.Driver(drv)),
		enttest.WithMigrateOptions(migrate.WithGlobalUniqueID(true)),
	)
	_, err = drv.DB().ExecContext(ctx,
		"CREATE VIEW bookmark_item_views AS SELECT id, 'Todo' AS kind, category_id, text, name FROM todos UNION ALL SELECT id, 'Project' AS kind, category_projects AS category_id, text, name FROM projects")
	require.NoError(t, err)
	gqlc := client.New(handler.NewDefaultServer(gen.NewSchema(ec)))

	cat1 := ec.Category.Create().SetText("cat1").SetStatus(category.StatusEnabled).SaveX(ctx)
	cat2 := ec.Category.Create().SetText("cat2").SetStatus(category.StatusEnabled).SaveX(ctx)
	// cat1: 2 todos + 1 project. cat2: 1 todo.
	ec.Todo.CreateBulk(
		ec.Todo.Create().SetText("t1").SetStatus(todo.StatusInProgress).SetCategory(cat1),
		ec.Todo.Create().SetText("t2").SetStatus(todo.StatusInProgress).SetCategory(cat1),
	).SaveX(ctx)
	ec.Project.Create().SetText("project").SetCategory(cat1).SaveX(ctx)
	ec.Todo.Create().SetText("t3").SetStatus(todo.StatusInProgress).SetCategory(cat2).SaveX(ctx)

	type page struct {
		BookmarkItems struct {
			TotalCount int
			Edges      []struct {
				Node struct {
					Owner *struct{ Text string }
				}
			}
		}
	}

	// where: filter to cat1's members only.
	var filtered page
	gqlc.MustPost(`query($cid: Int!) {
		bookmarkItems(where: {categoryID: $cid}) {
			totalCount
			edges { node { ... on BookmarkItem { owner { text } } } }
		}
	}`, &filtered, client.Var("cid", cat1.ID))
	require.Equal(t, 3, filtered.BookmarkItems.TotalCount, "cat1 has 2 todos + 1 project")
	require.Len(t, filtered.BookmarkItems.Edges, 3)
	for _, e := range filtered.BookmarkItems.Edges {
		require.Equal(t, "cat1", e.Node.Owner.Text)
	}

	// orderBy CATEGORY_ID DESC: cat2's member(s) come first.
	var ordered page
	gqlc.MustPost(`{
		bookmarkItems(orderBy: {field: CATEGORY_ID, direction: DESC}) {
			totalCount
			edges { node { ... on BookmarkItem { owner { text } } } }
		}
	}`, &ordered)
	require.Equal(t, 4, ordered.BookmarkItems.TotalCount)
	require.Len(t, ordered.BookmarkItems.Edges, 4)
	require.Equal(t, "cat2", ordered.BookmarkItems.Edges[0].Node.Owner.Text, "highest category_id first")
	require.Equal(t, "cat1", ordered.BookmarkItems.Edges[3].Node.Owner.Text)
}

// TestInterfaceViewConnectionFilterOrderPaginate exercises the generated
// filter (where), orderBy, and pagination on the global bookmarkItems connection,
// including their combination and multi-page cursor walks.
func TestInterfaceViewConnectionFilterOrderPaginate(t *testing.T) {
	ctx := context.Background()
	drv, err := sql.Open(dialect.SQLite,
		fmt.Sprintf("file:%s?mode=memory&cache=shared&_fk=1", t.Name()))
	require.NoError(t, err)
	ec := enttest.NewClient(t,
		enttest.WithOptions(ent.Driver(drv)),
		enttest.WithMigrateOptions(migrate.WithGlobalUniqueID(true)),
	)
	_, err = drv.DB().ExecContext(ctx,
		"CREATE VIEW bookmark_item_views AS SELECT id, 'Todo' AS kind, category_id, text, name FROM todos UNION ALL SELECT id, 'Project' AS kind, category_projects AS category_id, text, name FROM projects")
	require.NoError(t, err)
	gqlc := client.New(handler.NewDefaultServer(gen.NewSchema(ec)))

	// Three categories (ascending ids) with 4 todos + 2 projects = 6 members.
	//   cat1: t1, t2, p1   cat2: t3, p2   cat3: t4
	cat1 := ec.Category.Create().SetText("cat1").SetStatus(category.StatusEnabled).SaveX(ctx)
	cat2 := ec.Category.Create().SetText("cat2").SetStatus(category.StatusEnabled).SaveX(ctx)
	cat3 := ec.Category.Create().SetText("cat3").SetStatus(category.StatusEnabled).SaveX(ctx)
	ec.Todo.CreateBulk(
		ec.Todo.Create().SetText("t1").SetStatus(todo.StatusInProgress).SetCategory(cat1),
		ec.Todo.Create().SetText("t2").SetStatus(todo.StatusInProgress).SetCategory(cat1),
		ec.Todo.Create().SetText("t3").SetStatus(todo.StatusInProgress).SetCategory(cat2),
		ec.Todo.Create().SetText("t4").SetStatus(todo.StatusInProgress).SetCategory(cat3),
	).SaveX(ctx)
	ec.Project.CreateBulk(
		ec.Project.Create().SetText("p1").SetCategory(cat1),
		ec.Project.Create().SetText("p2").SetCategory(cat2),
	).SaveX(ctx)

	type node struct {
		Typename string `json:"__typename"`
		ID       string
		Text     string
		Owner    *struct{ Text string }
	}
	type page struct {
		BookmarkItems struct {
			TotalCount int
			PageInfo   struct {
				HasNextPage bool
				EndCursor   *string
			}
			Edges []struct {
				Cursor string
				Node   node
			}
		}
	}
	post := func(q string, opts ...client.Option) page {
		var p page
		gqlc.MustPost(q, &p, opts...)
		return p
	}

	t.Run("FilterByKind", func(t *testing.T) {
		todos := post(`{ bookmarkItems(where: {kind: "Todo"}) { totalCount edges { node { __typename } } } }`)
		require.Equal(t, 4, todos.BookmarkItems.TotalCount)
		require.Len(t, todos.BookmarkItems.Edges, 4)
		for _, e := range todos.BookmarkItems.Edges {
			require.Equal(t, "Todo", e.Node.Typename)
		}
		projects := post(`{ bookmarkItems(where: {kind: "Project"}) { totalCount edges { node { __typename } } } }`)
		require.Equal(t, 2, projects.BookmarkItems.TotalCount)
		for _, e := range projects.BookmarkItems.Edges {
			require.Equal(t, "Project", e.Node.Typename)
		}
	})

	t.Run("FilterByText", func(t *testing.T) {
		// hasPrefix "t" matches the 4 todos (projects are p1/p2).
		byPrefix := post(`{ bookmarkItems(where: {textHasPrefix: "t"}) { edges { node { __typename text } } } }`)
		require.Len(t, byPrefix.BookmarkItems.Edges, 4)
		for _, e := range byPrefix.BookmarkItems.Edges {
			require.Equal(t, "Todo", e.Node.Typename)
			require.True(t, strings.HasPrefix(e.Node.Text, "t"))
		}
		// exact match on a shared field.
		exact := post(`{ bookmarkItems(where: {text: "p1"}) { edges { node { __typename text } } } }`)
		require.Len(t, exact.BookmarkItems.Edges, 1)
		require.Equal(t, "Project", exact.BookmarkItems.Edges[0].Node.Typename)
		require.Equal(t, "p1", exact.BookmarkItems.Edges[0].Node.Text)
	})

	t.Run("FilterAndOrNot", func(t *testing.T) {
		// NOT Todo -> the 2 projects.
		notTodo := post(`{ bookmarkItems(where: {not: {kind: "Todo"}}) { totalCount edges { node { __typename } } } }`)
		require.Equal(t, 2, notTodo.BookmarkItems.TotalCount)
		for _, e := range notTodo.BookmarkItems.Edges {
			require.Equal(t, "Project", e.Node.Typename)
		}
		// Project OR text=t1 -> 2 projects + t1 = 3.
		orExpr := post(`{ bookmarkItems(where: {or: [{kind: "Project"}, {text: "t1"}]}) { totalCount } }`)
		require.Equal(t, 3, orExpr.BookmarkItems.TotalCount)
		// Todo AND category=cat1 -> t1, t2.
		andExpr := post(`query($cid: Int!) { bookmarkItems(where: {and: [{kind: "Todo"}, {categoryID: $cid}]}) { totalCount edges { node { text } } } }`,
			client.Var("cid", cat1.ID))
		require.Equal(t, 2, andExpr.BookmarkItems.TotalCount)
	})

	t.Run("OrderByCategoryWithPagination", func(t *testing.T) {
		// Walk all pages (first: 2) ordered by CATEGORY_ID ASC; the owner category
		// must be non-decreasing and every member must appear exactly once.
		const q = `query($first: Int, $after: Cursor) {
			bookmarkItems(first: $first, after: $after, orderBy: {field: CATEGORY_ID, direction: ASC}) {
				pageInfo { hasNextPage endCursor }
				edges { cursor node { ... on BookmarkItem { owner { text } } } }
			}
		}`
		var (
			owners  []string
			cursors = map[string]bool{}
			after   *string
		)
		for i := 0; i < 10; i++ { // bounded to avoid an infinite loop on a bug
			opts := []client.Option{client.Var("first", 2)}
			if after != nil {
				opts = append(opts, client.Var("after", *after))
			}
			p := post(q, opts...)
			require.LessOrEqual(t, len(p.BookmarkItems.Edges), 2)
			for _, e := range p.BookmarkItems.Edges {
				require.False(t, cursors[e.Cursor], "cursor %s repeated across pages", e.Cursor)
				cursors[e.Cursor] = true
				owners = append(owners, e.Node.Owner.Text)
			}
			if !p.BookmarkItems.PageInfo.HasNextPage {
				break
			}
			require.NotNil(t, p.BookmarkItems.PageInfo.EndCursor)
			after = p.BookmarkItems.PageInfo.EndCursor
		}
		require.Equal(t, []string{"cat1", "cat1", "cat1", "cat2", "cat2", "cat3"}, owners,
			"ordered by category_id ascending, paginated across pages")
	})

	t.Run("FilterOrderPaginateCombined", func(t *testing.T) {
		// Todos only, ordered by CATEGORY_ID desc, first 2 then the rest.
		const q = `query($first: Int, $after: Cursor) {
			bookmarkItems(first: $first, after: $after, where: {kind: "Todo"}, orderBy: {field: CATEGORY_ID, direction: DESC}) {
				totalCount
				pageInfo { hasNextPage endCursor }
				edges { node { __typename ... on BookmarkItem { owner { text } } } }
			}
		}`
		first := post(q, client.Var("first", 2))
		require.Equal(t, 4, first.BookmarkItems.TotalCount, "4 todos match the filter")
		require.Len(t, first.BookmarkItems.Edges, 2)
		require.True(t, first.BookmarkItems.PageInfo.HasNextPage)
		// Highest category first: cat3 (t4), then cat2 (t3).
		require.Equal(t, "cat3", first.BookmarkItems.Edges[0].Node.Owner.Text)
		require.Equal(t, "cat2", first.BookmarkItems.Edges[1].Node.Owner.Text)
		for _, e := range first.BookmarkItems.Edges {
			require.Equal(t, "Todo", e.Node.Typename)
		}

		rest := post(q, client.Var("first", 2), client.Var("after", *first.BookmarkItems.PageInfo.EndCursor))
		require.Len(t, rest.BookmarkItems.Edges, 2)
		require.False(t, rest.BookmarkItems.PageInfo.HasNextPage)
		// Remaining todos are both in cat1.
		for _, e := range rest.BookmarkItems.Edges {
			require.Equal(t, "Todo", e.Node.Typename)
			require.Equal(t, "cat1", e.Node.Owner.Text)
		}
	})
}

// TestInterfaceViewConnectionPreload verifies that resolving the global,
// view-backed BookmarkItem connection:
//   - preloads only the selected per-type fields, and
//   - does not issue N+1 queries: the view is queried once, each member type is
//     queried once (batched by id), and selected edges (owner) are eager-loaded
//     in batched queries rather than one per node.
func TestInterfaceViewConnectionPreload(t *testing.T) {
	ctx := context.Background()
	drv, err := sql.Open(dialect.SQLite,
		fmt.Sprintf("file:%s?mode=memory&cache=shared&_fk=1", t.Name()))
	require.NoError(t, err)
	rec := &queryRecorder{Driver: drv}
	ec := enttest.NewClient(t,
		enttest.WithOptions(ent.Driver(rec)),
		enttest.WithMigrateOptions(migrate.WithGlobalUniqueID(true)),
	)
	_, err = drv.DB().ExecContext(ctx,
		"CREATE VIEW bookmark_item_views AS SELECT id, 'Todo' AS kind, category_id, text, name FROM todos UNION ALL SELECT id, 'Project' AS kind, category_projects AS category_id, text, name FROM projects")
	require.NoError(t, err)
	gqlc := client.New(handler.NewDefaultServer(gen.NewSchema(ec)))

	// Fixed dataset (2 todos + 1 project) so the batched IN-queries have a
	// deterministic number of placeholders for the exact-match assertion below.
	cat := ec.Category.Create().SetText("cat1").SetStatus(category.StatusEnabled).SaveX(ctx)
	ec.Todo.CreateBulk(
		ec.Todo.Create().SetText("t1").SetStatus(todo.StatusInProgress).SetCategory(cat),
		ec.Todo.Create().SetText("t2").SetStatus(todo.StatusCompleted).SetCategory(cat),
	).SaveX(ctx)
	ec.Project.Create().SetText("p1").SetCategory(cat).SaveX(ctx)

	var rsp struct {
		BookmarkItems struct {
			Edges []struct {
				Node struct {
					Typename string `json:"__typename"`
					ID       string
					Text     string  // shared interface field
					Name     *string // shared interface field
					Owner    *struct{ Text string }
					Status   *string // Todo-only
				}
			}
		}
	}
	rec.reset()
	// The shared fields (id, text, name, owner) are queried directly on the
	// interface; status is selected via a Todo fragment.
	// language=GraphQL
	gqlc.MustPost(`{
		bookmarkItems {
			edges {
				node {
					__typename
					id
					text
					name
					owner { text }
					... on Todo { status }
				}
			}
		}
	}`, &rsp)

	// Preload: every item resolved to its concrete type with the selected fields.
	require.Len(t, rsp.BookmarkItems.Edges, 3)
	var todos, projects int
	for _, e := range rsp.BookmarkItems.Edges {
		require.NotNil(t, e.Node.Owner)
		require.Equal(t, "cat1", e.Node.Owner.Text)
		require.NotEmpty(t, e.Node.Text)
		switch e.Node.Typename {
		case "Todo":
			todos++
			require.NotNil(t, e.Node.Status)
		case "Project":
			projects++
			require.Equal(t, "p1", e.Node.Text)
		default:
			t.Fatalf("unexpected __typename %q", e.Node.Typename)
		}
	}
	require.Equal(t, 2, todos)
	require.Equal(t, 1, projects)

	// The exact set of SQL queries proves the resolution strategy:
	//   - one query against the view (ids + shared column, ordered);
	//   - one batched query per member type selecting only the requested columns
	//     (plus the FK needed for the owner edge);
	//   - one batched owner eager-load per member type.
	// No N+1 (each relation is queried once), and no count query since totalCount
	// is not selected.
	require.Equal(t, []string{
		"SELECT `bookmark_item_views`.`id`, `bookmark_item_views`.`kind`, `bookmark_item_views`.`category_id`, `bookmark_item_views`.`text`, `bookmark_item_views`.`name` FROM `bookmark_item_views` ORDER BY `bookmark_item_views`.`id`",
		"SELECT `projects`.`id`, `projects`.`text`, `projects`.`name`, `projects`.`category_projects` FROM `projects` WHERE `projects`.`id` IN (?)",
		"SELECT `categories`.`id`, `categories`.`text` FROM `categories` WHERE `categories`.`id` IN (?)",
		"SELECT `todos`.`id`, `todos`.`text`, `todos`.`name`, `todos`.`status`, `todos`.`category_id` FROM `todos` WHERE `todos`.`id` IN (?, ?)",
		"SELECT `categories`.`id`, `categories`.`text` FROM `categories` WHERE `categories`.`id` IN (?)",
	}, rec.queries)

	// `first: 0` returns an empty page (not the whole set) and issues only the
	// view query, never resolving concrete nodes.
	var zero struct {
		BookmarkItems struct {
			TotalCount int
			Edges      []struct{ Cursor string }
		}
	}
	rec.reset()
	gqlc.MustPost(`{ bookmarkItems(first: 0) { totalCount edges { cursor } } }`, &zero)
	require.Empty(t, zero.BookmarkItems.Edges, "first:0 must return no edges")
	require.Equal(t, 3, zero.BookmarkItems.TotalCount)
	require.Equal(t, []string{
		"SELECT COUNT(*) FROM `bookmark_item_views`",
	}, rec.queries, "first:0 only counts; no view rows or per-type queries")

	// A totalCount-only selection resolves no nodes (no per-type queries).
	var countOnly struct {
		BookmarkItems struct{ TotalCount int }
	}
	rec.reset()
	gqlc.MustPost(`{ bookmarkItems { totalCount } }`, &countOnly)
	require.Equal(t, 3, countOnly.BookmarkItems.TotalCount)
	require.Equal(t, []string{
		"SELECT COUNT(*) FROM `bookmark_item_views`",
	}, rec.queries, "totalCount-only must not fetch view rows or resolve nodes")

	// A selection fully covered by the view's columns (__typename + shared
	// scalars) is synthesized from the view rows, with no per-type queries.
	var synth struct {
		BookmarkItems struct {
			Edges []struct {
				Node struct {
					Typename string `json:"__typename"`
					ID       string
					Text     string
				}
			}
		}
	}
	rec.reset()
	gqlc.MustPost(`{ bookmarkItems { edges { node { __typename id text } } } }`, &synth)
	require.Len(t, synth.BookmarkItems.Edges, 3)
	for _, e := range synth.BookmarkItems.Edges {
		require.True(t, e.Node.Typename == "Todo" || e.Node.Typename == "Project", "concrete __typename")
		require.NotEmpty(t, e.Node.ID)
		require.NotEmpty(t, e.Node.Text)
	}
	require.Equal(t, []string{
		"SELECT `bookmark_item_views`.`id`, `bookmark_item_views`.`kind`, `bookmark_item_views`.`category_id`, `bookmark_item_views`.`text`, `bookmark_item_views`.`name` FROM `bookmark_item_views` ORDER BY `bookmark_item_views`.`id`",
	}, rec.queries, "a view-covered selection is served from the view; no per-type queries")
}

// TestInterfaceViewConnectionSynth exercises the view-row synthesis boundary: a
// selection touching only __typename plus the view's non-optional shared columns
// is built from the view rows (no per-type query), while touching an optional
// shared column or a type-specific field falls back to the per-type queries.
func TestInterfaceViewConnectionSynth(t *testing.T) {
	ctx := context.Background()
	drv, err := sql.Open(dialect.SQLite,
		fmt.Sprintf("file:%s?mode=memory&cache=shared&_fk=1", t.Name()))
	require.NoError(t, err)
	rec := &queryRecorder{Driver: drv}
	ec := enttest.NewClient(t,
		enttest.WithOptions(ent.Driver(rec)),
		enttest.WithMigrateOptions(migrate.WithGlobalUniqueID(true)),
	)
	_, err = drv.DB().ExecContext(ctx,
		"CREATE VIEW bookmark_item_views AS SELECT id, 'Todo' AS kind, category_id, text, name FROM todos UNION ALL SELECT id, 'Project' AS kind, category_projects AS category_id, text, name FROM projects")
	require.NoError(t, err)
	gqlc := client.New(handler.NewDefaultServer(gen.NewSchema(ec)))

	// 2 todos (no name) + 1 project (with name) so both synthesized and queried
	// paths have a deterministic, checkable dataset.
	cat := ec.Category.Create().SetText("cat1").SetStatus(category.StatusEnabled).SaveX(ctx)
	ec.Todo.CreateBulk(
		ec.Todo.Create().SetText("t1").SetStatus(todo.StatusInProgress).SetCategory(cat),
		ec.Todo.Create().SetText("t2").SetStatus(todo.StatusCompleted).SetCategory(cat),
	).SaveX(ctx)
	ec.Project.Create().SetText("p1").SetName("pname1").SetCategory(cat).SaveX(ctx)

	t.Run("SynthesizedFromView", func(t *testing.T) {
		// __typename + the non-optional shared columns (id, text) are all carried by
		// the view, so every node is synthesized from its view row.
		var rsp struct {
			BookmarkItems struct {
				Edges []struct {
					Node struct {
						Typename string `json:"__typename"`
						ID       string
						Text     string
					}
				}
			}
		}
		rec.reset()
		gqlc.MustPost(`{ bookmarkItems { edges { node { __typename id text } } } }`, &rsp)
		require.Len(t, rsp.BookmarkItems.Edges, 3)
		got := map[string]string{}
		for _, e := range rsp.BookmarkItems.Edges {
			require.NotEmpty(t, e.Node.ID)
			got[e.Node.Text] = e.Node.Typename
		}
		require.Equal(t, map[string]string{"t1": "Todo", "t2": "Todo", "p1": "Project"}, got,
			"synthesized nodes carry the correct concrete type and shared values")
		require.Equal(t, []string{
			"SELECT `bookmark_item_views`.`id`, `bookmark_item_views`.`kind`, `bookmark_item_views`.`category_id`, `bookmark_item_views`.`text`, `bookmark_item_views`.`name` FROM `bookmark_item_views` ORDER BY `bookmark_item_views`.`id`",
		}, rec.queries, "covered selection is synthesized from the view rows only")
	})

	t.Run("OptionalSharedFieldFallsBack", func(t *testing.T) {
		// `name` is a shared field but Optional, so it is not synthesizable (a NULL
		// name is indistinguishable from ""); selecting it forces the per-type queries
		// even though `text` on its own would have been covered.
		var rsp struct {
			BookmarkItems struct {
				Edges []struct {
					Node struct {
						Typename string `json:"__typename"`
						ID       string
						Text     string
						Name     *string
					}
				}
			}
		}
		rec.reset()
		gqlc.MustPost(`{ bookmarkItems { edges { node { __typename id text name } } } }`, &rsp)
		require.Len(t, rsp.BookmarkItems.Edges, 3)
		var projectName *string
		for _, e := range rsp.BookmarkItems.Edges {
			if e.Node.Typename == "Project" {
				projectName = e.Node.Name
			}
		}
		require.NotNil(t, projectName)
		require.Equal(t, "pname1", *projectName, "the project's name is loaded from its concrete row")
		require.Equal(t, []string{
			"SELECT `bookmark_item_views`.`id`, `bookmark_item_views`.`kind`, `bookmark_item_views`.`category_id`, `bookmark_item_views`.`text`, `bookmark_item_views`.`name` FROM `bookmark_item_views` ORDER BY `bookmark_item_views`.`id`",
			"SELECT `projects`.`id`, `projects`.`text`, `projects`.`name` FROM `projects` WHERE `projects`.`id` IN (?)",
			"SELECT `todos`.`id`, `todos`.`text`, `todos`.`name` FROM `todos` WHERE `todos`.`id` IN (?, ?)",
		}, rec.queries, "an optional shared field forces one batched query per member type")
	})

	t.Run("TypeSpecificFieldFallsBack", func(t *testing.T) {
		// The covering check is per member type. `... on Todo { status }` adds a
		// non-view field only to Todo rows, so only Todos are queried; the Project
		// selection is just __typename, so Projects are still synthesized.
		var rsp struct {
			BookmarkItems struct {
				Edges []struct {
					Node struct {
						Typename string `json:"__typename"`
						Status   *string
					}
				}
			}
		}
		rec.reset()
		gqlc.MustPost(`{ bookmarkItems { edges { node { __typename ... on Todo { status } } } } }`, &rsp)
		require.Len(t, rsp.BookmarkItems.Edges, 3)
		var statuses int
		for _, e := range rsp.BookmarkItems.Edges {
			if e.Node.Typename == "Todo" {
				require.NotNil(t, e.Node.Status)
				statuses++
			}
		}
		require.Equal(t, 2, statuses)
		require.Equal(t, []string{
			"SELECT `bookmark_item_views`.`id`, `bookmark_item_views`.`kind`, `bookmark_item_views`.`category_id`, `bookmark_item_views`.`text`, `bookmark_item_views`.`name` FROM `bookmark_item_views` ORDER BY `bookmark_item_views`.`id`",
			"SELECT `todos`.`id`, `todos`.`status` FROM `todos` WHERE `todos`.`id` IN (?, ?)",
		}, rec.queries, "only Todos are queried; Projects (only __typename) are synthesized")
	})
}

// TestBookmarkPolymorphicItem exercises the 1:1 polymorphic interface field:
// Bookmark.item resolves to a concrete Todo or Project (both implement the
// auto-generated BookmarkItem interface), and is null when neither edge is set.
func TestBookmarkPolymorphicItem(t *testing.T) {
	ctx := context.Background()
	ec := enttest.Open(t, dialect.SQLite,
		fmt.Sprintf("file:%s?mode=memory&cache=shared&_fk=1", t.Name()),
		enttest.WithMigrateOptions(migrate.WithGlobalUniqueID(true)),
	)
	gqlc := client.New(handler.NewDefaultServer(gen.NewSchema(ec)))

	cat := ec.Category.Create().SetText("cat").SetStatus(category.StatusEnabled).SaveX(ctx)
	td := ec.Todo.Create().SetText("t1").SetStatus(todo.StatusInProgress).SetName("todo-name").SetCategory(cat).SaveX(ctx)
	pr := ec.Project.Create().SetText("p1").SetName("proj-name").SetCategory(cat).SaveX(ctx)
	ec.Bookmark.Create().SetName("to-todo").SetTodo(td).SaveX(ctx)
	ec.Bookmark.Create().SetName("to-project").SetProject(pr).SaveX(ctx)
	ec.Bookmark.Create().SetName("empty").SaveX(ctx)

	var rsp struct {
		Bookmarks struct {
			Edges []struct {
				Node struct {
					Name string
					Item *struct {
						Typename string  `json:"__typename"`
						Name     *string // shared interface field
						Status   *string // Todo-only, via fragment
						Text     *string // via Project fragment
					}
				}
			}
		}
	}
	// language=GraphQL
	gqlc.MustPost(`{
		bookmarks {
			edges {
				node {
					name
					item {
						__typename
						name
						... on Todo { status }
						... on Project { text }
					}
				}
			}
		}
	}`, &rsp)

	require.Len(t, rsp.Bookmarks.Edges, 3)
	for _, e := range rsp.Bookmarks.Edges {
		switch e.Node.Name {
		case "to-todo":
			require.NotNil(t, e.Node.Item)
			require.Equal(t, "Todo", e.Node.Item.Typename)
			require.NotNil(t, e.Node.Item.Name)
			require.Equal(t, "todo-name", *e.Node.Item.Name)
			require.NotNil(t, e.Node.Item.Status, "Todo-specific field resolved via fragment")
		case "to-project":
			require.NotNil(t, e.Node.Item)
			require.Equal(t, "Project", e.Node.Item.Typename)
			require.Equal(t, "proj-name", *e.Node.Item.Name)
			require.NotNil(t, e.Node.Item.Text)
			require.Equal(t, "p1", *e.Node.Item.Text)
		case "empty":
			require.Nil(t, e.Node.Item, "item is null when neither edge is set")
		default:
			t.Fatalf("unexpected bookmark %q", e.Node.Name)
		}
	}
}

// TestBookmarkItemEagerLoad verifies how the polymorphic field is resolved under
// a list: a selection needing only __typename/id is served from the foreign key
// (no query to the target tables), while a real field forces one batched load per
// member type (no N+1).
func TestBookmarkItemEagerLoad(t *testing.T) {
	ctx := context.Background()
	drv, err := sql.Open(dialect.SQLite, fmt.Sprintf("file:%s?mode=memory&cache=shared&_fk=1", t.Name()))
	require.NoError(t, err)
	rec := &queryRecorder{Driver: drv}
	ec := enttest.NewClient(t,
		enttest.WithOptions(ent.Driver(rec)),
		enttest.WithMigrateOptions(migrate.WithGlobalUniqueID(true)),
	)
	gqlc := client.New(handler.NewDefaultServer(gen.NewSchema(ec)))

	cat := ec.Category.Create().SetText("cat").SetStatus(category.StatusEnabled).SaveX(ctx)
	td1 := ec.Todo.Create().SetText("t1").SetStatus(todo.StatusInProgress).SetCategory(cat).SaveX(ctx)
	td2 := ec.Todo.Create().SetText("t2").SetStatus(todo.StatusInProgress).SetCategory(cat).SaveX(ctx)
	pr := ec.Project.Create().SetText("p1").SetCategory(cat).SaveX(ctx)
	ec.Bookmark.Create().SetName("b1").SetTodo(td1).SaveX(ctx)
	ec.Bookmark.Create().SetName("b2").SetTodo(td2).SaveX(ctx)
	ec.Bookmark.Create().SetName("b3").SetProject(pr).SaveX(ctx)

	// Covered: only __typename is requested, so the concrete type is served from
	// the foreign key on the bookmark row — no query to the target tables.
	var covered struct {
		Bookmarks struct {
			Edges []struct {
				Node struct {
					Item *struct {
						Typename string `json:"__typename"`
					}
				}
			}
		}
	}
	rec.reset()
	gqlc.MustPost(`{ bookmarks { edges { node { item { __typename } } } } }`, &covered)
	require.Len(t, covered.Bookmarks.Edges, 3)
	var todos, projects int
	for _, e := range covered.Bookmarks.Edges {
		require.NotNil(t, e.Node.Item)
		switch e.Node.Item.Typename {
		case "Todo":
			todos++
		case "Project":
			projects++
		}
	}
	require.Equal(t, 2, todos)
	require.Equal(t, 1, projects)
	require.Equal(t, []string{
		"SELECT `bookmarks`.`id`, `bookmarks`.`bookmark_todo`, `bookmarks`.`bookmark_project` FROM `bookmarks` ORDER BY `bookmarks`.`id`",
	}, rec.queries, "__typename is served from the foreign key; no query to todos/projects")

	// Not covered: a real field (name) forces loading the concrete nodes, batched
	// once per member type (no N+1).
	var loaded struct {
		Bookmarks struct {
			Edges []struct {
				Node struct {
					Item *struct {
						Typename string `json:"__typename"`
						Name     *string
					}
				}
			}
		}
	}
	rec.reset()
	gqlc.MustPost(`{ bookmarks { edges { node { item { __typename name } } } } }`, &loaded)
	require.Len(t, loaded.Bookmarks.Edges, 3)
	require.Equal(t, []string{
		"SELECT `bookmarks`.`id`, `bookmarks`.`bookmark_todo`, `bookmarks`.`bookmark_project` FROM `bookmarks` ORDER BY `bookmarks`.`id`",
		"SELECT `todos`.`id`, `todos`.`name` FROM `todos` WHERE `todos`.`id` IN (?, ?)",
		"SELECT `projects`.`id`, `projects`.`name` FROM `projects` WHERE `projects`.`id` IN (?)",
	}, rec.queries, "a real field forces one batched load per member type (no N+1)")
}

func TestPaginate(t *testing.T) {
	ctx := context.Background()
	ec := enttest.Open(
		t, dialect.SQLite,
		fmt.Sprintf("file:%s?mode=memory&cache=shared&_fk=1", t.Name()),
		enttest.WithMigrateOptions(migrate.WithGlobalUniqueID(true)),
	)
	first := 1
	// Ensure that the pagination query compiles.
	_, err := ec.Todo.Query().
		Select(todo.FieldPriority, todo.FieldStatus).
		Paginate(ctx, nil, &first, nil, nil)
	require.NoError(t, err)
}

func TestPrivateFieldSelectionForPagination(t *testing.T) {
	ctx := context.Background()
	drv, err := sql.Open(dialect.SQLite, fmt.Sprintf("file:%s?mode=memory&cache=shared&_fk=1", t.Name()))
	require.NoError(t, err)
	rec := &queryRecorder{Driver: drv}
	ec := enttest.NewClient(t,
		enttest.WithOptions(ent.Driver(rec)),
		enttest.WithMigrateOptions(migrate.WithGlobalUniqueID(true)),
	)
	ec.Todo.CreateBulk(
		ec.Todo.Create().SetText("t0.1").SetStatus(todo.StatusInProgress),
		ec.Todo.Create().SetText("t0.2").SetStatus(todo.StatusInProgress),
		ec.Todo.Create().SetText("t0.3").SetStatus(todo.StatusCompleted),
		ec.Todo.Create().SetText("t0.4").SetStatus(todo.StatusCompleted),
		ec.Todo.Create().SetText("t0.5").SetStatus(todo.StatusCompleted),
	).SaveX(ctx)

	var (
		// language=GraphQL
		query = `query {
			todosWithJoins(first: 2, orderBy: [{direction: DESC, field: STATUS}]) {
				edges {
					cursor
					node {
						text
					}
				}
			}
		}`
		rsp struct {
			TodosWithJoins struct {
				Edges []struct {
					Cursor string
					Node   struct {
						Text string
					}
				}
			}
		}
		gqlc = client.New(handler.NewDefaultServer(gen.NewSchema(ec)))
	)
	rec.reset()
	gqlc.MustPost(query, &rsp)
	require.Equal(t, []string{
		"SELECT `todos`.`id`, `todos`.`text`, `todos`.`status` FROM `todos` LEFT JOIN `categories` AS `t1` ON `todos`.`category_id` = `t1`.`id` GROUP BY `todos`.`id` ORDER BY `todos`.`status` DESC, `todos`.`id` LIMIT 3",
	}, rec.queries)

	t.Log(rsp.TodosWithJoins)

	var (
		// language=GraphQL
		query2 = `query {
			todosWithJoins(first: 2, after: "gqFp0wAAAAYAAAACoXaRq0lOX1BST0dSRVNT", orderBy: [{direction: DESC, field: STATUS}]) {
				edges {
					cursor
					node {
						text
					}
				}
			}
		}`
		rsp2 struct {
			TodosWithJoins struct {
				Edges []struct {
					Cursor string
					Node   struct {
						Text string
					}
				}
			}
		}
	)
	rec.reset()
	gqlc.MustPost(query2, &rsp2)
	require.Equal(t, []string{
		// BEFORE: "SELECT `todos`.`id`, `todos`.`text`, `todos`.`status` FROM `todos` LEFT JOIN `categories` AS `t1` ON `todos`.`category_id` = `t1`.`id` WHERE `status` < ? OR (`status` = ? AND `id` > ?) GROUP BY `todos`.`id` ORDER BY `todos`.`status` DESC, `todos`.`id` LIMIT 3",
		"SELECT `todos`.`id`, `todos`.`text`, `todos`.`status` FROM `todos` LEFT JOIN `categories` AS `t1` ON `todos`.`category_id` = `t1`.`id` WHERE `todos`.`status` < ? OR (`todos`.`status` = ? AND `todos`.`id` > ?) GROUP BY `todos`.`id` ORDER BY `todos`.`status` DESC, `todos`.`id` LIMIT 3",
	}, rec.queries)
}
