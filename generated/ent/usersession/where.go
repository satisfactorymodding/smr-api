// Code generated by ent, DO NOT EDIT.

package usersession

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/satisfactorymodding/smr-api/generated/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id string) predicate.UserSession {
	return predicate.UserSession(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id string) predicate.UserSession {
	return predicate.UserSession(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id string) predicate.UserSession {
	return predicate.UserSession(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...string) predicate.UserSession {
	return predicate.UserSession(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...string) predicate.UserSession {
	return predicate.UserSession(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id string) predicate.UserSession {
	return predicate.UserSession(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id string) predicate.UserSession {
	return predicate.UserSession(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id string) predicate.UserSession {
	return predicate.UserSession(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id string) predicate.UserSession {
	return predicate.UserSession(sql.FieldLTE(FieldID, id))
}

// IDEqualFold applies the EqualFold predicate on the ID field.
func IDEqualFold(id string) predicate.UserSession {
	return predicate.UserSession(sql.FieldEqualFold(FieldID, id))
}

// IDContainsFold applies the ContainsFold predicate on the ID field.
func IDContainsFold(id string) predicate.UserSession {
	return predicate.UserSession(sql.FieldContainsFold(FieldID, id))
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.UserSession {
	return predicate.UserSession(sql.FieldEQ(FieldCreatedAt, v))
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.UserSession {
	return predicate.UserSession(sql.FieldEQ(FieldUpdatedAt, v))
}

// DeletedAt applies equality check predicate on the "deleted_at" field. It's identical to DeletedAtEQ.
func DeletedAt(v time.Time) predicate.UserSession {
	return predicate.UserSession(sql.FieldEQ(FieldDeletedAt, v))
}

// Token applies equality check predicate on the "token" field. It's identical to TokenEQ.
func Token(v string) predicate.UserSession {
	return predicate.UserSession(sql.FieldEQ(FieldToken, v))
}

// UserAgent applies equality check predicate on the "user_agent" field. It's identical to UserAgentEQ.
func UserAgent(v string) predicate.UserSession {
	return predicate.UserSession(sql.FieldEQ(FieldUserAgent, v))
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.UserSession {
	return predicate.UserSession(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.UserSession {
	return predicate.UserSession(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.UserSession {
	return predicate.UserSession(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.UserSession {
	return predicate.UserSession(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.UserSession {
	return predicate.UserSession(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.UserSession {
	return predicate.UserSession(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.UserSession {
	return predicate.UserSession(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.UserSession {
	return predicate.UserSession(sql.FieldLTE(FieldCreatedAt, v))
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.UserSession {
	return predicate.UserSession(sql.FieldEQ(FieldUpdatedAt, v))
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.UserSession {
	return predicate.UserSession(sql.FieldNEQ(FieldUpdatedAt, v))
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.UserSession {
	return predicate.UserSession(sql.FieldIn(FieldUpdatedAt, vs...))
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.UserSession {
	return predicate.UserSession(sql.FieldNotIn(FieldUpdatedAt, vs...))
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v time.Time) predicate.UserSession {
	return predicate.UserSession(sql.FieldGT(FieldUpdatedAt, v))
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.UserSession {
	return predicate.UserSession(sql.FieldGTE(FieldUpdatedAt, v))
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.UserSession {
	return predicate.UserSession(sql.FieldLT(FieldUpdatedAt, v))
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.UserSession {
	return predicate.UserSession(sql.FieldLTE(FieldUpdatedAt, v))
}

// DeletedAtEQ applies the EQ predicate on the "deleted_at" field.
func DeletedAtEQ(v time.Time) predicate.UserSession {
	return predicate.UserSession(sql.FieldEQ(FieldDeletedAt, v))
}

// DeletedAtNEQ applies the NEQ predicate on the "deleted_at" field.
func DeletedAtNEQ(v time.Time) predicate.UserSession {
	return predicate.UserSession(sql.FieldNEQ(FieldDeletedAt, v))
}

// DeletedAtIn applies the In predicate on the "deleted_at" field.
func DeletedAtIn(vs ...time.Time) predicate.UserSession {
	return predicate.UserSession(sql.FieldIn(FieldDeletedAt, vs...))
}

// DeletedAtNotIn applies the NotIn predicate on the "deleted_at" field.
func DeletedAtNotIn(vs ...time.Time) predicate.UserSession {
	return predicate.UserSession(sql.FieldNotIn(FieldDeletedAt, vs...))
}

// DeletedAtGT applies the GT predicate on the "deleted_at" field.
func DeletedAtGT(v time.Time) predicate.UserSession {
	return predicate.UserSession(sql.FieldGT(FieldDeletedAt, v))
}

// DeletedAtGTE applies the GTE predicate on the "deleted_at" field.
func DeletedAtGTE(v time.Time) predicate.UserSession {
	return predicate.UserSession(sql.FieldGTE(FieldDeletedAt, v))
}

// DeletedAtLT applies the LT predicate on the "deleted_at" field.
func DeletedAtLT(v time.Time) predicate.UserSession {
	return predicate.UserSession(sql.FieldLT(FieldDeletedAt, v))
}

// DeletedAtLTE applies the LTE predicate on the "deleted_at" field.
func DeletedAtLTE(v time.Time) predicate.UserSession {
	return predicate.UserSession(sql.FieldLTE(FieldDeletedAt, v))
}

// DeletedAtIsNil applies the IsNil predicate on the "deleted_at" field.
func DeletedAtIsNil() predicate.UserSession {
	return predicate.UserSession(sql.FieldIsNull(FieldDeletedAt))
}

// DeletedAtNotNil applies the NotNil predicate on the "deleted_at" field.
func DeletedAtNotNil() predicate.UserSession {
	return predicate.UserSession(sql.FieldNotNull(FieldDeletedAt))
}

// TokenEQ applies the EQ predicate on the "token" field.
func TokenEQ(v string) predicate.UserSession {
	return predicate.UserSession(sql.FieldEQ(FieldToken, v))
}

// TokenNEQ applies the NEQ predicate on the "token" field.
func TokenNEQ(v string) predicate.UserSession {
	return predicate.UserSession(sql.FieldNEQ(FieldToken, v))
}

// TokenIn applies the In predicate on the "token" field.
func TokenIn(vs ...string) predicate.UserSession {
	return predicate.UserSession(sql.FieldIn(FieldToken, vs...))
}

// TokenNotIn applies the NotIn predicate on the "token" field.
func TokenNotIn(vs ...string) predicate.UserSession {
	return predicate.UserSession(sql.FieldNotIn(FieldToken, vs...))
}

// TokenGT applies the GT predicate on the "token" field.
func TokenGT(v string) predicate.UserSession {
	return predicate.UserSession(sql.FieldGT(FieldToken, v))
}

// TokenGTE applies the GTE predicate on the "token" field.
func TokenGTE(v string) predicate.UserSession {
	return predicate.UserSession(sql.FieldGTE(FieldToken, v))
}

// TokenLT applies the LT predicate on the "token" field.
func TokenLT(v string) predicate.UserSession {
	return predicate.UserSession(sql.FieldLT(FieldToken, v))
}

// TokenLTE applies the LTE predicate on the "token" field.
func TokenLTE(v string) predicate.UserSession {
	return predicate.UserSession(sql.FieldLTE(FieldToken, v))
}

// TokenContains applies the Contains predicate on the "token" field.
func TokenContains(v string) predicate.UserSession {
	return predicate.UserSession(sql.FieldContains(FieldToken, v))
}

// TokenHasPrefix applies the HasPrefix predicate on the "token" field.
func TokenHasPrefix(v string) predicate.UserSession {
	return predicate.UserSession(sql.FieldHasPrefix(FieldToken, v))
}

// TokenHasSuffix applies the HasSuffix predicate on the "token" field.
func TokenHasSuffix(v string) predicate.UserSession {
	return predicate.UserSession(sql.FieldHasSuffix(FieldToken, v))
}

// TokenEqualFold applies the EqualFold predicate on the "token" field.
func TokenEqualFold(v string) predicate.UserSession {
	return predicate.UserSession(sql.FieldEqualFold(FieldToken, v))
}

// TokenContainsFold applies the ContainsFold predicate on the "token" field.
func TokenContainsFold(v string) predicate.UserSession {
	return predicate.UserSession(sql.FieldContainsFold(FieldToken, v))
}

// UserAgentEQ applies the EQ predicate on the "user_agent" field.
func UserAgentEQ(v string) predicate.UserSession {
	return predicate.UserSession(sql.FieldEQ(FieldUserAgent, v))
}

// UserAgentNEQ applies the NEQ predicate on the "user_agent" field.
func UserAgentNEQ(v string) predicate.UserSession {
	return predicate.UserSession(sql.FieldNEQ(FieldUserAgent, v))
}

// UserAgentIn applies the In predicate on the "user_agent" field.
func UserAgentIn(vs ...string) predicate.UserSession {
	return predicate.UserSession(sql.FieldIn(FieldUserAgent, vs...))
}

// UserAgentNotIn applies the NotIn predicate on the "user_agent" field.
func UserAgentNotIn(vs ...string) predicate.UserSession {
	return predicate.UserSession(sql.FieldNotIn(FieldUserAgent, vs...))
}

// UserAgentGT applies the GT predicate on the "user_agent" field.
func UserAgentGT(v string) predicate.UserSession {
	return predicate.UserSession(sql.FieldGT(FieldUserAgent, v))
}

// UserAgentGTE applies the GTE predicate on the "user_agent" field.
func UserAgentGTE(v string) predicate.UserSession {
	return predicate.UserSession(sql.FieldGTE(FieldUserAgent, v))
}

// UserAgentLT applies the LT predicate on the "user_agent" field.
func UserAgentLT(v string) predicate.UserSession {
	return predicate.UserSession(sql.FieldLT(FieldUserAgent, v))
}

// UserAgentLTE applies the LTE predicate on the "user_agent" field.
func UserAgentLTE(v string) predicate.UserSession {
	return predicate.UserSession(sql.FieldLTE(FieldUserAgent, v))
}

// UserAgentContains applies the Contains predicate on the "user_agent" field.
func UserAgentContains(v string) predicate.UserSession {
	return predicate.UserSession(sql.FieldContains(FieldUserAgent, v))
}

// UserAgentHasPrefix applies the HasPrefix predicate on the "user_agent" field.
func UserAgentHasPrefix(v string) predicate.UserSession {
	return predicate.UserSession(sql.FieldHasPrefix(FieldUserAgent, v))
}

// UserAgentHasSuffix applies the HasSuffix predicate on the "user_agent" field.
func UserAgentHasSuffix(v string) predicate.UserSession {
	return predicate.UserSession(sql.FieldHasSuffix(FieldUserAgent, v))
}

// UserAgentIsNil applies the IsNil predicate on the "user_agent" field.
func UserAgentIsNil() predicate.UserSession {
	return predicate.UserSession(sql.FieldIsNull(FieldUserAgent))
}

// UserAgentNotNil applies the NotNil predicate on the "user_agent" field.
func UserAgentNotNil() predicate.UserSession {
	return predicate.UserSession(sql.FieldNotNull(FieldUserAgent))
}

// UserAgentEqualFold applies the EqualFold predicate on the "user_agent" field.
func UserAgentEqualFold(v string) predicate.UserSession {
	return predicate.UserSession(sql.FieldEqualFold(FieldUserAgent, v))
}

// UserAgentContainsFold applies the ContainsFold predicate on the "user_agent" field.
func UserAgentContainsFold(v string) predicate.UserSession {
	return predicate.UserSession(sql.FieldContainsFold(FieldUserAgent, v))
}

// HasUser applies the HasEdge predicate on the "user" edge.
func HasUser() predicate.UserSession {
	return predicate.UserSession(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, UserTable, UserColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasUserWith applies the HasEdge predicate on the "user" edge with a given conditions (other predicates).
func HasUserWith(preds ...predicate.User) predicate.UserSession {
	return predicate.UserSession(func(s *sql.Selector) {
		step := newUserStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.UserSession) predicate.UserSession {
	return predicate.UserSession(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.UserSession) predicate.UserSession {
	return predicate.UserSession(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.UserSession) predicate.UserSession {
	return predicate.UserSession(sql.NotPredicates(p))
}
