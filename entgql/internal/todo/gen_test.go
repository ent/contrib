package todo_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/stretchr/testify/require"
)

func TestE(t *testing.T) {
	tempDir := t.TempDir()
	ex, err := entgql.NewExtension(
		entgql.WithConfigPath("gqlgen.yml"),
		entgql.WithSchemaGenerator(),
		entgql.WithSchemaPath(filepath.Join(tempDir, "ent.graphql")),
		entgql.WithWhereInputs(true),
		entgql.WithNodeDescriptor(true),
	)
	require.NoError(t, err)
	err = entc.Generate("./ent/schema", &gen.Config{
		Target: tempDir,
		Features: []gen.Feature{
			gen.FeatureModifier,
		},
	}, entc.Extensions(ex))
	require.NoError(t, err)
	fmt.Println("!")
	// compare file
}
