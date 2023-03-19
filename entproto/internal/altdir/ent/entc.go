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

//go:build ignore
// +build ignore

package main

import (
	"log"

	"entgo.io/contrib/entproto"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

func main() {
	extension, err := entproto.NewExtension(
		entproto.WithProtoDir("./v1/api"),
	)
	if err != nil {
		panic(err)
	}
	if err := entc.Generate("./schema",
		&gen.Config{},
		entc.Extensions(
			extension,
		),
	); err != nil {
		log.Fatal("running ent codegen:", err)
	}
}
