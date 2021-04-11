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

package entpb

import (
	"context"
	"strings"
	"testing"
	"time"

	"entgo.io/contrib/entproto/internal/todo/ent"
	"entgo.io/contrib/entproto/internal/todo/ent/enttest"
	"entgo.io/contrib/entproto/runtime"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestUserService_Create(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()
	svc := NewUserService(client)
	ctx := context.Background()
	group := client.Group.Create().SetName("managers").SaveX(ctx)
	attachment := client.Attachment.Create().SaveX(ctx)
	crmID, err := uuid.New().MarshalBinary()
	require.NoError(t, err)
	attachmentID, err := attachment.ID.MarshalBinary()
	require.NoError(t, err)
	inputUser := &User{
		UserName:   "rotemtam",
		Joined:     timestamppb.Now(),
		Exp:        100,
		Points:     1000,
		Status:     User_ACTIVE,
		ExternalId: 1,
		Group: &Group{
			Id: int32(group.ID),
		},
		CrmId:      crmID,
		Attachment: &Attachment{Id: attachmentID},
		Banned:     true,
	}
	created, err := svc.Create(ctx, &CreateUserRequest{
		User: inputUser,
	})
	require.NoError(t, err)
	require.EqualValues(t, created.Status, User_ACTIVE)

	fromDB := client.User.GetX(ctx, int(created.Id))
	require.EqualValues(t, inputUser.UserName, fromDB.UserName)
	require.EqualValues(t, inputUser.Joined.AsTime().Unix(), fromDB.Joined.Unix())
	require.EqualValues(t, inputUser.Exp, fromDB.Exp)
	require.EqualValues(t, inputUser.Points, fromDB.Points)
	require.EqualValues(t, inputUser.Status.String(), strings.ToUpper(string(fromDB.Status)))
	require.EqualValues(t, inputUser.Banned, fromDB.Banned)

	// preexisting user
	_, err = svc.Create(ctx, &CreateUserRequest{
		User: inputUser,
	})
	respStatus, ok := status.FromError(err)
	require.True(t, ok, "expected a gRPC status error")
	require.EqualValues(t, respStatus.Code(), codes.AlreadyExists)
}

func TestUserService_Get(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()
	svc := NewUserService(client)
	ctx := context.Background()
	created := client.User.Create().
		SetUserName("rotemtam").
		SetJoined(time.Now()).
		SetPoints(10).
		SetExp(1000).
		SetStatus("pending").
		SetExternalID(1).
		SetCrmID(uuid.New()).
		SetCustomPb(1).
		SaveX(ctx)
	get, err := svc.Get(ctx, &GetUserRequest{
		Id: int32(created.ID),
	})
	require.NoError(t, err)
	require.EqualValues(t, created.UserName, get.UserName)
	require.EqualValues(t, created.Exp, get.Exp)
	require.EqualValues(t, created.Joined.Unix(), get.Joined.AsTime().Unix())
	require.EqualValues(t, created.Points, get.Points)
	get, err = svc.Get(ctx, &GetUserRequest{
		Id: 1000,
	})
	require.Nil(t, get)
	respStatus, ok := status.FromError(err)
	require.True(t, ok, "expected a gRPC status error")
	require.EqualValues(t, respStatus.Code(), codes.NotFound)
}

func TestUserService_Delete(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()
	svc := NewUserService(client)
	ctx := context.Background()
	created := client.User.Create().
		SetUserName("rotemtam").
		SetJoined(time.Now()).
		SetPoints(10).
		SetExp(1000).
		SetStatus("pending").
		SetExternalID(1).
		SetCrmID(uuid.New()).
		SetCustomPb(1).
		SaveX(ctx)
	d, err := svc.Delete(ctx, &DeleteUserRequest{
		Id: int32(created.ID),
	})
	require.NoError(t, err)
	require.NotNil(t, d)
	_, err = client.User.Get(ctx, created.ID)
	require.True(t, ent.IsNotFound(err))

	d, err = svc.Delete(ctx, &DeleteUserRequest{
		Id: 1000,
	})
	require.Nil(t, d)
	respStatus, ok := status.FromError(err)
	require.True(t, ok, "expected a gRPC status error")
	require.EqualValues(t, respStatus.Code(), codes.NotFound)
}

func TestUserService_Update(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()
	svc := NewUserService(client)
	ctx := context.Background()
	attachment := client.Attachment.Create().SaveX(ctx)
	created := client.User.Create().
		SetUserName("rotemtam").
		SetJoined(time.Now()).
		SetPoints(10).
		SetExp(1000).
		SetStatus("pending").
		SetExternalID(1).
		SetCrmID(uuid.New()).
		SetCustomPb(1).
		SaveX(ctx)

	attachmentID, err := attachment.ID.MarshalBinary()
	require.NoError(t, err)
	group := client.Group.Create().SetName("managers").SaveX(ctx)
	inputUser := &User{
		Id:         int32(created.ID),
		UserName:   "rotemtam",
		Joined:     timestamppb.Now(),
		Exp:        999,
		Points:     999,
		ExternalId: 1,
		Status:     User_ACTIVE,
		Group: &Group{
			Id: int32(group.ID),
		},
		Attachment: &Attachment{
			Id: attachmentID,
		},
		CrmId: runtime.MustExtractUUIDBytes(created.CrmID),
	}
	updated, err := svc.Update(ctx, &UpdateUserRequest{
		User: inputUser,
	})
	require.NoError(t, err)
	require.EqualValues(t, inputUser.Exp, updated.Exp)

	afterUpd := client.User.GetX(ctx, created.ID)
	require.EqualValues(t, inputUser.Exp, afterUpd.Exp)
}
