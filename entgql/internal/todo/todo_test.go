// Copyright 2019-present Facebook
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
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
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

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
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/AlekSi/pointer"
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
	idOffset = 1 << 32
)

func (s *todoTestSuite) SetupTest() {
	s.ent = enttest.Open(s.T(), dialect.SQLite,
		fmt.Sprintf("file:%s-%d?mode=memory&cache=shared&_fk=1",
			s.T().Name(), time.Now().UnixNano(),
		),
		enttest.WithMigrateOptions(migrate.WithGlobalUniqueID(true)),
	)

	srv := handler.New(gen.NewSchema(s.ent))
	srv.AddTransport(transport.POST{})
	srv.Use(entgql.Transactioner{TxOpener: s.ent})
	s.Client = client.New(srv)

	const mutation = `mutation($priority: Int, $text: String!, $parent: ID) {
		createTodo(todo: {status: COMPLETED, priority: $priority, text: $text, parent: $parent}) {
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
				ID        string
				CreatedAt string
				Priority  int
				Status    todo.Status
				Text      string
				Parent    struct {
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
		query = `query($after: Cursor, $first: Int, $before: Cursor, $last: Int, $direction: OrderDirection!, $field: TodoOrderField) {
			todos(after: $after, first: $first, before: $before, last: $last, orderBy: { direction: $direction, field: $field }) {
				totalCount
				edges {
					node {
						id
						createdAt
						priority
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
				client.Var("field", "PRIORITY"),
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
				return rsp.Todos.Edges[i].Node.Priority < rsp.Todos.Edges[j].Node.Priority
			}))
			s.Require().NotNil(rsp.Todos.PageInfo.StartCursor)
			start := rsp.Todos.Edges[0]
			s.Require().Equal(*rsp.Todos.PageInfo.StartCursor, start.Cursor)
			s.Require().NotNil(rsp.Todos.PageInfo.EndCursor)
			end := rsp.Todos.Edges[len(rsp.Todos.Edges)-1]
			s.Require().Equal(*rsp.Todos.PageInfo.EndCursor, end.Cursor)
			if i > 0 {
				s.Require().Greater(startPriority, end.Node.Priority)
			}
			startPriority = start.Node.Priority
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
		query = `query($after: Cursor, $first: Int, $before: Cursor, $last: Int, $status: Status, $hasParent: Boolean, $hasCategory: Boolean) {
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
}

func (s *todoTestSuite) TestNode() {
	const (
		query = `query($id: ID!) {
			todo: node(id: $id) {
				... on Todo {
					priority
				}
			}
		}`
	)
	var rsp struct{ Todo struct{ Priority int } }
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
						text
						children {
							text
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
			Children []struct {
				Text     string
				Children []struct {
					Text string
				}
			}
		}
	}
	err := s.Post(query, &rsp, client.Var("id", idOffset+1))
	s.Require().NoError(err)
	s.Require().Nil(rsp.Todo.Parent)
	s.Require().Len(rsp.Todo.Children, maxTodos/2+1)
	s.Require().Condition(func() bool {
		for _, child := range rsp.Todo.Children {
			if child.Text == strconv.Itoa(idOffset+3) {
				s.Require().Len(child.Children, 1)
				s.Require().Equal(strconv.Itoa(idOffset+5), child.Children[0].Text)
				return true
			}
		}
		return false
	})

	err = s.Post(query, &rsp, client.Var("id", idOffset+4))
	s.Require().NoError(err)
	s.Require().NotNil(rsp.Todo.Parent)
	s.Require().Equal(strconv.Itoa(idOffset+1), rsp.Todo.Parent.Text)
	s.Require().Empty(rsp.Todo.Children)

	err = s.Post(query, &rsp, client.Var("id", strconv.Itoa(idOffset+5)))
	s.Require().NoError(err)
	s.Require().NotNil(rsp.Todo.Parent)
	s.Require().Equal(strconv.Itoa(idOffset+3), rsp.Todo.Parent.Text)
	s.Require().NotNil(rsp.Todo.Parent.Parent)
	s.Require().Equal(strconv.Itoa(idOffset+1), rsp.Todo.Parent.Parent.Text)
	s.Require().Len(rsp.Todo.Children, 1)
	s.Require().Equal(strconv.Itoa(idOffset+7), rsp.Todo.Children[0].Text)
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
							id
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
					Children []struct {
						ID string
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
			s.Require().Len(edge.Node.Children, maxTodos/2+1)
		case i%2 == 0:
			s.Require().NotNil(edge.Node.Parent)
			id, err := strconv.Atoi(edge.Node.Parent.ID)
			s.Require().NoError(err)
			s.Require().Equal(idOffset+i-1, id)
			if i < len(rsp.Todos.Edges)-2 {
				s.Require().Len(edge.Node.Children, 1)
			} else {
				s.Require().Empty(edge.Node.Children)
			}
		case i%2 != 0:
			s.Require().NotNil(edge.Node.Parent)
			s.Require().Equal(strconv.Itoa(idOffset+1), edge.Node.Parent.ID)
			s.Require().Empty(edge.Node.Children)
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
	err := s.Post(`mutation {
		createTodo(todo: { text: "OKE", parent: 4294967297 }) {
			parent {
				id
				text
			}
			text
		}
	}`, &rsp, client.Var("text", s.T().Name()))
	s.Require().NoError(err)
	s.Require().Equal("OKE", rsp.CreateTodo.Text)
	s.Require().Equal(strconv.Itoa(idOffset+1), rsp.CreateTodo.Parent.ID)
	s.Require().Equal(strconv.Itoa(idOffset+1), rsp.CreateTodo.Parent.Text)
}
