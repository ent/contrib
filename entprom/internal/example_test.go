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
