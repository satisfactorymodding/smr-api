// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/satisfactorymodding/smr-api/generated/ent/predicate"
	"github.com/satisfactorymodding/smr-api/generated/ent/user"
	"github.com/satisfactorymodding/smr-api/generated/ent/usergroup"
)

// UserGroupUpdate is the builder for updating UserGroup entities.
type UserGroupUpdate struct {
	config
	hooks     []Hook
	mutation  *UserGroupMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the UserGroupUpdate builder.
func (ugu *UserGroupUpdate) Where(ps ...predicate.UserGroup) *UserGroupUpdate {
	ugu.mutation.Where(ps...)
	return ugu
}

// SetUpdatedAt sets the "updated_at" field.
func (ugu *UserGroupUpdate) SetUpdatedAt(t time.Time) *UserGroupUpdate {
	ugu.mutation.SetUpdatedAt(t)
	return ugu
}

// SetDeletedAt sets the "deleted_at" field.
func (ugu *UserGroupUpdate) SetDeletedAt(t time.Time) *UserGroupUpdate {
	ugu.mutation.SetDeletedAt(t)
	return ugu
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (ugu *UserGroupUpdate) SetNillableDeletedAt(t *time.Time) *UserGroupUpdate {
	if t != nil {
		ugu.SetDeletedAt(*t)
	}
	return ugu
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (ugu *UserGroupUpdate) ClearDeletedAt() *UserGroupUpdate {
	ugu.mutation.ClearDeletedAt()
	return ugu
}

// SetUserID sets the "user_id" field.
func (ugu *UserGroupUpdate) SetUserID(s string) *UserGroupUpdate {
	ugu.mutation.SetUserID(s)
	return ugu
}

// SetNillableUserID sets the "user_id" field if the given value is not nil.
func (ugu *UserGroupUpdate) SetNillableUserID(s *string) *UserGroupUpdate {
	if s != nil {
		ugu.SetUserID(*s)
	}
	return ugu
}

// SetGroupID sets the "group_id" field.
func (ugu *UserGroupUpdate) SetGroupID(s string) *UserGroupUpdate {
	ugu.mutation.SetGroupID(s)
	return ugu
}

// SetNillableGroupID sets the "group_id" field if the given value is not nil.
func (ugu *UserGroupUpdate) SetNillableGroupID(s *string) *UserGroupUpdate {
	if s != nil {
		ugu.SetGroupID(*s)
	}
	return ugu
}

// SetUser sets the "user" edge to the User entity.
func (ugu *UserGroupUpdate) SetUser(u *User) *UserGroupUpdate {
	return ugu.SetUserID(u.ID)
}

// Mutation returns the UserGroupMutation object of the builder.
func (ugu *UserGroupUpdate) Mutation() *UserGroupMutation {
	return ugu.mutation
}

// ClearUser clears the "user" edge to the User entity.
func (ugu *UserGroupUpdate) ClearUser() *UserGroupUpdate {
	ugu.mutation.ClearUser()
	return ugu
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (ugu *UserGroupUpdate) Save(ctx context.Context) (int, error) {
	if err := ugu.defaults(); err != nil {
		return 0, err
	}
	return withHooks(ctx, ugu.sqlSave, ugu.mutation, ugu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (ugu *UserGroupUpdate) SaveX(ctx context.Context) int {
	affected, err := ugu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (ugu *UserGroupUpdate) Exec(ctx context.Context) error {
	_, err := ugu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ugu *UserGroupUpdate) ExecX(ctx context.Context) {
	if err := ugu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (ugu *UserGroupUpdate) defaults() error {
	if _, ok := ugu.mutation.UpdatedAt(); !ok {
		if usergroup.UpdateDefaultUpdatedAt == nil {
			return fmt.Errorf("ent: uninitialized usergroup.UpdateDefaultUpdatedAt (forgotten import ent/runtime?)")
		}
		v := usergroup.UpdateDefaultUpdatedAt()
		ugu.mutation.SetUpdatedAt(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (ugu *UserGroupUpdate) check() error {
	if v, ok := ugu.mutation.UserID(); ok {
		if err := usergroup.UserIDValidator(v); err != nil {
			return &ValidationError{Name: "user_id", err: fmt.Errorf(`ent: validator failed for field "UserGroup.user_id": %w`, err)}
		}
	}
	if v, ok := ugu.mutation.GroupID(); ok {
		if err := usergroup.GroupIDValidator(v); err != nil {
			return &ValidationError{Name: "group_id", err: fmt.Errorf(`ent: validator failed for field "UserGroup.group_id": %w`, err)}
		}
	}
	if ugu.mutation.UserCleared() && len(ugu.mutation.UserIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "UserGroup.user"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (ugu *UserGroupUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *UserGroupUpdate {
	ugu.modifiers = append(ugu.modifiers, modifiers...)
	return ugu
}

func (ugu *UserGroupUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := ugu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(usergroup.Table, usergroup.Columns, sqlgraph.NewFieldSpec(usergroup.FieldID, field.TypeString))
	if ps := ugu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ugu.mutation.UpdatedAt(); ok {
		_spec.SetField(usergroup.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := ugu.mutation.DeletedAt(); ok {
		_spec.SetField(usergroup.FieldDeletedAt, field.TypeTime, value)
	}
	if ugu.mutation.DeletedAtCleared() {
		_spec.ClearField(usergroup.FieldDeletedAt, field.TypeTime)
	}
	if value, ok := ugu.mutation.GroupID(); ok {
		_spec.SetField(usergroup.FieldGroupID, field.TypeString, value)
	}
	if ugu.mutation.UserCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   usergroup.UserTable,
			Columns: []string{usergroup.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ugu.mutation.UserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   usergroup.UserTable,
			Columns: []string{usergroup.UserColumn},
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
	_spec.AddModifiers(ugu.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, ugu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{usergroup.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	ugu.mutation.done = true
	return n, nil
}

// UserGroupUpdateOne is the builder for updating a single UserGroup entity.
type UserGroupUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *UserGroupMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetUpdatedAt sets the "updated_at" field.
func (uguo *UserGroupUpdateOne) SetUpdatedAt(t time.Time) *UserGroupUpdateOne {
	uguo.mutation.SetUpdatedAt(t)
	return uguo
}

// SetDeletedAt sets the "deleted_at" field.
func (uguo *UserGroupUpdateOne) SetDeletedAt(t time.Time) *UserGroupUpdateOne {
	uguo.mutation.SetDeletedAt(t)
	return uguo
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (uguo *UserGroupUpdateOne) SetNillableDeletedAt(t *time.Time) *UserGroupUpdateOne {
	if t != nil {
		uguo.SetDeletedAt(*t)
	}
	return uguo
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (uguo *UserGroupUpdateOne) ClearDeletedAt() *UserGroupUpdateOne {
	uguo.mutation.ClearDeletedAt()
	return uguo
}

// SetUserID sets the "user_id" field.
func (uguo *UserGroupUpdateOne) SetUserID(s string) *UserGroupUpdateOne {
	uguo.mutation.SetUserID(s)
	return uguo
}

// SetNillableUserID sets the "user_id" field if the given value is not nil.
func (uguo *UserGroupUpdateOne) SetNillableUserID(s *string) *UserGroupUpdateOne {
	if s != nil {
		uguo.SetUserID(*s)
	}
	return uguo
}

// SetGroupID sets the "group_id" field.
func (uguo *UserGroupUpdateOne) SetGroupID(s string) *UserGroupUpdateOne {
	uguo.mutation.SetGroupID(s)
	return uguo
}

// SetNillableGroupID sets the "group_id" field if the given value is not nil.
func (uguo *UserGroupUpdateOne) SetNillableGroupID(s *string) *UserGroupUpdateOne {
	if s != nil {
		uguo.SetGroupID(*s)
	}
	return uguo
}

// SetUser sets the "user" edge to the User entity.
func (uguo *UserGroupUpdateOne) SetUser(u *User) *UserGroupUpdateOne {
	return uguo.SetUserID(u.ID)
}

// Mutation returns the UserGroupMutation object of the builder.
func (uguo *UserGroupUpdateOne) Mutation() *UserGroupMutation {
	return uguo.mutation
}

// ClearUser clears the "user" edge to the User entity.
func (uguo *UserGroupUpdateOne) ClearUser() *UserGroupUpdateOne {
	uguo.mutation.ClearUser()
	return uguo
}

// Where appends a list predicates to the UserGroupUpdate builder.
func (uguo *UserGroupUpdateOne) Where(ps ...predicate.UserGroup) *UserGroupUpdateOne {
	uguo.mutation.Where(ps...)
	return uguo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (uguo *UserGroupUpdateOne) Select(field string, fields ...string) *UserGroupUpdateOne {
	uguo.fields = append([]string{field}, fields...)
	return uguo
}

// Save executes the query and returns the updated UserGroup entity.
func (uguo *UserGroupUpdateOne) Save(ctx context.Context) (*UserGroup, error) {
	if err := uguo.defaults(); err != nil {
		return nil, err
	}
	return withHooks(ctx, uguo.sqlSave, uguo.mutation, uguo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (uguo *UserGroupUpdateOne) SaveX(ctx context.Context) *UserGroup {
	node, err := uguo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (uguo *UserGroupUpdateOne) Exec(ctx context.Context) error {
	_, err := uguo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (uguo *UserGroupUpdateOne) ExecX(ctx context.Context) {
	if err := uguo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (uguo *UserGroupUpdateOne) defaults() error {
	if _, ok := uguo.mutation.UpdatedAt(); !ok {
		if usergroup.UpdateDefaultUpdatedAt == nil {
			return fmt.Errorf("ent: uninitialized usergroup.UpdateDefaultUpdatedAt (forgotten import ent/runtime?)")
		}
		v := usergroup.UpdateDefaultUpdatedAt()
		uguo.mutation.SetUpdatedAt(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (uguo *UserGroupUpdateOne) check() error {
	if v, ok := uguo.mutation.UserID(); ok {
		if err := usergroup.UserIDValidator(v); err != nil {
			return &ValidationError{Name: "user_id", err: fmt.Errorf(`ent: validator failed for field "UserGroup.user_id": %w`, err)}
		}
	}
	if v, ok := uguo.mutation.GroupID(); ok {
		if err := usergroup.GroupIDValidator(v); err != nil {
			return &ValidationError{Name: "group_id", err: fmt.Errorf(`ent: validator failed for field "UserGroup.group_id": %w`, err)}
		}
	}
	if uguo.mutation.UserCleared() && len(uguo.mutation.UserIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "UserGroup.user"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (uguo *UserGroupUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *UserGroupUpdateOne {
	uguo.modifiers = append(uguo.modifiers, modifiers...)
	return uguo
}

func (uguo *UserGroupUpdateOne) sqlSave(ctx context.Context) (_node *UserGroup, err error) {
	if err := uguo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(usergroup.Table, usergroup.Columns, sqlgraph.NewFieldSpec(usergroup.FieldID, field.TypeString))
	id, ok := uguo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "UserGroup.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := uguo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, usergroup.FieldID)
		for _, f := range fields {
			if !usergroup.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != usergroup.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := uguo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := uguo.mutation.UpdatedAt(); ok {
		_spec.SetField(usergroup.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := uguo.mutation.DeletedAt(); ok {
		_spec.SetField(usergroup.FieldDeletedAt, field.TypeTime, value)
	}
	if uguo.mutation.DeletedAtCleared() {
		_spec.ClearField(usergroup.FieldDeletedAt, field.TypeTime)
	}
	if value, ok := uguo.mutation.GroupID(); ok {
		_spec.SetField(usergroup.FieldGroupID, field.TypeString, value)
	}
	if uguo.mutation.UserCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   usergroup.UserTable,
			Columns: []string{usergroup.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uguo.mutation.UserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   usergroup.UserTable,
			Columns: []string{usergroup.UserColumn},
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
	_spec.AddModifiers(uguo.modifiers...)
	_node = &UserGroup{config: uguo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, uguo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{usergroup.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	uguo.mutation.done = true
	return _node, nil
}
