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
	"database/sql/driver"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"entgo.io/contrib/entoas/spec"
	"entgo.io/ent/dialect"
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/entc/load"
	entfield "entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestOASType(t *testing.T) {
	t.Parallel()
	for d, ex := range map[*entfield.Descriptor]*spec.Type{
		// Numeric
		entfield.Int("int").Descriptor():         _int32,
		entfield.Int8("int8").Descriptor():       _int32,
		entfield.Int16("int16").Descriptor():     _int32,
		entfield.Int32("int32").Descriptor():     _int32,
		entfield.Int64("int64").Descriptor():     _int64,
		entfield.Uint("uint").Descriptor():       _int32,
		entfield.Uint8("uint8").Descriptor():     _int32,
		entfield.Uint16("uint16").Descriptor():   _int32,
		entfield.Uint32("uint32").Descriptor():   _int32,
		entfield.Uint64("uint64").Descriptor():   _int64,
		entfield.Float32("float32").Descriptor(): _float,
		entfield.Float("float64").Descriptor():   _double,
		// Basic
		entfield.String("string").Descriptor():                  _string,
		entfield.Bool("bool").Descriptor():                      _bool,
		entfield.UUID("uuid", uuid.Nil).Descriptor():            _string,
		entfield.Time("time").Descriptor():                      _dateTime,
		entfield.Text("text").Descriptor():                      _string,
		entfield.Enum("state").Values("on", "off").Descriptor(): _string,
		// List
		entfield.Strings("strings").Descriptor(): arr(_string),
		entfield.Ints("ints").Descriptor():       arr(_int32),
		entfield.Floats("floats").Descriptor():   arr(_double),
		entfield.Bytes("bytes").Descriptor():     _bytes,
		// Custom
		entfield.JSON("nicknames", []string{}).Descriptor(): arr(_string),
		entfield.JSON("json_slice", []http.Dir{}).
			Annotations(OASType(arr(_string))).Descriptor(): arr(_string),
		entfield.JSON("json_obj", url.URL{}).
			Annotations(OASType(_string)).Descriptor(): _string,
		entfield.Other("other", &Link{}).
			SchemaType(map[string]string{dialect.Postgres: "varchar"}).
			Default(DefaultLink()).
			Annotations(OASType(_string)).
			Descriptor(): _string,
	} {
		t.Run(d.Name, func(t *testing.T) {
			f, err := load.NewField(d)
			require.NoError(t, err)
			gf := &gen.Field{Name: f.Name, Type: f.Info, Annotations: f.Annotations}
			ac, err := oasType(gf)
			if ex == _empty {
				require.Error(t, err)
				require.EqualError(t, err, fmt.Sprintf(
					"no OAS-type exists for type %q of field %s",
					gf.Type.String(),
					gf.StructField(),
				))
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

type Link struct {
	*url.URL
}

func DefaultLink() *Link {
	u, _ := url.Parse("127.0.0.1")
	return &Link{URL: u}
}

// Scan implements the Scanner interface.
func (l *Link) Scan(value interface{}) (err error) {
	switch v := value.(type) {
	case nil:
	case []byte:
		l.URL, err = url.Parse(string(v))
	case string:
		l.URL, err = url.Parse(v)
	default:
		err = fmt.Errorf("unexpected type %T", v)
	}
	return
}

// Value implements the driver Valuer interface.
func (l Link) Value() (driver.Value, error) {
	if l.URL == nil {
		return nil, nil
	}
	return l.String(), nil
}

func arr(t *spec.Type) *spec.Type { return &spec.Type{Type: "array", Items: t} }
