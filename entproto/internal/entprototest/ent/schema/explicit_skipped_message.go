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
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/mixin"

	"github.com/bionicstork/contrib/entproto"
)

// ExplicitSkippedMessage holds the schema definition for the ExplicitSkippedMessage entity.
type ExplicitSkippedMessage struct {
	ent.Schema
}

func (ExplicitSkippedMessage) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.SkipGen(),
	}
}

func (ExplicitSkippedMessage) Mixin() []ent.Mixin {
	return []ent.Mixin{
		ProtoMixin{},
	}
}

type ProtoMixin struct {
	mixin.Schema
}

func (ProtoMixin) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(),
	}
}
