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

package todo

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.43

import (
	"context"
	"fmt"

	"entgo.io/contrib/entgql"
	"entgo.io/contrib/entgql/internal/todogotype/ent"
)

// TodosCount is the resolver for the todosCount field.
func (r *categoryResolver) TodosCount(ctx context.Context, obj *ent.Category) (*int, error) {
	panic(fmt.Errorf("not implemented"))
}

// CreateCategory is the resolver for the createCategory field.
func (r *mutationResolver) CreateCategory(ctx context.Context, input ent.CreateCategoryInput) (*ent.Category, error) {
	panic(fmt.Errorf("not implemented"))
}

// CreateTodo is the resolver for the createTodo field.
func (r *mutationResolver) CreateTodo(ctx context.Context, input ent.CreateTodoInput) (*ent.Todo, error) {
	panic(fmt.Errorf("not implemented"))
}

// UpdateTodo is the resolver for the updateTodo field.
func (r *mutationResolver) UpdateTodo(ctx context.Context, id string, input ent.UpdateTodoInput) (*ent.Todo, error) {
	panic(fmt.Errorf("not implemented"))
}

// ClearTodos is the resolver for the clearTodos field.
func (r *mutationResolver) ClearTodos(ctx context.Context) (int, error) {
	panic(fmt.Errorf("not implemented"))
}

// UpdateFriendship is the resolver for the updateFriendship field.
func (r *mutationResolver) UpdateFriendship(ctx context.Context, id string, input UpdateFriendshipInput) (*ent.Friendship, error) {
	panic(fmt.Errorf("not implemented"))
}

// Ping is the resolver for the ping field.
func (r *queryResolver) Ping(ctx context.Context) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

// TodosWithJoins is the resolver for the todosWithJoins field.
func (r *queryResolver) TodosWithJoins(ctx context.Context, after *entgql.Cursor[string], first *int, before *entgql.Cursor[string], last *int, orderBy []*ent.TodoOrder, where *ent.TodoWhereInput) (*ent.TodoConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

// ExtendedField is the resolver for the extendedField field.
func (r *todoResolver) ExtendedField(ctx context.Context, obj *ent.Todo) (*string, error) {
	panic(fmt.Errorf("not implemented"))
}

// CreateTodos is the resolver for the createTodos field.
func (r *createCategoryInputResolver) CreateTodos(ctx context.Context, obj *ent.CreateCategoryInput, data []*ent.CreateTodoInput) error {
	panic(fmt.Errorf("not implemented"))
}

// CreatedToday is the resolver for the createdToday field.
func (r *todoWhereInputResolver) CreatedToday(ctx context.Context, obj *ent.TodoWhereInput, data *bool) error {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
