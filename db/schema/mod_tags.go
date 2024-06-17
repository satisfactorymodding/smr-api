package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type ModTag struct {
	ent.Schema
}

func (ModTag) Annotations() []schema.Annotation {
	return []schema.Annotation{
		field.ID("mod_id", "tag_id"),
	}
}

func (ModTag) Fields() []ent.Field {
	return []ent.Field{
		field.String("mod_id"),
		field.String("tag_id"),
	}
}

func (ModTag) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("mod", Mod.Type).
			Unique().
			Required().
			Field("mod_id"),
		edge.To("tag", Tag.Type).
			Unique().
			Required().
			Field("tag_id"),
	}
}
