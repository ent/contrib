// Code generated by entc, DO NOT EDIT.

package ent

import (
	"github.com/bionicstork/contrib/schemast/internal/mutatetest/ent/schema"
	"github.com/bionicstork/contrib/schemast/internal/mutatetest/ent/withmodifiedfield"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	withmodifiedfieldFields := schema.WithModifiedField{}.Fields()
	_ = withmodifiedfieldFields
	// withmodifiedfieldDescName is the schema descriptor for name field.
	withmodifiedfieldDescName := withmodifiedfieldFields[0].Descriptor()
	// withmodifiedfield.NameValidator is a validator for the "name" field. It is called by the builders before save.
	withmodifiedfield.NameValidator = func() func(string) error {
		validators := withmodifiedfieldDescName.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(name string) error {
			for _, fn := range fns {
				if err := fn(name); err != nil {
					return err
				}
			}
			return nil
		}
	}()
}
