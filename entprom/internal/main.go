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
	"fmt"
	"log"
	"net/http"

	"entgo.io/contrib/entprom"
	"entgo.io/contrib/entprom/internal/ent"
	_ "entgo.io/contrib/entprom/internal/ent/runtime"
	_ "github.com/mattn/go-sqlite3"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var dbClient *ent.Client

func initDB() *ent.Client {
	client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	ctx := context.Background()
	// Run the auto migration tool.
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	client.Use(entprom.Hook(
		entprom.Labels(map[string]string{"environment": "dev"}),
	))

	return client
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	// Run operations.
	a8m, err := dbClient.User.Create().SetName("a8m").Save(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	root, err := dbClient.File.Create().SetName("/").SetOwner(a8m).Save(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = dbClient.File.Create().SetName("dev").SetParent(root).Save(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "check the metrics")
}

func main() {
	dbClient = initDB()
	http.HandleFunc("/", handler)
	http.Handle("/metrics", promhttp.Handler())
	log.Println("server starting on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
