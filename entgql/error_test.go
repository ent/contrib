// Copyright 2019-present Facebook Inc. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package entgql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/facebookincubator/ent-contrib/entgql"
	"github.com/stretchr/testify/assert"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func TestDefaultErrorPresenter(t *testing.T) {
	err := fmt.Errorf("wrapping gqlerr: %w", gqlerror.Errorf("gqlerr"))
	gqlerr := entgql.DefaultErrorPresenter(context.Background(), err)
	assert.Equal(t, "gqlerr", gqlerr.Message)
}

func TestErrNodeNotFound(t *testing.T) {
	err := entgql.ErrNodeNotFound(42)
	assert.EqualError(t, err, "input: Could not resolve to a node with the global id of '42'")
	assert.Equal(t, "NOT_FOUND", err.Extensions["code"])
}