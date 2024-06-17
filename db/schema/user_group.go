package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type UserGroup struct {
	ent.Schema
}

func (UserGroup) Mixin() []ent.Mixin {
	return []ent.Mixin{
		IDMixin{},
		TimeMixin{},
		SoftDeleteMixin{},
	}
}

func (UserGroup) Fields() []ent.Field {
	return []ent.Field{
		field.String("user_id").MaxLen(14),
		field.String("group_id").MaxLen(14),
	}
}

func (UserGroup) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "group_id").Unique(),
	}
}

func (UserGroup) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("groups").
			Field("user_id").
			Unique().
			Required(),
	}
}
