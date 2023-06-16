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

package schema

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"math/big"
	"strings"

	"entgo.io/contrib/entproto"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/descriptorpb"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(),
		entproto.Service(),
	}
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Uint32("id").StorageKey("user_id").Annotations(entproto.Field(1)),
		field.String("user_name").
			Unique().
			Annotations(entproto.Field(2)),
		field.Time("joined").
			Immutable().
			Annotations(entproto.Field(3)),
		field.Uint("points").
			Annotations(entproto.Field(4)),
		field.Uint64("exp").
			Annotations(entproto.Field(5)),
		field.Enum("status").
			Values("pending", "active").
			Annotations(
				entproto.Field(6),
				entproto.Enum(map[string]int32{
					"pending": 1,
					"active":  2,
				}),
			),
		field.Int("external_id").
			Unique().
			Annotations(entproto.Field(8)),
		field.UUID("crm_id", uuid.New()).
			Annotations(entproto.Field(9)),
		field.Bool("banned").
			Default(false).
			Annotations(entproto.Field(10)),
		field.Uint8("custom_pb").
			Annotations(
				entproto.Field(12,
					entproto.Type(descriptorpb.FieldDescriptorProto_TYPE_UINT64),
				),
			),
		field.Int("opt_num").
			Optional().
			Annotations(entproto.Field(13)),
		field.String("opt_str").
			Optional().
			Annotations(entproto.Field(14)),
		field.Bool("opt_bool").
			Optional().
			Annotations(entproto.Field(15)),
		field.Int("big_int").
			Optional().
			GoType(BigInt{}).
			Annotations(entproto.Field(
				17,
				entproto.Type(descriptorpb.FieldDescriptorProto_TYPE_MESSAGE),
				entproto.TypeName("google.protobuf.StringValue"),
			)),
		field.Int("b_user_1").
			Optional().
			Unique().
			Annotations(entproto.Field(18)),
		field.Float32("height_in_cm").
			Default(0.0).
			Annotations(entproto.Field(19)),
		field.Float("account_balance").
			Default(0.0).
			Annotations(entproto.Field(20)),
		field.String("unnecessary").
			Optional().
			Annotations(entproto.Skip()),
		field.String("type").
			Optional().
			Annotations(
				entproto.Field(23),
			),
		field.Strings("labels").
			Optional().
			Annotations(
				entproto.Field(24),
			),
		field.JSON("int32s", []int32{}).
			Optional().
			Annotations(
				entproto.Field(25),
			),
		field.JSON("int64s", []int64{}).
			Optional().
			Annotations(
				entproto.Field(26),
			),
		field.JSON("uint32s", []uint32{}).
			Optional().
			Annotations(
				entproto.Field(27),
			),
		field.JSON("uint64s", []uint64{}).
			Optional().
			Annotations(
				entproto.Field(28),
			),
		field.Enum("device_type").
			Values("GLOWY9000", "SPEEDY300").
			Default("GLOWY9000").
			Annotations(
				entproto.Field(100),
				entproto.Enum(map[string]int32{
					"GLOWY9000": 0,
					"SPEEDY300": 1,
				}),
			),
		field.Enum("omit_prefix").
			Values("foo", "bar").
			Annotations(
				entproto.Field(103),
				entproto.Enum(
					map[string]int32{
						"foo": 1,
						"bar": 2,
					},
					entproto.OmitFieldPrefix(),
				),
			),
		field.Enum("mime_type").
			NamedValues(
				"png", "image/png",
				"svg", "image/xml+svg",
			).
			Annotations(
				entproto.Field(104),
				entproto.Enum(
					map[string]int32{
						"image/png":     1,
						"image/xml+svg": 2,
					},
				),
			),
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("group", Group.Type).
			Unique().
			Annotations(
				entproto.Field(7),
			),
		edge.To("attachment", Attachment.Type).
			Unique().
			Annotations(
				entproto.Field(11),
			),
		edge.From("received_1", Attachment.Type).
			Ref("recipients").
			Annotations(entproto.Field(16)),
		edge.To("pet", Pet.Type).
			Unique().
			Annotations(entproto.Field(21)),
		edge.To("skip_edge", SkipEdgeExample.Type).
			Unique().
			Annotations(entproto.Skip()),
	}
}

type BigInt struct {
	*big.Int
}

func NewBigInt(i int64) BigInt {
	return BigInt{Int: big.NewInt(i)}
}

func (b *BigInt) Scan(src interface{}) error {
	var i sql.NullString
	if err := i.Scan(src); err != nil {
		return err
	}
	if !i.Valid {
		return nil
	}
	if b.Int == nil {
		b.Int = big.NewInt(0)
	}
	// Value came in a floating point format.
	if strings.ContainsAny(i.String, ".+e") {
		f := big.NewFloat(0)
		if _, err := fmt.Sscan(i.String, f); err != nil {
			return err
		}
		b.Int, _ = f.Int(b.Int)
	} else if _, err := fmt.Sscan(i.String, b.Int); err != nil {
		return err
	}
	return nil
}

func (b BigInt) Value() (driver.Value, error) {
	return b.String(), nil
}
