// Code generated by ent, DO NOT EDIT.

package usergroup

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/satisfactorymodding/smr-api/generated/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id string) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id string) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id string) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...string) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...string) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id string) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id string) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id string) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id string) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldLTE(FieldID, id))
}

// IDEqualFold applies the EqualFold predicate on the ID field.
func IDEqualFold(id string) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldEqualFold(FieldID, id))
}

// IDContainsFold applies the ContainsFold predicate on the ID field.
func IDContainsFold(id string) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldContainsFold(FieldID, id))
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldEQ(FieldCreatedAt, v))
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldEQ(FieldUpdatedAt, v))
}

// DeletedAt applies equality check predicate on the "deleted_at" field. It's identical to DeletedAtEQ.
func DeletedAt(v time.Time) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldEQ(FieldDeletedAt, v))
}

// UserID applies equality check predicate on the "user_id" field. It's identical to UserIDEQ.
func UserID(v string) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldEQ(FieldUserID, v))
}

// GroupID applies equality check predicate on the "group_id" field. It's identical to GroupIDEQ.
func GroupID(v string) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldEQ(FieldGroupID, v))
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldLTE(FieldCreatedAt, v))
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldEQ(FieldUpdatedAt, v))
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldNEQ(FieldUpdatedAt, v))
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldIn(FieldUpdatedAt, vs...))
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldNotIn(FieldUpdatedAt, vs...))
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v time.Time) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldGT(FieldUpdatedAt, v))
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldGTE(FieldUpdatedAt, v))
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldLT(FieldUpdatedAt, v))
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldLTE(FieldUpdatedAt, v))
}

// DeletedAtEQ applies the EQ predicate on the "deleted_at" field.
func DeletedAtEQ(v time.Time) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldEQ(FieldDeletedAt, v))
}

// DeletedAtNEQ applies the NEQ predicate on the "deleted_at" field.
func DeletedAtNEQ(v time.Time) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldNEQ(FieldDeletedAt, v))
}

// DeletedAtIn applies the In predicate on the "deleted_at" field.
func DeletedAtIn(vs ...time.Time) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldIn(FieldDeletedAt, vs...))
}

// DeletedAtNotIn applies the NotIn predicate on the "deleted_at" field.
func DeletedAtNotIn(vs ...time.Time) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldNotIn(FieldDeletedAt, vs...))
}

// DeletedAtGT applies the GT predicate on the "deleted_at" field.
func DeletedAtGT(v time.Time) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldGT(FieldDeletedAt, v))
}

// DeletedAtGTE applies the GTE predicate on the "deleted_at" field.
func DeletedAtGTE(v time.Time) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldGTE(FieldDeletedAt, v))
}

// DeletedAtLT applies the LT predicate on the "deleted_at" field.
func DeletedAtLT(v time.Time) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldLT(FieldDeletedAt, v))
}

// DeletedAtLTE applies the LTE predicate on the "deleted_at" field.
func DeletedAtLTE(v time.Time) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldLTE(FieldDeletedAt, v))
}

// DeletedAtIsNil applies the IsNil predicate on the "deleted_at" field.
func DeletedAtIsNil() predicate.UserGroup {
	return predicate.UserGroup(sql.FieldIsNull(FieldDeletedAt))
}

// DeletedAtNotNil applies the NotNil predicate on the "deleted_at" field.
func DeletedAtNotNil() predicate.UserGroup {
	return predicate.UserGroup(sql.FieldNotNull(FieldDeletedAt))
}

// UserIDEQ applies the EQ predicate on the "user_id" field.
func UserIDEQ(v string) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldEQ(FieldUserID, v))
}

// UserIDNEQ applies the NEQ predicate on the "user_id" field.
func UserIDNEQ(v string) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldNEQ(FieldUserID, v))
}

// UserIDIn applies the In predicate on the "user_id" field.
func UserIDIn(vs ...string) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldIn(FieldUserID, vs...))
}

// UserIDNotIn applies the NotIn predicate on the "user_id" field.
func UserIDNotIn(vs ...string) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldNotIn(FieldUserID, vs...))
}

