package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Version struct {
	ent.Schema
}

func (Version) Mixin() []ent.Mixin {
	return []ent.Mixin{
		IDMixin{},
		TimeMixin{},
		SoftDeleteMixin{},
	}
}

func (Version) Fields() []ent.Field {
	return []ent.Field{
		field.String("version").MaxLen(16),
		field.String("sml_version").MaxLen(16),
		field.String("changelog"),
		field.Uint("downloads"),
		field.String("key"),
		field.Enum("stability").Values("alpha", "beta", "release"),
		field.Bool("approved").Default(false),
		field.Uint("hotness"),
		field.Bool("denied").Default(false),
		field.String("metadata"),
		field.String("mod_reference").MaxLen(32),
		field.Int("version_major"),
		field.Int("version_minor"),
		field.Int("version_patch"),
		field.Int64("size"),
		field.String("hash").MinLen(64).MaxLen(64),
	}
}

func (Version) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("mod", Mod.Type).
			Ref("versions").
			Unique().
			Required(),
		edge.To("dependencies", Mod.Type).
			Through("version_dependencies", VersionDependency.Type),
	}
}
