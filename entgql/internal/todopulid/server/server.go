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

package main

import (
	"context"
	"net/http"
	"time"

	"entgo.io/contrib/entgql"
	todopulid "entgo.io/contrib/entgql/internal/todopulid"
	"entgo.io/contrib/entgql/internal/todopulid/ent"
	"entgo.io/contrib/entgql/internal/todopulid/ent/todo"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/debug"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/alecthomas/kong"
	"go.uber.org/zap"

	_ "entgo.io/contrib/entgql/internal/todopulid/ent/runtime"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var cli struct {
		Addr  string `name:"address" default:":8081" help:"Address to listen on."`
		Debug bool   `name:"debug" help:"Enable debugging mode."`
	}
	kong.Parse(&cli)

	log, _ := zap.NewDevelopment()
	client, err := ent.Open(
		"sqlite3",
		"file:ent?mode=memory&cache=shared&_fk=1",
	)
	client = client.Debug()
	if err != nil {
		log.Fatal("opening ent client", zap.Error(err))
	}
	ctx := context.Background()
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatal("running schema migration", zap.Error(err))
	}

	if _, err := client.Todo.Create().
		SetStatus(todo.StatusInProgress).
		SetText("todo 1").
		Save(ctx); err != nil {
		log.Fatal("creating todo", zap.Error(err))
	}

	if _, err := client.Todo.Create().
		SetStatus(todo.StatusInProgress).
		SetText("todo 2").
		Save(ctx); err != nil {
		log.Fatal("creating todo", zap.Error(err))
	}

	srv := handler.NewDefaultServer(todopulid.NewSchema(client))
	srv.Use(entgql.Transactioner{TxOpener: client})
	if cli.Debug {
		srv.Use(&debug.Tracer{})
	}

	http.Handle("/",
		playground.Handler("Todo", "/query"),
	)
	http.Handle("/query", srv)

	log.Info("listening on", zap.String("address", cli.Addr))
	server := &http.Server{
		Addr:              cli.Addr,
		ReadHeaderTimeout: 30 * time.Second,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Error("http server terminated", zap.Error(err))
	}
}
