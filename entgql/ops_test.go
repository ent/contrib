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
