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
	"testing"

	"entgo.io/contrib/entproto"
	"entgo.io/contrib/entproto/internal/todo/ent/enttest"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestPonyService_BatchCreate(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()
	svc := NewPonyService(client)
	ctx := context.Background()

	// Create requests
	var requests []*CreatePonyRequest
	for i := 0; i < (entproto.MaxBatchCreateSize*2)+5; i++ {
		request := &CreatePonyRequest{
			Pony: &Pony{
				Name: fmt.Sprintf("Pony%d", i),
			},
		}
		requests = append(requests, request)
	}

	// Valid request
	resp, err := svc.BatchCreate(ctx, &BatchCreatePoniesRequest{
		Requests: requests[:entproto.MaxBatchCreateSize],
	})
	require.NoError(t, err)
	// Check number of entities returned. Should be max batch create size
	require.EqualValues(t, entproto.MaxBatchCreateSize, len(resp.Ponies))
	// Check unique values of returned entities
	for i, entry := range resp.Ponies {
		require.EqualValues(t, fmt.Sprintf("Pony%d", i), entry.Name)
	}

	// Invalid batch size
	resp, err = svc.BatchCreate(ctx, &BatchCreatePoniesRequest{
		Requests: requests,
	})
	require.Nil(t, resp)
	respStatus, ok := status.FromError(err)
	require.True(t, ok, "expected a gRPC status error")
	require.EqualValues(t, respStatus.Code(), codes.InvalidArgument)
}
