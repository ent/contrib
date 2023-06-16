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

package todo

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
)

func HasPermission() func(context.Context, interface{}, graphql.Resolver, []string) (interface{}, error) {
	return func(
		ctx context.Context,
		obj interface{},
		next graphql.Resolver,
		permissions []string,
	) (res interface{}, err error) {
		// you can do your thing here for permissions
		return next(ctx)
	}
}
