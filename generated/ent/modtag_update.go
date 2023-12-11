// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/satisfactorymodding/smr-api/generated/ent/mod"
	"github.com/satisfactorymodding/smr-api/generated/ent/modtag"
	"github.com/satisfactorymodding/smr-api/generated/ent/predicate"
	"github.com/satisfactorymodding/smr-api/generated/ent/tag"
)

// ModTagUpdate is the builder for updating ModTag entities.
type ModTagUpdate struct {
	config
	hooks     []Hook
	mutation  *ModTagMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the ModTagUpdate builder.
func (mtu *ModTagUpdate) Where(ps ...predicate.ModTag) *ModTagUpdate {
	mtu.mutation.Where(ps...)
	return mtu
}

// SetModID sets the "mod_id" field.
func (mtu *ModTagUpdate) SetModID(s string) *ModTagUpdate {
	mtu.mutation.SetModID(s)
	return mtu
}

// SetTagID sets the "tag_id" field.
func (mtu *ModTagUpdate) SetTagID(s string) *ModTagUpdate {
	mtu.mutation.SetTagID(s)
	return mtu
}

// SetMod sets the "mod" edge to the Mod entity.
func (mtu *ModTagUpdate) SetMod(m *Mod) *ModTagUpdate {
	return mtu.SetModID(m.ID)
}

// SetTag sets the "tag" edge to the Tag entity.
func (mtu *ModTagUpdate) SetTag(t *Tag) *ModTagUpdate {
	return mtu.SetTagID(t.ID)
}

// Mutation returns the ModTagMutation object of the builder.
func (mtu *ModTagUpdate) Mutation() *ModTagMutation {
	return mtu.mutation
}

// ClearMod clears the "mod" edge to the Mod entity.
func (mtu *ModTagUpdate) ClearMod() *ModTagUpdate {
	mtu.mutation.ClearMod()
	return mtu
}

// ClearTag clears the "tag" edge to the Tag entity.
func (mtu *ModTagUpdate) ClearTag() *ModTagUpdate {
	mtu.mutation.ClearTag()
	return mtu
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (mtu *ModTagUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, mtu.sqlSave, mtu.mutation, mtu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (mtu *ModTagUpdate) SaveX(ctx context.Context) int {
	affected, err := mtu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (mtu *ModTagUpdate) Exec(ctx context.Context) error {
	_, err := mtu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (mtu *ModTagUpdate) ExecX(ctx context.Context) {
	if err := mtu.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (mtu *ModTagUpdate) check() error {
	if _, ok := mtu.mutation.ModID(); mtu.mutation.ModCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "ModTag.mod"`)
	}
	if _, ok := mtu.mutation.TagID(); mtu.mutation.TagCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "ModTag.tag"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (mtu *ModTagUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *ModTagUpdate {
	mtu.modifiers = append(mtu.modifiers, modifiers...)
	return mtu
}

func (mtu *ModTagUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := mtu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(modtag.Table, modtag.Columns, sqlgraph.NewFieldSpec(modtag.FieldModID, field.TypeString), sqlgraph.NewFieldSpec(modtag.FieldTagID, field.TypeString))
	if ps := mtu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if mtu.mutation.ModCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   modtag.ModTable,
			Columns: []string{modtag.ModColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(mod.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := mtu.mutation.ModIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   modtag.ModTable,
			Columns: []string{modtag.ModColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(mod.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if mtu.mutation.TagCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   modtag.TagTable,
			Columns: []string{modtag.TagColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(tag.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := mtu.mutation.TagIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   modtag.TagTable,
			Columns: []string{modtag.TagColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(tag.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.AddModifiers(mtu.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, mtu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{modtag.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	mtu.mutation.done = true
	return n, nil
}

// ModTagUpdateOne is the builder for updating a single ModTag entity.
type ModTagUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *ModTagMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetModID sets the "mod_id" field.
func (mtuo *ModTagUpdateOne) SetModID(s string) *ModTagUpdateOne {
	mtuo.mutation.SetModID(s)
	return mtuo
}

// SetTagID sets the "tag_id" field.
func (mtuo *ModTagUpdateOne) SetTagID(s string) *ModTagUpdateOne {
	mtuo.mutation.SetTagID(s)
	return mtuo
}

// SetMod sets the "mod" edge to the Mod entity.
func (mtuo *ModTagUpdateOne) SetMod(m *Mod) *ModTagUpdateOne {
	return mtuo.SetModID(m.ID)
}

// SetTag sets the "tag" edge to the Tag entity.
func (mtuo *ModTagUpdateOne) SetTag(t *Tag) *ModTagUpdateOne {
	return mtuo.SetTagID(t.ID)
}

// Mutation returns the ModTagMutation object of the builder.
func (mtuo *ModTagUpdateOne) Mutation() *ModTagMutation {
	return mtuo.mutation
}

// ClearMod clears the "mod" edge to the Mod entity.
func (mtuo *ModTagUpdateOne) ClearMod() *ModTagUpdateOne {
	mtuo.mutation.ClearMod()
	return mtuo
}

// ClearTag clears the "tag" edge to the Tag entity.
func (mtuo *ModTagUpdateOne) ClearTag() *ModTagUpdateOne {
	mtuo.mutation.ClearTag()
	return mtuo
}

// Where appends a list predicates to the ModTagUpdate builder.
func (mtuo *ModTagUpdateOne) Where(ps ...predicate.ModTag) *ModTagUpdateOne {
	mtuo.mutation.Where(ps...)
	return mtuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (mtuo *ModTagUpdateOne) Select(field string, fields ...string) *ModTagUpdateOne {
	mtuo.fields = append([]string{field}, fields...)
	return mtuo
}

// Save executes the query and returns the updated ModTag entity.
func (mtuo *ModTagUpdateOne) Save(ctx context.Context) (*ModTag, error) {
	return withHooks(ctx, mtuo.sqlSave, mtuo.mutation, mtuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (mtuo *ModTagUpdateOne) SaveX(ctx context.Context) *ModTag {
	node, err := mtuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (mtuo *ModTagUpdateOne) Exec(ctx context.Context) error {
	_, err := mtuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (mtuo *ModTagUpdateOne) ExecX(ctx context.Context) {
	if err := mtuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (mtuo *ModTagUpdateOne) check() error {
	if _, ok := mtuo.mutation.ModID(); mtuo.mutation.ModCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "ModTag.mod"`)
	}
	if _, ok := mtuo.mutation.TagID(); mtuo.mutation.TagCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "ModTag.tag"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (mtuo *ModTagUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *ModTagUpdateOne {
	mtuo.modifiers = append(mtuo.modifiers, modifiers...)
	return mtuo
}

func (mtuo *ModTagUpdateOne) sqlSave(ctx context.Context) (_node *ModTag, err error) {
	if err := mtuo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(modtag.Table, modtag.Columns, sqlgraph.NewFieldSpec(modtag.FieldModID, field.TypeString), sqlgraph.NewFieldSpec(modtag.FieldTagID, field.TypeString))
	if id, ok := mtuo.mutation.ModID(); !ok {
		return nil, &ValidationError{Name: "mod_id", err: errors.New(`ent: missing "ModTag.mod_id" for update`)}
	} else {
		_spec.Node.CompositeID[0].Value = id
	}
	if id, ok := mtuo.mutation.TagID(); !ok {
		return nil, &ValidationError{Name: "tag_id", err: errors.New(`ent: missing "ModTag.tag_id" for update`)}
	} else {
		_spec.Node.CompositeID[1].Value = id
	}
	if fields := mtuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, len(fields))
		for i, f := range fields {
			if !modtag.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			_spec.Node.Columns[i] = f
		}
	}
	if ps := mtuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if mtuo.mutation.ModCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   modtag.ModTable,
			Columns: []string{modtag.ModColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(mod.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := mtuo.mutation.ModIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   modtag.ModTable,
			Columns: []string{modtag.ModColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(mod.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if mtuo.mutation.TagCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   modtag.TagTable,
			Columns: []string{modtag.TagColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(tag.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := mtuo.mutation.TagIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   modtag.TagTable,
			Columns: []string{modtag.TagColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(tag.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.AddModifiers(mtuo.modifiers...)
	_node = &ModTag{config: mtuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, mtuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{modtag.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	mtuo.mutation.done = true
	return _node, nil
}