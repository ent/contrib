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

package renamedtype_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"entgo.io/contrib/entgql"
	gen "entgo.io/contrib/entgql/internal/renamedtype"
	"entgo.io/contrib/entgql/internal/renamedtype/ent"
	"entgo.io/contrib/entgql/internal/renamedtype/ent/enttest"
	"entgo.io/contrib/entgql/internal/renamedtype/ent/migrate"
	"entgo.io/ent/dialect"
	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/stretchr/testify/suite"

	_ "github.com/mattn/go-sqlite3"
)

type todoTestSuite struct {
	suite.Suite
	*client.Client
	ent *ent.Client
}

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
}

func TestRenamedType(t *testing.T) {
	ctx := context.Background()
	ec := enttest.Open(
		t, dialect.SQLite,
		fmt.Sprintf("file:%s?mode=memory&cache=shared&_fk=1", t.Name()),
		enttest.WithMigrateOptions(migrate.WithGlobalUniqueID(true)),
	)
	gqlc := client.New(handler.NewDefaultServer(gen.NewSchema(ec)))
	c1 := ec.ClashingText.Create().SetContent("c1").SaveX(ctx)
	var (
		// language=GraphQL
		query = `query Node($id: ID!) {
			text: node(id: $id) {
				id
				... on NotClashingText {
					content
				}
			}
		}`
		rsp struct {
			Text struct {
				ID, Content string
			}
		}
	)

	gqlc.MustPost(query, &rsp, client.Var("id", c1.ID))
	require.Equal(t, "c1", rsp.Text.Content)
}
