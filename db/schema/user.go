package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type User struct {
	ent.Schema
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		IDMixin{},
		TimeMixin{},
		SoftDeleteMixin{},
	}
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("email").MaxLen(256).Unique(),
		field.String("username").MaxLen(32),
		field.String("avatar").Optional(),
		field.String("joined_from").Optional(),
		field.Bool("banned").Default(false),
		field.Int("rank").Default(1),
		field.String("github_id").MaxLen(16).Unique().Optional(),
		field.String("google_id").MaxLen(16).Unique().Optional(),
		field.String("facebook_id").MaxLen(16).Unique().Optional(),
	}
}

func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("email"),
		index.Fields("github_id"),
		index.Fields("google_id"),
		index.Fields("facebook_id"),
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("guides", Guide.Type),
		edge.To("sessions", UserSession.Type).
			StorageKey(edge.Column("user_id")),
		edge.To("mods", Mod.Type).
			Through("user_mods", UserMod.Type),
		edge.To("groups", UserGroup.Type),
	}
}
