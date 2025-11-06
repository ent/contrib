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

// Package schematype provides custom types for ent/schema.
package schematype

import (
	"database/sql/driver"
	"encoding/json"
)

// CategoryConfig implements the field.ValueScanner interface.
type CategoryConfig struct {
	MaxMembers int `json:"maxMembers,omitempty"`
}

func (t *CategoryConfig) Scan(v interface{}) (err error) {
	switch v := v.(type) {
	case string:
		err = json.Unmarshal([]byte(v), t)
	case []byte:
		err = json.Unmarshal(v, t)
	}
	return
}

func (t *CategoryConfig) Value() (driver.Value, error) {
	return json.Marshal(t)
}

// CategoryTypes is a simple JSON type.
type CategoryTypes struct {
	Public bool `json:"public,omitempty"`
}
