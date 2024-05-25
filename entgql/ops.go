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

package entgql

import "entgo.io/ent/entc/gen"

// Ops represents a bitwise flag of predicate operations for schema fields.
type Ops int

// List of all builtin predicates as bitwise flags.
const (
	OpsEQ           Ops = 1 << iota // =
	OpsNEQ                          // <>
	OpsGT                           // >
	OpsGTE                          // >=
	OpsLT                           // <
	OpsLTE                          // <=
	OpsIsNil                        // IS NULL / has
	OpsNotNil                       // IS NOT NULL / hasNot
	OpsIn                           // within
	OpsNotIn                        // without
	OpsEqualFold                    // equals case-insensitive
	OpsContains                     // containing
	OpsContainsFold                 // containing case-insensitive
	OpsHasPrefix                    // startingWith
	OpsHasSuffix                    // endingWith
)

var (
	// OpsALL
	/*
		OpsALL is assumed when the Ops value is 0.

		It is also provided as an explicit value here to make it easier to allow all operations except one or more.

		Example:
		OpsALL &^ OpsEQ removes the EQ operation.
	*/
	OpsALL = OpsEQ | OpsNEQ | OpsGT | OpsGTE | OpsLT | OpsLTE | OpsIsNil | OpsNotNil | OpsIn | OpsNotIn | OpsEqualFold | OpsContains | OpsContainsFold | OpsHasPrefix | OpsHasSuffix
)

func (o Ops) hasGenOp(value gen.Op) bool {
	switch value {
	case gen.EQ:
		return o&OpsEQ != 0
	case gen.NEQ:
		return o&OpsNEQ != 0
	case gen.GT:
		return o&OpsGT != 0
	case gen.GTE:
		return o&OpsGTE != 0
	case gen.LT:
		return o&OpsLT != 0
	case gen.LTE:
		return o&OpsLTE != 0
	case gen.IsNil:
		return o&OpsIsNil != 0
	case gen.NotNil:
		return o&OpsNotNil != 0
	case gen.In:
		return o&OpsIn != 0
	case gen.NotIn:
		return o&OpsNotIn != 0
	case gen.EqualFold:
		return o&OpsEqualFold != 0
	case gen.Contains:
		return o&OpsContains != 0
	case gen.ContainsFold:
		return o&OpsContainsFold != 0
	case gen.HasPrefix:
		return o&OpsHasPrefix != 0
	case gen.HasSuffix:
		return o&OpsHasSuffix != 0
	default:
		// We don't return false here because we want our
		// tests to fail if gen.Op and Ops are out of sync.
		panic("unknown operation")
	}
}
