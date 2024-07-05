package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type UserSession struct {
	ent.Schema
}

func (UserSession) Mixin() []ent.Mixin {
	return []ent.Mixin{
		IDMixin{},
		TimeMixin{},
		SoftDeleteMixin{},
	}
}

func (UserSession) Fields() []ent.Field {
	return []ent.Field{
		field.String("token").MaxLen(512).Unique(),
		field.String("user_agent").Optional(),
	}
}

func (UserSession) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("token").Unique().StorageKey("uix_user_sessions_token"),
	}
}

func (UserSession) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("sessions").
			Unique().
			Required(),
	}
}
