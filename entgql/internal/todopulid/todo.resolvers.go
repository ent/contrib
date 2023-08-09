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

package todopulid

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"entgo.io/contrib/entgql"
	"entgo.io/contrib/entgql/internal/todopulid/ent"
	"entgo.io/contrib/entgql/internal/todopulid/ent/schema/pulid"
	"entgo.io/contrib/entgql/internal/todopulid/ent/todo"
)

func (r *categoryResolver) TodosCount(ctx context.Context, obj *ent.Category) (*int, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateCategory(ctx context.Context, input ent.CreateCategoryInput) (*ent.Category, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateTodo(ctx context.Context, input ent.CreateTodoInput) (*ent.Todo, error) {
	return ent.FromContext(ctx).Todo.
		Create().
		SetInput(input).
		Save(ctx)
}

func (r *mutationResolver) UpdateTodo(ctx context.Context, id pulid.ID, input ent.UpdateTodoInput) (*ent.Todo, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ClearTodos(ctx context.Context) (int, error) {
	client := ent.FromContext(ctx)
	return client.Todo.
		Delete().
		Exec(ctx)
}

func (r *mutationResolver) UpdateFriendship(ctx context.Context, id pulid.ID, input UpdateFriendshipInput) (*ent.Friendship, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Ping(ctx context.Context) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) TodosWithJoins(ctx context.Context, after *entgql.Cursor[pulid.ID], first *int, before *entgql.Cursor[pulid.ID], last *int, orderBy []*ent.TodoOrder, where *ent.TodoWhereInput) (*ent.TodoConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *todoResolver) ExtendedField(ctx context.Context, obj *ent.Todo) (*string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *createCategoryInputResolver) CreateTodos(ctx context.Context, obj *ent.CreateCategoryInput, data []*ent.CreateTodoInput) error {
	panic(fmt.Errorf("not implemented"))
}

func (r *todoWhereInputResolver) CreatedToday(ctx context.Context, obj *ent.TodoWhereInput, data *bool) error {
	if data == nil {
		return nil
	}

	startOfDay := time.Now().Truncate(24 * time.Hour)
	endOfDay := startOfDay.Add(24*time.Hour - 1)
	if *data {
		obj.AddPredicates(todo.And(todo.CreatedAtGTE(startOfDay), todo.CreatedAtLTE(endOfDay)))
	} else {
		obj.AddPredicates(todo.Or(todo.CreatedAtLT(startOfDay), todo.CreatedAtGT(endOfDay)))
	}

	return nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
