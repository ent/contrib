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

package entproto

import (
	"entgo.io/ent/schema"
)

const SkipAnnotation = "ProtoSkip"

type skipped struct{}

// Skip annotates an ent.Schema to specify that this field will be skipped during .proto generation.
func Skip() schema.Annotation {
	return skipped{}
}

func (f skipped) Name() string {
	return SkipAnnotation
}
