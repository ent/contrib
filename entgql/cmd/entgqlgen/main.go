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
	"flag"
	"fmt"
	"log"
	"os"

	"entgo.io/contrib/entgql"
	"entgo.io/contrib/entgql/plugin"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen/config"
)

func main() {
	var (
		schemaPath = flag.String("path", "", "path to schema directory")
	)
	flag.Parse()
	if *schemaPath == "" {
		log.Fatal("entgqlgen: must specify schema path. use entgqlgen -path ./ent/schema")
	}

	graph, err := entc.LoadGraph(*schemaPath, &gen.Config{
		Annotations: map[string]interface{}{
			"EntGQL": entgql.Annotation{
				GQLScalarMappings: map[string]string{
					"Time":     "Time",
					"Date":     "Date",
					"Duration": "Duration",
					"Uint64":   "Uint64",
				},
			},
		},
	})
	if err != nil {
		log.Fatalf("entproto: failed loading ent graph: %v", err)
	}

	cfg, err := config.LoadConfigFromDefaultLocations()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load config", err.Error())
		os.Exit(2)
	}

	entgqlPlugin, err := plugin.New(graph)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to create entgql plugin", err.Error())
		os.Exit(2)
	}
	err = api.Generate(cfg,
		api.PrependPlugin(entgqlPlugin),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(3)
	}
}
