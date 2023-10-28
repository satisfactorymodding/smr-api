// Code generated by ent, DO NOT EDIT.

package modtag

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/satisfactorymodding/smr-api/generated/ent/predicate"
)

// ModID applies equality check predicate on the "mod_id" field. It's identical to ModIDEQ.
func ModID(v string) predicate.ModTag {
	return predicate.ModTag(sql.FieldEQ(FieldModID, v))
}

// TagID applies equality check predicate on the "tag_id" field. It's identical to TagIDEQ.
func TagID(v string) predicate.ModTag {
	return predicate.ModTag(sql.FieldEQ(FieldTagID, v))
}

// ModIDEQ applies the EQ predicate on the "mod_id" field.
func ModIDEQ(v string) predicate.ModTag {
	return predicate.ModTag(sql.FieldEQ(FieldModID, v))
}

// ModIDNEQ applies the NEQ predicate on the "mod_id" field.
func ModIDNEQ(v string) predicate.ModTag {
	return predicate.ModTag(sql.FieldNEQ(FieldModID, v))
}

// ModIDIn applies the In predicate on the "mod_id" field.
func ModIDIn(vs ...string) predicate.ModTag {
	return predicate.ModTag(sql.FieldIn(FieldModID, vs...))
}

// ModIDNotIn applies the NotIn predicate on the "mod_id" field.
func ModIDNotIn(vs ...string) predicate.ModTag {
	return predicate.ModTag(sql.FieldNotIn(FieldModID, vs...))
}

// ModIDGT applies the GT predicate on the "mod_id" field.
func ModIDGT(v string) predicate.ModTag {
	return predicate.ModTag(sql.FieldGT(FieldModID, v))
}

// ModIDGTE applies the GTE predicate on the "mod_id" field.
func ModIDGTE(v string) predicate.ModTag {
	return predicate.ModTag(sql.FieldGTE(FieldModID, v))
}

// ModIDLT applies the LT predicate on the "mod_id" field.
func ModIDLT(v string) predicate.ModTag {
	return predicate.ModTag(sql.FieldLT(FieldModID, v))
}

// ModIDLTE applies the LTE predicate on the "mod_id" field.
func ModIDLTE(v string) predicate.ModTag {
	return predicate.ModTag(sql.FieldLTE(FieldModID, v))
}

// ModIDContains applies the Contains predicate on the "mod_id" field.
func ModIDContains(v string) predicate.ModTag {
	return predicate.ModTag(sql.FieldContains(FieldModID, v))
}

// ModIDHasPrefix applies the HasPrefix predicate on the "mod_id" field.
func ModIDHasPrefix(v string) predicate.ModTag {
	return predicate.ModTag(sql.FieldHasPrefix(FieldModID, v))
}

// ModIDHasSuffix applies the HasSuffix predicate on the "mod_id" field.
func ModIDHasSuffix(v string) predicate.ModTag {
	return predicate.ModTag(sql.FieldHasSuffix(FieldModID, v))
}

// ModIDEqualFold applies the EqualFold predicate on the "mod_id" field.
func ModIDEqualFold(v string) predicate.ModTag {
	return predicate.ModTag(sql.FieldEqualFold(FieldModID, v))
}

// ModIDContainsFold applies the ContainsFold predicate on the "mod_id" field.
func ModIDContainsFold(v string) predicate.ModTag {
	return predicate.ModTag(sql.FieldContainsFold(FieldModID, v))
}

// TagIDEQ applies the EQ predicate on the "tag_id" field.
func TagIDEQ(v string) predicate.ModTag {
	return predicate.ModTag(sql.FieldEQ(FieldTagID, v))
}

// TagIDNEQ applies the NEQ predicate on the "tag_id" field.
func TagIDNEQ(v string) predicate.ModTag {
	return predicate.ModTag(sql.FieldNEQ(FieldTagID, v))
}

// TagIDIn applies the In predicate on the "tag_id" field.
func TagIDIn(vs ...string) predicate.ModTag {
	return predicate.ModTag(sql.FieldIn(FieldTagID, vs...))
}

// TagIDNotIn applies the NotIn predicate on the "tag_id" field.
func TagIDNotIn(vs ...string) predicate.ModTag {
	return predicate.ModTag(sql.FieldNotIn(FieldTagID, vs...))
}

// TagIDGT applies the GT predicate on the "tag_id" field.
func TagIDGT(v string) predicate.ModTag {
	return predicate.ModTag(sql.FieldGT(FieldTagID, v))
}

// TagIDGTE applies the GTE predicate on the "tag_id" field.
func TagIDGTE(v string) predicate.ModTag {
	return predicate.ModTag(sql.FieldGTE(FieldTagID, v))
}

// TagIDLT applies the LT predicate on the "tag_id" field.
func TagIDLT(v string) predicate.ModTag {
	return predicate.ModTag(sql.FieldLT(FieldTagID, v))
}

// TagIDLTE applies the LTE predicate on the "tag_id" field.
func TagIDLTE(v string) predicate.ModTag {
	return predicate.ModTag(sql.FieldLTE(FieldTagID, v))
}

// TagIDContains applies the Contains predicate on the "tag_id" field.
func TagIDContains(v string) predicate.ModTag {
	return predicate.ModTag(sql.FieldContains(FieldTagID, v))
}

// TagIDHasPrefix applies the HasPrefix predicate on the "tag_id" field.
func TagIDHasPrefix(v string) predicate.ModTag {
	return predicate.ModTag(sql.FieldHasPrefix(FieldTagID, v))
}

// TagIDHasSuffix applies the HasSuffix predicate on the "tag_id" field.
func TagIDHasSuffix(v string) predicate.ModTag {
	return predicate.ModTag(sql.FieldHasSuffix(FieldTagID, v))
}

// TagIDEqualFold applies the EqualFold predicate on the "tag_id" field.
func TagIDEqualFold(v string) predicate.ModTag {
	return predicate.ModTag(sql.FieldEqualFold(FieldTagID, v))
}

// TagIDContainsFold applies the ContainsFold predicate on the "tag_id" field.
func TagIDContainsFold(v string) predicate.ModTag {
	return predicate.ModTag(sql.FieldContainsFold(FieldTagID, v))
}

// HasMod applies the HasEdge predicate on the "mod" edge.
func HasMod() predicate.ModTag {
	return predicate.ModTag(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, ModColumn),
			sqlgraph.Edge(sqlgraph.M2O, false, ModTable, ModColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasModWith applies the HasEdge predicate on the "mod" edge with a given conditions (other predicates).
func HasModWith(preds ...predicate.Mod) predicate.ModTag {
	return predicate.ModTag(func(s *sql.Selector) {
		step := newModStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasTag applies the HasEdge predicate on the "tag" edge.
func HasTag() predicate.ModTag {
	return predicate.ModTag(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, TagColumn),
			sqlgraph.Edge(sqlgraph.M2O, false, TagTable, TagColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasTagWith applies the HasEdge predicate on the "tag" edge with a given conditions (other predicates).
func HasTagWith(preds ...predicate.Tag) predicate.ModTag {
	return predicate.ModTag(func(s *sql.Selector) {
		step := newTagStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.ModTag) predicate.ModTag {
	return predicate.ModTag(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.ModTag) predicate.ModTag {
	return predicate.ModTag(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.ModTag) predicate.ModTag {
	return predicate.ModTag(sql.NotPredicates(p))
}
