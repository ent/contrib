// Copyright 2019-present Facebook
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package entgql_test

import (
	"testing"

	"entgo.io/contrib/entgql"
	"github.com/stretchr/testify/require"
)

func TestAnnotation(t *testing.T) {
	t.Parallel()
	annotation := entgql.OrderField("foo")
	require.Equal(t, "foo", annotation.OrderField)

	annotation = entgql.Bind()
	require.False(t, annotation.Unbind)
	annotation = entgql.Unbind()
	require.True(t, annotation.Unbind)
	require.Empty(t, annotation.Mapping)

	names := []string{"foo", "bar", "baz"}
	annotation = entgql.MapsTo(names...)
	require.True(t, annotation.Unbind)
	require.ElementsMatch(t, names, annotation.Mapping)

	descr := "This is a test description"
	annotation = entgql.Description(descr)
	require.Equal(t, descr, annotation.Description)
}

func TestAnnotationDecode(t *testing.T) {
	ann := &entgql.Annotation{}
	err := ann.Decode(map[string]interface{}{})
	require.NoError(t, err)
	require.Equal(t, ann, &entgql.Annotation{})
	ann = &entgql.Annotation{}
	err = ann.Decode(map[string]interface{}{
		"OrderField":  "NAME",
		"Unbind":      true,
		"Mapping":     []string{"f1", "f2"},
		"Skip":        entgql.SkipAll,
		"Description": "This is a test description",
	})
	require.NoError(t, err)
	require.Equal(t, ann, &entgql.Annotation{
		OrderField:  "NAME",
		Unbind:      true,
		Mapping:     []string{"f1", "f2"},
		Skip:        entgql.SkipAll,
		Description: "This is a test description",
	})
	err = ann.Decode("invalid")
	require.NotNil(t, err)
	require.Equal(t, err.Error(), "json: cannot unmarshal string into Go value of type entgql.Annotation")
}
