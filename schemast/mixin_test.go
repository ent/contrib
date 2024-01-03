package schemast

import (
	"bytes"
	"go/printer"
	"go/token"
	"testing"

	"github.com/stretchr/testify/require"

	"entgo.io/contrib/schemast/internal/mixintest"
)

func TestMixin(t *testing.T) {
	entMixin := mixintest.UUID{}

	result, _, err := Mixin(entMixin)
	require.NoError(t, err)

	var buf bytes.Buffer
	fst := token.NewFileSet()
	err = printer.Fprint(&buf, fst, result)
	require.NoError(t, err)
	require.EqualValues(t, "mixintest.UUID{}", buf.String())
}
