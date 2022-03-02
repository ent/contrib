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
	"fmt"
	"github.com/bionicstork/contrib/entproto"
	"log"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

func main() {
	var (
		schemaPath      = flag.String("path", "", "path to schema directory")
		protoTarget     = flag.String("target", "", "proto schema target")
		targetGoPackage = flag.String("target-go-package", "", "pb.go package path, used for setting the \"go_package\" in the proto files")
		idType          = flag.String("idtype", "", "type of the table's primary key")
	)
	flag.Parse()
	if *schemaPath == "" {
		log.Fatal("entproto: must specify schema path. use entproto -path ./ent/schema")
	}

	idConfigType := IDType(field.TypeInt)
	if *idType != "" {
		err := idConfigType.Set(*idType)
		if err != nil {
			log.Fatal(err)
		}
	}
	graph, err := entc.LoadGraph(*schemaPath, &gen.Config{Target: *protoTarget, IDType: &field.TypeInfo{Type: field.Type(idConfigType)}})
	if err != nil {
		log.Fatalf("entproto: failed loading ent graph: %v", err)
	}
	if err := entproto.Generate(graph, *targetGoPackage); err != nil {
		log.Fatalf("entproto: failed generating protos: %s", err)
	}
}

// IDType is a custom ID implementation for pflag.
type IDType field.Type

// Set implements the Set method of the flag.Value interface.
func (t *IDType) Set(s string) error {
	switch s {
	case field.TypeInt.String():
		*t = IDType(field.TypeInt)
	case field.TypeInt64.String():
		*t = IDType(field.TypeInt64)
	case field.TypeUint.String():
		*t = IDType(field.TypeUint)
	case field.TypeUint64.String():
		*t = IDType(field.TypeUint64)
	case field.TypeString.String():
		*t = IDType(field.TypeString)
	default:
		return fmt.Errorf("invalid type %q", s)
	}
	return nil
}
