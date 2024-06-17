package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type VersionDependency struct {
	ent.Schema
}

func (VersionDependency) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
		SoftDeleteMixin{},
	}
}

func (VersionDependency) Annotations() []schema.Annotation {
	return []schema.Annotation{
		field.ID("version_id", "mod_id"),
	}
}

func (VersionDependency) Fields() []ent.Field {
	return []ent.Field{
		field.String("version_id"),
		field.String("mod_id"),
		field.String("condition").MaxLen(64),
		field.Bool("optional"),
	}
}

func (VersionDependency) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("version", Version.Type).
			Unique().
			Required().
			Field("version_id"),
		edge.To("mod", Mod.Type).
			Unique().
			Required().
			Field("mod_id"),
	}
}
