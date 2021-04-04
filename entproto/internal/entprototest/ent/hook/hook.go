// Code generated by entc, DO NOT EDIT.

package hook

import (
	"context"
	"fmt"

	"entgo.io/contrib/entproto/internal/entprototest/ent"
)

// The BlogPostFunc type is an adapter to allow the use of ordinary
// function as BlogPost mutator.
type BlogPostFunc func(context.Context, *ent.BlogPostMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f BlogPostFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.BlogPostMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.BlogPostMutation", m)
	}
	return f(ctx, mv)
}

// The CategoryFunc type is an adapter to allow the use of ordinary
// function as Category mutator.
type CategoryFunc func(context.Context, *ent.CategoryMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f CategoryFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.CategoryMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.CategoryMutation", m)
	}
	return f(ctx, mv)
}

// The DependsOnSkippedFunc type is an adapter to allow the use of ordinary
// function as DependsOnSkipped mutator.
type DependsOnSkippedFunc func(context.Context, *ent.DependsOnSkippedMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f DependsOnSkippedFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.DependsOnSkippedMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.DependsOnSkippedMutation", m)
	}
	return f(ctx, mv)
}

// The DuplicateNumberMessageFunc type is an adapter to allow the use of ordinary
// function as DuplicateNumberMessage mutator.
type DuplicateNumberMessageFunc func(context.Context, *ent.DuplicateNumberMessageMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f DuplicateNumberMessageFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.DuplicateNumberMessageMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.DuplicateNumberMessageMutation", m)
	}
	return f(ctx, mv)
}

// The ExplicitSkippedMessageFunc type is an adapter to allow the use of ordinary
// function as ExplicitSkippedMessage mutator.
type ExplicitSkippedMessageFunc func(context.Context, *ent.ExplicitSkippedMessageMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f ExplicitSkippedMessageFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.ExplicitSkippedMessageMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.ExplicitSkippedMessageMutation", m)
	}
	return f(ctx, mv)
}

// The ImageFunc type is an adapter to allow the use of ordinary
// function as Image mutator.
type ImageFunc func(context.Context, *ent.ImageMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f ImageFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.ImageMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.ImageMutation", m)
	}
	return f(ctx, mv)
}

// The ImplicitSkippedMessageFunc type is an adapter to allow the use of ordinary
// function as ImplicitSkippedMessage mutator.
type ImplicitSkippedMessageFunc func(context.Context, *ent.ImplicitSkippedMessageMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f ImplicitSkippedMessageFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.ImplicitSkippedMessageMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.ImplicitSkippedMessageMutation", m)
	}
	return f(ctx, mv)
}

// The InvalidFieldMessageFunc type is an adapter to allow the use of ordinary
// function as InvalidFieldMessage mutator.
type InvalidFieldMessageFunc func(context.Context, *ent.InvalidFieldMessageMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f InvalidFieldMessageFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.InvalidFieldMessageMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.InvalidFieldMessageMutation", m)
	}
	return f(ctx, mv)
}

// The MessageWithEnumFunc type is an adapter to allow the use of ordinary
// function as MessageWithEnum mutator.
type MessageWithEnumFunc func(context.Context, *ent.MessageWithEnumMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f MessageWithEnumFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.MessageWithEnumMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.MessageWithEnumMutation", m)
	}
	return f(ctx, mv)
}

// The MessageWithFieldOneFunc type is an adapter to allow the use of ordinary
// function as MessageWithFieldOne mutator.
type MessageWithFieldOneFunc func(context.Context, *ent.MessageWithFieldOneMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f MessageWithFieldOneFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.MessageWithFieldOneMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.MessageWithFieldOneMutation", m)
	}
	return f(ctx, mv)
}

// The MessageWithIDFunc type is an adapter to allow the use of ordinary
// function as MessageWithID mutator.
type MessageWithIDFunc func(context.Context, *ent.MessageWithIDMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f MessageWithIDFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.MessageWithIDMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.MessageWithIDMutation", m)
	}
	return f(ctx, mv)
}

// The MessageWithPackageNameFunc type is an adapter to allow the use of ordinary
// function as MessageWithPackageName mutator.
type MessageWithPackageNameFunc func(context.Context, *ent.MessageWithPackageNameMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f MessageWithPackageNameFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.MessageWithPackageNameMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.MessageWithPackageNameMutation", m)
	}
	return f(ctx, mv)
}

// The PortalFunc type is an adapter to allow the use of ordinary
// function as Portal mutator.
type PortalFunc func(context.Context, *ent.PortalMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f PortalFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.PortalMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.PortalMutation", m)
	}
	return f(ctx, mv)
}

// The UserFunc type is an adapter to allow the use of ordinary
// function as User mutator.
type UserFunc func(context.Context, *ent.UserMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f UserFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.UserMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.UserMutation", m)
	}
	return f(ctx, mv)
}

// The ValidMessageFunc type is an adapter to allow the use of ordinary
// function as ValidMessage mutator.
type ValidMessageFunc func(context.Context, *ent.ValidMessageMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f ValidMessageFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.ValidMessageMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.ValidMessageMutation", m)
	}
	return f(ctx, mv)
}

