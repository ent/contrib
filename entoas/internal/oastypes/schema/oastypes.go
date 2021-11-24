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

package schema

import (
	"database/sql/driver"
	"fmt"
	"net/http"
	"net/url"

	"entgo.io/contrib/entoas"
	"entgo.io/contrib/entoas/spec"
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// OASTypes holds the schema definition for the OASTypes entity.
type OASTypes struct {
	ent.Schema
}

// Fields of the OASTypes.
func (OASTypes) Fields() []ent.Field {
	return []ent.Field{
		// Numeric
		field.Int("int"),
		field.Int8("int8"),
		field.Int16("int16"),
		field.Int32("int32"),
		field.Int64("int64"),
		field.Uint("uint"),
		field.Uint8("uint8"),
		field.Uint16("uint16"),
		field.Uint32("uint32"),
		field.Uint64("uint64"),
		field.Float32("float32"),
		field.Float("float64"),
		// Basic
		field.String("string"),
		field.Bool("bool"),
		field.UUID("uuid", uuid.Nil).
			Default(uuid.New),
		field.Time("time"),
		field.Text("text"),
		field.Enum("state").Values("on", "off"),
		// List
		field.Strings("strings"),
		field.Ints("ints"),
		field.Floats("floats"),
		field.Bytes("bytes"),
		// Custom
		field.JSON("nicknames", []string{}),
		field.JSON("json_slice", []http.Dir{}).
			Annotations(entoas.OASType(&spec.Type{Type: "array", Items: &spec.Type{Type: "string"}})),
		field.JSON("json_obj", url.URL{}).
			Annotations(entoas.OASType(&spec.Type{Type: "string"})),
		field.Other("other", &Link{}).
			SchemaType(map[string]string{dialect.Postgres: "varchar"}).
			Default(DefaultLink()).
			Annotations(entoas.OASType(&spec.Type{Type: "string"})),
	}
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
