package entgql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithParallelWhereInputFiles_DefaultDisabled(t *testing.T) {
	t.Parallel()

	ex, err := NewExtension(WithSplitGoFiles(true), WithWhereInputs(true))
	require.NoError(t, err)
	require.False(t, ex.parallelWhereInputFiles)
}

func TestWithParallelWhereInputFiles_Enabled(t *testing.T) {
	t.Parallel()

	ex, err := NewExtension(
		WithSplitGoFiles(true),
		WithWhereInputs(true),
		WithParallelWhereInputFiles(true),
	)
	require.NoError(t, err)
	require.True(t, ex.parallelWhereInputFiles)
}
