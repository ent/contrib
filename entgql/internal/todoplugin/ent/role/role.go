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

import "io"

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

func (r Role) MarshalGQL(w io.Writer) {
	panic("implement me")
}

func (r Role) UnmarshalGQL(v interface{}) error {
	panic("implement me")
}
