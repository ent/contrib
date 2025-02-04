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

package entpb

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"entgo.io/contrib/entproto"
	"entgo.io/contrib/entproto/internal/todo/ent/user"

	"entgo.io/contrib/entproto/internal/todo/ent"
	"entgo.io/contrib/entproto/internal/todo/ent/enttest"
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
		Status:     User_STATUS_ACTIVE,
		ExternalId: 1,
		Group: &Group{
			Id: int64(group.ID),
		},
		CrmId:          crmID,
		Attachment:     &Attachment{Id: attachmentID},
		Banned:         true,
		HeightInCm:     170.18,
		AccountBalance: 2000.50,
		Labels:         []string{"member", "production"},
		OmitPrefix:     User_BAR,
		MimeType:       User_MIME_TYPE_IMAGE_XML_SVG,
		Int32S:         []int32{1, 2, 3},
		Int64S:         []int64{1, 2, 3},
	}
	created, err := svc.Create(ctx, &CreateUserRequest{
		User: inputUser,
	})
	require.NoError(t, err)
	require.EqualValues(t, created.Status, User_STATUS_ACTIVE)

	fromDB := client.User.GetX(ctx, created.Id)
	require.EqualValues(t, inputUser.UserName, fromDB.UserName)
	require.EqualValues(t, inputUser.Joined.AsTime().Unix(), fromDB.Joined.Unix())
	require.EqualValues(t, inputUser.Exp, fromDB.Exp)
	require.EqualValues(t, inputUser.Points, fromDB.Points)
	require.EqualValues(t, inputUser.Status.String(), strings.ToUpper("STATUS_"+string(fromDB.Status)))
	require.EqualValues(t, inputUser.Banned, fromDB.Banned)
	require.EqualValues(t, inputUser.HeightInCm, fromDB.HeightInCm)
	require.EqualValues(t, inputUser.AccountBalance, fromDB.AccountBalance)
	require.EqualValues(t, inputUser.Labels, fromDB.Labels)
	require.EqualValues(t, inputUser.Int64S, fromDB.Int64s)
	require.EqualValues(t, inputUser.Int32S, fromDB.Int32s)
	require.EqualValues(t, inputUser.MimeType.String(), strings.ToUpper("MIME_TYPE_"+regexp.MustCompile("[^a-zA-Z0-9_]+").ReplaceAllString(string(fromDB.MimeType), "_")))

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
		SetHeightInCm(170.18).
		SetAccountBalance(2000.50).
		SetLabels([]string{"on", "off"}).
		SetOmitPrefix(user.OmitPrefixBar).
		SetMimeType(user.MimeTypeSvg).
		SaveX(ctx)
	get, err := svc.Get(ctx, &GetUserRequest{
		Id: created.ID,
	})
	require.NoError(t, err)
	require.EqualValues(t, created.UserName, get.UserName)
	require.EqualValues(t, created.Exp, get.Exp)
	require.EqualValues(t, created.Joined.Unix(), get.Joined.AsTime().Unix())
	require.EqualValues(t, created.Points, get.Points)
	require.EqualValues(t, created.HeightInCm, get.HeightInCm)
	require.EqualValues(t, created.AccountBalance, get.AccountBalance)
	require.EqualValues(t, created.Labels, get.Labels)
	require.EqualValues(t, User_BAR, get.OmitPrefix)
	require.EqualValues(t, User_MIME_TYPE_IMAGE_XML_SVG, get.MimeType)
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
		SetOmitPrefix(user.OmitPrefixBar).
		SetMimeType(user.MimeTypeSvg).
		SaveX(ctx)
	d, err := svc.Delete(ctx, &DeleteUserRequest{
		Id: created.ID,
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
		SetHeightInCm(170.18).
		SetAccountBalance(2000.50).
		SetLabels(nil).
		SetOmitPrefix(user.OmitPrefixFoo).
		SetMimeType(user.MimeTypeSvg).
		SaveX(ctx)

	attachmentID, err := attachment.ID.MarshalBinary()
	require.NoError(t, err)
	group := client.Group.Create().SetName("managers").SaveX(ctx)
	crmID, err := created.CrmID.MarshalBinary()
	require.NoError(t, err, "Converting UUID to Bytes: %v", crmID)

	inputUser := &User{
		Id:         created.ID,
		UserName:   "rotemtam",
		Joined:     timestamppb.Now(),
		Exp:        999,
		Points:     999,
		ExternalId: 1,
		Status:     User_STATUS_ACTIVE,
		Group: &Group{
			Id: int64(group.ID),
		},
		Attachment: &Attachment{
			Id: attachmentID,
		},
		CrmId:          crmID,
		HeightInCm:     175.18,
		AccountBalance: 5000.75,
		OmitPrefix:     User_FOO,
		MimeType:       User_MIME_TYPE_IMAGE_PNG,
	}
	updated, err := svc.Update(ctx, &UpdateUserRequest{
		User: inputUser,
	})
	require.NoError(t, err)
	require.EqualValues(t, inputUser.Exp, updated.Exp)

	afterUpd := client.User.GetX(ctx, created.ID)
	require.EqualValues(t, inputUser.Exp, afterUpd.Exp)
	require.EqualValues(t, user.OmitPrefixFoo, afterUpd.OmitPrefix)
	require.EqualValues(t, user.MimeTypePng, afterUpd.MimeType)
}

func TestUserService_List(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()
	svc := NewUserService(client)
	ctx := context.Background()

	// Create test entries
	for i := 0; i < (entproto.MaxPageSize*2)+5; i++ {
		_ = client.User.Create().
			SetUserName(fmt.Sprintf("User%d", i)).
			SetExternalID(i).
			SetJoined(time.Now()).
			SetExp(1000).
			SetPoints(10).
			SetStatus("pending").
			SetCrmID(uuid.New()).
			SetCustomPb(1).
			SetLabels(nil).
			SetOmitPrefix(user.OmitPrefixBar).
			SetMimeType(user.MimeTypeSvg).
			SaveX(ctx)
	}

	// First page
	resp, err := svc.List(ctx, &ListUserRequest{
		PageSize: entproto.MaxPageSize * 2,
	})
	require.NoError(t, err)
	// Check number of entities returned. Should be max page size
	require.EqualValues(t, entproto.MaxPageSize, len(resp.UserList))
	// Check unique values of returned entities
	for entryIdx, entry := range resp.UserList {
		entityID := ((entproto.MaxPageSize * 2) + 5) - (entryIdx + 1)
		require.EqualValues(t, fmt.Sprintf("User%d", entityID), entry.UserName)
		require.EqualValues(t, entityID, entry.ExternalId)
	}

	// Second page
	resp, err = svc.List(ctx, &ListUserRequest{
		PageToken: resp.NextPageToken,
	})
	require.NoError(t, err)
	// Check number of entities returned. Should be max page size which is the default
	require.EqualValues(t, entproto.MaxPageSize, len(resp.UserList))
	// Check that we actually got values from the second page
	for entryIdx, entry := range resp.UserList {
		entityID := (entproto.MaxPageSize + 5) - (entryIdx + 1)
		require.EqualValues(t, fmt.Sprintf("User%d", entityID), entry.UserName)
		require.EqualValues(t, entityID, entry.ExternalId)
	}

	// Final page
	resp, err = svc.List(ctx, &ListUserRequest{
		PageToken: resp.NextPageToken,
	})
	require.NoError(t, err)
	// Check number of entities returned
	require.EqualValues(t, 5, len(resp.UserList))
	// Check that no next page token was returned
	require.EqualValues(t, "", resp.NextPageToken)

	// Invalid page size
	resp, err = svc.List(ctx, &ListUserRequest{
		PageSize: -1,
	})
	require.Nil(t, resp)
	respStatus, ok := status.FromError(err)
	require.True(t, ok, "expected a gRPC status error")
	require.EqualValues(t, respStatus.Code(), codes.InvalidArgument)

	// Invalid page token
	resp, err = svc.List(ctx, &ListUserRequest{
		PageToken: "INVALID PAGE TOKEN",
	})
	require.Nil(t, resp)
	respStatus, ok = status.FromError(err)
	require.True(t, ok, "expected a gRPC status error")
	require.EqualValues(t, respStatus.Code(), codes.InvalidArgument)
}

func TestUserService_BatchCreate(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()
	svc := NewUserService(client)
	ctx := context.Background()

	// Create requests
	var requests []*CreateUserRequest
	for i := 0; i < (entproto.MaxBatchCreateSize*2)+5; i++ {
		crmid, _ := uuid.New().MarshalBinary()
		request := &CreateUserRequest{
			User: &User{
				UserName:   fmt.Sprintf("User%d", i),
				ExternalId: int64(i),
				Joined:     timestamppb.Now(),
				Exp:        1000,
				Points:     10,
				CrmId:      crmid,
				CustomPb:   1,
				Labels:     nil,
				Status:     User_STATUS_ACTIVE,
				OmitPrefix: User_BAR,
				MimeType:   User_MIME_TYPE_IMAGE_PNG,
			},
		}
		requests = append(requests, request)
	}

	// Valid request
	resp, err := svc.BatchCreate(ctx, &BatchCreateUsersRequest{
		Requests: requests[:entproto.MaxBatchCreateSize],
	})
	require.NoError(t, err)
	// Check number of entities returned. Should be max batch create size
	require.EqualValues(t, entproto.MaxBatchCreateSize, len(resp.Users))
	// Check unique values of returned entities
	for i, entry := range resp.Users {
		require.EqualValues(t, fmt.Sprintf("User%d", i), entry.UserName)
		require.EqualValues(t, i, entry.ExternalId)
	}

	// Invalid batch size
	resp, err = svc.BatchCreate(ctx, &BatchCreateUsersRequest{
		Requests: requests,
	})
	require.Nil(t, resp)
	respStatus, ok := status.FromError(err)
	require.True(t, ok, "expected a gRPC status error")
	require.EqualValues(t, respStatus.Code(), codes.InvalidArgument)
}
