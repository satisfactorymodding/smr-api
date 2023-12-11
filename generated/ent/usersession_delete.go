// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/satisfactorymodding/smr-api/generated/ent/predicate"
	"github.com/satisfactorymodding/smr-api/generated/ent/usersession"
)

// UserSessionDelete is the builder for deleting a UserSession entity.
type UserSessionDelete struct {
	config
	hooks    []Hook
	mutation *UserSessionMutation
}

// Where appends a list predicates to the UserSessionDelete builder.
func (usd *UserSessionDelete) Where(ps ...predicate.UserSession) *UserSessionDelete {
	usd.mutation.Where(ps...)
	return usd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (usd *UserSessionDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, usd.sqlExec, usd.mutation, usd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (usd *UserSessionDelete) ExecX(ctx context.Context) int {
	n, err := usd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (usd *UserSessionDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(usersession.Table, sqlgraph.NewFieldSpec(usersession.FieldID, field.TypeString))
	if ps := usd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, usd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	usd.mutation.done = true
	return affected, err
}

// UserSessionDeleteOne is the builder for deleting a single UserSession entity.
type UserSessionDeleteOne struct {
	usd *UserSessionDelete
}

// Where appends a list predicates to the UserSessionDelete builder.
func (usdo *UserSessionDeleteOne) Where(ps ...predicate.UserSession) *UserSessionDeleteOne {
	usdo.usd.mutation.Where(ps...)
	return usdo
}

// Exec executes the deletion query.
func (usdo *UserSessionDeleteOne) Exec(ctx context.Context) error {
	n, err := usdo.usd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{usersession.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (usdo *UserSessionDeleteOne) ExecX(ctx context.Context) {
	if err := usdo.Exec(ctx); err != nil {
		panic(err)
	}
}