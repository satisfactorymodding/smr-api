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
		field.String("name").MaxLen(32),
		field.String("short_description").MaxLen(128),
		field.String("guide"),
		field.Int("views"),
	}
}

func (Guide) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("guides").
			Unique().
			Required(),
		edge.To("tags", Tag.Type).
			Through("guide_tags", GuideTag.Type),
	}
}
