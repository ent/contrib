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
//
// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"github.com/bionicstork/contrib/entgql/internal/todopulid/ent/schema/pulid"
	"github.com/bionicstork/contrib/entgql/internal/todopulid/ent/todo"
)

// prefixMap maps PULID prefixes to table names.
var prefixMap = map[pulid.ID]string{
	"TD": todo.Table,
}

// IDToType maps a pulid.ID to the underlying table.
func IDToType(ctx context.Context, id pulid.ID) (string, error) {
	if len(id) < 2 {
		return "", fmt.Errorf("IDToType: id too short")
	}
	prefix := id[:2]
	typ := prefixMap[prefix]
	if typ == "" {
		return "", fmt.Errorf("IDToType: could not map prefix '%s' to a type", prefix)
	}
	return typ, nil
}
