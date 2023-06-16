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

package todo

import (
	"context"
	"fmt"
	"testing"

	"entgo.io/contrib/entgql/internal/todogotype/ent/enttest"
	"entgo.io/contrib/entgql/internal/todogotype/ent/todo"

	"entgo.io/ent/dialect"
	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

func TestSanity(t *testing.T) {
	ctx := context.Background()
	ec := enttest.Open(
		t, dialect.SQLite,
		fmt.Sprintf("file:%s?mode=memory&cache=shared&_fk=1", t.Name()),
	)
	srv := handler.NewDefaultServer(NewSchema(ec))
	gqlc := client.New(srv)

	todos := ec.Todo.CreateBulk(
		ec.Todo.Create().SetID("todos/1").SetText("1").SetStatus(todo.StatusInProgress),
		ec.Todo.Create().SetID("todos/2").SetText("2").SetStatus(todo.StatusInProgress),
		ec.Todo.Create().SetID("todos/3").SetText("3").SetStatus(todo.StatusInProgress),
	).SaveX(ctx)

	for i := range todos {
		var rsp struct {
			Node struct {
				Text string
			}
		}
		err := gqlc.Post(`query node($id: ID!) {
	    	node(id: $id) {
	    		... on Todo {
					text
				}
			}
		}`, &rsp, client.Var("id", todos[i].ID))
		require.NoError(t, err)
		require.Equal(t, todos[i].Text, rsp.Node.Text)
	}
}
