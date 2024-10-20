// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/satisfactorymodding/smr-api/generated/ent/guidetag"
	"github.com/satisfactorymodding/smr-api/generated/ent/predicate"
)

// GuideTagDelete is the builder for deleting a GuideTag entity.
type GuideTagDelete struct {
	config
	hooks    []Hook
	mutation *GuideTagMutation
}

// Where appends a list predicates to the GuideTagDelete builder.
func (gtd *GuideTagDelete) Where(ps ...predicate.GuideTag) *GuideTagDelete {
	gtd.mutation.Where(ps...)
	return gtd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (gtd *GuideTagDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, gtd.sqlExec, gtd.mutation, gtd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (gtd *GuideTagDelete) ExecX(ctx context.Context) int {
	n, err := gtd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (gtd *GuideTagDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(guidetag.Table, nil)
	if ps := gtd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, gtd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	gtd.mutation.done = true
	return affected, err
}

// GuideTagDeleteOne is the builder for deleting a single GuideTag entity.
type GuideTagDeleteOne struct {
	gtd *GuideTagDelete
}

// Where appends a list predicates to the GuideTagDelete builder.
func (gtdo *GuideTagDeleteOne) Where(ps ...predicate.GuideTag) *GuideTagDeleteOne {
	gtdo.gtd.mutation.Where(ps...)
	return gtdo
}

// Exec executes the deletion query.
func (gtdo *GuideTagDeleteOne) Exec(ctx context.Context) error {
	n, err := gtdo.gtd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{guidetag.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (gtdo *GuideTagDeleteOne) ExecX(ctx context.Context) {
	if err := gtdo.Exec(ctx); err != nil {
		panic(err)
	}
}
