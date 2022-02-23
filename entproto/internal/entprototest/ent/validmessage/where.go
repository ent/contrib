// Code generated by entc, DO NOT EDIT.

package validmessage

import (
	"time"

	"github.com/bionicstork/bionicstork/pkg/entproto/internal/entprototest/ent/predicate"
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
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
func IDNotIn(ids ...int) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
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
func IDGT(id int) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldName), v))
	})
}

// Ts applies equality check predicate on the "ts" field. It's identical to TsEQ.
func Ts(v time.Time) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldTs), v))
	})
}

// UUID applies equality check predicate on the "uuid" field. It's identical to UUIDEQ.
func UUID(v uuid.UUID) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUUID), v))
	})
}

// U8 applies equality check predicate on the "u8" field. It's identical to U8EQ.
func U8(v uint8) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldU8), v))
	})
}

// Opti8 applies equality check predicate on the "opti8" field. It's identical to Opti8EQ.
func Opti8(v int8) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldOpti8), v))
	})
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldName), v))
	})
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldName), v))
	})
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.ValidMessage {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.ValidMessage(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldName), v...))
	})
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.ValidMessage {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.ValidMessage(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldName), v...))
	})
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldName), v))
	})
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldName), v))
	})
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldName), v))
	})
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldName), v))
	})
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldName), v))
	})
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldName), v))
	})
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldName), v))
	})
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldName), v))
	})
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldName), v))
	})
}

// TsEQ applies the EQ predicate on the "ts" field.
func TsEQ(v time.Time) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldTs), v))
	})
}

// TsNEQ applies the NEQ predicate on the "ts" field.
func TsNEQ(v time.Time) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldTs), v))
	})
}

// TsIn applies the In predicate on the "ts" field.
func TsIn(vs ...time.Time) predicate.ValidMessage {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.ValidMessage(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldTs), v...))
	})
}

// TsNotIn applies the NotIn predicate on the "ts" field.
func TsNotIn(vs ...time.Time) predicate.ValidMessage {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.ValidMessage(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldTs), v...))
	})
}

// TsGT applies the GT predicate on the "ts" field.
func TsGT(v time.Time) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldTs), v))
	})
}

// TsGTE applies the GTE predicate on the "ts" field.
func TsGTE(v time.Time) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldTs), v))
	})
}

// TsLT applies the LT predicate on the "ts" field.
func TsLT(v time.Time) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldTs), v))
	})
}

// TsLTE applies the LTE predicate on the "ts" field.
func TsLTE(v time.Time) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldTs), v))
	})
}

// UUIDEQ applies the EQ predicate on the "uuid" field.
func UUIDEQ(v uuid.UUID) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUUID), v))
	})
}

// UUIDNEQ applies the NEQ predicate on the "uuid" field.
func UUIDNEQ(v uuid.UUID) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldUUID), v))
	})
}

// UUIDIn applies the In predicate on the "uuid" field.
func UUIDIn(vs ...uuid.UUID) predicate.ValidMessage {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.ValidMessage(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldUUID), v...))
	})
}

// UUIDNotIn applies the NotIn predicate on the "uuid" field.
func UUIDNotIn(vs ...uuid.UUID) predicate.ValidMessage {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.ValidMessage(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldUUID), v...))
	})
}

// UUIDGT applies the GT predicate on the "uuid" field.
func UUIDGT(v uuid.UUID) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldUUID), v))
	})
}

// UUIDGTE applies the GTE predicate on the "uuid" field.
func UUIDGTE(v uuid.UUID) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldUUID), v))
	})
}

// UUIDLT applies the LT predicate on the "uuid" field.
func UUIDLT(v uuid.UUID) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldUUID), v))
	})
}

// UUIDLTE applies the LTE predicate on the "uuid" field.
func UUIDLTE(v uuid.UUID) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldUUID), v))
	})
}

// U8EQ applies the EQ predicate on the "u8" field.
func U8EQ(v uint8) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldU8), v))
	})
}

// U8NEQ applies the NEQ predicate on the "u8" field.
func U8NEQ(v uint8) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldU8), v))
	})
}

// U8In applies the In predicate on the "u8" field.
func U8In(vs ...uint8) predicate.ValidMessage {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.ValidMessage(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldU8), v...))
	})
}

// U8NotIn applies the NotIn predicate on the "u8" field.
func U8NotIn(vs ...uint8) predicate.ValidMessage {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.ValidMessage(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldU8), v...))
	})
}

// U8GT applies the GT predicate on the "u8" field.
func U8GT(v uint8) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldU8), v))
	})
}

// U8GTE applies the GTE predicate on the "u8" field.
func U8GTE(v uint8) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldU8), v))
	})
}

// U8LT applies the LT predicate on the "u8" field.
func U8LT(v uint8) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldU8), v))
	})
}

// U8LTE applies the LTE predicate on the "u8" field.
func U8LTE(v uint8) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldU8), v))
	})
}

// Opti8EQ applies the EQ predicate on the "opti8" field.
func Opti8EQ(v int8) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldOpti8), v))
	})
}

// Opti8NEQ applies the NEQ predicate on the "opti8" field.
func Opti8NEQ(v int8) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldOpti8), v))
	})
}

// Opti8In applies the In predicate on the "opti8" field.
func Opti8In(vs ...int8) predicate.ValidMessage {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.ValidMessage(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldOpti8), v...))
	})
}

// Opti8NotIn applies the NotIn predicate on the "opti8" field.
func Opti8NotIn(vs ...int8) predicate.ValidMessage {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.ValidMessage(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldOpti8), v...))
	})
}

// Opti8GT applies the GT predicate on the "opti8" field.
func Opti8GT(v int8) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldOpti8), v))
	})
}

// Opti8GTE applies the GTE predicate on the "opti8" field.
func Opti8GTE(v int8) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldOpti8), v))
	})
}

// Opti8LT applies the LT predicate on the "opti8" field.
func Opti8LT(v int8) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldOpti8), v))
	})
}

// Opti8LTE applies the LTE predicate on the "opti8" field.
func Opti8LTE(v int8) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldOpti8), v))
	})
}

// Opti8IsNil applies the IsNil predicate on the "opti8" field.
func Opti8IsNil() predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldOpti8)))
	})
}

// Opti8NotNil applies the NotNil predicate on the "opti8" field.
func Opti8NotNil() predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldOpti8)))
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.ValidMessage) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.ValidMessage) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
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
func Not(p predicate.ValidMessage) predicate.ValidMessage {
	return predicate.ValidMessage(func(s *sql.Selector) {
		p(s.Not())
	})
}
