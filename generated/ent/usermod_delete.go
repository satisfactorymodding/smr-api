// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/satisfactorymodding/smr-api/generated/ent/predicate"
	"github.com/satisfactorymodding/smr-api/generated/ent/usermod"
)

// UserModDelete is the builder for deleting a UserMod entity.
type UserModDelete struct {
	config
	hooks    []Hook
	mutation *UserModMutation
}

// Where appends a list predicates to the UserModDelete builder.
func (umd *UserModDelete) Where(ps ...predicate.UserMod) *UserModDelete {
	umd.mutation.Where(ps...)
	return umd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (umd *UserModDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, umd.sqlExec, umd.mutation, umd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (umd *UserModDelete) ExecX(ctx context.Context) int {
	n, err := umd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (umd *UserModDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(usermod.Table, nil)
	if ps := umd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, umd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	umd.mutation.done = true
	return affected, err
}

// UserModDeleteOne is the builder for deleting a single UserMod entity.
type UserModDeleteOne struct {
	umd *UserModDelete
}

// Where appends a list predicates to the UserModDelete builder.
func (umdo *UserModDeleteOne) Where(ps ...predicate.UserMod) *UserModDeleteOne {
	umdo.umd.mutation.Where(ps...)
	return umdo
}

// Exec executes the deletion query.
func (umdo *UserModDeleteOne) Exec(ctx context.Context) error {
	n, err := umdo.umd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{usermod.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (umdo *UserModDeleteOne) ExecX(ctx context.Context) {
	if err := umdo.Exec(ctx); err != nil {
		panic(err)
	}
}
