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
	"entgo.io/ent/schema/field"
	"flag"
	"log"

	"entgo.io/contrib/entproto"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

func main() {
	var (
		schemaPath = flag.String("path", "", "path to schema directory")
		idtype     = flag.String("idtype", "int", "path to schema directory")
	)
	flag.Parse()
	if *schemaPath == "" {
		log.Fatal("entproto: must specify schema path. use entproto -path ./ent/schema")
	}
	cfg := &gen.Config{}
	cfg.IDType = &field.TypeInfo{Type: str2Type(idtype)}
	graph, err := entc.LoadGraph(*schemaPath, cfg)
	if err != nil {
		log.Fatalf("entproto: failed loading ent graph: %v", err)
	}
	if err := entproto.Generate(graph); err != nil {
		log.Fatalf("entproto: failed generating protos: %s", err)
	}
}

func str2Type(t *string) field.Type {
	switch *t {
	case field.TypeInt.String():
		return field.TypeInt
	case field.TypeInt64.String():
		return field.TypeInt64
	case field.TypeString.String():
		return field.TypeString
	default:
		return field.TypeInvalid
	}
}
