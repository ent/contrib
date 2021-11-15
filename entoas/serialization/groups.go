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

package serialization

import (
	"hash/fnv"
	"sort"
	"strconv"
)

// Groups are used to determine what properties to load and serialize.
type Groups []string

// Add adds a group to the groups. If the group is already present it does nothing.
func (gs *Groups) Add(g ...string) {
	for _, g1 := range g {
		if !gs.HasGroup(g1) {
			*gs = append(*gs, g1)
		}
	}
}

// HasGroup checks if the given group is present.
func (gs Groups) HasGroup(g string) bool {
	for _, e := range gs {
		if e == g {
			return true
		}
	}
	return false
}

// Match check if at least one of the given Groups is present.
func (gs Groups) Match(other Groups) bool {
	for _, g := range other {
		if gs.HasGroup(g) {
			return true
		}
	}
	return false
}

// Equal reports if two Groups have the same entries.
func (gs Groups) Equal(other Groups) bool {
	if len(gs) != len(other) {
		return false
	}
	for _, g := range other {
		if !gs.HasGroup(g) {
			return false
		}
	}
	return true
}

// Hash returns a hash value for a Groups.
func (gs Groups) Hash() uint32 {
	sort.Strings(gs)
	h := fnv.New32a()
	for _, g := range gs {
		_, _ = h.Write([]byte(g))
	}
	_, _ = h.Write([]byte(strconv.Itoa(len(gs))))
	return h.Sum32()
}
