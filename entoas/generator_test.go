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

package entoas

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"entgo.io/ent/dialect"
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/entc/load"
	entfield "entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/ogen-go/ogen"
	"github.com/stretchr/testify/require"
)

func TestOgenSchema(t *testing.T) {
	t.Parallel()
	for d, ex := range map[*entfield.Descriptor]*ogen.Schema{
		// Numeric
		entfield.Int("int").Descriptor():         ogen.Int(),
		entfield.Int8("int8").Descriptor():       ogen.Int32().SetMinimum(&min8).SetMaximum(&max8),
		entfield.Int16("int16").Descriptor():     ogen.Int32().SetMinimum(&min16).SetMaximum(&max16),
		entfield.Int32("int32").Descriptor():     ogen.Int32(),
		entfield.Int64("int64").Descriptor():     ogen.Int64(),
		entfield.Uint("uint").Descriptor():       ogen.Int64().SetMinimum(&zero).SetMaximum(&maxu32),
		entfield.Uint8("uint8").Descriptor():     ogen.Int32().SetMinimum(&zero).SetMaximum(&maxu8),
		entfield.Uint16("uint16").Descriptor():   ogen.Int32().SetMinimum(&zero).SetMaximum(&maxu16),
		entfield.Uint32("uint32").Descriptor():   ogen.Int64().SetMinimum(&zero).SetMaximum(&maxu32),
		entfield.Uint64("uint64").Descriptor():   ogen.Int64().SetMinimum(&zero),
		entfield.Float32("float32").Descriptor(): ogen.Float(),
		entfield.Float("float64").Descriptor():   ogen.Double(),
		// Basic
		entfield.String("string").Descriptor():       ogen.String(),
		entfield.Bool("bool").Descriptor():           ogen.Bool(),
		entfield.UUID("uuid", uuid.Nil).Descriptor(): ogen.UUID(),
		entfield.Time("time").Descriptor():           ogen.DateTime(),
		entfield.Text("text").Descriptor():           ogen.String(),
		entfield.Enum("state").
			Values("on", "off").
			Descriptor(): ogen.String().AsEnum(nil, json.RawMessage(`"on"`), json.RawMessage(`"off"`)),
		// List
		entfield.Strings("strings").Descriptor(): ogen.String().AsArray(),
		entfield.Ints("ints").Descriptor():       ogen.Int().AsArray(),
		entfield.Floats("floats").Descriptor():   ogen.Double().AsArray(),
		entfield.Bytes("bytes").Descriptor():     ogen.Bytes(),
		// Custom
		entfield.JSON("nicknames", []string{}).Descriptor(): ogen.String().AsArray(),
		entfield.JSON("json_slice", []http.Dir{}).
			Annotations(Schema(ogen.String().AsArray())).Descriptor(): ogen.String().AsArray(),
		entfield.JSON("json_obj", url.URL{}).
			Annotations(Schema(ogen.String())).Descriptor(): ogen.String(),
		entfield.Other("other", &Link{}).
			SchemaType(map[string]string{dialect.Postgres: "varchar"}).
			Default(DefaultLink()).
			Annotations(Schema(ogen.String())).
			Descriptor(): ogen.String(),
	} {
		t.Run(d.Name, func(t *testing.T) {
			f, err := load.NewField(d)
			require.NoError(t, err)
			ens := make([]gen.Enum, len(f.Enums))
			for i, e := range f.Enums {
				ens[i] = gen.Enum{Name: e.N, Value: e.V}
			}
			gf := &gen.Field{
				Name:        f.Name,
				Type:        f.Info,
				Annotations: f.Annotations,
				Enums:       ens,
			}
			ac, err := OgenSchema(gf)
			if ex == nil {
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

func TestOgenSchema_Example(t *testing.T) {
	t.Parallel()
	entFields := map[*entfield.Descriptor]*ogen.Schema{
		entfield.String("name").
			Annotations(Example("name")).Descriptor(): func() *ogen.Schema {
			schema := ogen.String()
			v, err := json.Marshal("name")
			require.NoError(t, err)
			schema.Example = v
			return schema
		}(),
		entfield.Float32("total").
			Annotations(Example("total")).Descriptor(): func() *ogen.Schema {
			schema := ogen.Float()
			v, err := json.Marshal("total")
			require.NoError(t, err)
			schema.Example = v
			return schema
		}(),
	}

	for d, ex := range entFields {
		t.Run(d.Name, func(t *testing.T) {
			f, err := load.NewField(d)
			require.NoError(t, err)
			gf := &gen.Field{
				Name:        f.Name,
				Type:        f.Info,
				Annotations: f.Annotations,
			}
			ant, err := FieldAnnotation(gf)
			require.NoError(t, err)
			require.Equal(t, ant.Example, gf.Name)

			ac, err := OgenSchema(gf)
			require.NoError(t, err)
			require.Equal(t, ex, ac)
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
