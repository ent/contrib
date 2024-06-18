package schemast

import (
	"entgo.io/ent"
	"fmt"
	"go/ast"
	"reflect"
	"strings"
)

// AppendMixin adds a mixin to the returned values of the Mixin method of type typeName.
func (c *Context) AppendMixin(typeName string, mixin ent.Mixin) error {
	newMixin, importPkg, err := Mixin(mixin)
	if err != nil {
		return err
	}
	c.appendImport(typeName, importPkg)
	return c.appendReturnItem(kindMixin, typeName, newMixin)
}

func Mixin(mixin ent.Mixin) (ast.Expr, string, error) {
	t := reflect.TypeOf(mixin)
	ident := t.String()
	pkgPath := t.PkgPath()
	typeInfo := strings.Split(ident, ".")
	if len(typeInfo) != 2 {
		return nil, "", fmt.Errorf("schemast: expected mixin type to be of the form {package}.{mixinName}")
	}

	return &ast.CompositeLit{
		Type: selectorLit(typeInfo[0], typeInfo[1]),
	}, pkgPath, nil
}
