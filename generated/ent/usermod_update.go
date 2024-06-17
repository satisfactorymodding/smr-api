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
	"github.com/satisfactorymodding/smr-api/generated/ent/predicate"
	"github.com/satisfactorymodding/smr-api/generated/ent/user"
	"github.com/satisfactorymodding/smr-api/generated/ent/usermod"
)

// UserModUpdate is the builder for updating UserMod entities.
type UserModUpdate struct {
	config
	hooks     []Hook
	mutation  *UserModMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the UserModUpdate builder.
func (umu *UserModUpdate) Where(ps ...predicate.UserMod) *UserModUpdate {
	umu.mutation.Where(ps...)
	return umu
}

// SetUserID sets the "user_id" field.
func (umu *UserModUpdate) SetUserID(s string) *UserModUpdate {
	umu.mutation.SetUserID(s)
	return umu
}

// SetNillableUserID sets the "user_id" field if the given value is not nil.
func (umu *UserModUpdate) SetNillableUserID(s *string) *UserModUpdate {
	if s != nil {
		umu.SetUserID(*s)
	}
	return umu
}

// SetModID sets the "mod_id" field.
func (umu *UserModUpdate) SetModID(s string) *UserModUpdate {
	umu.mutation.SetModID(s)
	return umu
}

// SetNillableModID sets the "mod_id" field if the given value is not nil.
func (umu *UserModUpdate) SetNillableModID(s *string) *UserModUpdate {
	if s != nil {
		umu.SetModID(*s)
	}
	return umu
}

// SetRole sets the "role" field.
func (umu *UserModUpdate) SetRole(s string) *UserModUpdate {
	umu.mutation.SetRole(s)
	return umu
}

// SetNillableRole sets the "role" field if the given value is not nil.
func (umu *UserModUpdate) SetNillableRole(s *string) *UserModUpdate {
	if s != nil {
		umu.SetRole(*s)
	}
	return umu
}

// SetUser sets the "user" edge to the User entity.
func (umu *UserModUpdate) SetUser(u *User) *UserModUpdate {
	return umu.SetUserID(u.ID)
}

// SetMod sets the "mod" edge to the Mod entity.
func (umu *UserModUpdate) SetMod(m *Mod) *UserModUpdate {
	return umu.SetModID(m.ID)
}

// Mutation returns the UserModMutation object of the builder.
func (umu *UserModUpdate) Mutation() *UserModMutation {
	return umu.mutation
}

// ClearUser clears the "user" edge to the User entity.
func (umu *UserModUpdate) ClearUser() *UserModUpdate {
	umu.mutation.ClearUser()
	return umu
}

// ClearMod clears the "mod" edge to the Mod entity.
func (umu *UserModUpdate) ClearMod() *UserModUpdate {
	umu.mutation.ClearMod()
	return umu
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (umu *UserModUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, umu.sqlSave, umu.mutation, umu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (umu *UserModUpdate) SaveX(ctx context.Context) int {
	affected, err := umu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (umu *UserModUpdate) Exec(ctx context.Context) error {
	_, err := umu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (umu *UserModUpdate) ExecX(ctx context.Context) {
	if err := umu.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (umu *UserModUpdate) check() error {
	if _, ok := umu.mutation.UserID(); umu.mutation.UserCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "UserMod.user"`)
	}
	if _, ok := umu.mutation.ModID(); umu.mutation.ModCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "UserMod.mod"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (umu *UserModUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *UserModUpdate {
	umu.modifiers = append(umu.modifiers, modifiers...)
	return umu
}

func (umu *UserModUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := umu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(usermod.Table, usermod.Columns, sqlgraph.NewFieldSpec(usermod.FieldUserID, field.TypeString), sqlgraph.NewFieldSpec(usermod.FieldModID, field.TypeString))
	if ps := umu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := umu.mutation.Role(); ok {
		_spec.SetField(usermod.FieldRole, field.TypeString, value)
	}
	if umu.mutation.UserCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   usermod.UserTable,
			Columns: []string{usermod.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := umu.mutation.UserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   usermod.UserTable,
			Columns: []string{usermod.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if umu.mutation.ModCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   usermod.ModTable,
			Columns: []string{usermod.ModColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(mod.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := umu.mutation.ModIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   usermod.ModTable,
			Columns: []string{usermod.ModColumn},
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
	_spec.AddModifiers(umu.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, umu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{usermod.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	umu.mutation.done = true
	return n, nil
}

// UserModUpdateOne is the builder for updating a single UserMod entity.
type UserModUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *UserModMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetUserID sets the "user_id" field.
func (umuo *UserModUpdateOne) SetUserID(s string) *UserModUpdateOne {
	umuo.mutation.SetUserID(s)
	return umuo
}

// SetNillableUserID sets the "user_id" field if the given value is not nil.
func (umuo *UserModUpdateOne) SetNillableUserID(s *string) *UserModUpdateOne {
	if s != nil {
		umuo.SetUserID(*s)
	}
	return umuo
}

// SetModID sets the "mod_id" field.
func (umuo *UserModUpdateOne) SetModID(s string) *UserModUpdateOne {
	umuo.mutation.SetModID(s)
	return umuo
}

// SetNillableModID sets the "mod_id" field if the given value is not nil.
func (umuo *UserModUpdateOne) SetNillableModID(s *string) *UserModUpdateOne {
	if s != nil {
		umuo.SetModID(*s)
	}
	return umuo
}

// SetRole sets the "role" field.
func (umuo *UserModUpdateOne) SetRole(s string) *UserModUpdateOne {
	umuo.mutation.SetRole(s)
	return umuo
}

// SetNillableRole sets the "role" field if the given value is not nil.
func (umuo *UserModUpdateOne) SetNillableRole(s *string) *UserModUpdateOne {
	if s != nil {
		umuo.SetRole(*s)
	}
	return umuo
}

// SetUser sets the "user" edge to the User entity.
func (umuo *UserModUpdateOne) SetUser(u *User) *UserModUpdateOne {
	return umuo.SetUserID(u.ID)
}

// SetMod sets the "mod" edge to the Mod entity.
func (umuo *UserModUpdateOne) SetMod(m *Mod) *UserModUpdateOne {
	return umuo.SetModID(m.ID)
}

// Mutation returns the UserModMutation object of the builder.
func (umuo *UserModUpdateOne) Mutation() *UserModMutation {
	return umuo.mutation
}

// ClearUser clears the "user" edge to the User entity.
func (umuo *UserModUpdateOne) ClearUser() *UserModUpdateOne {
	umuo.mutation.ClearUser()
	return umuo
}

// ClearMod clears the "mod" edge to the Mod entity.
func (umuo *UserModUpdateOne) ClearMod() *UserModUpdateOne {
	umuo.mutation.ClearMod()
	return umuo
}

// Where appends a list predicates to the UserModUpdate builder.
func (umuo *UserModUpdateOne) Where(ps ...predicate.UserMod) *UserModUpdateOne {
	umuo.mutation.Where(ps...)
	return umuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (umuo *UserModUpdateOne) Select(field string, fields ...string) *UserModUpdateOne {
	umuo.fields = append([]string{field}, fields...)
	return umuo
}

// Save executes the query and returns the updated UserMod entity.
func (umuo *UserModUpdateOne) Save(ctx context.Context) (*UserMod, error) {
	return withHooks(ctx, umuo.sqlSave, umuo.mutation, umuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (umuo *UserModUpdateOne) SaveX(ctx context.Context) *UserMod {
	node, err := umuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (umuo *UserModUpdateOne) Exec(ctx context.Context) error {
	_, err := umuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (umuo *UserModUpdateOne) ExecX(ctx context.Context) {
	if err := umuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (umuo *UserModUpdateOne) check() error {
	if _, ok := umuo.mutation.UserID(); umuo.mutation.UserCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "UserMod.user"`)
	}
	if _, ok := umuo.mutation.ModID(); umuo.mutation.ModCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "UserMod.mod"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (umuo *UserModUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *UserModUpdateOne {
	umuo.modifiers = append(umuo.modifiers, modifiers...)
	return umuo
}

func (umuo *UserModUpdateOne) sqlSave(ctx context.Context) (_node *UserMod, err error) {
	if err := umuo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(usermod.Table, usermod.Columns, sqlgraph.NewFieldSpec(usermod.FieldUserID, field.TypeString), sqlgraph.NewFieldSpec(usermod.FieldModID, field.TypeString))
	if id, ok := umuo.mutation.UserID(); !ok {
		return nil, &ValidationError{Name: "user_id", err: errors.New(`ent: missing "UserMod.user_id" for update`)}
	} else {
		_spec.Node.CompositeID[0].Value = id
	}
	if id, ok := umuo.mutation.ModID(); !ok {
		return nil, &ValidationError{Name: "mod_id", err: errors.New(`ent: missing "UserMod.mod_id" for update`)}
	} else {
		_spec.Node.CompositeID[1].Value = id
	}
	if fields := umuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, len(fields))
		for i, f := range fields {
			if !usermod.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			_spec.Node.Columns[i] = f
		}
	}
	if ps := umuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := umuo.mutation.Role(); ok {
		_spec.SetField(usermod.FieldRole, field.TypeString, value)
	}
	if umuo.mutation.UserCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   usermod.UserTable,
			Columns: []string{usermod.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := umuo.mutation.UserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   usermod.UserTable,
			Columns: []string{usermod.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if umuo.mutation.ModCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   usermod.ModTable,
			Columns: []string{usermod.ModColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(mod.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := umuo.mutation.ModIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   usermod.ModTable,
			Columns: []string{usermod.ModColumn},
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
	_spec.AddModifiers(umuo.modifiers...)
	_node = &UserMod{config: umuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, umuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{usermod.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	umuo.mutation.done = true
	return _node, nil
}
