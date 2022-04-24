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

package todo

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"time"

	"entgo.io/contrib/entgql/internal/todo/ent"
	"entgo.io/contrib/entgql/internal/todo/ent/todo"
)

func (r *mutationResolver) CreateTodo(ctx context.Context, todo TodoInput) (*ent.Todo, error) {
	client := ent.FromContext(ctx)
	return client.Todo.
		Create().
		SetStatus(todo.Status).
		SetNillablePriority(todo.Priority).
		SetText(todo.Text).
		SetNillableParentID(todo.Parent).
		Save(ctx)
}

func (r *mutationResolver) ClearTodos(ctx context.Context) (int, error) {
	client := ent.FromContext(ctx)
	return client.Todo.
		Delete().
		Exec(ctx)
}

func (r *queryResolver) Node(ctx context.Context, id int) (ent.Noder, error) {
	return r.client.Noder(ctx, id)
}

func (r *queryResolver) Nodes(ctx context.Context, ids []int) ([]ent.Noder, error) {
	return r.client.Noders(ctx, ids)
}

func (r *queryResolver) Todos(ctx context.Context, after *ent.Cursor, first *int, before *ent.Cursor, last *int, orderBy *ent.TodoOrder, where *ent.TodoWhereInput) (*ent.TodoConnection, error) {
	return r.client.Todo.Query().
		Paginate(ctx, after, first, before, last,
			ent.WithTodoOrder(orderBy),
			ent.WithTodoFilter(where.Filter),
		)
}

func (r *queryResolver) Users(ctx context.Context, after *ent.Cursor, first *int, before *ent.Cursor, last *int, where *ent.UserWhereInput) (*ent.UserConnection, error) {
	return r.client.User.Query().
		Paginate(ctx, after, first, before, last,
			ent.WithUserFilter(where.Filter),
		)
}

func (r *queryResolver) Groups(ctx context.Context, after *ent.Cursor, first *int, before *ent.Cursor, last *int, where *ent.GroupWhereInput) (*ent.GroupConnection, error) {
	return r.client.Group.Query().
		Paginate(ctx, after, first, before, last,
			ent.WithGroupFilter(where.Filter),
		)
}

func (r *todoWhereInputResolver) CreatedToday(ctx context.Context, obj *ent.TodoWhereInput, data *bool) error {
	if data == nil {
		return nil
	}

	startOfDay := time.Now().Truncate(24 * time.Hour)
	endOfDay := startOfDay.Add(24*time.Hour - 1)
	if *data {
		obj.Predicates = append(obj.Predicates, todo.And(todo.CreatedAtGTE(startOfDay), todo.CreatedAtLTE(endOfDay)))
	} else {
		obj.Predicates = append(obj.Predicates, todo.Or(todo.CreatedAtLT(startOfDay), todo.CreatedAtGT(endOfDay)))
	}

	return nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
