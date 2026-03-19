package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"github.com/satisfactorymodding/smr-api/util"
)

type Modpack struct {
	ent.Schema
}

func (Modpack) Mixin() []ent.Mixin {
	return []ent.Mixin{
		IDMixin{},
		TimeMixin{},
		SoftDeleteMixin{},
	}
}

func (Modpack) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("short_description").MaxLen(128),
		field.String("full_description"),
		field.String("logo").Optional(),
		field.String("logo_thumbhash").Optional(),
		field.String("creator_id"),
		field.Uint("views").Default(0),
		field.Uint("hotness").Default(0),
		field.Uint("installs").Default(0),
		field.Uint("popularity").Default(0),
		field.Bool("hidden").Default(false),
		field.JSON("compatibility", &util.CompatibilityInfo{}).Optional(),
		field.String("parent_id").Optional().Immutable(),
	}
}

func (Modpack) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("children", Modpack.Type).
			Annotations(entsql.OnDelete(entsql.Restrict)),
		edge.From("parent", Modpack.Type).
			Ref("children").
			Field("parent_id").
			Unique().
			Immutable().
			Annotations(entsql.OnDelete(entsql.Restrict)),
		edge.To("targets", ModpackTarget.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("releases", ModpackRelease.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("mods", Mod.Type).
			Through("modpack_mods", ModpackMod.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.From("authors", User.Type).
			Ref("modpacks").
			Through("user_modpacks", UserModpack.Type),
		edge.To("tags", Tag.Type).
			Through("modpack_tags", ModpackTag.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}
