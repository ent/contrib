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

package altdir

import (
	"context"
	"testing"

	"entgo.io/contrib/entproto/internal/altdir/ent/enttest"
	"entgo.io/contrib/entproto/internal/altdir/ent/v1/api/entpb"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

func TestService(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()
	svc := entpb.NewUserService(client)
	ctx := context.Background()
	_, err := svc.Create(ctx, &entpb.CreateUserRequest{
		User: &entpb.User{
			Name: "a8m",
		},
	})
	require.NoError(t, err)
	uc := client.User.Query().CountX(ctx)
	require.Equal(t, 1, uc)
}
