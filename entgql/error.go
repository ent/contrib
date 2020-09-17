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

package entgql

import (
	"context"
	"errors"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/errcode"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// DefaultErrorPresenter adds error unwrapping to graphql.DefaultErrorPresenter.
func DefaultErrorPresenter(ctx context.Context, err error) (gqlerr *gqlerror.Error) {
	if errors.As(err, &gqlerr) {
		if gqlerr.Path == nil {
			gqlerr.Path = graphql.GetFieldContext(ctx).Path()
		}
		return gqlerr
	}
	gqlerr = &gqlerror.Error{
		Message: err.Error(),
		Path:    graphql.GetFieldContext(ctx).Path(),
	}
	var ee graphql.ExtendedError
	if errors.As(err, &ee) {
		gqlerr.Extensions = ee.Extensions()
	}
	return gqlerr
}

// ErrNodeNotFound creates a node not found graphql error.
func ErrNodeNotFound(id interface{}) *gqlerror.Error {
	err := gqlerror.Errorf("Could not resolve to a node with the global id of '%v'", id)
	errcode.Set(err, "NOT_FOUND")
	return err
}
