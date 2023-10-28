package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/satisfactorymodding/smr-api/db/postgres"
)

type Mod struct {
	ent.Schema
}

func (Mod) Mixin() []ent.Mixin {
	return []ent.Mixin{
		IDMixin{},
		TimeMixin{},
		SoftDeleteMixin{},
	}
}

func (Mod) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").MaxLen(32),
		field.String("short_description").MaxLen(128),
		field.String("full_description"),
		field.String("logo"),
		field.String("source_url"),
		field.String("creator_id"),
		field.Bool("approved").Default(false),
		field.Uint("views"),
		field.Uint("hotness"),
		field.Uint("popularity"),
		field.Uint("downloads"),
		field.Bool("denied").Default(false),
		field.Time("last_version_date"),
		field.String("mod_reference").MaxLen(32).Unique(),
		field.Bool("hidden"),
		field.JSON("compatibility", &postgres.CompatibilityInfo{}),
	}
}

func (Mod) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("last_version_date"),
	}
}

func (Mod) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("versions", Version.Type).
			StorageKey(edge.Column("mod_id")),
		edge.From("authors", User.Type).
			Ref("mods").
			Through("user_mods", UserMod.Type),
		edge.To("tags", Tag.Type).
			Through("mod_tags", ModTag.Type),
		edge.From("dependents", Version.Type).
			Ref("dependencies").
			Through("version_dependencies", VersionDependency.Type),
	}
}
