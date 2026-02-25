package entgql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithMax_SetsMaxPageSize(t *testing.T) {
	t.Parallel()

	ex, err := NewExtension(WithMax(500))
	require.NoError(t, err)
	require.Equal(t, 500, ex.maxPageSize)
}

func TestWithMax_ZeroIsNoop(t *testing.T) {
	t.Parallel()

	ex, err := NewExtension(WithMax(0))
	require.NoError(t, err)
	require.Equal(t, 0, ex.maxPageSize)
}

func TestWithMax_NegativeReturnsError(t *testing.T) {
	t.Parallel()

	_, err := NewExtension(WithMax(-1))
	require.Error(t, err)
}

func TestWithMax_DefaultIsZero(t *testing.T) {
	t.Parallel()

	ex, err := NewExtension()
	require.NoError(t, err)
	require.Equal(t, 0, ex.maxPageSize)
}
