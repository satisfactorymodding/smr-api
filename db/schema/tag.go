package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Tag struct {
	ent.Schema
}

func (Tag) Mixin() []ent.Mixin {
	return []ent.Mixin{
		IDMixin{},
		TimeMixin{},
		SoftDeleteMixin{},
	}
}

func (Tag) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").MaxLen(24).Unique(),
	}
}

func (Tag) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("mods", Mod.Type).
			Ref("tags").
			Through("mod_tags", ModTag.Type),
		edge.From("guides", Guide.Type).
			Ref("tags").
			Through("guide_tags", GuideTag.Type),
	}
}
