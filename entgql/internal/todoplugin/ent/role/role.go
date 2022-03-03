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

package role

import (
	"fmt"
	"io"
	"strconv"
)

type Role string

const (
	Admin   Role = "ADMIN"
	User    Role = "USER"
	Unknown Role = "UNKNOWN"
)

func (Role) Values() (roles []string) {
	for _, r := range []Role{Admin, User, Unknown} {
		roles = append(roles, string(r))
	}
	return
}

func (e Role) String() string {
	return string(e)
}

// Validator is a validator for the "role" field enum values. It is called by the builders before save.
func Validator(e Role) error {
	for _, v := range e.Values() {
		if v == e.String() {
			return nil
		}
	}

	return fmt.Errorf("role: invalid enum value for role field: %q", e)
}

// MarshalGQL implements graphql.Marshaler interface.
func (e Role) MarshalGQL(w io.Writer) {
	_, _ = io.WriteString(w, strconv.Quote(e.String()))
}

// UnmarshalGQL implements graphql.Unmarshaler interface.
func (e *Role) UnmarshalGQL(val interface{}) error {
	str, ok := val.(string)
	if !ok {
		return fmt.Errorf("enum %T must be a string", val)
	}
	*e = Role(str)
	if err := Validator(*e); err != nil {
		return fmt.Errorf("%s is not a valid Role", str)
	}
	return nil
}
