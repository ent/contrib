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

package entoas

import (
	"testing"

	"entgo.io/ent/entc/gen"
	"github.com/ogen-go/ogen"
	"github.com/stretchr/testify/require"

	"entgo.io/contrib/entoas/serialization"
)

func TestAnnotation(t *testing.T) {
	t.Parallel()

	a := ReadOnly(true)
	require.Equal(t, true, a.ReadOnly)

	a = Groups("create", "groups")
	require.Equal(t, serialization.Groups{"create", "groups"}, a.Groups)

	a = CreateOperation(OperationGroups("create", "groups"), OperationPolicy(PolicyExpose), RequestGroups("reqCreate", "groups"))
	require.Equal(t, OperationConfig{PolicyExpose, serialization.Groups{"create", "groups"}, serialization.Groups{"reqCreate", "groups"}}, a.Create)

	a = ReadOperation(OperationGroups("read", "groups"), OperationPolicy(PolicyExpose), RequestGroups("reqRead", "groups"))
	require.Equal(t, OperationConfig{PolicyExpose, serialization.Groups{"read", "groups"}, serialization.Groups{"reqRead", "groups"}}, a.Read)

	a = UpdateOperation(OperationGroups("update", "groups"), OperationPolicy(PolicyExpose), RequestGroups("reqUpdate", "groups"))
	require.Equal(t, OperationConfig{PolicyExpose, serialization.Groups{"update", "groups"}, serialization.Groups{"reqUpdate", "groups"}}, a.Update)

	a = DeleteOperation(OperationGroups("delete", "groups"), OperationPolicy(PolicyExpose), RequestGroups("reqDelete", "groups"))
	require.Equal(t, OperationConfig{PolicyExpose, serialization.Groups{"delete", "groups"}, serialization.Groups{"reqDelete", "groups"}}, a.Delete)

	a = ListOperation(OperationGroups("list", "groups"), OperationPolicy(PolicyExpose), RequestGroups("reqList", "groups"))
	require.Equal(t, OperationConfig{PolicyExpose, serialization.Groups{"list", "groups"}, serialization.Groups{"reqList", "groups"}}, a.List)

	b := Example("example")
	require.Equal(t, "example", b.Example)

	c := Schema(ogen.Binary())
	require.Equal(t, ogen.Binary(), c.Schema)

	a = a.Merge(b).(Annotation).Merge(c).(Annotation)
	ex := Annotation{
		Example: "example",
		Schema:  ogen.Binary(),
		List: OperationConfig{
			Groups:        serialization.Groups{"list", "groups"},
			RequestGroups: serialization.Groups{"reqList", "groups"},
			Policy:        PolicyExpose,
		},
	}
	require.Equal(t, ex, a)

	a = a.Merge(ReadOnly(true)).(Annotation)
	ex.ReadOnly = true
	require.Equal(t, ex, a)

	crOp := CreateOperation(OperationPolicy(PolicyExpose))
	dlOp := DeleteOperation(OperationPolicy(PolicyExclude))
	crdlEx := Annotation{
		Create: OperationConfig{
			Policy: PolicyExpose,
		},
		Delete: OperationConfig{
			Policy: PolicyExclude,
		},
	}
	crOp = crOp.Merge(dlOp).(Annotation)
	require.Equal(t, crdlEx, crOp)

	ac, err := SchemaAnnotation(new(gen.Type))
	require.NoError(t, err)
	require.NotNil(t, ac)
	ac, err = SchemaAnnotation(&gen.Type{Annotations: gen.Annotations{a.Name(): ex}})
	require.NoError(t, err)
	require.Equal(t, &ex, ac)

	ac, err = FieldAnnotation(new(gen.Field))
	require.NoError(t, err)
	require.NotNil(t, ac)
	ac, err = FieldAnnotation(&gen.Field{Annotations: gen.Annotations{a.Name(): ex}})
	require.NoError(t, err)
	require.Equal(t, &ex, ac)

	ac, err = EdgeAnnotation(new(gen.Edge))
	require.NoError(t, err)
	require.NotNil(t, ac)
	ac, err = EdgeAnnotation(&gen.Edge{Annotations: gen.Annotations{a.Name(): ex}})
	require.NoError(t, err)
	require.Equal(t, &ex, ac)
}
