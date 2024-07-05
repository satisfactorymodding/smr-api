// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/satisfactorymodding/smr-api/generated/ent/smlversion"
	"github.com/satisfactorymodding/smr-api/generated/ent/smlversiontarget"
	"github.com/satisfactorymodding/smr-api/util"
)

// SmlVersionCreate is the builder for creating a SmlVersion entity.
type SmlVersionCreate struct {
	config
	mutation *SmlVersionMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetCreatedAt sets the "created_at" field.
func (svc *SmlVersionCreate) SetCreatedAt(t time.Time) *SmlVersionCreate {
	svc.mutation.SetCreatedAt(t)
	return svc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (svc *SmlVersionCreate) SetNillableCreatedAt(t *time.Time) *SmlVersionCreate {
	if t != nil {
		svc.SetCreatedAt(*t)
	}
	return svc
}

// SetUpdatedAt sets the "updated_at" field.
func (svc *SmlVersionCreate) SetUpdatedAt(t time.Time) *SmlVersionCreate {
	svc.mutation.SetUpdatedAt(t)
	return svc
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (svc *SmlVersionCreate) SetNillableUpdatedAt(t *time.Time) *SmlVersionCreate {
	if t != nil {
		svc.SetUpdatedAt(*t)
	}
	return svc
}

// SetDeletedAt sets the "deleted_at" field.
func (svc *SmlVersionCreate) SetDeletedAt(t time.Time) *SmlVersionCreate {
	svc.mutation.SetDeletedAt(t)
	return svc
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (svc *SmlVersionCreate) SetNillableDeletedAt(t *time.Time) *SmlVersionCreate {
	if t != nil {
		svc.SetDeletedAt(*t)
	}
	return svc
}

// SetVersion sets the "version" field.
func (svc *SmlVersionCreate) SetVersion(s string) *SmlVersionCreate {
	svc.mutation.SetVersion(s)
	return svc
}

// SetSatisfactoryVersion sets the "satisfactory_version" field.
func (svc *SmlVersionCreate) SetSatisfactoryVersion(i int) *SmlVersionCreate {
	svc.mutation.SetSatisfactoryVersion(i)
	return svc
}

// SetStability sets the "stability" field.
func (svc *SmlVersionCreate) SetStability(u util.Stability) *SmlVersionCreate {
	svc.mutation.SetStability(u)
	return svc
}

// SetDate sets the "date" field.
func (svc *SmlVersionCreate) SetDate(t time.Time) *SmlVersionCreate {
	svc.mutation.SetDate(t)
	return svc
}

// SetLink sets the "link" field.
func (svc *SmlVersionCreate) SetLink(s string) *SmlVersionCreate {
	svc.mutation.SetLink(s)
	return svc
}

// SetChangelog sets the "changelog" field.
func (svc *SmlVersionCreate) SetChangelog(s string) *SmlVersionCreate {
	svc.mutation.SetChangelog(s)
	return svc
}

// SetBootstrapVersion sets the "bootstrap_version" field.
func (svc *SmlVersionCreate) SetBootstrapVersion(s string) *SmlVersionCreate {
	svc.mutation.SetBootstrapVersion(s)
	return svc
}

// SetNillableBootstrapVersion sets the "bootstrap_version" field if the given value is not nil.
func (svc *SmlVersionCreate) SetNillableBootstrapVersion(s *string) *SmlVersionCreate {
	if s != nil {
		svc.SetBootstrapVersion(*s)
	}
	return svc
}

// SetEngineVersion sets the "engine_version" field.
func (svc *SmlVersionCreate) SetEngineVersion(s string) *SmlVersionCreate {
	svc.mutation.SetEngineVersion(s)
	return svc
}

// SetNillableEngineVersion sets the "engine_version" field if the given value is not nil.
func (svc *SmlVersionCreate) SetNillableEngineVersion(s *string) *SmlVersionCreate {
	if s != nil {
		svc.SetEngineVersion(*s)
	}
	return svc
}

// SetID sets the "id" field.
func (svc *SmlVersionCreate) SetID(s string) *SmlVersionCreate {
	svc.mutation.SetID(s)
	return svc
}

// SetNillableID sets the "id" field if the given value is not nil.
func (svc *SmlVersionCreate) SetNillableID(s *string) *SmlVersionCreate {
	if s != nil {
		svc.SetID(*s)
	}
	return svc
}

// AddTargetIDs adds the "targets" edge to the SmlVersionTarget entity by IDs.
func (svc *SmlVersionCreate) AddTargetIDs(ids ...string) *SmlVersionCreate {
	svc.mutation.AddTargetIDs(ids...)
	return svc
}

// AddTargets adds the "targets" edges to the SmlVersionTarget entity.
func (svc *SmlVersionCreate) AddTargets(s ...*SmlVersionTarget) *SmlVersionCreate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return svc.AddTargetIDs(ids...)
}

// Mutation returns the SmlVersionMutation object of the builder.
func (svc *SmlVersionCreate) Mutation() *SmlVersionMutation {
	return svc.mutation
}

// Save creates the SmlVersion in the database.
func (svc *SmlVersionCreate) Save(ctx context.Context) (*SmlVersion, error) {
	if err := svc.defaults(); err != nil {
		return nil, err
	}
	return withHooks(ctx, svc.sqlSave, svc.mutation, svc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (svc *SmlVersionCreate) SaveX(ctx context.Context) *SmlVersion {
	v, err := svc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (svc *SmlVersionCreate) Exec(ctx context.Context) error {
	_, err := svc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (svc *SmlVersionCreate) ExecX(ctx context.Context) {
	if err := svc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (svc *SmlVersionCreate) defaults() error {
	if _, ok := svc.mutation.CreatedAt(); !ok {
		if smlversion.DefaultCreatedAt == nil {
			return fmt.Errorf("ent: uninitialized smlversion.DefaultCreatedAt (forgotten import ent/runtime?)")
		}
		v := smlversion.DefaultCreatedAt()
		svc.mutation.SetCreatedAt(v)
	}
	if _, ok := svc.mutation.UpdatedAt(); !ok {
		if smlversion.DefaultUpdatedAt == nil {
			return fmt.Errorf("ent: uninitialized smlversion.DefaultUpdatedAt (forgotten import ent/runtime?)")
		}
		v := smlversion.DefaultUpdatedAt()
		svc.mutation.SetUpdatedAt(v)
	}
	if _, ok := svc.mutation.EngineVersion(); !ok {
		v := smlversion.DefaultEngineVersion
		svc.mutation.SetEngineVersion(v)
	}
	if _, ok := svc.mutation.ID(); !ok {
		if smlversion.DefaultID == nil {
			return fmt.Errorf("ent: uninitialized smlversion.DefaultID (forgotten import ent/runtime?)")
		}
		v := smlversion.DefaultID()
		svc.mutation.SetID(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (svc *SmlVersionCreate) check() error {
	if _, ok := svc.mutation.Version(); !ok {
		return &ValidationError{Name: "version", err: errors.New(`ent: missing required field "SmlVersion.version"`)}
	}
	if v, ok := svc.mutation.Version(); ok {
		if err := smlversion.VersionValidator(v); err != nil {
			return &ValidationError{Name: "version", err: fmt.Errorf(`ent: validator failed for field "SmlVersion.version": %w`, err)}
		}
	}
	if _, ok := svc.mutation.SatisfactoryVersion(); !ok {
		return &ValidationError{Name: "satisfactory_version", err: errors.New(`ent: missing required field "SmlVersion.satisfactory_version"`)}
	}
	if _, ok := svc.mutation.Stability(); !ok {
		return &ValidationError{Name: "stability", err: errors.New(`ent: missing required field "SmlVersion.stability"`)}
	}
	if v, ok := svc.mutation.Stability(); ok {
		if err := smlversion.StabilityValidator(v); err != nil {
			return &ValidationError{Name: "stability", err: fmt.Errorf(`ent: validator failed for field "SmlVersion.stability": %w`, err)}
		}
	}
	if _, ok := svc.mutation.Date(); !ok {
		return &ValidationError{Name: "date", err: errors.New(`ent: missing required field "SmlVersion.date"`)}
	}
	if _, ok := svc.mutation.Link(); !ok {
		return &ValidationError{Name: "link", err: errors.New(`ent: missing required field "SmlVersion.link"`)}
	}
	if _, ok := svc.mutation.Changelog(); !ok {
		return &ValidationError{Name: "changelog", err: errors.New(`ent: missing required field "SmlVersion.changelog"`)}
	}
	if v, ok := svc.mutation.BootstrapVersion(); ok {
		if err := smlversion.BootstrapVersionValidator(v); err != nil {
			return &ValidationError{Name: "bootstrap_version", err: fmt.Errorf(`ent: validator failed for field "SmlVersion.bootstrap_version": %w`, err)}
		}
	}
	if _, ok := svc.mutation.EngineVersion(); !ok {
		return &ValidationError{Name: "engine_version", err: errors.New(`ent: missing required field "SmlVersion.engine_version"`)}
	}
	if v, ok := svc.mutation.EngineVersion(); ok {
		if err := smlversion.EngineVersionValidator(v); err != nil {
			return &ValidationError{Name: "engine_version", err: fmt.Errorf(`ent: validator failed for field "SmlVersion.engine_version": %w`, err)}
		}
	}
	return nil
}

func (svc *SmlVersionCreate) sqlSave(ctx context.Context) (*SmlVersion, error) {
	if err := svc.check(); err != nil {
		return nil, err
	}
	_node, _spec := svc.createSpec()
	if err := sqlgraph.CreateNode(ctx, svc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(string); ok {
			_node.ID = id
		} else {
			return nil, fmt.Errorf("unexpected SmlVersion.ID type: %T", _spec.ID.Value)
		}
	}
	svc.mutation.id = &_node.ID
	svc.mutation.done = true
	return _node, nil
}

func (svc *SmlVersionCreate) createSpec() (*SmlVersion, *sqlgraph.CreateSpec) {
	var (
		_node = &SmlVersion{config: svc.config}
		_spec = sqlgraph.NewCreateSpec(smlversion.Table, sqlgraph.NewFieldSpec(smlversion.FieldID, field.TypeString))
	)
	_spec.OnConflict = svc.conflict
	if id, ok := svc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = id
	}
	if value, ok := svc.mutation.CreatedAt(); ok {
		_spec.SetField(smlversion.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if value, ok := svc.mutation.UpdatedAt(); ok {
		_spec.SetField(smlversion.FieldUpdatedAt, field.TypeTime, value)
		_node.UpdatedAt = value
	}
	if value, ok := svc.mutation.DeletedAt(); ok {
		_spec.SetField(smlversion.FieldDeletedAt, field.TypeTime, value)
		_node.DeletedAt = value
	}
	if value, ok := svc.mutation.Version(); ok {
		_spec.SetField(smlversion.FieldVersion, field.TypeString, value)
		_node.Version = value
	}
	if value, ok := svc.mutation.SatisfactoryVersion(); ok {
		_spec.SetField(smlversion.FieldSatisfactoryVersion, field.TypeInt, value)
		_node.SatisfactoryVersion = value
	}
	if value, ok := svc.mutation.Stability(); ok {
		_spec.SetField(smlversion.FieldStability, field.TypeEnum, value)
		_node.Stability = value
	}
	if value, ok := svc.mutation.Date(); ok {
		_spec.SetField(smlversion.FieldDate, field.TypeTime, value)
		_node.Date = value
	}
	if value, ok := svc.mutation.Link(); ok {
		_spec.SetField(smlversion.FieldLink, field.TypeString, value)
		_node.Link = value
	}
	if value, ok := svc.mutation.Changelog(); ok {
		_spec.SetField(smlversion.FieldChangelog, field.TypeString, value)
		_node.Changelog = value
	}
	if value, ok := svc.mutation.BootstrapVersion(); ok {
		_spec.SetField(smlversion.FieldBootstrapVersion, field.TypeString, value)
		_node.BootstrapVersion = value
	}
	if value, ok := svc.mutation.EngineVersion(); ok {
		_spec.SetField(smlversion.FieldEngineVersion, field.TypeString, value)
		_node.EngineVersion = value
	}
	if nodes := svc.mutation.TargetsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   smlversion.TargetsTable,
			Columns: []string{smlversion.TargetsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(smlversiontarget.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.SmlVersion.Create().
//		SetCreatedAt(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.SmlVersionUpsert) {
//			SetCreatedAt(v+v).
//		}).
//		Exec(ctx)
func (svc *SmlVersionCreate) OnConflict(opts ...sql.ConflictOption) *SmlVersionUpsertOne {
	svc.conflict = opts
	return &SmlVersionUpsertOne{
		create: svc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.SmlVersion.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (svc *SmlVersionCreate) OnConflictColumns(columns ...string) *SmlVersionUpsertOne {
	svc.conflict = append(svc.conflict, sql.ConflictColumns(columns...))
	return &SmlVersionUpsertOne{
		create: svc,
	}
}

type (
	// SmlVersionUpsertOne is the builder for "upsert"-ing
	//  one SmlVersion node.
	SmlVersionUpsertOne struct {
		create *SmlVersionCreate
	}

	// SmlVersionUpsert is the "OnConflict" setter.
	SmlVersionUpsert struct {
		*sql.UpdateSet
	}
)

// SetUpdatedAt sets the "updated_at" field.
func (u *SmlVersionUpsert) SetUpdatedAt(v time.Time) *SmlVersionUpsert {
	u.Set(smlversion.FieldUpdatedAt, v)
	return u
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *SmlVersionUpsert) UpdateUpdatedAt() *SmlVersionUpsert {
	u.SetExcluded(smlversion.FieldUpdatedAt)
	return u
}

// ClearUpdatedAt clears the value of the "updated_at" field.
func (u *SmlVersionUpsert) ClearUpdatedAt() *SmlVersionUpsert {
	u.SetNull(smlversion.FieldUpdatedAt)
	return u
}

// SetDeletedAt sets the "deleted_at" field.
func (u *SmlVersionUpsert) SetDeletedAt(v time.Time) *SmlVersionUpsert {
	u.Set(smlversion.FieldDeletedAt, v)
	return u
}

// UpdateDeletedAt sets the "deleted_at" field to the value that was provided on create.
func (u *SmlVersionUpsert) UpdateDeletedAt() *SmlVersionUpsert {
	u.SetExcluded(smlversion.FieldDeletedAt)
	return u
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (u *SmlVersionUpsert) ClearDeletedAt() *SmlVersionUpsert {
	u.SetNull(smlversion.FieldDeletedAt)
	return u
}

// SetVersion sets the "version" field.
func (u *SmlVersionUpsert) SetVersion(v string) *SmlVersionUpsert {
	u.Set(smlversion.FieldVersion, v)
	return u
}

// UpdateVersion sets the "version" field to the value that was provided on create.
func (u *SmlVersionUpsert) UpdateVersion() *SmlVersionUpsert {
	u.SetExcluded(smlversion.FieldVersion)
	return u
}

// SetSatisfactoryVersion sets the "satisfactory_version" field.
func (u *SmlVersionUpsert) SetSatisfactoryVersion(v int) *SmlVersionUpsert {
	u.Set(smlversion.FieldSatisfactoryVersion, v)
	return u
}

// UpdateSatisfactoryVersion sets the "satisfactory_version" field to the value that was provided on create.
func (u *SmlVersionUpsert) UpdateSatisfactoryVersion() *SmlVersionUpsert {
	u.SetExcluded(smlversion.FieldSatisfactoryVersion)
	return u
}

// AddSatisfactoryVersion adds v to the "satisfactory_version" field.
func (u *SmlVersionUpsert) AddSatisfactoryVersion(v int) *SmlVersionUpsert {
	u.Add(smlversion.FieldSatisfactoryVersion, v)
	return u
}

// SetStability sets the "stability" field.
func (u *SmlVersionUpsert) SetStability(v util.Stability) *SmlVersionUpsert {
	u.Set(smlversion.FieldStability, v)
	return u
}

// UpdateStability sets the "stability" field to the value that was provided on create.
func (u *SmlVersionUpsert) UpdateStability() *SmlVersionUpsert {
	u.SetExcluded(smlversion.FieldStability)
	return u
}

// SetDate sets the "date" field.
func (u *SmlVersionUpsert) SetDate(v time.Time) *SmlVersionUpsert {
	u.Set(smlversion.FieldDate, v)
	return u
}

// UpdateDate sets the "date" field to the value that was provided on create.
func (u *SmlVersionUpsert) UpdateDate() *SmlVersionUpsert {
	u.SetExcluded(smlversion.FieldDate)
	return u
}

// SetLink sets the "link" field.
func (u *SmlVersionUpsert) SetLink(v string) *SmlVersionUpsert {
	u.Set(smlversion.FieldLink, v)
	return u
}

// UpdateLink sets the "link" field to the value that was provided on create.
func (u *SmlVersionUpsert) UpdateLink() *SmlVersionUpsert {
	u.SetExcluded(smlversion.FieldLink)
	return u
}

// SetChangelog sets the "changelog" field.
func (u *SmlVersionUpsert) SetChangelog(v string) *SmlVersionUpsert {
	u.Set(smlversion.FieldChangelog, v)
	return u
}

// UpdateChangelog sets the "changelog" field to the value that was provided on create.
func (u *SmlVersionUpsert) UpdateChangelog() *SmlVersionUpsert {
	u.SetExcluded(smlversion.FieldChangelog)
	return u
}

// SetBootstrapVersion sets the "bootstrap_version" field.
func (u *SmlVersionUpsert) SetBootstrapVersion(v string) *SmlVersionUpsert {
	u.Set(smlversion.FieldBootstrapVersion, v)
	return u
}

// UpdateBootstrapVersion sets the "bootstrap_version" field to the value that was provided on create.
func (u *SmlVersionUpsert) UpdateBootstrapVersion() *SmlVersionUpsert {
	u.SetExcluded(smlversion.FieldBootstrapVersion)
	return u
}

// ClearBootstrapVersion clears the value of the "bootstrap_version" field.
func (u *SmlVersionUpsert) ClearBootstrapVersion() *SmlVersionUpsert {
	u.SetNull(smlversion.FieldBootstrapVersion)
	return u
}

// SetEngineVersion sets the "engine_version" field.
func (u *SmlVersionUpsert) SetEngineVersion(v string) *SmlVersionUpsert {
	u.Set(smlversion.FieldEngineVersion, v)
	return u
}

// UpdateEngineVersion sets the "engine_version" field to the value that was provided on create.
func (u *SmlVersionUpsert) UpdateEngineVersion() *SmlVersionUpsert {
	u.SetExcluded(smlversion.FieldEngineVersion)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create except the ID field.
// Using this option is equivalent to using:
//
//	client.SmlVersion.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(smlversion.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *SmlVersionUpsertOne) UpdateNewValues() *SmlVersionUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.ID(); exists {
			s.SetIgnore(smlversion.FieldID)
		}
		if _, exists := u.create.mutation.CreatedAt(); exists {
			s.SetIgnore(smlversion.FieldCreatedAt)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.SmlVersion.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *SmlVersionUpsertOne) Ignore() *SmlVersionUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *SmlVersionUpsertOne) DoNothing() *SmlVersionUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the SmlVersionCreate.OnConflict
// documentation for more info.
func (u *SmlVersionUpsertOne) Update(set func(*SmlVersionUpsert)) *SmlVersionUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&SmlVersionUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdatedAt sets the "updated_at" field.
func (u *SmlVersionUpsertOne) SetUpdatedAt(v time.Time) *SmlVersionUpsertOne {
	return u.Update(func(s *SmlVersionUpsert) {
		s.SetUpdatedAt(v)
	})
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *SmlVersionUpsertOne) UpdateUpdatedAt() *SmlVersionUpsertOne {
	return u.Update(func(s *SmlVersionUpsert) {
		s.UpdateUpdatedAt()
	})
}

// ClearUpdatedAt clears the value of the "updated_at" field.
func (u *SmlVersionUpsertOne) ClearUpdatedAt() *SmlVersionUpsertOne {
	return u.Update(func(s *SmlVersionUpsert) {
		s.ClearUpdatedAt()
	})
}

// SetDeletedAt sets the "deleted_at" field.
func (u *SmlVersionUpsertOne) SetDeletedAt(v time.Time) *SmlVersionUpsertOne {
	return u.Update(func(s *SmlVersionUpsert) {
		s.SetDeletedAt(v)
	})
}

// UpdateDeletedAt sets the "deleted_at" field to the value that was provided on create.
func (u *SmlVersionUpsertOne) UpdateDeletedAt() *SmlVersionUpsertOne {
	return u.Update(func(s *SmlVersionUpsert) {
		s.UpdateDeletedAt()
	})
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (u *SmlVersionUpsertOne) ClearDeletedAt() *SmlVersionUpsertOne {
	return u.Update(func(s *SmlVersionUpsert) {
		s.ClearDeletedAt()
	})
}

// SetVersion sets the "version" field.
func (u *SmlVersionUpsertOne) SetVersion(v string) *SmlVersionUpsertOne {
	return u.Update(func(s *SmlVersionUpsert) {
		s.SetVersion(v)
	})
}

// UpdateVersion sets the "version" field to the value that was provided on create.
func (u *SmlVersionUpsertOne) UpdateVersion() *SmlVersionUpsertOne {
	return u.Update(func(s *SmlVersionUpsert) {
		s.UpdateVersion()
	})
}

// SetSatisfactoryVersion sets the "satisfactory_version" field.
func (u *SmlVersionUpsertOne) SetSatisfactoryVersion(v int) *SmlVersionUpsertOne {
	return u.Update(func(s *SmlVersionUpsert) {
		s.SetSatisfactoryVersion(v)
	})
}

// AddSatisfactoryVersion adds v to the "satisfactory_version" field.
func (u *SmlVersionUpsertOne) AddSatisfactoryVersion(v int) *SmlVersionUpsertOne {
	return u.Update(func(s *SmlVersionUpsert) {
		s.AddSatisfactoryVersion(v)
	})
}

// UpdateSatisfactoryVersion sets the "satisfactory_version" field to the value that was provided on create.
func (u *SmlVersionUpsertOne) UpdateSatisfactoryVersion() *SmlVersionUpsertOne {
	return u.Update(func(s *SmlVersionUpsert) {
		s.UpdateSatisfactoryVersion()
	})
}

// SetStability sets the "stability" field.
func (u *SmlVersionUpsertOne) SetStability(v util.Stability) *SmlVersionUpsertOne {
	return u.Update(func(s *SmlVersionUpsert) {
		s.SetStability(v)
	})
}

// UpdateStability sets the "stability" field to the value that was provided on create.
func (u *SmlVersionUpsertOne) UpdateStability() *SmlVersionUpsertOne {
	return u.Update(func(s *SmlVersionUpsert) {
		s.UpdateStability()
	})
}

// SetDate sets the "date" field.
func (u *SmlVersionUpsertOne) SetDate(v time.Time) *SmlVersionUpsertOne {
	return u.Update(func(s *SmlVersionUpsert) {
		s.SetDate(v)
	})
}

// UpdateDate sets the "date" field to the value that was provided on create.
func (u *SmlVersionUpsertOne) UpdateDate() *SmlVersionUpsertOne {
	return u.Update(func(s *SmlVersionUpsert) {
		s.UpdateDate()
	})
}

// SetLink sets the "link" field.
func (u *SmlVersionUpsertOne) SetLink(v string) *SmlVersionUpsertOne {
	return u.Update(func(s *SmlVersionUpsert) {
		s.SetLink(v)
	})
}

// UpdateLink sets the "link" field to the value that was provided on create.
func (u *SmlVersionUpsertOne) UpdateLink() *SmlVersionUpsertOne {
	return u.Update(func(s *SmlVersionUpsert) {
		s.UpdateLink()
	})
}

// SetChangelog sets the "changelog" field.
func (u *SmlVersionUpsertOne) SetChangelog(v string) *SmlVersionUpsertOne {
	return u.Update(func(s *SmlVersionUpsert) {
		s.SetChangelog(v)
	})
}

// UpdateChangelog sets the "changelog" field to the value that was provided on create.
func (u *SmlVersionUpsertOne) UpdateChangelog() *SmlVersionUpsertOne {
	return u.Update(func(s *SmlVersionUpsert) {
		s.UpdateChangelog()
	})
}

// SetBootstrapVersion sets the "bootstrap_version" field.
func (u *SmlVersionUpsertOne) SetBootstrapVersion(v string) *SmlVersionUpsertOne {
	return u.Update(func(s *SmlVersionUpsert) {
		s.SetBootstrapVersion(v)
	})
}

// UpdateBootstrapVersion sets the "bootstrap_version" field to the value that was provided on create.
func (u *SmlVersionUpsertOne) UpdateBootstrapVersion() *SmlVersionUpsertOne {
	return u.Update(func(s *SmlVersionUpsert) {
		s.UpdateBootstrapVersion()
	})
}

// ClearBootstrapVersion clears the value of the "bootstrap_version" field.
func (u *SmlVersionUpsertOne) ClearBootstrapVersion() *SmlVersionUpsertOne {
	return u.Update(func(s *SmlVersionUpsert) {
		s.ClearBootstrapVersion()
	})
}

// SetEngineVersion sets the "engine_version" field.
func (u *SmlVersionUpsertOne) SetEngineVersion(v string) *SmlVersionUpsertOne {
	return u.Update(func(s *SmlVersionUpsert) {
		s.SetEngineVersion(v)
	})
}

// UpdateEngineVersion sets the "engine_version" field to the value that was provided on create.
func (u *SmlVersionUpsertOne) UpdateEngineVersion() *SmlVersionUpsertOne {
	return u.Update(func(s *SmlVersionUpsert) {
		s.UpdateEngineVersion()
	})
}

// Exec executes the query.
func (u *SmlVersionUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for SmlVersionCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *SmlVersionUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *SmlVersionUpsertOne) ID(ctx context.Context) (id string, err error) {
	if u.create.driver.Dialect() == dialect.MySQL {
		// In case of "ON CONFLICT", there is no way to get back non-numeric ID
		// fields from the database since MySQL does not support the RETURNING clause.
		return id, errors.New("ent: SmlVersionUpsertOne.ID is not supported by MySQL driver. Use SmlVersionUpsertOne.Exec instead")
	}
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *SmlVersionUpsertOne) IDX(ctx context.Context) string {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// SmlVersionCreateBulk is the builder for creating many SmlVersion entities in bulk.
type SmlVersionCreateBulk struct {
	config
	err      error
	builders []*SmlVersionCreate
	conflict []sql.ConflictOption
}

// Save creates the SmlVersion entities in the database.
func (svcb *SmlVersionCreateBulk) Save(ctx context.Context) ([]*SmlVersion, error) {
	if svcb.err != nil {
		return nil, svcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(svcb.builders))
	nodes := make([]*SmlVersion, len(svcb.builders))
	mutators := make([]Mutator, len(svcb.builders))
	for i := range svcb.builders {
		func(i int, root context.Context) {
			builder := svcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*SmlVersionMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, svcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = svcb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, svcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, svcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (svcb *SmlVersionCreateBulk) SaveX(ctx context.Context) []*SmlVersion {
	v, err := svcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (svcb *SmlVersionCreateBulk) Exec(ctx context.Context) error {
	_, err := svcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (svcb *SmlVersionCreateBulk) ExecX(ctx context.Context) {
	if err := svcb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.SmlVersion.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.SmlVersionUpsert) {
//			SetCreatedAt(v+v).
//		}).
//		Exec(ctx)
func (svcb *SmlVersionCreateBulk) OnConflict(opts ...sql.ConflictOption) *SmlVersionUpsertBulk {
	svcb.conflict = opts
	return &SmlVersionUpsertBulk{
		create: svcb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.SmlVersion.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (svcb *SmlVersionCreateBulk) OnConflictColumns(columns ...string) *SmlVersionUpsertBulk {
	svcb.conflict = append(svcb.conflict, sql.ConflictColumns(columns...))
	return &SmlVersionUpsertBulk{
		create: svcb,
	}
}

// SmlVersionUpsertBulk is the builder for "upsert"-ing
// a bulk of SmlVersion nodes.
type SmlVersionUpsertBulk struct {
	create *SmlVersionCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.SmlVersion.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(smlversion.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *SmlVersionUpsertBulk) UpdateNewValues() *SmlVersionUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.ID(); exists {
				s.SetIgnore(smlversion.FieldID)
			}
			if _, exists := b.mutation.CreatedAt(); exists {
				s.SetIgnore(smlversion.FieldCreatedAt)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.SmlVersion.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *SmlVersionUpsertBulk) Ignore() *SmlVersionUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *SmlVersionUpsertBulk) DoNothing() *SmlVersionUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the SmlVersionCreateBulk.OnConflict
// documentation for more info.
func (u *SmlVersionUpsertBulk) Update(set func(*SmlVersionUpsert)) *SmlVersionUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&SmlVersionUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdatedAt sets the "updated_at" field.
func (u *SmlVersionUpsertBulk) SetUpdatedAt(v time.Time) *SmlVersionUpsertBulk {
	return u.Update(func(s *SmlVersionUpsert) {
		s.SetUpdatedAt(v)
	})
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *SmlVersionUpsertBulk) UpdateUpdatedAt() *SmlVersionUpsertBulk {
	return u.Update(func(s *SmlVersionUpsert) {
		s.UpdateUpdatedAt()
	})
}

// ClearUpdatedAt clears the value of the "updated_at" field.
func (u *SmlVersionUpsertBulk) ClearUpdatedAt() *SmlVersionUpsertBulk {
	return u.Update(func(s *SmlVersionUpsert) {
		s.ClearUpdatedAt()
	})
}

// SetDeletedAt sets the "deleted_at" field.
func (u *SmlVersionUpsertBulk) SetDeletedAt(v time.Time) *SmlVersionUpsertBulk {
	return u.Update(func(s *SmlVersionUpsert) {
		s.SetDeletedAt(v)
	})
}

// UpdateDeletedAt sets the "deleted_at" field to the value that was provided on create.
func (u *SmlVersionUpsertBulk) UpdateDeletedAt() *SmlVersionUpsertBulk {
	return u.Update(func(s *SmlVersionUpsert) {
		s.UpdateDeletedAt()
	})
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (u *SmlVersionUpsertBulk) ClearDeletedAt() *SmlVersionUpsertBulk {
	return u.Update(func(s *SmlVersionUpsert) {
		s.ClearDeletedAt()
	})
}

// SetVersion sets the "version" field.
func (u *SmlVersionUpsertBulk) SetVersion(v string) *SmlVersionUpsertBulk {
	return u.Update(func(s *SmlVersionUpsert) {
		s.SetVersion(v)
	})
}

// UpdateVersion sets the "version" field to the value that was provided on create.
func (u *SmlVersionUpsertBulk) UpdateVersion() *SmlVersionUpsertBulk {
	return u.Update(func(s *SmlVersionUpsert) {
		s.UpdateVersion()
	})
}

// SetSatisfactoryVersion sets the "satisfactory_version" field.
func (u *SmlVersionUpsertBulk) SetSatisfactoryVersion(v int) *SmlVersionUpsertBulk {
	return u.Update(func(s *SmlVersionUpsert) {
		s.SetSatisfactoryVersion(v)
	})
}

// AddSatisfactoryVersion adds v to the "satisfactory_version" field.
func (u *SmlVersionUpsertBulk) AddSatisfactoryVersion(v int) *SmlVersionUpsertBulk {
	return u.Update(func(s *SmlVersionUpsert) {
		s.AddSatisfactoryVersion(v)
	})
}

// UpdateSatisfactoryVersion sets the "satisfactory_version" field to the value that was provided on create.
func (u *SmlVersionUpsertBulk) UpdateSatisfactoryVersion() *SmlVersionUpsertBulk {
	return u.Update(func(s *SmlVersionUpsert) {
		s.UpdateSatisfactoryVersion()
	})
}

// SetStability sets the "stability" field.
func (u *SmlVersionUpsertBulk) SetStability(v util.Stability) *SmlVersionUpsertBulk {
	return u.Update(func(s *SmlVersionUpsert) {
		s.SetStability(v)
	})
}

// UpdateStability sets the "stability" field to the value that was provided on create.
func (u *SmlVersionUpsertBulk) UpdateStability() *SmlVersionUpsertBulk {
	return u.Update(func(s *SmlVersionUpsert) {
		s.UpdateStability()
	})
}

// SetDate sets the "date" field.
func (u *SmlVersionUpsertBulk) SetDate(v time.Time) *SmlVersionUpsertBulk {
	return u.Update(func(s *SmlVersionUpsert) {
		s.SetDate(v)
	})
}

// UpdateDate sets the "date" field to the value that was provided on create.
func (u *SmlVersionUpsertBulk) UpdateDate() *SmlVersionUpsertBulk {
	return u.Update(func(s *SmlVersionUpsert) {
		s.UpdateDate()
	})
}

// SetLink sets the "link" field.
func (u *SmlVersionUpsertBulk) SetLink(v string) *SmlVersionUpsertBulk {
	return u.Update(func(s *SmlVersionUpsert) {
		s.SetLink(v)
	})
}

// UpdateLink sets the "link" field to the value that was provided on create.
func (u *SmlVersionUpsertBulk) UpdateLink() *SmlVersionUpsertBulk {
	return u.Update(func(s *SmlVersionUpsert) {
		s.UpdateLink()
	})
}

// SetChangelog sets the "changelog" field.
func (u *SmlVersionUpsertBulk) SetChangelog(v string) *SmlVersionUpsertBulk {
	return u.Update(func(s *SmlVersionUpsert) {
		s.SetChangelog(v)
	})
}

// UpdateChangelog sets the "changelog" field to the value that was provided on create.
func (u *SmlVersionUpsertBulk) UpdateChangelog() *SmlVersionUpsertBulk {
	return u.Update(func(s *SmlVersionUpsert) {
		s.UpdateChangelog()
	})
}

// SetBootstrapVersion sets the "bootstrap_version" field.
func (u *SmlVersionUpsertBulk) SetBootstrapVersion(v string) *SmlVersionUpsertBulk {
	return u.Update(func(s *SmlVersionUpsert) {
		s.SetBootstrapVersion(v)
	})
}

// UpdateBootstrapVersion sets the "bootstrap_version" field to the value that was provided on create.
func (u *SmlVersionUpsertBulk) UpdateBootstrapVersion() *SmlVersionUpsertBulk {
	return u.Update(func(s *SmlVersionUpsert) {
		s.UpdateBootstrapVersion()
	})
}

// ClearBootstrapVersion clears the value of the "bootstrap_version" field.
func (u *SmlVersionUpsertBulk) ClearBootstrapVersion() *SmlVersionUpsertBulk {
	return u.Update(func(s *SmlVersionUpsert) {
		s.ClearBootstrapVersion()
	})
}

// SetEngineVersion sets the "engine_version" field.
func (u *SmlVersionUpsertBulk) SetEngineVersion(v string) *SmlVersionUpsertBulk {
	return u.Update(func(s *SmlVersionUpsert) {
		s.SetEngineVersion(v)
	})
}

// UpdateEngineVersion sets the "engine_version" field to the value that was provided on create.
func (u *SmlVersionUpsertBulk) UpdateEngineVersion() *SmlVersionUpsertBulk {
	return u.Update(func(s *SmlVersionUpsert) {
		s.UpdateEngineVersion()
	})
}

// Exec executes the query.
func (u *SmlVersionUpsertBulk) Exec(ctx context.Context) error {
	if u.create.err != nil {
		return u.create.err
	}
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the SmlVersionCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for SmlVersionCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *SmlVersionUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
