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
	"entgo.io/contrib/entgql/internal/todouuid/ent/todo"
	"github.com/google/uuid"
)

// CreateTodoInput represents a mutation input for creating todos.
type CreateTodoInput struct {
	Status     todo.Status
	Priority   *int
	Text       string
	ParentID   *uuid.UUID
	ChildIDs   []uuid.UUID
	CategoryID *uuid.UUID
	SecretID   *uuid.UUID
}

// Mutate applies the CreateTodoInput on the TodoMutation builder.
func (i *CreateTodoInput) Mutate(m *TodoMutation) {
	m.SetStatus(i.Status)
	if v := i.Priority; v != nil {
		m.SetPriority(*v)
	}
	m.SetText(i.Text)
	if v := i.ParentID; v != nil {
		m.SetParentID(*v)
	}
	if v := i.ChildIDs; len(v) > 0 {
		m.AddChildIDs(v...)
	}
	if v := i.CategoryID; v != nil {
		m.SetCategoryID(*v)
	}
	if v := i.SecretID; v != nil {
		m.SetSecretID(*v)
	}
}

// SetInput applies the change-set in the CreateTodoInput on the TodoCreate builder.
func (c *TodoCreate) SetInput(i CreateTodoInput) *TodoCreate {
	i.Mutate(c.Mutation())
	return c
}
