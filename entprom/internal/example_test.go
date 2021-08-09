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


package main

import (
	"context"
	"log"

	"entgo.io/contrib/entprom"
	"entgo.io/contrib/entprom/internal/ent"
	_ "github.com/mattn/go-sqlite3"
)

func Example_prometheusHook() {
	client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	ctx := context.Background()
	// Run the auto migration tool.
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	// Add Global Hook
	client.Use(entprom.Hook())

	// Run operations.
	a8m := client.User.Create().SetName("a8m").SaveX(ctx)
	root := client.File.Create().SetName("/").SetOwner(a8m).SaveX(ctx)
	client.File.Create().SetName("dev").SetParent(root).SaveX(ctx)
}
