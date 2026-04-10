package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type ModpackRelease struct {
	ent.Schema
}

func (ModpackRelease) Mixin() []ent.Mixin {
	return []ent.Mixin{
		IDMixin{},
		TimeMixin{},
	}
}

func (ModpackRelease) Fields() []ent.Field {
	return []ent.Field{
		field.String("modpack_id"),
		field.String("version"),
		field.String("changelog"),
		field.String("lockfile"),
	}
}

func (ModpackRelease) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("modpack", Modpack.Type).
			Ref("releases").
			Field("modpack_id").
			Unique().
			Required().
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("targets", ModpackTarget.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

func (ModpackRelease) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("modpack_id", "version").Unique(),
	}
}
