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
//
// Code generated by entc, DO NOT EDIT.

package ent

import (
	"time"

	"github.com/bionicstork/contrib/entgql/internal/todouuid/ent/category"
	"github.com/bionicstork/contrib/entgql/internal/todouuid/ent/schema"
	"github.com/bionicstork/contrib/entgql/internal/todouuid/ent/todo"
	"github.com/bionicstork/contrib/entgql/internal/todouuid/ent/verysecret"
	"github.com/google/uuid"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	categoryMixin := schema.Category{}.Mixin()
	categoryMixinFields0 := categoryMixin[0].Fields()
	_ = categoryMixinFields0
	categoryFields := schema.Category{}.Fields()
	_ = categoryFields
	// categoryDescText is the schema descriptor for text field.
	categoryDescText := categoryMixinFields0[0].Descriptor()
	// category.TextValidator is a validator for the "text" field. It is called by the builders before save.
	category.TextValidator = categoryDescText.Validators[0].(func(string) error)
	// categoryDescID is the schema descriptor for id field.
	categoryDescID := categoryFields[0].Descriptor()
	// category.DefaultID holds the default value on creation for the id field.
	category.DefaultID = categoryDescID.Default.(func() uuid.UUID)
	todoMixin := schema.Todo{}.Mixin()
	todoMixinFields0 := todoMixin[0].Fields()
	_ = todoMixinFields0
	todoFields := schema.Todo{}.Fields()
	_ = todoFields
	// todoDescCreatedAt is the schema descriptor for created_at field.
	todoDescCreatedAt := todoMixinFields0[0].Descriptor()
	// todo.DefaultCreatedAt holds the default value on creation for the created_at field.
	todo.DefaultCreatedAt = todoDescCreatedAt.Default.(func() time.Time)
	// todoDescPriority is the schema descriptor for priority field.
	todoDescPriority := todoMixinFields0[2].Descriptor()
	// todo.DefaultPriority holds the default value on creation for the priority field.
	todo.DefaultPriority = todoDescPriority.Default.(int)
	// todoDescText is the schema descriptor for text field.
	todoDescText := todoMixinFields0[3].Descriptor()
	// todo.TextValidator is a validator for the "text" field. It is called by the builders before save.
	todo.TextValidator = todoDescText.Validators[0].(func(string) error)
	// todoDescID is the schema descriptor for id field.
	todoDescID := todoFields[0].Descriptor()
	// todo.DefaultID holds the default value on creation for the id field.
	todo.DefaultID = todoDescID.Default.(func() uuid.UUID)
	verysecretFields := schema.VerySecret{}.Fields()
	_ = verysecretFields
	// verysecretDescID is the schema descriptor for id field.
	verysecretDescID := verysecretFields[0].Descriptor()
	// verysecret.DefaultID holds the default value on creation for the id field.
	verysecret.DefaultID = verysecretDescID.Default.(func() uuid.UUID)
}
