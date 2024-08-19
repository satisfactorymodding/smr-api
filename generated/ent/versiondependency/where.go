// Code generated by ent, DO NOT EDIT.

package versiondependency

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/satisfactorymodding/smr-api/generated/ent/predicate"
)

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldEQ(FieldCreatedAt, v))
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldEQ(FieldUpdatedAt, v))
}

// DeletedAt applies equality check predicate on the "deleted_at" field. It's identical to DeletedAtEQ.
func DeletedAt(v time.Time) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldEQ(FieldDeletedAt, v))
}

// VersionID applies equality check predicate on the "version_id" field. It's identical to VersionIDEQ.
func VersionID(v string) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldEQ(FieldVersionID, v))
}

// ModID applies equality check predicate on the "mod_id" field. It's identical to ModIDEQ.
func ModID(v string) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldEQ(FieldModID, v))
}

// Condition applies equality check predicate on the "condition" field. It's identical to ConditionEQ.
func Condition(v string) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldEQ(FieldCondition, v))
}

// Optional applies equality check predicate on the "optional" field. It's identical to OptionalEQ.
func Optional(v bool) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldEQ(FieldOptional, v))
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldLTE(FieldCreatedAt, v))
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldEQ(FieldUpdatedAt, v))
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldNEQ(FieldUpdatedAt, v))
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldIn(FieldUpdatedAt, vs...))
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldNotIn(FieldUpdatedAt, vs...))
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v time.Time) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldGT(FieldUpdatedAt, v))
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldGTE(FieldUpdatedAt, v))
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldLT(FieldUpdatedAt, v))
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldLTE(FieldUpdatedAt, v))
}

// DeletedAtEQ applies the EQ predicate on the "deleted_at" field.
func DeletedAtEQ(v time.Time) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldEQ(FieldDeletedAt, v))
}

// DeletedAtNEQ applies the NEQ predicate on the "deleted_at" field.
func DeletedAtNEQ(v time.Time) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldNEQ(FieldDeletedAt, v))
}

// DeletedAtIn applies the In predicate on the "deleted_at" field.
func DeletedAtIn(vs ...time.Time) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldIn(FieldDeletedAt, vs...))
}

// DeletedAtNotIn applies the NotIn predicate on the "deleted_at" field.
func DeletedAtNotIn(vs ...time.Time) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldNotIn(FieldDeletedAt, vs...))
}

// DeletedAtGT applies the GT predicate on the "deleted_at" field.
func DeletedAtGT(v time.Time) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldGT(FieldDeletedAt, v))
}

// DeletedAtGTE applies the GTE predicate on the "deleted_at" field.
func DeletedAtGTE(v time.Time) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldGTE(FieldDeletedAt, v))
}

// DeletedAtLT applies the LT predicate on the "deleted_at" field.
func DeletedAtLT(v time.Time) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldLT(FieldDeletedAt, v))
}

// DeletedAtLTE applies the LTE predicate on the "deleted_at" field.
func DeletedAtLTE(v time.Time) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldLTE(FieldDeletedAt, v))
}

// DeletedAtIsNil applies the IsNil predicate on the "deleted_at" field.
func DeletedAtIsNil() predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldIsNull(FieldDeletedAt))
}

// DeletedAtNotNil applies the NotNil predicate on the "deleted_at" field.
func DeletedAtNotNil() predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldNotNull(FieldDeletedAt))
}

// VersionIDEQ applies the EQ predicate on the "version_id" field.
func VersionIDEQ(v string) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldEQ(FieldVersionID, v))
}

// VersionIDNEQ applies the NEQ predicate on the "version_id" field.
func VersionIDNEQ(v string) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldNEQ(FieldVersionID, v))
}

// VersionIDIn applies the In predicate on the "version_id" field.
func VersionIDIn(vs ...string) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldIn(FieldVersionID, vs...))
}

// VersionIDNotIn applies the NotIn predicate on the "version_id" field.
func VersionIDNotIn(vs ...string) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldNotIn(FieldVersionID, vs...))
}

// VersionIDGT applies the GT predicate on the "version_id" field.
func VersionIDGT(v string) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldGT(FieldVersionID, v))
}

// VersionIDGTE applies the GTE predicate on the "version_id" field.
func VersionIDGTE(v string) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldGTE(FieldVersionID, v))
}

// VersionIDLT applies the LT predicate on the "version_id" field.
func VersionIDLT(v string) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldLT(FieldVersionID, v))
}

// VersionIDLTE applies the LTE predicate on the "version_id" field.
func VersionIDLTE(v string) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldLTE(FieldVersionID, v))
}

// VersionIDContains applies the Contains predicate on the "version_id" field.
func VersionIDContains(v string) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldContains(FieldVersionID, v))
}

// VersionIDHasPrefix applies the HasPrefix predicate on the "version_id" field.
func VersionIDHasPrefix(v string) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldHasPrefix(FieldVersionID, v))
}

// VersionIDHasSuffix applies the HasSuffix predicate on the "version_id" field.
func VersionIDHasSuffix(v string) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldHasSuffix(FieldVersionID, v))
}

// VersionIDEqualFold applies the EqualFold predicate on the "version_id" field.
func VersionIDEqualFold(v string) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldEqualFold(FieldVersionID, v))
}

// VersionIDContainsFold applies the ContainsFold predicate on the "version_id" field.
func VersionIDContainsFold(v string) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldContainsFold(FieldVersionID, v))
}

