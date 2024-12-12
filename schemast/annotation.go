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

package schemast

import (
	"fmt"
	"go/ast"
	"sort"

	"entgo.io/contrib/entproto"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"github.com/go-viper/mapstructure/v2"
	"google.golang.org/protobuf/types/descriptorpb"
)

type UnsupportedAnnotationError struct {
	annot schema.Annotation
}

func (e *UnsupportedAnnotationError) Error() string {
	return fmt.Sprintf("schemast: no Annotator configured for annotation %q", e.annot.Name())
}

// Annotator takes a schema.Annotation and returns an AST node that can be used
// to build that annotation. In addition, Annotator reports whether an AST node should be created or not.
// If the passed schema.Annotation is not supported by the Annotator, it returns UnsupportedAnnotationError.
type Annotator func(schema.Annotation) (ast.Expr, bool, error)

// Annotation is an Annotator that searches a map of well-known ent annotation (entproto, entsql, etc.) and
// invokes that Annotator if found.
func Annotation(annot schema.Annotation) (ast.Expr, bool, error) {
	annotators := map[string]Annotator{
		entproto.MessageAnnotation: protoMsg,
		entproto.ServiceAnnotation: protoSvc,
		entproto.FieldAnnotation:   protoField,
		entproto.EnumAnnotation:    protoEnum,
		"EntSQL":                   entSQL,
	}
	fn, ok := annotators[annot.Name()]
	if !ok {
		return nil, false, &UnsupportedAnnotationError{annot: annot}
	}
	return fn(annot)
}

func (c *Context) AppendTypeAnnotation(typeName string, annot schema.Annotation) error {
	newAnnot, shouldAdd, err := Annotation(annot)
	if err != nil {
		return err
	}
	if !shouldAdd {
		return nil
	}
	return c.appendReturnItem(kindAnnot, typeName, newAnnot)
}

func protoMsg(annot schema.Annotation) (ast.Expr, bool, error) {
	var m struct {
		Generate bool
		Package  string
	}
	if err := mapstructure.Decode(annot, &m); err != nil {
		return nil, false, err
	}
	if !m.Generate {
		return fnCall(selectorLit("entproto", "SkipGen")), true, nil
	}
	c := fnCall(selectorLit("entproto", "Message"))
	if m.Package != "entpb" {
		c.Args = []ast.Expr{fnCall(selectorLit("entproto", "PackageName"), strLit(m.Package))}
	}
	return c, true, nil
}

func protoSvc(annot schema.Annotation) (ast.Expr, bool, error) {
	var m struct {
		Generate bool
	}
	if err := mapstructure.Decode(annot, &m); err != nil {
		return nil, false, err
	}
	if !m.Generate {
		return nil, false, nil
	}
	return fnCall(selectorLit("entproto", "Service")), true, nil
}

func protoField(annot schema.Annotation) (ast.Expr, bool, error) {
	var m struct {
		Number   int
		Type     descriptorpb.FieldDescriptorProto_Type
		TypeName string
	}
	if err := mapstructure.Decode(annot, &m); err != nil {
		return nil, false, err
	}
	c := fnCall(selectorLit("entproto", "Field"), intLit(m.Number))
	if m.Type > 0 {
		pbt := selectorLit("descriptorpb", "FieldDescriptorProto_"+m.Type.String())
		c.Args = append(c.Args, fnCall(selectorLit("entproto", "Type"), pbt))
	}
	if m.TypeName != "" {
		c.Args = append(c.Args, fnCall(selectorLit("entproto", "TypeName"), strLit(m.TypeName)))
	}
	return c, true, nil
}

func protoEnum(annot schema.Annotation) (ast.Expr, bool, error) {
	var m struct {
		Options map[string]int32
	}
	if err := mapstructure.Decode(annot, &m); err != nil {
		return nil, false, err
	}
	opts := &ast.CompositeLit{
		Type: &ast.MapType{
			Key:   ast.NewIdent("string"),
			Value: ast.NewIdent("int32"),
		},
	}
	for k, v := range m.Options {
		opts.Elts = append(opts.Elts, &ast.KeyValueExpr{
			Key:   strLit(k),
			Value: intLit(int(v)),
		})
	}
	sort.Slice(opts.Elts, func(i, j int) bool {
		return opts.Elts[i].(*ast.KeyValueExpr).Value.(*ast.BasicLit).Value < opts.Elts[j].(*ast.KeyValueExpr).Value.(*ast.BasicLit).Value
	})
	return fnCall(selectorLit("entproto", "Enum"), opts), true, nil
}

func entSQL(annot schema.Annotation) (ast.Expr, bool, error) {
	m := &entsql.Annotation{}
	if err := mapstructure.Decode(annot, m); err != nil {
		return nil, false, err
	}
	c := &ast.CompositeLit{
		Type: selectorLit("entsql", "Annotation"),
	}
	if m.Table != "" {
		c.Elts = append(c.Elts, structAttr("Table", strLit(m.Table)))
	}
	if m.Charset != "" {
		c.Elts = append(c.Elts, structAttr("Charset", strLit(m.Charset)))
	}
	if m.Collation != "" {
		c.Elts = append(c.Elts, structAttr("Collation", strLit(m.Collation)))
	}
	if m.Default != "" {
		c.Elts = append(c.Elts, structAttr("Default", strLit(m.Default)))
	}
	if m.Size > 0 {
		c.Elts = append(c.Elts, structAttr("Size", intLit(int(m.Size))))
	}
	if m.OnDelete != "" {
		switch m.OnDelete {
		case entsql.NoAction:
			c.Elts = append(c.Elts, structAttr("OnDelete", selectorLit("entsql", "NoAction")))
		case entsql.Restrict:
			c.Elts = append(c.Elts, structAttr("OnDelete", selectorLit("entsql", "Restrict")))
		case entsql.Cascade:
			c.Elts = append(c.Elts, structAttr("OnDelete", selectorLit("entsql", "Cascade")))
		case entsql.SetNull:
			c.Elts = append(c.Elts, structAttr("OnDelete", selectorLit("entsql", "SetNull")))
		case entsql.SetDefault:
			c.Elts = append(c.Elts, structAttr("OnDelete", selectorLit("entsql", "SetDefault")))
		default:
			return nil, false, fmt.Errorf("schemast: unknown entsql ReferenceOption: %q", m.OnDelete)
		}
	}
	// TODO(rotemtam): support m.Incremental (it is a *bool)
	return c, true, nil
}

func toAnnotASTs(annots []schema.Annotation) ([]ast.Expr, error) {
	out := make([]ast.Expr, 0, len(annots))
	for _, annot := range annots {
		a, shouldAdd, err := Annotation(annot)
		if err != nil {
			return nil, err
		}
		if !shouldAdd {
			continue
		}
		out = append(out, a)
	}
	return out, nil
}
