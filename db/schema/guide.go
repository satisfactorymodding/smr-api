package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Guide struct {
	ent.Schema
}

func (Guide) Mixin() []ent.Mixin {
	return []ent.Mixin{
		IDMixin{},
		TimeMixin{},
		SoftDeleteMixin{},
	}
}

func (Guide) Fields() []ent.Field {
	return []ent.Field{
		field.String("user_id").Optional(),
		field.String("name").MaxLen(32),
		field.String("short_description").MaxLen(128),
		field.String("guide"),
		field.Int("views").Default(0),
	}
}

func (Guide) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("guides").
			Field("user_id").
			Unique(),
		edge.To("tags", Tag.Type).
			Through("guide_tags", GuideTag.Type),
	}
}