// Condition is a hook condition function.
type Condition func(context.Context, ent.Mutation) bool

// And groups conditions with the AND operator.
func And(first, second Condition, rest ...Condition) Condition {
	return func(ctx context.Context, m ent.Mutation) bool {
		if !first(ctx, m) || !second(ctx, m) {
			return false
		}
		for _, cond := range rest {
			if !cond(ctx, m) {
				return false
			}
		}
		return true
	}
}

// Or groups conditions with the OR operator.
func Or(first, second Condition, rest ...Condition) Condition {
	return func(ctx context.Context, m ent.Mutation) bool {
		if first(ctx, m) || second(ctx, m) {
			return true
		}
		for _, cond := range rest {
			if cond(ctx, m) {
				return true
			}
		}
		return false
	}
}

// Not negates a given condition.
func Not(cond Condition) Condition {
	return func(ctx context.Context, m ent.Mutation) bool {
		return !cond(ctx, m)
	}
}

// HasOp is a condition testing mutation operation.
func HasOp(op ent.Op) Condition {
	return func(_ context.Context, m ent.Mutation) bool {
		return m.Op().Is(op)
	}
}

// HasAddedFields is a condition validating `.AddedField` on fields.
func HasAddedFields(field string, fields ...string) Condition {
	return func(_ context.Context, m ent.Mutation) bool {
		if _, exists := m.AddedField(field); !exists {
			return false
		}
		for _, field := range fields {
			if _, exists := m.AddedField(field); !exists {
				return false
			}
		}
		return true
	}
}

// HasClearedFields is a condition validating `.FieldCleared` on fields.
func HasClearedFields(field string, fields ...string) Condition {
	return func(_ context.Context, m ent.Mutation) bool {
		if exists := m.FieldCleared(field); !exists {
			return false
		}
		for _, field := range fields {
			if exists := m.FieldCleared(field); !exists {
				return false
			}
		}
		return true
	}
}

// HasFields is a condition validating `.Field` on fields.
func HasFields(field string, fields ...string) Condition {
	return func(_ context.Context, m ent.Mutation) bool {
		if _, exists := m.Field(field); !exists {
			return false
		}
		for _, field := range fields {
			if _, exists := m.Field(field); !exists {
				return false
			}
		}
		return true
	}
}

// If executes the given hook under condition.
//
//	hook.If(ComputeAverage, And(HasFields(...), HasAddedFields(...)))
//
func If(hk ent.Hook, cond Condition) ent.Hook {
	return func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			if cond(ctx, m) {
				return hk(next).Mutate(ctx, m)
			}
			return next.Mutate(ctx, m)
		})
	}
}

// On executes the given hook only for the given operation.
//
//	hook.On(Log, ent.Delete|ent.Create)
//
func On(hk ent.Hook, op ent.Op) ent.Hook {
	return If(hk, HasOp(op))
}

// Unless skips the given hook only for the given operation.
//
//	hook.Unless(Log, ent.Update|ent.UpdateOne)
//
func Unless(hk ent.Hook, op ent.Op) ent.Hook {
	return If(hk, Not(HasOp(op)))
}

// FixedError is a hook returning a fixed error.
func FixedError(err error) ent.Hook {
	return func(ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(context.Context, ent.Mutation) (ent.Value, error) {
			return nil, err
		})
	}
}

// Reject returns a hook that rejects all operations that match op.
//
//	func (T) Hooks() []ent.Hook {
//		return []ent.Hook{
//			Reject(ent.Delete|ent.Update),
//		}
//	}
//
func Reject(op ent.Op) ent.Hook {
	hk := FixedError(fmt.Errorf("%s operation is not allowed", op))
	return On(hk, op)
}

// Chain acts as a list of hooks and is effectively immutable.
// Once created, it will always hold the same set of hooks in the same order.
type Chain struct {
	hooks []ent.Hook
}

// NewChain creates a new chain of hooks.
func NewChain(hooks ...ent.Hook) Chain {
	return Chain{append([]ent.Hook(nil), hooks...)}
}

// Hook chains the list of hooks and returns the final hook.
func (c Chain) Hook() ent.Hook {
	return func(mutator ent.Mutator) ent.Mutator {
		for i := len(c.hooks) - 1; i >= 0; i-- {
			mutator = c.hooks[i](mutator)
		}
		return mutator
	}
}

// Append extends a chain, adding the specified hook
// as the last ones in the mutation flow.
func (c Chain) Append(hooks ...ent.Hook) Chain {
	newHooks := make([]ent.Hook, 0, len(c.hooks)+len(hooks))
	newHooks = append(newHooks, c.hooks...)
	newHooks = append(newHooks, hooks...)
	return Chain{newHooks}
}

// Extend extends a chain, adding the specified chain
// as the last ones in the mutation flow.
func (c Chain) Extend(chain Chain) Chain {
	return c.Append(chain.hooks...)
}
