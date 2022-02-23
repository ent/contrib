// Code generated by entc, DO NOT EDIT.

package nilexample

import (
	"time"

	"github.com/bionicstork/bionicstork/pkg/entproto/internal/todo/ent/predicate"
	"entgo.io/ent/dialect/sql"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.NilExample {
	return predicate.NilExample(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.NilExample {
	return predicate.NilExample(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.NilExample {
	return predicate.NilExample(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.NilExample {
	return predicate.NilExample(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(ids) == 0 {
			s.Where(sql.False())
			return
		}
		v := make([]interface{}, len(ids))
		for i := range v {
			v[i] = ids[i]
		}
		s.Where(sql.In(s.C(FieldID), v...))
	})
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.NilExample {
	return predicate.NilExample(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(ids) == 0 {
			s.Where(sql.False())
			return
		}
		v := make([]interface{}, len(ids))
		for i := range v {
			v[i] = ids[i]
		}
		s.Where(sql.NotIn(s.C(FieldID), v...))
	})
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.NilExample {
	return predicate.NilExample(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.NilExample {
	return predicate.NilExample(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.NilExample {
	return predicate.NilExample(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.NilExample {
	return predicate.NilExample(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// StrNil applies equality check predicate on the "str_nil" field. It's identical to StrNilEQ.
func StrNil(v string) predicate.NilExample {
	return predicate.NilExample(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldStrNil), v))
	})
}

// TimeNil applies equality check predicate on the "time_nil" field. It's identical to TimeNilEQ.
func TimeNil(v time.Time) predicate.NilExample {
	return predicate.NilExample(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldTimeNil), v))
	})
}

// StrNilEQ applies the EQ predicate on the "str_nil" field.
func StrNilEQ(v string) predicate.NilExample {
	return predicate.NilExample(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldStrNil), v))
	})
}

// StrNilNEQ applies the NEQ predicate on the "str_nil" field.
func StrNilNEQ(v string) predicate.NilExample {
	return predicate.NilExample(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldStrNil), v))
	})
}

// StrNilIn applies the In predicate on the "str_nil" field.
func StrNilIn(vs ...string) predicate.NilExample {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.NilExample(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldStrNil), v...))
	})
}

// StrNilNotIn applies the NotIn predicate on the "str_nil" field.
func StrNilNotIn(vs ...string) predicate.NilExample {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.NilExample(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldStrNil), v...))
	})
}

// StrNilGT applies the GT predicate on the "str_nil" field.
func StrNilGT(v string) predicate.NilExample {
	return predicate.NilExample(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldStrNil), v))
	})
}

// StrNilGTE applies the GTE predicate on the "str_nil" field.
func StrNilGTE(v string) predicate.NilExample {
	return predicate.NilExample(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldStrNil), v))
	})
}

// StrNilLT applies the LT predicate on the "str_nil" field.
func StrNilLT(v string) predicate.NilExample {
	return predicate.NilExample(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldStrNil), v))
	})
}

// StrNilLTE applies the LTE predicate on the "str_nil" field.
func StrNilLTE(v string) predicate.NilExample {
	return predicate.NilExample(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldStrNil), v))
	})
}

// StrNilContains applies the Contains predicate on the "str_nil" field.
func StrNilContains(v string) predicate.NilExample {
	return predicate.NilExample(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldStrNil), v))
	})
}

// StrNilHasPrefix applies the HasPrefix predicate on the "str_nil" field.
func StrNilHasPrefix(v string) predicate.NilExample {
	return predicate.NilExample(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldStrNil), v))
	})
}

// StrNilHasSuffix applies the HasSuffix predicate on the "str_nil" field.
func StrNilHasSuffix(v string) predicate.NilExample {
	return predicate.NilExample(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldStrNil), v))
	})
}

// StrNilIsNil applies the IsNil predicate on the "str_nil" field.
func StrNilIsNil() predicate.NilExample {
	return predicate.NilExample(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldStrNil)))
	})
}

// StrNilNotNil applies the NotNil predicate on the "str_nil" field.
func StrNilNotNil() predicate.NilExample {
	return predicate.NilExample(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldStrNil)))
	})
}

// StrNilEqualFold applies the EqualFold predicate on the "str_nil" field.
func StrNilEqualFold(v string) predicate.NilExample {
	return predicate.NilExample(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldStrNil), v))
	})
}

// StrNilContainsFold applies the ContainsFold predicate on the "str_nil" field.
func StrNilContainsFold(v string) predicate.NilExample {
	return predicate.NilExample(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldStrNil), v))
	})
}

// TimeNilEQ applies the EQ predicate on the "time_nil" field.
func TimeNilEQ(v time.Time) predicate.NilExample {
	return predicate.NilExample(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldTimeNil), v))
	})
}

// TimeNilNEQ applies the NEQ predicate on the "time_nil" field.
func TimeNilNEQ(v time.Time) predicate.NilExample {
	return predicate.NilExample(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldTimeNil), v))
	})
}

// TimeNilIn applies the In predicate on the "time_nil" field.
func TimeNilIn(vs ...time.Time) predicate.NilExample {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.NilExample(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldTimeNil), v...))
	})
}

// TimeNilNotIn applies the NotIn predicate on the "time_nil" field.
func TimeNilNotIn(vs ...time.Time) predicate.NilExample {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.NilExample(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldTimeNil), v...))
	})
}

// TimeNilGT applies the GT predicate on the "time_nil" field.
func TimeNilGT(v time.Time) predicate.NilExample {
	return predicate.NilExample(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldTimeNil), v))
	})
}

// TimeNilGTE applies the GTE predicate on the "time_nil" field.
func TimeNilGTE(v time.Time) predicate.NilExample {
	return predicate.NilExample(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldTimeNil), v))
	})
}

// TimeNilLT applies the LT predicate on the "time_nil" field.
func TimeNilLT(v time.Time) predicate.NilExample {
	return predicate.NilExample(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldTimeNil), v))
	})
}

// TimeNilLTE applies the LTE predicate on the "time_nil" field.
func TimeNilLTE(v time.Time) predicate.NilExample {
	return predicate.NilExample(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldTimeNil), v))
	})
}

// TimeNilIsNil applies the IsNil predicate on the "time_nil" field.
func TimeNilIsNil() predicate.NilExample {
	return predicate.NilExample(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldTimeNil)))
	})
}

// TimeNilNotNil applies the NotNil predicate on the "time_nil" field.
func TimeNilNotNil() predicate.NilExample {
	return predicate.NilExample(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldTimeNil)))
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.NilExample) predicate.NilExample {
	return predicate.NilExample(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.NilExample) predicate.NilExample {
	return predicate.NilExample(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for i, p := range predicates {
			if i > 0 {
				s1.Or()
			}
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Not applies the not operator on the given predicate.
func Not(p predicate.NilExample) predicate.NilExample {
	return predicate.NilExample(func(s *sql.Selector) {
		p(s.Not())
	})
}
