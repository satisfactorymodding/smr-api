package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
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
		field.String("token").MaxLen(256).Unique(),
		field.String("user_agent").Optional(),
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
