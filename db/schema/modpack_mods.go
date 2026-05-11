package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type ModpackMod struct {
	ent.Schema
}

func (ModpackMod) Annotations() []schema.Annotation {
	return []schema.Annotation{
		field.ID("modpack_id", "mod_id"),
	}
}

func (ModpackMod) Fields() []ent.Field {
	return []ent.Field{
		field.String("modpack_id"),
		field.String("mod_id"),
		field.String("version_constraint"),
	}
}

func (ModpackMod) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("modpack", Modpack.Type).
			Unique().
			Required().
			Field("modpack_id").
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("mod", Mod.Type).
			Unique().
			Required().
			Field("mod_id").
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}
