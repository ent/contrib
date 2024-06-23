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

import (
	"math"
	"testing"

	"entgo.io/ent/entc/gen"
)

// TestOpsHasAllGenOpsMappings
/*
Loops through all available gen.Op integer values and ensures
we have a mapping for each one in OpsALL.
*/
func TestOpsHasAllGenOpsMappings(t *testing.T) {
	genOp := 0
	breakOnName := gen.Op(math.MaxInt).Name()
	validator := OpsALL

	for ; ; genOp++ {
		name := gen.Op(genOp).Name()
		if name == breakOnName {
			break
		}

		testValue := gen.Op(genOp)
		if !validator.hasGenOp(testValue) {
			t.Errorf("OpsALL.hasGenOp(%s) = false, want true", name)
		}
	}
}

// TestOpsExcludeFromAll is just a sanity check for the &^ operator.
func TestOpsExcludeFromAll(t *testing.T) {
	ops := OpsALL &^ OpsEQ
	if ops&OpsEQ != 0 {
		t.Error("OpsALL &^ OpsEQ should remove OpsEQ")
	}
	if ops&OpsNEQ == 0 {
		t.Error("OpsALL &^ OpsEQ should keep OpsNEQ")
	}
}
