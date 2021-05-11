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

package plugin

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema/field"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/vektah/gqlparser/v2/ast"
	"testing"
)

func TestTypeFields(t *testing.T) {
	e, err := New(&gen.Graph{
		Config: &gen.Config{},
	})
	require.NoError(t, err)
	fields, err := e.typeFields(&gen.Type{
		ID: &gen.Field{
			Name: "Id",
			Type: &field.TypeInfo{
				Type: field.TypeInt,
			},
		},
		Fields: []*gen.Field{
			{
				Name: "name",
				Type: &field.TypeInfo{
					Type: field.TypeString,
				},
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, ast.FieldList{
		{
			Name: "id",
			Type: namedType("ID", false),
		},
		{
			Name: "name",
			Type: namedType("String", false),
		},
	}, fields)
}

func TestFields(t *testing.T) {
	testCases := []struct {
		name         string
		fieldType    field.Type
		userDefined  string
		expectedType string
		err          error
	}{
		{"firstname", field.TypeString, "", "String", nil},
		{"age", field.TypeInt, "", "Int", nil},
		{"f", field.TypeFloat64, "", "Float", nil},
		{"f", field.TypeFloat32, "", "Float", nil},
		{"status", field.TypeEnum, "", "Status", nil},
		{"status", field.TypeEnum, "StatusEnum", "StatusEnum", nil},
		{"status", field.TypeEnum, "StatusEnum", "StatusEnum", nil},
		{"timestamp", field.TypeTime, "", "Time", nil},
		{"active", field.TypeBool, "", "Boolean", nil},
		{"data", field.TypeBytes, "", "", fmt.Errorf("bytes type not implemented")},
		{"json", field.TypeJSON, "", "", fmt.Errorf("json type not implemented")},
		{"other", field.TypeOther, "", "Invalid", fmt.Errorf("other type must have typed defined")},
	}
	e, err := New(&gen.Graph{
		Config: &gen.Config{
			Annotations: map[string]interface{}{
				annotationName: entgql.Annotation{
					GqlScalarMappings: map[string]string{
						"Time": "Time",
					}},
			},
		},
	})
	require.NoError(t, err)
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s(%s)", tc.name, tc.fieldType.ConstName()), func(t *testing.T) {
			f, err := e.entTypToGqlType(&gen.Field{
				Name: tc.name,
				Type: &field.TypeInfo{
					Type: tc.fieldType,
				},
				Nillable: true,
			}, false, tc.userDefined)
			require.Equal(t, tc.err, err)
			if tc.err == nil {
				require.Equal(t, tc.expectedType, f.String())
			}
			f, err = e.entTypToGqlType(&gen.Field{
				Name: tc.name,
				Type: &field.TypeInfo{
					Type: tc.fieldType,
				},
			}, false, tc.userDefined)
			require.Equal(t, tc.err, err)
			if tc.err == nil {
				require.Equal(t, tc.expectedType+"!", f.String())
			}
		})
	}
}

func TestIdField(t *testing.T) {
	e, err := New(&gen.Graph{
		Config: &gen.Config{},
	})
	require.NoError(t, err)
	f, err := e.entTypToGqlType(&gen.Field{
		Name: "id",
		Type: &field.TypeInfo{
			Type: field.TypeInt,
		},
	}, true, "")
	require.NoError(t, err)
	require.Equal(t, "ID!", f.String())
}
