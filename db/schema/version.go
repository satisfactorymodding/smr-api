package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"github.com/satisfactorymodding/smr-api/util"
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
		field.String("mod_id"),
		field.String("version").MaxLen(16),
		field.String("sml_version").MaxLen(16),
		field.String("changelog").Optional(),
		field.Uint("downloads").Default(0),
		field.String("key").Optional(),
		field.Enum("stability").GoType(util.Stability("")),
		field.Bool("approved").Default(false),
		field.Uint("hotness").Default(0),
		field.Bool("denied").Default(false),
		field.String("metadata").Optional(),
		field.String("mod_reference").MaxLen(32),
		field.Int("version_major").Optional(),
		field.Int("version_minor").Optional(),
		field.Int("version_patch").Optional(),
		field.Int64("size").Optional(),
		field.String("hash").MinLen(64).MaxLen(64).Optional(),
	}
}

func (Version) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("mod", Mod.Type).
			Ref("versions").
			Field("mod_id").
			Unique().
			Required(),
		edge.To("dependencies", Mod.Type).
			Through("version_dependencies", VersionDependency.Type),
		edge.To("targets", VersionTarget.Type),
	}
}
