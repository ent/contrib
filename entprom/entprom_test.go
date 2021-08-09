package entprom

import (
	"context"
	"testing"

	"entgo.io/ent/examples/fs/ent"
	_ "github.com/mattn/go-sqlite3"
	"github.com/prometheus/tsdb/testutil"
	"github.com/stretchr/testify/require"
)

func TestHookDefault(t *testing.T) {
	client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	require.NoError(t, err)
	defer client.Close()
	ctx := context.Background()
	require.NoError(t, client.Schema.Create(ctx))

	ph := &hook{}
	client.Use(newHook(ph))
	client.File.Create().SetName("a8m").SaveX(ctx)

	require.Equal(t, 1, testutil.CollectAndCount(ph.opsProcessedTotal))
}

func TestExtraLabelsOption(t *testing.T) {
	ph := &hook{}
	expectedLabels := map[string]string{"dog": "fortuna"}
	Labels(expectedLabels)(ph)
	require.EqualValues(t, expectedLabels, ph.extraLabels)
}
