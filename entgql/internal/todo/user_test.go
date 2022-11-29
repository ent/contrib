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
	"fmt"
	"strconv"
	"testing"
	"time"

	"entgo.io/contrib/entgql"
	gen "entgo.io/contrib/entgql/internal/todo"
	"entgo.io/contrib/entgql/internal/todo/ent"
	"entgo.io/contrib/entgql/internal/todo/ent/enttest"
	"entgo.io/contrib/entgql/internal/todo/ent/migrate"
	"entgo.io/contrib/entgql/internal/todo/ent/schema"
	"entgo.io/ent/dialect"
	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/stretchr/testify/suite"
)

type userTestSuite struct {
	suite.Suite
	*client.Client
	ent       *ent.Client
	famousIDs []string
}

const (
	queryAllUsers = `query {
		users {
			totalCount
			edges {
				node {
					id
					name
					friends {
						id
					}
					friendships {
						id
						createdAt
						userID
						user {
							id
						}
						friendID
						friend {
							id
						}
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
	maxUsers       = 32
	maxFamousUsers = 5
	idOffsetUser   = 4 << 32
)

func (s *userTestSuite) SetupTest() {
	schema.TimeNow = func() time.Time {
		return time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	}
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

	const mutation = `mutation($name: String!, $friendIDs: [ID!]) {
		createUser(input: {name: $name, friendIDs: $friendIDs}) {
			id
		}
	}`
	var (
		rsp struct {
			CreateUser struct {
				ID string
			}
		}
		friendIDs = make([]string, 0, maxFamousUsers)
	)
	for i := 1; i <= maxFamousUsers; i++ {
		id := strconv.Itoa(idOffsetUser + i)
		err := s.Post(mutation, &rsp,
			client.Var("name", id),
		)
		s.Require().NoError(err)
		s.Require().Equal(id, rsp.CreateUser.ID)
		friendIDs = append(friendIDs, rsp.CreateUser.ID)
	}
	for i := 1; i <= maxUsers; i++ {
		schema.TimeNow = func() time.Time {
			return time.Date(2022, 1, i, 0, 0, 0, 0, time.UTC)
		}
		id := strconv.Itoa(idOffsetUser + i + maxFamousUsers)
		opts := []client.Option{client.Var("name", id)}
		if i%2 == 1 {
			opts = append(opts, client.Var("friendIDs", friendIDs))
		}

		err := s.Post(mutation, &rsp, opts...)
		s.Require().NoError(err)
		s.Require().Equal(id, rsp.CreateUser.ID)
	}
	s.famousIDs = friendIDs
}

func TestUser(t *testing.T) {
	suite.Run(t, new(userTestSuite))
}

type responseUser struct {
	Users struct {
		TotalCount int
		Edges      []struct {
			Node struct {
				ID      string
				Name    string
				Friends []struct {
					ID string
				}
				Friendships []struct {
					ID        string
					CreatedAt string
					UserID    string
					User      struct {
						ID string
					}
					FriendID string
					Friend   struct {
						ID string
					}
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

func (s *userTestSuite) TestQueryEmpty() {
	{
		var rsp struct{ ClearFriendships int }
		err := s.Post(`mutation { clearFriendships }`, &rsp)
		s.Require().NoError(err)
		s.Require().Equal(maxFamousUsers*maxUsers, rsp.ClearFriendships)
	}
	{
		var rsp struct{ ClearUsers int }
		err := s.Post(`mutation { clearUsers }`, &rsp)
		s.Require().NoError(err)
		s.Require().Equal(maxFamousUsers+maxUsers, rsp.ClearUsers)
	}
	var rsp responseUser
	err := s.Post(queryAllUsers, &rsp)
	s.Require().NoError(err)
	s.Require().Zero(rsp.Users.TotalCount)
	s.Require().Empty(rsp.Users.Edges)
	s.Require().False(rsp.Users.PageInfo.HasNextPage)
	s.Require().False(rsp.Users.PageInfo.HasPreviousPage)
	s.Require().Nil(rsp.Users.PageInfo.StartCursor)
	s.Require().Nil(rsp.Users.PageInfo.EndCursor)
}

func (s *userTestSuite) TestQueryAll() {
	totalUsers := maxUsers + maxFamousUsers
	var rsp responseUser
	err := s.Post(queryAllUsers, &rsp)
	s.Require().NoError(err)

	s.Require().Equal(totalUsers, rsp.Users.TotalCount)
	s.Require().Len(rsp.Users.Edges, totalUsers)
	s.Require().False(rsp.Users.PageInfo.HasNextPage)
	s.Require().False(rsp.Users.PageInfo.HasPreviousPage)
	s.Require().Equal(
		rsp.Users.Edges[0].Cursor,
		*rsp.Users.PageInfo.StartCursor,
	)
	s.Require().Equal(
		rsp.Users.Edges[len(rsp.Users.Edges)-1].Cursor,
		*rsp.Users.PageInfo.EndCursor,
	)

	for i, edge := range rsp.Users.Edges {
		s.Require().Equal(strconv.Itoa(idOffsetUser+i+1), edge.Node.ID)
		s.Require().NotEmpty(edge.Cursor)
		for idx, friend := range edge.Node.Friends {
			fs := edge.Node.Friendships[idx]

			s.Require().NotEmpty(fs.CreatedAt)
			s.Require().Equal(fs.Friend.ID, friend.ID)
			s.Require().Equal(fs.FriendID, friend.ID)
			s.Require().Equal(fs.User.ID, edge.Node.ID)
			s.Require().Equal(fs.UserID, edge.Node.ID)
		}
	}
}

func (s *userTestSuite) TestEdgeSchemaFiltering() {
	famousCheck := map[string]struct{}{}
	for _, id := range s.famousIDs {
		famousCheck[id] = struct{}{}
	}

	s.Run("Friendship on Jan-07", func() {
		const (
			query = `query($friendDate: Time) {
				users (where:{hasFriendshipsWith:{createdAt: $friendDate}}) {
					totalCount
					edges {
						node {
							id
							friendships {
								id
							}
						}
					}
				}
			}`
		)

		var rsp responseUser
		err := s.Post(query, &rsp,
			client.Var("friendDate", "2022-01-07T00:00:00Z"),
		)
		s.NoError(err)
		s.Equal(maxFamousUsers+1, rsp.Users.TotalCount)

		for _, edge := range rsp.Users.Edges {
			if len(edge.Node.Friendships) == maxFamousUsers {
				s.Equal(strconv.Itoa(idOffsetUser+maxFamousUsers+7), edge.Node.ID)
			} else {
				s.Equal(maxUsers/2, len(edge.Node.Friendships))
				_, ok := famousCheck[edge.Node.ID]
				s.True(ok)
			}
		}
	})
}
