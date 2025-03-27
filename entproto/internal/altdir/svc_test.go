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
	require := require.New(t)

	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	userSvc := entpb.NewUserService(client)
	accountSvc := entpb.NewAccountService(client)
	ctx := context.Background()

	owner, err := userSvc.Create(ctx, &entpb.CreateUserRequest{
		User: &entpb.User{
			Name: "a8m",
		},
	})
	require.NoError(err)
	uc := client.User.Query().CountX(ctx)
	require.Equal(1, uc)

	other, err := userSvc.Create(ctx, &entpb.CreateUserRequest{
		User: &entpb.User{
			Name: "Anton Chigurh",
		},
	})
	require.NoError(err)

	account, err := accountSvc.Create(ctx, &entpb.CreateAccountRequest{
		Account: &entpb.Account{
			Owner: owner,
		},
	})
	require.NoError(err)

	_, err = accountSvc.Update(ctx, &entpb.UpdateAccountRequest{
		Account: &entpb.Account{
			Id:    account.Id,
			Owner: other,
		},
	})
	require.NoError(err)

	account, err = accountSvc.Get(ctx, &entpb.GetAccountRequest{
		Id:   account.Id,
		View: entpb.GetAccountRequest_WITH_EDGE_IDS,
	})
	require.NoError(err)
	require.Equal(owner.Id, account.Owner.Id)
}
