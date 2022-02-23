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

package schema

import (
	"github.com/bionicstork/bionicstork/pkg/entproto"
	"entgo.io/ent"
	"entgo.io/ent/schema"
)

// AllMethodsService holds the schema definition for the AllMethodsService entity.
type AllMethodsService struct {
	ent.Schema
}

func (AllMethodsService) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(),
		entproto.Service(
			entproto.Methods(entproto.MethodAll),
		),
	}
}
