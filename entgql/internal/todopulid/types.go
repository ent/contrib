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

package todo

import (
	"context"
	"fmt"

	"github.com/facebookincubator/ent-contrib/entgql/internal/todopulid/ent/todo"
)

// TypeMap maps PULID prefixes to database table names.
var TypeMap = map[string]string{
	"TO": todo.Table,
}

// IDToType returns the type name associated with a PULID id.
func IDToType(ctx context.Context, id string) (string, error) {
	fmt.Println("IDTOTYPE:", id)
	if len(id) < 2 {
		return "", fmt.Errorf("idtotype: id too short")
	}
	prefix := id[:2]
	typ := TypeMap[prefix]
	fmt.Println("  IDTOTYPE:", prefix, typ)
	if typ == "" {
		return "", fmt.Errorf("idtotype: could not map prefix '%s' to a type", prefix)
	}
	return typ, nil
}