// ModIDEQ applies the EQ predicate on the "mod_id" field.
func ModIDEQ(v string) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldEQ(FieldModID, v))
}

// ModIDNEQ applies the NEQ predicate on the "mod_id" field.
func ModIDNEQ(v string) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldNEQ(FieldModID, v))
}

// ModIDIn applies the In predicate on the "mod_id" field.
func ModIDIn(vs ...string) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldIn(FieldModID, vs...))
}

// ModIDNotIn applies the NotIn predicate on the "mod_id" field.
func ModIDNotIn(vs ...string) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldNotIn(FieldModID, vs...))
}

// ModIDGT applies the GT predicate on the "mod_id" field.
func ModIDGT(v string) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldGT(FieldModID, v))
}

// ModIDGTE applies the GTE predicate on the "mod_id" field.
func ModIDGTE(v string) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldGTE(FieldModID, v))
}

// ModIDLT applies the LT predicate on the "mod_id" field.
func ModIDLT(v string) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldLT(FieldModID, v))
}

// ModIDLTE applies the LTE predicate on the "mod_id" field.
func ModIDLTE(v string) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldLTE(FieldModID, v))
}

// ModIDContains applies the Contains predicate on the "mod_id" field.
func ModIDContains(v string) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldContains(FieldModID, v))
}

// ModIDHasPrefix applies the HasPrefix predicate on the "mod_id" field.
func ModIDHasPrefix(v string) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldHasPrefix(FieldModID, v))
}

// ModIDHasSuffix applies the HasSuffix predicate on the "mod_id" field.
func ModIDHasSuffix(v string) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldHasSuffix(FieldModID, v))
}

// ModIDEqualFold applies the EqualFold predicate on the "mod_id" field.
func ModIDEqualFold(v string) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldEqualFold(FieldModID, v))
}

// ModIDContainsFold applies the ContainsFold predicate on the "mod_id" field.
func ModIDContainsFold(v string) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldContainsFold(FieldModID, v))
}

// ConditionEQ applies the EQ predicate on the "condition" field.
func ConditionEQ(v string) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldEQ(FieldCondition, v))
}

// ConditionNEQ applies the NEQ predicate on the "condition" field.
func ConditionNEQ(v string) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldNEQ(FieldCondition, v))
}

// ConditionIn applies the In predicate on the "condition" field.
func ConditionIn(vs ...string) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldIn(FieldCondition, vs...))
}

// ConditionNotIn applies the NotIn predicate on the "condition" field.
func ConditionNotIn(vs ...string) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldNotIn(FieldCondition, vs...))
}

// ConditionGT applies the GT predicate on the "condition" field.
func ConditionGT(v string) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldGT(FieldCondition, v))
}

// ConditionGTE applies the GTE predicate on the "condition" field.
func ConditionGTE(v string) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldGTE(FieldCondition, v))
}

// ConditionLT applies the LT predicate on the "condition" field.
func ConditionLT(v string) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldLT(FieldCondition, v))
}

// ConditionLTE applies the LTE predicate on the "condition" field.
func ConditionLTE(v string) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldLTE(FieldCondition, v))
}

// ConditionContains applies the Contains predicate on the "condition" field.
func ConditionContains(v string) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldContains(FieldCondition, v))
}

// ConditionHasPrefix applies the HasPrefix predicate on the "condition" field.
func ConditionHasPrefix(v string) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldHasPrefix(FieldCondition, v))
}

// ConditionHasSuffix applies the HasSuffix predicate on the "condition" field.
func ConditionHasSuffix(v string) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldHasSuffix(FieldCondition, v))
}

// ConditionEqualFold applies the EqualFold predicate on the "condition" field.
func ConditionEqualFold(v string) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldEqualFold(FieldCondition, v))
}

// ConditionContainsFold applies the ContainsFold predicate on the "condition" field.
func ConditionContainsFold(v string) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldContainsFold(FieldCondition, v))
}

// OptionalEQ applies the EQ predicate on the "optional" field.
func OptionalEQ(v bool) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldEQ(FieldOptional, v))
}

// OptionalNEQ applies the NEQ predicate on the "optional" field.
func OptionalNEQ(v bool) predicate.VersionDependency {
	return predicate.VersionDependency(sql.FieldNEQ(FieldOptional, v))
}

// HasVersion applies the HasEdge predicate on the "version" edge.
func HasVersion() predicate.VersionDependency {
	return predicate.VersionDependency(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, VersionColumn),
			sqlgraph.Edge(sqlgraph.M2O, false, VersionTable, VersionColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasVersionWith applies the HasEdge predicate on the "version" edge with a given conditions (other predicates).
func HasVersionWith(preds ...predicate.Version) predicate.VersionDependency {
	return predicate.VersionDependency(func(s *sql.Selector) {
		step := newVersionStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasMod applies the HasEdge predicate on the "mod" edge.
func HasMod() predicate.VersionDependency {
	return predicate.VersionDependency(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, ModColumn),
			sqlgraph.Edge(sqlgraph.M2O, false, ModTable, ModColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasModWith applies the HasEdge predicate on the "mod" edge with a given conditions (other predicates).
func HasModWith(preds ...predicate.Mod) predicate.VersionDependency {
	return predicate.VersionDependency(func(s *sql.Selector) {
		step := newModStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.VersionDependency) predicate.VersionDependency {
	return predicate.VersionDependency(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.VersionDependency) predicate.VersionDependency {
	return predicate.VersionDependency(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.VersionDependency) predicate.VersionDependency {
	return predicate.VersionDependency(sql.NotPredicates(p))
}
