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
	"testing"
	"time"

	"entgo.io/contrib/entproto/internal/todo/ent/enttest"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

func TestMapping(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	ts := time.Now()
	created := client.User.Create().
		SetUserName("rotemtam").
		SetExp(100).
		SetPoints(1000).
		SetStatus("active").
		SetJoined(ts).
		SaveX(context.Background())

	pbUser := toProtoUser(created)
	require.NotNil(t, pbUser)
	require.EqualValues(t, "rotemtam", pbUser.UserName)
	require.EqualValues(t, 100, pbUser.Exp)
	require.EqualValues(t, 1000, pbUser.Points)
	require.EqualValues(t, User_ACTIVE, pbUser.Status)
	require.EqualValues(t, ts.Unix(), pbUser.Joined.AsTime().Unix())
}
