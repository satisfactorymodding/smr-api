package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type ModpackTarget struct {
	ent.Schema
}

func (ModpackTarget) Mixin() []ent.Mixin {
	return []ent.Mixin{
		IDMixin{},
	}
}

func (ModpackTarget) Fields() []ent.Field {
	return []ent.Field{
		field.String("modpack_id"),
		field.String("target_name"),
	}
}

func (ModpackTarget) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("modpack", Modpack.Type).
			Ref("targets").
			Field("modpack_id").
			Unique().
			Required().
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

func (ModpackTarget) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("modpack_id", "target_name").Unique(),
	}
}
