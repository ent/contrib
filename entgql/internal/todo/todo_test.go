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
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/AlekSi/pointer"
	"github.com/facebook/ent/dialect"
	"github.com/facebookincubator/ent-contrib/entgql"
	gen "github.com/facebookincubator/ent-contrib/entgql/internal/todo"
	"github.com/facebookincubator/ent-contrib/entgql/internal/todo/ent/enttest"
	"github.com/facebookincubator/ent-contrib/entgql/internal/todo/ent/migrate"
	"github.com/facebookincubator/ent-contrib/entgql/internal/todo/ent/todo"
	"github.com/stretchr/testify/suite"
	"github.com/vektah/gqlparser/v2/gqlerror"

	_ "github.com/mattn/go-sqlite3"
)

type todoTestSuite struct {
	suite.Suite
	*client.Client
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
)

func (s *todoTestSuite) SetupTest() {
	ec := enttest.Open(s.T(), dialect.SQLite,
		fmt.Sprintf("file:%s-%d?mode=memory&cache=shared&_fk=1",
			s.T().Name(), time.Now().UnixNano(),
		),
		enttest.WithMigrateOptions(migrate.WithGlobalUniqueID(true)),
	)

	srv := handler.New(gen.NewSchema(ec))
	srv.AddTransport(transport.POST{})
	srv.SetErrorPresenter(entgql.DefaultErrorPresenter)
	srv.Use(entgql.Transactioner{TxOpener: ec})
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
		root = 1
	)
	for i := 1; i <= maxTodos; i++ {
		id := strconv.Itoa(i)
		var parent *int
		if i != root {
			if i%2 != 0 {
				parent = pointer.ToInt(i - 2)
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
	s.Assert().Zero(rsp.Todos.TotalCount)
	s.Assert().Empty(rsp.Todos.Edges)
	s.Assert().False(rsp.Todos.PageInfo.HasNextPage)
	s.Assert().False(rsp.Todos.PageInfo.HasPreviousPage)
	s.Assert().Nil(rsp.Todos.PageInfo.StartCursor)
	s.Assert().Nil(rsp.Todos.PageInfo.EndCursor)
}

func (s *todoTestSuite) TestQueryAll() {
	var rsp response
	err := s.Post(queryAll, &rsp)
	s.Require().NoError(err)

	s.Assert().Equal(maxTodos, rsp.Todos.TotalCount)
	s.Require().Len(rsp.Todos.Edges, maxTodos)
	s.Assert().False(rsp.Todos.PageInfo.HasNextPage)
	s.Assert().False(rsp.Todos.PageInfo.HasPreviousPage)
	s.Assert().Equal(
		rsp.Todos.Edges[0].Cursor,
		*rsp.Todos.PageInfo.StartCursor,
	)
	s.Assert().Equal(
		rsp.Todos.Edges[len(rsp.Todos.Edges)-1].Cursor,
		*rsp.Todos.PageInfo.EndCursor,
	)
	for i, edge := range rsp.Todos.Edges {
		s.Assert().Equal(strconv.Itoa(i+1), edge.Node.ID)
		s.Assert().EqualValues(todo.StatusCompleted, edge.Node.Status)
		s.Assert().NotEmpty(edge.Cursor)
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
		id    = 1
	)
	for i := 0; i < maxTodos/first; i++ {
		err := s.Post(query, &rsp,
			client.Var("after", after),
			client.Var("first", first),
		)
		s.Require().NoError(err)
		s.Require().Equal(maxTodos, rsp.Todos.TotalCount)
		s.Require().Len(rsp.Todos.Edges, first)
		s.Assert().True(rsp.Todos.PageInfo.HasNextPage)
		s.Assert().NotEmpty(rsp.Todos.PageInfo.EndCursor)

		for _, edge := range rsp.Todos.Edges {
			s.Assert().Equal(strconv.Itoa(id), edge.Node.ID)
			s.Assert().NotEmpty(edge.Cursor)
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
	s.Assert().Len(rsp.Todos.Edges, maxTodos%first)
	s.Assert().False(rsp.Todos.PageInfo.HasNextPage)
	s.Assert().NotEmpty(rsp.Todos.PageInfo.EndCursor)

	for _, edge := range rsp.Todos.Edges {
		s.Assert().Equal(strconv.Itoa(id), edge.Node.ID)
		s.Assert().NotEmpty(edge.Cursor)
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
	s.Assert().Empty(rsp.Todos.Edges)
	s.Assert().Empty(rsp.Todos.PageInfo.EndCursor)
	s.Assert().False(rsp.Todos.PageInfo.HasNextPage)
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
		id     = maxTodos
	)
	for i := 0; i < maxTodos/last; i++ {
		err := s.Post(query, &rsp,
			client.Var("before", before),
			client.Var("last", last),
		)
		s.Require().NoError(err)
		s.Require().Equal(maxTodos, rsp.Todos.TotalCount)
		s.Require().Len(rsp.Todos.Edges, last)
		s.Assert().True(rsp.Todos.PageInfo.HasPreviousPage)
		s.Assert().NotEmpty(rsp.Todos.PageInfo.StartCursor)

		for i := len(rsp.Todos.Edges) - 1; i >= 0; i-- {
			edge := &rsp.Todos.Edges[i]
			s.Assert().Equal(strconv.Itoa(id), edge.Node.ID)
			s.Assert().NotEmpty(edge.Cursor)
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
	s.Assert().Len(rsp.Todos.Edges, maxTodos%last)
	s.Assert().False(rsp.Todos.PageInfo.HasPreviousPage)
	s.Assert().NotEmpty(rsp.Todos.PageInfo.StartCursor)

	for i := len(rsp.Todos.Edges) - 1; i >= 0; i-- {
		edge := &rsp.Todos.Edges[i]
		s.Assert().Equal(strconv.Itoa(id), edge.Node.ID)
		s.Assert().NotEmpty(edge.Cursor)
		id--
	}
	s.Assert().Zero(id)

	before = rsp.Todos.PageInfo.StartCursor
	rsp = response{}
	err = s.Post(query, &rsp,
		client.Var("before", before),
		client.Var("last", last),
	)
	s.Require().NoError(err)
	s.Require().Equal(maxTodos, rsp.Todos.TotalCount)
	s.Assert().Empty(rsp.Todos.Edges)
	s.Assert().Empty(rsp.Todos.PageInfo.StartCursor)
	s.Assert().False(rsp.Todos.PageInfo.HasPreviousPage)
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
			s.Assert().Equal(maxTodos, rsp.Todos.TotalCount)
			if i < steps-1 {
				s.Require().Len(rsp.Todos.Edges, step)
				s.Assert().True(rsp.Todos.PageInfo.HasNextPage)
			} else {
				s.Require().Len(rsp.Todos.Edges, maxTodos%step)
				s.Assert().False(rsp.Todos.PageInfo.HasNextPage)
			}
			s.Assert().True(sort.SliceIsSorted(rsp.Todos.Edges, func(i, j int) bool {
				return rsp.Todos.Edges[i].Node.Text < rsp.Todos.Edges[j].Node.Text
			}))
			s.Require().NotNil(rsp.Todos.PageInfo.StartCursor)
			s.Assert().Equal(*rsp.Todos.PageInfo.StartCursor, rsp.Todos.Edges[0].Cursor)
			s.Require().NotNil(rsp.Todos.PageInfo.EndCursor)
			end := rsp.Todos.Edges[len(rsp.Todos.Edges)-1]
			s.Assert().Equal(*rsp.Todos.PageInfo.EndCursor, end.Cursor)
			if i > 0 {
				s.Assert().Less(endText, rsp.Todos.Edges[0].Node.Text)
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
			s.Assert().Equal(maxTodos, rsp.Todos.TotalCount)
			if i < steps-1 {
				s.Require().Len(rsp.Todos.Edges, step)
				s.Assert().True(rsp.Todos.PageInfo.HasNextPage)
			} else {
				s.Require().Len(rsp.Todos.Edges, maxTodos%step)
				s.Assert().False(rsp.Todos.PageInfo.HasNextPage)
			}
			s.Assert().True(sort.SliceIsSorted(rsp.Todos.Edges, func(i, j int) bool {
				left, _ := strconv.Atoi(rsp.Todos.Edges[i].Node.ID)
				right, _ := strconv.Atoi(rsp.Todos.Edges[j].Node.ID)
				return left > right
			}))
			s.Require().NotNil(rsp.Todos.PageInfo.StartCursor)
			s.Assert().Equal(*rsp.Todos.PageInfo.StartCursor, rsp.Todos.Edges[0].Cursor)
			s.Require().NotNil(rsp.Todos.PageInfo.EndCursor)
			end := rsp.Todos.Edges[len(rsp.Todos.Edges)-1]
			s.Assert().Equal(*rsp.Todos.PageInfo.EndCursor, end.Cursor)
			if i > 0 {
				id, _ := strconv.Atoi(rsp.Todos.Edges[0].Node.ID)
				s.Assert().Greater(endID, id)
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
			s.Assert().Equal(maxTodos, rsp.Todos.TotalCount)
			if i < steps-1 {
				s.Require().Len(rsp.Todos.Edges, step)
				s.Assert().True(rsp.Todos.PageInfo.HasPreviousPage)
			} else {
				s.Require().Len(rsp.Todos.Edges, maxTodos%step)
				s.Assert().False(rsp.Todos.PageInfo.HasPreviousPage)
			}
			s.Assert().True(sort.SliceIsSorted(rsp.Todos.Edges, func(i, j int) bool {
				return rsp.Todos.Edges[i].Node.Priority < rsp.Todos.Edges[j].Node.Priority
			}))
			s.Require().NotNil(rsp.Todos.PageInfo.StartCursor)
			start := rsp.Todos.Edges[0]
			s.Assert().Equal(*rsp.Todos.PageInfo.StartCursor, start.Cursor)
			s.Require().NotNil(rsp.Todos.PageInfo.EndCursor)
			end := rsp.Todos.Edges[len(rsp.Todos.Edges)-1]
			s.Assert().Equal(*rsp.Todos.PageInfo.EndCursor, end.Cursor)
			if i > 0 {
				s.Assert().Greater(startPriority, end.Node.Priority)
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
			s.Assert().Equal(maxTodos, rsp.Todos.TotalCount)
			if i < steps-1 {
				s.Require().Len(rsp.Todos.Edges, step)
				s.Assert().True(rsp.Todos.PageInfo.HasPreviousPage)
			} else {
				s.Require().Len(rsp.Todos.Edges, maxTodos%step)
				s.Assert().False(rsp.Todos.PageInfo.HasPreviousPage)
			}
			s.Assert().True(sort.SliceIsSorted(rsp.Todos.Edges, func(i, j int) bool {
				left, _ := time.Parse(time.RFC3339, rsp.Todos.Edges[i].Node.CreatedAt)
				right, _ := time.Parse(time.RFC3339, rsp.Todos.Edges[j].Node.CreatedAt)
				return left.After(right)
			}))
			s.Require().NotNil(rsp.Todos.PageInfo.StartCursor)
			start := rsp.Todos.Edges[0]
			s.Assert().Equal(*rsp.Todos.PageInfo.StartCursor, start.Cursor)
			s.Require().NotNil(rsp.Todos.PageInfo.EndCursor)
			end := rsp.Todos.Edges[len(rsp.Todos.Edges)-1]
			s.Assert().Equal(*rsp.Todos.PageInfo.EndCursor, end.Cursor)
			if i > 0 {
				endCreatedAt, _ := time.Parse(time.RFC3339, end.Node.CreatedAt)
				s.Assert().True(startCreatedAt.Before(endCreatedAt) || startCreatedAt.Equal(endCreatedAt))
			}
			startCreatedAt, _ = time.Parse(time.RFC3339, start.Node.CreatedAt)
		}
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
	err := s.Post(query, &rsp, client.Var("id", maxTodos))
	s.Require().NoError(err)
	jerr, ok := s.Post(query, &rsp, client.Var("id", -1)).(client.RawJsonError)
	s.Require().True(ok)
	var errs gqlerror.List
	err = json.Unmarshal(jerr.RawMessage, &errs)
	s.Require().NoError(err)
	s.Require().Len(errs, 1)
	s.Assert().Equal("Could not resolve to a node with the global id of '-1'", errs[0].Message)
	s.Assert().Equal("NOT_FOUND", errs[0].Extensions["code"])
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
	err := s.Post(query, &rsp, client.Var("id", 1))
	s.Require().NoError(err)
	s.Assert().Nil(rsp.Todo.Parent)
	s.Assert().Len(rsp.Todo.Children, maxTodos/2+1)
	s.Assert().Condition(func() bool {
		for _, child := range rsp.Todo.Children {
			if child.Text == "3" {
				s.Require().Len(child.Children, 1)
				s.Assert().Equal("5", child.Children[0].Text)
				return true
			}
		}
		return false
	})

	err = s.Post(query, &rsp, client.Var("id", 4))
	s.Require().NoError(err)
	s.Require().NotNil(rsp.Todo.Parent)
	s.Assert().Equal("1", rsp.Todo.Parent.Text)
	s.Assert().Empty(rsp.Todo.Children)

	err = s.Post(query, &rsp, client.Var("id", 5))
	s.Require().NoError(err)
	s.Require().NotNil(rsp.Todo.Parent)
	s.Assert().Equal("3", rsp.Todo.Parent.Text)
	s.Require().NotNil(rsp.Todo.Parent.Parent)
	s.Assert().Equal("1", rsp.Todo.Parent.Parent.Text)
	s.Require().Len(rsp.Todo.Children, 1)
	s.Assert().Equal("7", rsp.Todo.Children[0].Text)
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
			s.Assert().Nil(edge.Node.Parent)
			s.Assert().Len(edge.Node.Children, maxTodos/2+1)
		case i%2 == 0:
			s.Require().NotNil(edge.Node.Parent)
			id, err := strconv.Atoi(edge.Node.Parent.ID)
			s.Require().NoError(err)
			s.Assert().Equal(i-1, id)
			if i < len(rsp.Todos.Edges)-2 {
				s.Assert().Len(edge.Node.Children, 1)
			} else {
				s.Assert().Empty(edge.Node.Children)
			}
		case i%2 != 0:
			s.Require().NotNil(edge.Node.Parent)
			s.Assert().Equal("1", edge.Node.Parent.ID)
			s.Assert().Empty(edge.Node.Children)
		}
	}
}

func (s *todoTestSuite) TestEnumEncoding() {
	s.Run("Encode", func() {
		const status = todo.StatusCompleted
		s.Assert().Implements((*graphql.Marshaler)(nil), status)
		var b strings.Builder
		status.MarshalGQL(&b)
		str := b.String()
		const quote = `"`
		s.Assert().Equal(quote, str[:1])
		s.Assert().Equal(quote, str[len(str)-1:])
		str = str[1 : len(str)-1]
		s.Assert().EqualValues(status, str)
	})
	s.Run("Decode", func() {
		const want = todo.StatusInProgress
		var got todo.Status
		s.Assert().Implements((*graphql.Unmarshaler)(nil), &got)
		err := got.UnmarshalGQL(want.String())
		s.Assert().NoError(err)
		s.Assert().Equal(want, got)
	})
}
