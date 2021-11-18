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

package entoas

import (
	"testing"

	"entgo.io/contrib/entoas/spec"
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/entc/load"
	entfield "entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestOASType(t *testing.T) {
	t.Parallel()
	for d, ex := range map[*entfield.Descriptor]*spec.Type{
		entfield.Bool("bool").Descriptor():           _bool,
		entfield.Bool("bool").Descriptor():           _bool,
		entfield.Time("time").Descriptor():           _dateTime,
		entfield.UUID("uuid", uuid.Nil).Descriptor(): _string,
		entfield.Bytes("bytes").Descriptor():         _empty,
		entfield.Enum("enum").Descriptor():           _string,
		entfield.String("string").Descriptor():       _string,
		entfield.Int8("int8").Descriptor():           _int32,
		entfield.Int16("int16").Descriptor():         _int32,
		entfield.Int32("int32").Descriptor():         _int32,
		entfield.Int("int").Descriptor():             _int32,
		entfield.Int64("int64").Descriptor():         _int64,
		entfield.Uint8("uint8").Descriptor():         _int32,
		entfield.Uint16("uint16").Descriptor():       _int32,
		entfield.Uint32("uint32").Descriptor():       _int32,
		entfield.Uint("uint").Descriptor():           _int32,
		entfield.Uint64("uint64").Descriptor():       _int64,
		entfield.Float32("float32").Descriptor():     _float,
		entfield.Float("float64").Descriptor():       _double,
	} {
		t.Run(d.Name, func(t *testing.T) {
			f, err := load.NewField(d)
			require.NoError(t, err)
			ac, err := oasType(&gen.Field{Type: f.Info})
			if ex == _empty {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, ex, ac)
			}
		})
	}
}

func TestOperation_Title(t *testing.T) {
	t.Parallel()
	require.Equal(t, "Create", OpCreate.Title())
	require.Equal(t, "Read", OpRead.Title())
	require.Equal(t, "Update", OpUpdate.Title())
	require.Equal(t, "Delete", OpDelete.Title())
	require.Equal(t, "List", OpList.Title())
}
