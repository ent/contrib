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
//
// Code generated by entc, DO NOT EDIT.

package directiveexample

import (
	"entgo.io/ent/dialect/sql"
)

const (
	// Label holds the string label denoting the directiveexample type in the database.
	Label = "directive_example"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldOnTypeField holds the string denoting the on_type_field field in the database.
	FieldOnTypeField = "on_type_field"
	// FieldOnMutationFields holds the string denoting the on_mutation_fields field in the database.
	FieldOnMutationFields = "on_mutation_fields"
	// FieldOnMutationCreate holds the string denoting the on_mutation_create field in the database.
	FieldOnMutationCreate = "on_mutation_create"
	// FieldOnMutationUpdate holds the string denoting the on_mutation_update field in the database.
	FieldOnMutationUpdate = "on_mutation_update"
	// FieldOnAllFields holds the string denoting the on_all_fields field in the database.
	FieldOnAllFields = "on_all_fields"
	// Table holds the table name of the directiveexample in the database.
	Table = "directive_examples"
)

// Columns holds all SQL columns for directiveexample fields.
var Columns = []string{
	FieldID,
	FieldOnTypeField,
	FieldOnMutationFields,
	FieldOnMutationCreate,
	FieldOnMutationUpdate,
	FieldOnAllFields,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

// OrderOption defines the ordering options for the DirectiveExample queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByOnTypeField orders the results by the on_type_field field.
func ByOnTypeField(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldOnTypeField, opts...).ToFunc()
}

// ByOnMutationFields orders the results by the on_mutation_fields field.
func ByOnMutationFields(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldOnMutationFields, opts...).ToFunc()
}

// ByOnMutationCreate orders the results by the on_mutation_create field.
func ByOnMutationCreate(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldOnMutationCreate, opts...).ToFunc()
}

// ByOnMutationUpdate orders the results by the on_mutation_update field.
func ByOnMutationUpdate(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldOnMutationUpdate, opts...).ToFunc()
}

// ByOnAllFields orders the results by the on_all_fields field.
func ByOnAllFields(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldOnAllFields, opts...).ToFunc()
}
