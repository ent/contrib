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

package todoplugin

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"entgo.io/contrib/entgql/internal/todoplugin/ent"
	"entgo.io/contrib/entgql/internal/todoplugin/ent/category"
	"entgo.io/contrib/entgql/internal/todoplugin/ent/todo"
)

func (r *mutationResolver) CreateTodo(ctx context.Context, input TodoInput) (*ent.Todo, error) {
	client := ent.FromContext(ctx)
	return client.Todo.
		Create().
		SetStatus(input.Status).
		SetNillablePriority(input.Priority).
		SetText(input.Text).
		SetNillableParentID(input.Parent).
		Save(ctx)
}

func (r *mutationResolver) ClearTodos(ctx context.Context) (int, error) {
	client := ent.FromContext(ctx)
	return client.Todo.
		Delete().
		Exec(ctx)
}

func (r *queryResolver) Ping(ctx context.Context) (string, error) {
	return "pong", nil
}

func (r *todoResolver) Category(ctx context.Context, obj *ent.Todo) (*ent.Category, error) {
	e, err := r.client.Category.
		Query().
		Where(category.HasTodosWith(todo.ID(obj.ID))).
		First(ctx)
	return e, ent.MaskNotFound(err)
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
