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
	"github.com/satisfactorymodding/smr-api/generated/ent/tag"
)

// ModTagCreate is the builder for creating a ModTag entity.
type ModTagCreate struct {
	config
	mutation *ModTagMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetModID sets the "mod_id" field.
func (mtc *ModTagCreate) SetModID(s string) *ModTagCreate {
	mtc.mutation.SetModID(s)
	return mtc
}

// SetTagID sets the "tag_id" field.
func (mtc *ModTagCreate) SetTagID(s string) *ModTagCreate {
	mtc.mutation.SetTagID(s)
	return mtc
}

// SetMod sets the "mod" edge to the Mod entity.
func (mtc *ModTagCreate) SetMod(m *Mod) *ModTagCreate {
	return mtc.SetModID(m.ID)
}

// SetTag sets the "tag" edge to the Tag entity.
func (mtc *ModTagCreate) SetTag(t *Tag) *ModTagCreate {
	return mtc.SetTagID(t.ID)
}

// Mutation returns the ModTagMutation object of the builder.
func (mtc *ModTagCreate) Mutation() *ModTagMutation {
	return mtc.mutation
}

// Save creates the ModTag in the database.
func (mtc *ModTagCreate) Save(ctx context.Context) (*ModTag, error) {
	return withHooks(ctx, mtc.sqlSave, mtc.mutation, mtc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (mtc *ModTagCreate) SaveX(ctx context.Context) *ModTag {
	v, err := mtc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (mtc *ModTagCreate) Exec(ctx context.Context) error {
	_, err := mtc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (mtc *ModTagCreate) ExecX(ctx context.Context) {
	if err := mtc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (mtc *ModTagCreate) check() error {
	if _, ok := mtc.mutation.ModID(); !ok {
		return &ValidationError{Name: "mod_id", err: errors.New(`ent: missing required field "ModTag.mod_id"`)}
	}
	if _, ok := mtc.mutation.TagID(); !ok {
		return &ValidationError{Name: "tag_id", err: errors.New(`ent: missing required field "ModTag.tag_id"`)}
	}
	if _, ok := mtc.mutation.ModID(); !ok {
		return &ValidationError{Name: "mod", err: errors.New(`ent: missing required edge "ModTag.mod"`)}
	}
	if _, ok := mtc.mutation.TagID(); !ok {
		return &ValidationError{Name: "tag", err: errors.New(`ent: missing required edge "ModTag.tag"`)}
	}
	return nil
}

func (mtc *ModTagCreate) sqlSave(ctx context.Context) (*ModTag, error) {
	if err := mtc.check(); err != nil {
		return nil, err
	}
	_node, _spec := mtc.createSpec()
	if err := sqlgraph.CreateNode(ctx, mtc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	return _node, nil
}

func (mtc *ModTagCreate) createSpec() (*ModTag, *sqlgraph.CreateSpec) {
	var (
		_node = &ModTag{config: mtc.config}
		_spec = sqlgraph.NewCreateSpec(modtag.Table, nil)
	)
	_spec.OnConflict = mtc.conflict
	if nodes := mtc.mutation.ModIDs(); len(nodes) > 0 {
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
		_node.ModID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := mtc.mutation.TagIDs(); len(nodes) > 0 {
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
		_node.TagID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.ModTag.Create().
//		SetModID(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.ModTagUpsert) {
//			SetModID(v+v).
//		}).
//		Exec(ctx)
func (mtc *ModTagCreate) OnConflict(opts ...sql.ConflictOption) *ModTagUpsertOne {
	mtc.conflict = opts
	return &ModTagUpsertOne{
		create: mtc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.ModTag.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (mtc *ModTagCreate) OnConflictColumns(columns ...string) *ModTagUpsertOne {
	mtc.conflict = append(mtc.conflict, sql.ConflictColumns(columns...))
	return &ModTagUpsertOne{
		create: mtc,
	}
}

type (
	// ModTagUpsertOne is the builder for "upsert"-ing
	//  one ModTag node.
	ModTagUpsertOne struct {
		create *ModTagCreate
	}

	// ModTagUpsert is the "OnConflict" setter.
	ModTagUpsert struct {
		*sql.UpdateSet
	}
)

// SetModID sets the "mod_id" field.
func (u *ModTagUpsert) SetModID(v string) *ModTagUpsert {
	u.Set(modtag.FieldModID, v)
	return u
}

// UpdateModID sets the "mod_id" field to the value that was provided on create.
func (u *ModTagUpsert) UpdateModID() *ModTagUpsert {
	u.SetExcluded(modtag.FieldModID)
	return u
}

// SetTagID sets the "tag_id" field.
func (u *ModTagUpsert) SetTagID(v string) *ModTagUpsert {
	u.Set(modtag.FieldTagID, v)
	return u
}

// UpdateTagID sets the "tag_id" field to the value that was provided on create.
func (u *ModTagUpsert) UpdateTagID() *ModTagUpsert {
	u.SetExcluded(modtag.FieldTagID)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create.
// Using this option is equivalent to using:
//
//	client.ModTag.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *ModTagUpsertOne) UpdateNewValues() *ModTagUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.ModTag.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *ModTagUpsertOne) Ignore() *ModTagUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *ModTagUpsertOne) DoNothing() *ModTagUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the ModTagCreate.OnConflict
// documentation for more info.
func (u *ModTagUpsertOne) Update(set func(*ModTagUpsert)) *ModTagUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&ModTagUpsert{UpdateSet: update})
	}))
	return u
}

// SetModID sets the "mod_id" field.
func (u *ModTagUpsertOne) SetModID(v string) *ModTagUpsertOne {
	return u.Update(func(s *ModTagUpsert) {
		s.SetModID(v)
	})
}

// UpdateModID sets the "mod_id" field to the value that was provided on create.
func (u *ModTagUpsertOne) UpdateModID() *ModTagUpsertOne {
	return u.Update(func(s *ModTagUpsert) {
		s.UpdateModID()
	})
}

// SetTagID sets the "tag_id" field.
func (u *ModTagUpsertOne) SetTagID(v string) *ModTagUpsertOne {
	return u.Update(func(s *ModTagUpsert) {
		s.SetTagID(v)
	})
}

// UpdateTagID sets the "tag_id" field to the value that was provided on create.
func (u *ModTagUpsertOne) UpdateTagID() *ModTagUpsertOne {
	return u.Update(func(s *ModTagUpsert) {
		s.UpdateTagID()
	})
}

// Exec executes the query.
func (u *ModTagUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for ModTagCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *ModTagUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// ModTagCreateBulk is the builder for creating many ModTag entities in bulk.
type ModTagCreateBulk struct {
	config
	err      error
	builders []*ModTagCreate
	conflict []sql.ConflictOption
}

// Save creates the ModTag entities in the database.
func (mtcb *ModTagCreateBulk) Save(ctx context.Context) ([]*ModTag, error) {
	if mtcb.err != nil {
		return nil, mtcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(mtcb.builders))
	nodes := make([]*ModTag, len(mtcb.builders))
	mutators := make([]Mutator, len(mtcb.builders))
	for i := range mtcb.builders {
		func(i int, root context.Context) {
			builder := mtcb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*ModTagMutation)
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
					_, err = mutators[i+1].Mutate(root, mtcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = mtcb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, mtcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
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
		if _, err := mutators[0].Mutate(ctx, mtcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (mtcb *ModTagCreateBulk) SaveX(ctx context.Context) []*ModTag {
	v, err := mtcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (mtcb *ModTagCreateBulk) Exec(ctx context.Context) error {
	_, err := mtcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (mtcb *ModTagCreateBulk) ExecX(ctx context.Context) {
	if err := mtcb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.ModTag.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.ModTagUpsert) {
//			SetModID(v+v).
//		}).
//		Exec(ctx)
func (mtcb *ModTagCreateBulk) OnConflict(opts ...sql.ConflictOption) *ModTagUpsertBulk {
	mtcb.conflict = opts
	return &ModTagUpsertBulk{
		create: mtcb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.ModTag.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (mtcb *ModTagCreateBulk) OnConflictColumns(columns ...string) *ModTagUpsertBulk {
	mtcb.conflict = append(mtcb.conflict, sql.ConflictColumns(columns...))
	return &ModTagUpsertBulk{
		create: mtcb,
	}
}

// ModTagUpsertBulk is the builder for "upsert"-ing
// a bulk of ModTag nodes.
type ModTagUpsertBulk struct {
	create *ModTagCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.ModTag.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *ModTagUpsertBulk) UpdateNewValues() *ModTagUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.ModTag.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *ModTagUpsertBulk) Ignore() *ModTagUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *ModTagUpsertBulk) DoNothing() *ModTagUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the ModTagCreateBulk.OnConflict
// documentation for more info.
func (u *ModTagUpsertBulk) Update(set func(*ModTagUpsert)) *ModTagUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&ModTagUpsert{UpdateSet: update})
	}))
	return u
}

// SetModID sets the "mod_id" field.
func (u *ModTagUpsertBulk) SetModID(v string) *ModTagUpsertBulk {
	return u.Update(func(s *ModTagUpsert) {
		s.SetModID(v)
	})
}

// UpdateModID sets the "mod_id" field to the value that was provided on create.
func (u *ModTagUpsertBulk) UpdateModID() *ModTagUpsertBulk {
	return u.Update(func(s *ModTagUpsert) {
		s.UpdateModID()
	})
}

// SetTagID sets the "tag_id" field.
func (u *ModTagUpsertBulk) SetTagID(v string) *ModTagUpsertBulk {
	return u.Update(func(s *ModTagUpsert) {
		s.SetTagID(v)
	})
}

// UpdateTagID sets the "tag_id" field to the value that was provided on create.
func (u *ModTagUpsertBulk) UpdateTagID() *ModTagUpsertBulk {
	return u.Update(func(s *ModTagUpsert) {
		s.UpdateTagID()
	})
}

// Exec executes the query.
func (u *ModTagUpsertBulk) Exec(ctx context.Context) error {
	if u.create.err != nil {
		return u.create.err
	}
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the ModTagCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for ModTagCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *ModTagUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
