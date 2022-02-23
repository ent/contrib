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

package main

import (
	"flag"
	"fmt"

	entopts "github.com/bionicstork/contrib/entproto/cmd/protoc-gen-ent/options/ent"
	"entgo.io/contrib/schemast"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

var schemaDir *string

func main() {
	var flags flag.FlagSet
	schemaDir = flags.String("schemadir", "./ent/schema", "path to ent schema dir")
	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {
		return printSchemas(*schemaDir, gen)
	})
}

func printSchemas(schemaDir string, gen *protogen.Plugin) error {
	ctx, err := schemast.Load(schemaDir)
	if err != nil {
		return err
	}
	var mutations []schemast.Mutator
	for _, f := range gen.Files {
		if !f.Generate {
			continue
		}
		// TODO(rotemtam): handle nested messages recursively?
		for _, msg := range f.Messages {
			opts, ok := schemaOpts(msg)
			if !ok || !opts.GetGen() {
				continue
			}
			schema, err := toSchema(msg, opts)
			if err != nil {
				return err
			}
			mutations = append(mutations, schema)
		}
	}
	if err := schemast.Mutate(ctx, mutations...); err != nil {
		return err
	}
	if err := ctx.Print(schemaDir, schemast.Header("File updated by protoc-gen-ent.")); err != nil {
		return err
	}
	return nil
}

func schemaOpts(msg *protogen.Message) (*entopts.Schema, bool) {
	opts, ok := msg.Desc.Options().(*descriptorpb.MessageOptions)
	if !ok {
		return nil, false
	}
	extension := proto.GetExtension(opts, entopts.E_Schema)
	mop, ok := extension.(*entopts.Schema)
	return mop, ok
}

func fieldOpts(fld *protogen.Field) (*entopts.Field, bool) {
	opts, ok := fld.Desc.Options().(*descriptorpb.FieldOptions)
	if !ok {
		return nil, false
	}
	extension := proto.GetExtension(opts, entopts.E_Field)
	fop, ok := extension.(*entopts.Field)
	return fop, ok
}

func edgeOpts(fld *protogen.Field) (*entopts.Edge, bool) {
	opts, ok := fld.Desc.Options().(*descriptorpb.FieldOptions)
	if !ok || opts == nil {
		return nil, false
	}
	extension := proto.GetExtension(opts, entopts.E_Edge)
	eop, ok := extension.(*entopts.Edge)
	return eop, ok
}

func toSchema(m *protogen.Message, opts *entopts.Schema) (*schemast.UpsertSchema, error) {
	name := string(m.Desc.Name())
	if opts.Name != nil {
		name = opts.GetName()
	}
	out := &schemast.UpsertSchema{
		Name: name,
	}
	for _, f := range m.Fields {
		if isEdge(f) {
			edg, err := toEdge(f)
			if err != nil {
				return nil, err
			}
			out.Edges = append(out.Edges, edg)
			continue
		}
		fld, err := toField(f)
		if err != nil {
			return nil, err
		}
		out.Fields = append(out.Fields, fld)
	}
	return out, nil
}

func isEdge(f *protogen.Field) bool {
	return f.Desc.Kind() == protoreflect.MessageKind
}

func toEdge(f *protogen.Field) (ent.Edge, error) {
	name := string(f.Desc.Name())
	msgType := string(f.Desc.Message().Name())
	opts, ok := edgeOpts(f)
	if !ok {
		return nil, fmt.Errorf("protoc-gen-ent: expected ent.edge option on field %q", name)
	}
	var e ent.Edge
	switch {
	// TODO(rotemtam): handle O2O/M2M same type
	case opts.Ref != nil:
		e = edge.From(name, placeholder.Type)
	default:
		e = edge.To(name, placeholder.Type)
	}
	e = withType(e, msgType)
	applyEdgeOpts(e, opts)
	return e, nil
}

func toField(f *protogen.Field) (ent.Field, error) {
	name := string(f.Desc.Name())
	var fld ent.Field
	switch f.Desc.Kind() {
	case protoreflect.StringKind:
		fld = field.String(name)
	case protoreflect.BoolKind:
		fld = field.Bool(name)
	case protoreflect.Sint32Kind:
		fld = field.Int32(name)
	case protoreflect.Uint32Kind:
		fld = field.Uint32(name)
	case protoreflect.Int64Kind:
		fld = field.Int64(name)
	case protoreflect.Sint64Kind:
		fld = field.Int64(name)
	case protoreflect.Uint64Kind:
		fld = field.Uint64(name)
	case protoreflect.Sfixed32Kind:
		fld = field.Int32(name)
	case protoreflect.Fixed32Kind:
		fld = field.Int32(name)
	case protoreflect.FloatKind:
		fld = field.Float(name)
	case protoreflect.Sfixed64Kind:
		fld = field.Int64(name)
	case protoreflect.Fixed64Kind:
		fld = field.Int64(name)
	case protoreflect.DoubleKind:
		fld = field.Float(name)
	case protoreflect.BytesKind:
		fld = field.Bytes(name)
	case protoreflect.Int32Kind:
		fld = field.Int32(name)
	case protoreflect.EnumKind:
		pbEnum := f.Desc.Enum().Values()
		values := make([]string, 0, pbEnum.Len())
		for i := 0; i < pbEnum.Len(); i++ {
			values = append(values, string(pbEnum.Get(i).Name()))
		}
		fld = field.Enum(name).Values(values...)
	default:
		return nil, fmt.Errorf("protoc-gen-ent: unsupported kind %q", f.Desc.Kind())
	}
	if opts, ok := fieldOpts(f); ok {
		applyFieldOpts(fld, opts)
	}
	return fld, nil
}

func applyFieldOpts(fld ent.Field, opts *entopts.Field) {
	d := fld.Descriptor()
	d.Nillable = opts.GetNillable()
	d.Optional = opts.GetOptional()
	d.Unique = opts.GetUnique()
	d.Sensitive = opts.GetSensitive()
	d.Immutable = opts.GetImmutable()
	d.Comment = opts.GetComment()
	d.Tag = opts.GetStructTag()
	d.StorageKey = opts.GetStorageKey()
	d.SchemaType = opts.GetSchemaType()
}

func applyEdgeOpts(edg ent.Edge, opts *entopts.Edge) {
	d := edg.Descriptor()
	d.Unique = opts.GetUnique()
	d.RefName = opts.GetRef()
	d.Required = opts.GetRequired()
	d.Field = opts.GetField()
	d.Tag = opts.GetStructTag()
	if sk := opts.StorageKey; sk != nil {
		d.StorageKey = &edge.StorageKey{
			Table:   sk.GetTable(),
			Columns: sk.GetColumns(),
		}
	}
}

type placeholder struct {
}

func (placeholder) Type() {

}

func withType(edg ent.Edge, tn string) ent.Edge {
	edg.Descriptor().Type = tn
	return edg
}
