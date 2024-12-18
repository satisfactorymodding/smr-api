package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type GuideTag struct {
	ent.Schema
}

func (GuideTag) Annotations() []schema.Annotation {
	return []schema.Annotation{
		field.ID("guide_id", "tag_id"),
	}
}

func (GuideTag) Fields() []ent.Field {
	return []ent.Field{
		field.String("guide_id"),
		field.String("tag_id"),
	}
}

func (GuideTag) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("guide", Guide.Type).
			Unique().
			Required().
			Field("guide_id"),
		edge.To("tag", Tag.Type).
			Unique().
			Required().
			Field("tag_id"),
	}
}
