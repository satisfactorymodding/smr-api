// Code generated by ent, DO NOT EDIT.

package smlversiontarget

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/satisfactorymodding/smr-api/generated/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldLTE(FieldID, id))
}

// IDEqualFold applies the EqualFold predicate on the ID field.
func IDEqualFold(id string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldEqualFold(FieldID, id))
}

// IDContainsFold applies the ContainsFold predicate on the ID field.
func IDContainsFold(id string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldContainsFold(FieldID, id))
}

// VersionID applies equality check predicate on the "version_id" field. It's identical to VersionIDEQ.
func VersionID(v string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldEQ(FieldVersionID, v))
}

// TargetName applies equality check predicate on the "target_name" field. It's identical to TargetNameEQ.
func TargetName(v string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldEQ(FieldTargetName, v))
}

// Link applies equality check predicate on the "link" field. It's identical to LinkEQ.
func Link(v string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldEQ(FieldLink, v))
}

// VersionIDEQ applies the EQ predicate on the "version_id" field.
func VersionIDEQ(v string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldEQ(FieldVersionID, v))
}

// VersionIDNEQ applies the NEQ predicate on the "version_id" field.
func VersionIDNEQ(v string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldNEQ(FieldVersionID, v))
}

// VersionIDIn applies the In predicate on the "version_id" field.
func VersionIDIn(vs ...string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldIn(FieldVersionID, vs...))
}

// VersionIDNotIn applies the NotIn predicate on the "version_id" field.
func VersionIDNotIn(vs ...string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldNotIn(FieldVersionID, vs...))
}

// VersionIDGT applies the GT predicate on the "version_id" field.
func VersionIDGT(v string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldGT(FieldVersionID, v))
}

// VersionIDGTE applies the GTE predicate on the "version_id" field.
func VersionIDGTE(v string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldGTE(FieldVersionID, v))
}

// VersionIDLT applies the LT predicate on the "version_id" field.
func VersionIDLT(v string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldLT(FieldVersionID, v))
}

// VersionIDLTE applies the LTE predicate on the "version_id" field.
func VersionIDLTE(v string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldLTE(FieldVersionID, v))
}

// VersionIDContains applies the Contains predicate on the "version_id" field.
func VersionIDContains(v string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldContains(FieldVersionID, v))
}

// VersionIDHasPrefix applies the HasPrefix predicate on the "version_id" field.
func VersionIDHasPrefix(v string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldHasPrefix(FieldVersionID, v))
}

// VersionIDHasSuffix applies the HasSuffix predicate on the "version_id" field.
func VersionIDHasSuffix(v string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldHasSuffix(FieldVersionID, v))
}

// VersionIDEqualFold applies the EqualFold predicate on the "version_id" field.
func VersionIDEqualFold(v string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldEqualFold(FieldVersionID, v))
}

// VersionIDContainsFold applies the ContainsFold predicate on the "version_id" field.
func VersionIDContainsFold(v string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldContainsFold(FieldVersionID, v))
}

// TargetNameEQ applies the EQ predicate on the "target_name" field.
func TargetNameEQ(v string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldEQ(FieldTargetName, v))
}

// TargetNameNEQ applies the NEQ predicate on the "target_name" field.
func TargetNameNEQ(v string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldNEQ(FieldTargetName, v))
}

// TargetNameIn applies the In predicate on the "target_name" field.
func TargetNameIn(vs ...string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldIn(FieldTargetName, vs...))
}

// TargetNameNotIn applies the NotIn predicate on the "target_name" field.
func TargetNameNotIn(vs ...string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldNotIn(FieldTargetName, vs...))
}

// TargetNameGT applies the GT predicate on the "target_name" field.
func TargetNameGT(v string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldGT(FieldTargetName, v))
}

// TargetNameGTE applies the GTE predicate on the "target_name" field.
func TargetNameGTE(v string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldGTE(FieldTargetName, v))
}

// TargetNameLT applies the LT predicate on the "target_name" field.
func TargetNameLT(v string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldLT(FieldTargetName, v))
}

// TargetNameLTE applies the LTE predicate on the "target_name" field.
func TargetNameLTE(v string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldLTE(FieldTargetName, v))
}

// TargetNameContains applies the Contains predicate on the "target_name" field.
func TargetNameContains(v string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldContains(FieldTargetName, v))
}

// TargetNameHasPrefix applies the HasPrefix predicate on the "target_name" field.
func TargetNameHasPrefix(v string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldHasPrefix(FieldTargetName, v))
}

// TargetNameHasSuffix applies the HasSuffix predicate on the "target_name" field.
func TargetNameHasSuffix(v string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldHasSuffix(FieldTargetName, v))
}

// TargetNameEqualFold applies the EqualFold predicate on the "target_name" field.
func TargetNameEqualFold(v string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldEqualFold(FieldTargetName, v))
}

// TargetNameContainsFold applies the ContainsFold predicate on the "target_name" field.
func TargetNameContainsFold(v string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldContainsFold(FieldTargetName, v))
}

// LinkEQ applies the EQ predicate on the "link" field.
func LinkEQ(v string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldEQ(FieldLink, v))
}

// LinkNEQ applies the NEQ predicate on the "link" field.
func LinkNEQ(v string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldNEQ(FieldLink, v))
}

// LinkIn applies the In predicate on the "link" field.
func LinkIn(vs ...string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldIn(FieldLink, vs...))
}

// LinkNotIn applies the NotIn predicate on the "link" field.
func LinkNotIn(vs ...string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldNotIn(FieldLink, vs...))
}

// LinkGT applies the GT predicate on the "link" field.
func LinkGT(v string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldGT(FieldLink, v))
}

// LinkGTE applies the GTE predicate on the "link" field.
func LinkGTE(v string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldGTE(FieldLink, v))
}

// LinkLT applies the LT predicate on the "link" field.
func LinkLT(v string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldLT(FieldLink, v))
}

// LinkLTE applies the LTE predicate on the "link" field.
func LinkLTE(v string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldLTE(FieldLink, v))
}

// LinkContains applies the Contains predicate on the "link" field.
func LinkContains(v string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldContains(FieldLink, v))
}

// LinkHasPrefix applies the HasPrefix predicate on the "link" field.
func LinkHasPrefix(v string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldHasPrefix(FieldLink, v))
}

// LinkHasSuffix applies the HasSuffix predicate on the "link" field.
func LinkHasSuffix(v string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldHasSuffix(FieldLink, v))
}

// LinkEqualFold applies the EqualFold predicate on the "link" field.
func LinkEqualFold(v string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldEqualFold(FieldLink, v))
}

// LinkContainsFold applies the ContainsFold predicate on the "link" field.
func LinkContainsFold(v string) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.FieldContainsFold(FieldLink, v))
}

// HasSmlVersion applies the HasEdge predicate on the "sml_version" edge.
func HasSmlVersion() predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, SmlVersionTable, SmlVersionColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasSmlVersionWith applies the HasEdge predicate on the "sml_version" edge with a given conditions (other predicates).
func HasSmlVersionWith(preds ...predicate.SmlVersion) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(func(s *sql.Selector) {
		step := newSmlVersionStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.SmlVersionTarget) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.SmlVersionTarget) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.SmlVersionTarget) predicate.SmlVersionTarget {
	return predicate.SmlVersionTarget(sql.NotPredicates(p))
}
