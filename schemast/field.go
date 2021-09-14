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

package schemast

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"entgo.io/ent/schema/field"
)

// Field converts a *field.Descriptor back into an *ast.CallExpr of the ent field package that can be used
// to construct it.
func Field(desc *field.Descriptor) (*ast.CallExpr, error) {
	switch t := desc.Info.Type; {
	case t.Numeric(), t == field.TypeString, t == field.TypeBool, t == field.TypeTime:
		return fromSimpleType(desc)
	case t == field.TypeEnum:
		return fromEnumType(desc)
	default:
		return nil, fmt.Errorf("schemast: unsupported type %s", t.ConstName())
	}
}

// AppendField adds a field to the returned values of the Fields method of type typeName.
func (c *Context) AppendField(typeName string, desc *field.Descriptor) error {
	newField, err := Field(desc)
	if err != nil {
		return err
	}
	return c.appendReturnItem(kindField, typeName, newField)
}

// RemoveField removes a field from the returned values of the Fields method of type typeName.
func (c *Context) RemoveField(typeName string, fieldName string) error {
	stmt, err := c.returnStmt(typeName, "Fields")
	if err != nil {
		return err
	}
	returned, ok := stmt.Results[0].(*ast.CompositeLit)
	if !ok {
		return fmt.Errorf("schemast: unexpected AST component type %T", stmt.Results[0])
	}
	for i, item := range returned.Elts {
		call, ok := item.(*ast.CallExpr)
		if !ok {
			return fmt.Errorf("schemast: expected return statement elements to be call expressions")
		}
		name, err := extractFieldName(call)
		if err != nil {
			return err
		}
		if name == fieldName {
			returned.Elts = append(returned.Elts[:i], returned.Elts[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("schemast: could not find field %q in type %q", fieldName, typeName)
}

func newFieldCall(desc *field.Descriptor) *builderCall {
	return &builderCall{
		curr: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("field"),
				Sel: ast.NewIdent(fieldConstructor(desc)),
			},
			Args: []ast.Expr{
				strLit(desc.Name),
			},
		},
	}
}

func fromEnumType(desc *field.Descriptor) (*ast.CallExpr, error) {
	call, err := fromSimpleType(desc)
	if err != nil {
		return nil, err
	}
	modifier := "Values"
	for _, pair := range desc.Enums {
		if pair.N != pair.V {
			modifier = "NamedValues"
			break
		}
	}
	args := make([]ast.Expr, 0, len(desc.Enums))
	for _, pair := range desc.Enums {
		args = append(args, strLit(pair.N))
		if modifier == "NamedValues" {
			args = append(args, strLit(pair.V))
		}
	}
	builder := &builderCall{curr: call}
	builder.method(modifier, args...)
	return builder.curr, nil
}

func fromSimpleType(desc *field.Descriptor) (*ast.CallExpr, error) {
	builder := newFieldCall(desc)
	if desc.Nillable {
		builder.method("Nillable")
	}
	if desc.Optional {
		builder.method("Optional")
	}
	if desc.Unique {
		builder.method("Unique")
	}
	if desc.Sensitive {
		builder.method("Sensitive")
	}
	if desc.Immutable {
		builder.method("Immutable")
	}
	if desc.Comment != "" {
		builder.method("Comment", strLit(desc.Comment))
	}
	if desc.Tag != "" {
		builder.method("StructTag", strLit(desc.Tag))
	}
	if desc.StorageKey != "" {
		builder.method("StorageKey", strLit(desc.StorageKey))
	}
	if len(desc.SchemaType) > 0 {
		builder.method("SchemaType", strMapLit(desc.SchemaType))
	}
	if len(desc.Annotations) != 0 {
		annots, err := toAnnotASTs(desc.Annotations)
		if err != nil {
			return nil, err
		}
		builder.annotate(annots...)
	}
	if desc.Default != nil {
		expr, err := defaultExpr(desc.Default)
		if err != nil {
			return nil, err
		}
		builder.method("Default", expr)
	}
	// Unsupported features
	var unsupported error
	if len(desc.Validators) != 0 {
		unsupported = combineUnsupported(unsupported, "Descriptor.Validators")
	}
	if desc.UpdateDefault != nil {
		unsupported = combineUnsupported(unsupported, "Descriptor.UpdateDefault")
	}
	if unsupported != nil {
		return nil, unsupported
	}
	return builder.curr, nil
}

func fieldConstructor(dsc *field.Descriptor) string {
	return strings.TrimPrefix(dsc.Info.ConstName(), "Type")
}

func defaultExpr(d interface{}) (ast.Expr, error) {
	v := reflect.ValueOf(d)
	switch v.Kind() {
	case reflect.String:
		return strLit(d.(string)), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		lit := &ast.BasicLit{
			Kind:  token.INT,
			Value: fmt.Sprintf("%d", d),
		}
		return lit, nil
	case reflect.Float32, reflect.Float64:
		lit := &ast.BasicLit{
			Kind:  token.FLOAT,
			Value: fmt.Sprintf("%#v", d),
		}
		return lit, nil
	case reflect.Bool:
		lit := &ast.BasicLit{
			Kind:  token.STRING,
			Value: strconv.FormatBool(d.(bool)),
		}
		return lit, nil
	case reflect.Func:
		f := runtime.FuncForPC(v.Pointer()).Name()
		parts := strings.Split(f, ".")
		if len(parts) != 2 {
			return nil, errors.New("schemast: only selector exprs are supported for default func")
		}
		return selectorLit(parts[0], parts[1]), nil
	default:
		return nil, fmt.Errorf("schemast: unsupported default field kind: %q", v.Kind())
	}
}

func extractFieldName(fd *ast.CallExpr) (string, error) {
	sel, ok := fd.Fun.(*ast.SelectorExpr)
	if !ok {
		return "", fmt.Errorf("schemast: unexpected type %T", fd.Fun)
	}
	if inner, ok := sel.X.(*ast.CallExpr); ok {
		return extractFieldName(inner)
	}
	if final, ok := sel.X.(*ast.Ident); ok && final.Name != "field" {
		return "", fmt.Errorf(`schemast: expected field AST to be of form field.<Type>("name")`)
	}
	if len(fd.Args) == 0 {
		return "", fmt.Errorf("schemast: expected field constructor to have at least name arg")
	}
	name, ok := fd.Args[0].(*ast.BasicLit)
	if !ok && name.Kind == token.STRING {
		return "", fmt.Errorf("schemast: expected field name to be a string literal")
	}
	return strconv.Unquote(name.Value)
}
