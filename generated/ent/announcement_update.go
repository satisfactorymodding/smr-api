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
	"github.com/satisfactorymodding/smr-api/generated/ent/announcement"
	"github.com/satisfactorymodding/smr-api/generated/ent/predicate"
)

// AnnouncementUpdate is the builder for updating Announcement entities.
type AnnouncementUpdate struct {
	config
	hooks     []Hook
	mutation  *AnnouncementMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the AnnouncementUpdate builder.
func (au *AnnouncementUpdate) Where(ps ...predicate.Announcement) *AnnouncementUpdate {
	au.mutation.Where(ps...)
	return au
}

// SetUpdatedAt sets the "updated_at" field.
func (au *AnnouncementUpdate) SetUpdatedAt(t time.Time) *AnnouncementUpdate {
	au.mutation.SetUpdatedAt(t)
	return au
}

// SetDeletedAt sets the "deleted_at" field.
func (au *AnnouncementUpdate) SetDeletedAt(t time.Time) *AnnouncementUpdate {
	au.mutation.SetDeletedAt(t)
	return au
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (au *AnnouncementUpdate) SetNillableDeletedAt(t *time.Time) *AnnouncementUpdate {
	if t != nil {
		au.SetDeletedAt(*t)
	}
	return au
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (au *AnnouncementUpdate) ClearDeletedAt() *AnnouncementUpdate {
	au.mutation.ClearDeletedAt()
	return au
}

// SetMessage sets the "message" field.
func (au *AnnouncementUpdate) SetMessage(s string) *AnnouncementUpdate {
	au.mutation.SetMessage(s)
	return au
}

// SetNillableMessage sets the "message" field if the given value is not nil.
func (au *AnnouncementUpdate) SetNillableMessage(s *string) *AnnouncementUpdate {
	if s != nil {
		au.SetMessage(*s)
	}
	return au
}

// SetImportance sets the "importance" field.
func (au *AnnouncementUpdate) SetImportance(s string) *AnnouncementUpdate {
	au.mutation.SetImportance(s)
	return au
}

// SetNillableImportance sets the "importance" field if the given value is not nil.
func (au *AnnouncementUpdate) SetNillableImportance(s *string) *AnnouncementUpdate {
	if s != nil {
		au.SetImportance(*s)
	}
	return au
}

// Mutation returns the AnnouncementMutation object of the builder.
func (au *AnnouncementUpdate) Mutation() *AnnouncementMutation {
	return au.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (au *AnnouncementUpdate) Save(ctx context.Context) (int, error) {
	if err := au.defaults(); err != nil {
		return 0, err
	}
	return withHooks(ctx, au.sqlSave, au.mutation, au.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (au *AnnouncementUpdate) SaveX(ctx context.Context) int {
	affected, err := au.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (au *AnnouncementUpdate) Exec(ctx context.Context) error {
	_, err := au.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (au *AnnouncementUpdate) ExecX(ctx context.Context) {
	if err := au.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (au *AnnouncementUpdate) defaults() error {
	if _, ok := au.mutation.UpdatedAt(); !ok {
		if announcement.UpdateDefaultUpdatedAt == nil {
			return fmt.Errorf("ent: uninitialized announcement.UpdateDefaultUpdatedAt (forgotten import ent/runtime?)")
		}
		v := announcement.UpdateDefaultUpdatedAt()
		au.mutation.SetUpdatedAt(v)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (au *AnnouncementUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *AnnouncementUpdate {
	au.modifiers = append(au.modifiers, modifiers...)
	return au
}

func (au *AnnouncementUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(announcement.Table, announcement.Columns, sqlgraph.NewFieldSpec(announcement.FieldID, field.TypeString))
	if ps := au.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := au.mutation.UpdatedAt(); ok {
		_spec.SetField(announcement.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := au.mutation.DeletedAt(); ok {
		_spec.SetField(announcement.FieldDeletedAt, field.TypeTime, value)
	}
	if au.mutation.DeletedAtCleared() {
		_spec.ClearField(announcement.FieldDeletedAt, field.TypeTime)
	}
	if value, ok := au.mutation.Message(); ok {
		_spec.SetField(announcement.FieldMessage, field.TypeString, value)
	}
	if value, ok := au.mutation.Importance(); ok {
		_spec.SetField(announcement.FieldImportance, field.TypeString, value)
	}
	_spec.AddModifiers(au.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, au.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{announcement.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	au.mutation.done = true
	return n, nil
}

// AnnouncementUpdateOne is the builder for updating a single Announcement entity.
type AnnouncementUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *AnnouncementMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetUpdatedAt sets the "updated_at" field.
func (auo *AnnouncementUpdateOne) SetUpdatedAt(t time.Time) *AnnouncementUpdateOne {
	auo.mutation.SetUpdatedAt(t)
	return auo
}

// SetDeletedAt sets the "deleted_at" field.
func (auo *AnnouncementUpdateOne) SetDeletedAt(t time.Time) *AnnouncementUpdateOne {
	auo.mutation.SetDeletedAt(t)
	return auo
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (auo *AnnouncementUpdateOne) SetNillableDeletedAt(t *time.Time) *AnnouncementUpdateOne {
	if t != nil {
		auo.SetDeletedAt(*t)
	}
	return auo
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (auo *AnnouncementUpdateOne) ClearDeletedAt() *AnnouncementUpdateOne {
	auo.mutation.ClearDeletedAt()
	return auo
}

// SetMessage sets the "message" field.
func (auo *AnnouncementUpdateOne) SetMessage(s string) *AnnouncementUpdateOne {
	auo.mutation.SetMessage(s)
	return auo
}

// SetNillableMessage sets the "message" field if the given value is not nil.
func (auo *AnnouncementUpdateOne) SetNillableMessage(s *string) *AnnouncementUpdateOne {
	if s != nil {
		auo.SetMessage(*s)
	}
	return auo
}

// SetImportance sets the "importance" field.
func (auo *AnnouncementUpdateOne) SetImportance(s string) *AnnouncementUpdateOne {
	auo.mutation.SetImportance(s)
	return auo
}

// SetNillableImportance sets the "importance" field if the given value is not nil.
func (auo *AnnouncementUpdateOne) SetNillableImportance(s *string) *AnnouncementUpdateOne {
	if s != nil {
		auo.SetImportance(*s)
	}
	return auo
}

// Mutation returns the AnnouncementMutation object of the builder.
func (auo *AnnouncementUpdateOne) Mutation() *AnnouncementMutation {
	return auo.mutation
}

// Where appends a list predicates to the AnnouncementUpdate builder.
func (auo *AnnouncementUpdateOne) Where(ps ...predicate.Announcement) *AnnouncementUpdateOne {
	auo.mutation.Where(ps...)
	return auo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (auo *AnnouncementUpdateOne) Select(field string, fields ...string) *AnnouncementUpdateOne {
	auo.fields = append([]string{field}, fields...)
	return auo
}

// Save executes the query and returns the updated Announcement entity.
func (auo *AnnouncementUpdateOne) Save(ctx context.Context) (*Announcement, error) {
	if err := auo.defaults(); err != nil {
		return nil, err
	}
	return withHooks(ctx, auo.sqlSave, auo.mutation, auo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (auo *AnnouncementUpdateOne) SaveX(ctx context.Context) *Announcement {
	node, err := auo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (auo *AnnouncementUpdateOne) Exec(ctx context.Context) error {
	_, err := auo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (auo *AnnouncementUpdateOne) ExecX(ctx context.Context) {
	if err := auo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (auo *AnnouncementUpdateOne) defaults() error {
	if _, ok := auo.mutation.UpdatedAt(); !ok {
		if announcement.UpdateDefaultUpdatedAt == nil {
			return fmt.Errorf("ent: uninitialized announcement.UpdateDefaultUpdatedAt (forgotten import ent/runtime?)")
		}
		v := announcement.UpdateDefaultUpdatedAt()
		auo.mutation.SetUpdatedAt(v)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (auo *AnnouncementUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *AnnouncementUpdateOne {
	auo.modifiers = append(auo.modifiers, modifiers...)
	return auo
}

func (auo *AnnouncementUpdateOne) sqlSave(ctx context.Context) (_node *Announcement, err error) {
	_spec := sqlgraph.NewUpdateSpec(announcement.Table, announcement.Columns, sqlgraph.NewFieldSpec(announcement.FieldID, field.TypeString))
	id, ok := auo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Announcement.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := auo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, announcement.FieldID)
		for _, f := range fields {
			if !announcement.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != announcement.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := auo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := auo.mutation.UpdatedAt(); ok {
		_spec.SetField(announcement.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := auo.mutation.DeletedAt(); ok {
		_spec.SetField(announcement.FieldDeletedAt, field.TypeTime, value)
	}
	if auo.mutation.DeletedAtCleared() {
		_spec.ClearField(announcement.FieldDeletedAt, field.TypeTime)
	}
	if value, ok := auo.mutation.Message(); ok {
		_spec.SetField(announcement.FieldMessage, field.TypeString, value)
	}
	if value, ok := auo.mutation.Importance(); ok {
		_spec.SetField(announcement.FieldImportance, field.TypeString, value)
	}
	_spec.AddModifiers(auo.modifiers...)
	_node = &Announcement{config: auo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, auo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{announcement.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	auo.mutation.done = true
	return _node, nil
}
