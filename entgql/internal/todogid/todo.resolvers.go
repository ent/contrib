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

package todogid

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"entgo.io/contrib/entgql/internal/todogid/ent"
)

func (r *mutationResolver) CreateTodo(ctx context.Context, input *CreateTodoInput) (*CreateTodoPayload, error) {
	// Extract the "real" user-id in case the global-id is used in inputs.
	id, err := ent.IntFromGlobalID(input.User)
	if err != nil {
		return nil, err
	}
	t, err := ent.FromContext(ctx).Todo.Create().SetText(input.Text).SetOwnerID(id).Save(ctx)
	if err != nil {
		return nil, err
	}
	return &CreateTodoPayload{
		Todo:             t,
		ClientMutationID: input.ClientMutationID,
	}, nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, input *CreateUserInput) (*CreateUserPayload, error) {
	u, err := ent.FromContext(ctx).User.Create().SetName(input.Name).Save(ctx)
	if err != nil {
		return nil, err
	}
	return &CreateUserPayload{
		User:             u,
		ClientMutationID: input.ClientMutationID,
	}, nil
}

func (r *mutationResolver) CreatePost(ctx context.Context, input *CreatePostInput) (*CreatePostPayload, error) {
	p, err := ent.FromContext(ctx).Post.Create().SetText(input.Text).Save(ctx)
	if err != nil {
		return nil, err
	}
	return &CreatePostPayload{
		Post:             p,
		ClientMutationID: input.ClientMutationID,
	}, nil
}

func (r *queryResolver) Node(ctx context.Context, id string) (ent.Noder, error) {
	return r.client.FromGlobalID(ctx, id)
}

func (r *queryResolver) Nodes(ctx context.Context, ids []string) ([]ent.Noder, error) {
	return r.client.FromGlobalIDs(ctx, ids)
}

func (r *queryResolver) Todos(ctx context.Context) ([]*ent.Todo, error) {
	return r.client.Todo.Query().All(ctx)
}

func (r *queryResolver) Users(ctx context.Context) ([]*ent.User, error) {
	return r.client.User.Query().All(ctx)
}

func (r *queryResolver) Posts(ctx context.Context) ([]*ent.Post, error) {
	return r.client.Post.Query().All(ctx)
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
