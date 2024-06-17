// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/satisfactorymodding/smr-api/generated/ent/predicate"
	"github.com/satisfactorymodding/smr-api/generated/ent/smlversion"
	"github.com/satisfactorymodding/smr-api/generated/ent/smlversiontarget"
)

// SmlVersionTargetUpdate is the builder for updating SmlVersionTarget entities.
type SmlVersionTargetUpdate struct {
	config
	hooks     []Hook
	mutation  *SmlVersionTargetMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the SmlVersionTargetUpdate builder.
func (svtu *SmlVersionTargetUpdate) Where(ps ...predicate.SmlVersionTarget) *SmlVersionTargetUpdate {
	svtu.mutation.Where(ps...)
	return svtu
}

// SetVersionID sets the "version_id" field.
func (svtu *SmlVersionTargetUpdate) SetVersionID(s string) *SmlVersionTargetUpdate {
	svtu.mutation.SetVersionID(s)
	return svtu
}

// SetNillableVersionID sets the "version_id" field if the given value is not nil.
func (svtu *SmlVersionTargetUpdate) SetNillableVersionID(s *string) *SmlVersionTargetUpdate {
	if s != nil {
		svtu.SetVersionID(*s)
	}
	return svtu
}

// SetTargetName sets the "target_name" field.
func (svtu *SmlVersionTargetUpdate) SetTargetName(s string) *SmlVersionTargetUpdate {
	svtu.mutation.SetTargetName(s)
	return svtu
}

// SetNillableTargetName sets the "target_name" field if the given value is not nil.
func (svtu *SmlVersionTargetUpdate) SetNillableTargetName(s *string) *SmlVersionTargetUpdate {
	if s != nil {
		svtu.SetTargetName(*s)
	}
	return svtu
}

// SetLink sets the "link" field.
func (svtu *SmlVersionTargetUpdate) SetLink(s string) *SmlVersionTargetUpdate {
	svtu.mutation.SetLink(s)
	return svtu
}

// SetNillableLink sets the "link" field if the given value is not nil.
func (svtu *SmlVersionTargetUpdate) SetNillableLink(s *string) *SmlVersionTargetUpdate {
	if s != nil {
		svtu.SetLink(*s)
	}
	return svtu
}

// SetSmlVersionID sets the "sml_version" edge to the SmlVersion entity by ID.
func (svtu *SmlVersionTargetUpdate) SetSmlVersionID(id string) *SmlVersionTargetUpdate {
	svtu.mutation.SetSmlVersionID(id)
	return svtu
}

// SetSmlVersion sets the "sml_version" edge to the SmlVersion entity.
func (svtu *SmlVersionTargetUpdate) SetSmlVersion(s *SmlVersion) *SmlVersionTargetUpdate {
	return svtu.SetSmlVersionID(s.ID)
}

// Mutation returns the SmlVersionTargetMutation object of the builder.
func (svtu *SmlVersionTargetUpdate) Mutation() *SmlVersionTargetMutation {
	return svtu.mutation
}

// ClearSmlVersion clears the "sml_version" edge to the SmlVersion entity.
func (svtu *SmlVersionTargetUpdate) ClearSmlVersion() *SmlVersionTargetUpdate {
	svtu.mutation.ClearSmlVersion()
	return svtu
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (svtu *SmlVersionTargetUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, svtu.sqlSave, svtu.mutation, svtu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (svtu *SmlVersionTargetUpdate) SaveX(ctx context.Context) int {
	affected, err := svtu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (svtu *SmlVersionTargetUpdate) Exec(ctx context.Context) error {
	_, err := svtu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (svtu *SmlVersionTargetUpdate) ExecX(ctx context.Context) {
	if err := svtu.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (svtu *SmlVersionTargetUpdate) check() error {
	if _, ok := svtu.mutation.SmlVersionID(); svtu.mutation.SmlVersionCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "SmlVersionTarget.sml_version"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (svtu *SmlVersionTargetUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *SmlVersionTargetUpdate {
	svtu.modifiers = append(svtu.modifiers, modifiers...)
	return svtu
}

func (svtu *SmlVersionTargetUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := svtu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(smlversiontarget.Table, smlversiontarget.Columns, sqlgraph.NewFieldSpec(smlversiontarget.FieldID, field.TypeString))
	if ps := svtu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := svtu.mutation.TargetName(); ok {
		_spec.SetField(smlversiontarget.FieldTargetName, field.TypeString, value)
	}
	if value, ok := svtu.mutation.Link(); ok {
		_spec.SetField(smlversiontarget.FieldLink, field.TypeString, value)
	}
	if svtu.mutation.SmlVersionCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   smlversiontarget.SmlVersionTable,
			Columns: []string{smlversiontarget.SmlVersionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(smlversion.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := svtu.mutation.SmlVersionIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   smlversiontarget.SmlVersionTable,
			Columns: []string{smlversiontarget.SmlVersionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(smlversion.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.AddModifiers(svtu.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, svtu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{smlversiontarget.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	svtu.mutation.done = true
	return n, nil
}

// SmlVersionTargetUpdateOne is the builder for updating a single SmlVersionTarget entity.
type SmlVersionTargetUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *SmlVersionTargetMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetVersionID sets the "version_id" field.
func (svtuo *SmlVersionTargetUpdateOne) SetVersionID(s string) *SmlVersionTargetUpdateOne {
	svtuo.mutation.SetVersionID(s)
	return svtuo
}

// SetNillableVersionID sets the "version_id" field if the given value is not nil.
func (svtuo *SmlVersionTargetUpdateOne) SetNillableVersionID(s *string) *SmlVersionTargetUpdateOne {
	if s != nil {
		svtuo.SetVersionID(*s)
	}
	return svtuo
}

// SetTargetName sets the "target_name" field.
func (svtuo *SmlVersionTargetUpdateOne) SetTargetName(s string) *SmlVersionTargetUpdateOne {
	svtuo.mutation.SetTargetName(s)
	return svtuo
}

// SetNillableTargetName sets the "target_name" field if the given value is not nil.
func (svtuo *SmlVersionTargetUpdateOne) SetNillableTargetName(s *string) *SmlVersionTargetUpdateOne {
	if s != nil {
		svtuo.SetTargetName(*s)
	}
	return svtuo
}

// SetLink sets the "link" field.
func (svtuo *SmlVersionTargetUpdateOne) SetLink(s string) *SmlVersionTargetUpdateOne {
	svtuo.mutation.SetLink(s)
	return svtuo
}

// SetNillableLink sets the "link" field if the given value is not nil.
func (svtuo *SmlVersionTargetUpdateOne) SetNillableLink(s *string) *SmlVersionTargetUpdateOne {
	if s != nil {
		svtuo.SetLink(*s)
	}
	return svtuo
}

// SetSmlVersionID sets the "sml_version" edge to the SmlVersion entity by ID.
func (svtuo *SmlVersionTargetUpdateOne) SetSmlVersionID(id string) *SmlVersionTargetUpdateOne {
	svtuo.mutation.SetSmlVersionID(id)
	return svtuo
}

// SetSmlVersion sets the "sml_version" edge to the SmlVersion entity.
func (svtuo *SmlVersionTargetUpdateOne) SetSmlVersion(s *SmlVersion) *SmlVersionTargetUpdateOne {
	return svtuo.SetSmlVersionID(s.ID)
}

// Mutation returns the SmlVersionTargetMutation object of the builder.
func (svtuo *SmlVersionTargetUpdateOne) Mutation() *SmlVersionTargetMutation {
	return svtuo.mutation
}

// ClearSmlVersion clears the "sml_version" edge to the SmlVersion entity.
func (svtuo *SmlVersionTargetUpdateOne) ClearSmlVersion() *SmlVersionTargetUpdateOne {
	svtuo.mutation.ClearSmlVersion()
	return svtuo
}

// Where appends a list predicates to the SmlVersionTargetUpdate builder.
func (svtuo *SmlVersionTargetUpdateOne) Where(ps ...predicate.SmlVersionTarget) *SmlVersionTargetUpdateOne {
	svtuo.mutation.Where(ps...)
	return svtuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (svtuo *SmlVersionTargetUpdateOne) Select(field string, fields ...string) *SmlVersionTargetUpdateOne {
	svtuo.fields = append([]string{field}, fields...)
	return svtuo
}

// Save executes the query and returns the updated SmlVersionTarget entity.
func (svtuo *SmlVersionTargetUpdateOne) Save(ctx context.Context) (*SmlVersionTarget, error) {
	return withHooks(ctx, svtuo.sqlSave, svtuo.mutation, svtuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (svtuo *SmlVersionTargetUpdateOne) SaveX(ctx context.Context) *SmlVersionTarget {
	node, err := svtuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (svtuo *SmlVersionTargetUpdateOne) Exec(ctx context.Context) error {
	_, err := svtuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (svtuo *SmlVersionTargetUpdateOne) ExecX(ctx context.Context) {
	if err := svtuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (svtuo *SmlVersionTargetUpdateOne) check() error {
	if _, ok := svtuo.mutation.SmlVersionID(); svtuo.mutation.SmlVersionCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "SmlVersionTarget.sml_version"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (svtuo *SmlVersionTargetUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *SmlVersionTargetUpdateOne {
	svtuo.modifiers = append(svtuo.modifiers, modifiers...)
	return svtuo
}

func (svtuo *SmlVersionTargetUpdateOne) sqlSave(ctx context.Context) (_node *SmlVersionTarget, err error) {
	if err := svtuo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(smlversiontarget.Table, smlversiontarget.Columns, sqlgraph.NewFieldSpec(smlversiontarget.FieldID, field.TypeString))
	id, ok := svtuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "SmlVersionTarget.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := svtuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, smlversiontarget.FieldID)
		for _, f := range fields {
			if !smlversiontarget.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != smlversiontarget.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := svtuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := svtuo.mutation.TargetName(); ok {
		_spec.SetField(smlversiontarget.FieldTargetName, field.TypeString, value)
	}
	if value, ok := svtuo.mutation.Link(); ok {
		_spec.SetField(smlversiontarget.FieldLink, field.TypeString, value)
	}
	if svtuo.mutation.SmlVersionCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   smlversiontarget.SmlVersionTable,
			Columns: []string{smlversiontarget.SmlVersionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(smlversion.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := svtuo.mutation.SmlVersionIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   smlversiontarget.SmlVersionTable,
			Columns: []string{smlversiontarget.SmlVersionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(smlversion.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.AddModifiers(svtuo.modifiers...)
	_node = &SmlVersionTarget{config: svtuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, svtuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{smlversiontarget.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	svtuo.mutation.done = true
	return _node, nil
}