// UserIDGT applies the GT predicate on the "user_id" field.
func UserIDGT(v string) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldGT(FieldUserID, v))
}

// UserIDGTE applies the GTE predicate on the "user_id" field.
func UserIDGTE(v string) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldGTE(FieldUserID, v))
}

// UserIDLT applies the LT predicate on the "user_id" field.
func UserIDLT(v string) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldLT(FieldUserID, v))
}

// UserIDLTE applies the LTE predicate on the "user_id" field.
func UserIDLTE(v string) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldLTE(FieldUserID, v))
}

// UserIDContains applies the Contains predicate on the "user_id" field.
func UserIDContains(v string) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldContains(FieldUserID, v))
}

// UserIDHasPrefix applies the HasPrefix predicate on the "user_id" field.
func UserIDHasPrefix(v string) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldHasPrefix(FieldUserID, v))
}

// UserIDHasSuffix applies the HasSuffix predicate on the "user_id" field.
func UserIDHasSuffix(v string) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldHasSuffix(FieldUserID, v))
}

// UserIDEqualFold applies the EqualFold predicate on the "user_id" field.
func UserIDEqualFold(v string) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldEqualFold(FieldUserID, v))
}

// UserIDContainsFold applies the ContainsFold predicate on the "user_id" field.
func UserIDContainsFold(v string) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldContainsFold(FieldUserID, v))
}

// GroupIDEQ applies the EQ predicate on the "group_id" field.
func GroupIDEQ(v string) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldEQ(FieldGroupID, v))
}

// GroupIDNEQ applies the NEQ predicate on the "group_id" field.
func GroupIDNEQ(v string) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldNEQ(FieldGroupID, v))
}

// GroupIDIn applies the In predicate on the "group_id" field.
func GroupIDIn(vs ...string) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldIn(FieldGroupID, vs...))
}

// GroupIDNotIn applies the NotIn predicate on the "group_id" field.
func GroupIDNotIn(vs ...string) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldNotIn(FieldGroupID, vs...))
}

// GroupIDGT applies the GT predicate on the "group_id" field.
func GroupIDGT(v string) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldGT(FieldGroupID, v))
}

// GroupIDGTE applies the GTE predicate on the "group_id" field.
func GroupIDGTE(v string) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldGTE(FieldGroupID, v))
}

// GroupIDLT applies the LT predicate on the "group_id" field.
func GroupIDLT(v string) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldLT(FieldGroupID, v))
}

// GroupIDLTE applies the LTE predicate on the "group_id" field.
func GroupIDLTE(v string) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldLTE(FieldGroupID, v))
}

// GroupIDContains applies the Contains predicate on the "group_id" field.
func GroupIDContains(v string) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldContains(FieldGroupID, v))
}

// GroupIDHasPrefix applies the HasPrefix predicate on the "group_id" field.
func GroupIDHasPrefix(v string) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldHasPrefix(FieldGroupID, v))
}

// GroupIDHasSuffix applies the HasSuffix predicate on the "group_id" field.
func GroupIDHasSuffix(v string) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldHasSuffix(FieldGroupID, v))
}

// GroupIDEqualFold applies the EqualFold predicate on the "group_id" field.
func GroupIDEqualFold(v string) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldEqualFold(FieldGroupID, v))
}

// GroupIDContainsFold applies the ContainsFold predicate on the "group_id" field.
func GroupIDContainsFold(v string) predicate.UserGroup {
	return predicate.UserGroup(sql.FieldContainsFold(FieldGroupID, v))
}

// HasUser applies the HasEdge predicate on the "user" edge.
func HasUser() predicate.UserGroup {
	return predicate.UserGroup(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, UserTable, UserColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasUserWith applies the HasEdge predicate on the "user" edge with a given conditions (other predicates).
func HasUserWith(preds ...predicate.User) predicate.UserGroup {
	return predicate.UserGroup(func(s *sql.Selector) {
		step := newUserStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.UserGroup) predicate.UserGroup {
	return predicate.UserGroup(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.UserGroup) predicate.UserGroup {
	return predicate.UserGroup(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.UserGroup) predicate.UserGroup {
	return predicate.UserGroup(sql.NotPredicates(p))
}
