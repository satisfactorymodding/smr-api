package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type ModpackTag struct {
	ent.Schema
}

func (ModpackTag) Annotations() []schema.Annotation {
	return []schema.Annotation{
		field.ID("modpack_id", "tag_id"),
	}
}

func (ModpackTag) Fields() []ent.Field {
	return []ent.Field{
		field.String("modpack_id"),
		field.String("tag_id"),
	}
}

func (ModpackTag) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("modpack", Modpack.Type).
			Unique().
			Required().
			Field("modpack_id").
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("tag", Tag.Type).
			Unique().
			Required().
			Field("tag_id").
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}
