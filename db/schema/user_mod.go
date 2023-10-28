package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type UserMod struct {
	ent.Schema
}

func (UserMod) Annotations() []schema.Annotation {
	return []schema.Annotation{
		field.ID("user_id", "mod_id"),
	}
}

func (UserMod) Fields() []ent.Field {
	return []ent.Field{
		field.String("user_id"),
		field.String("mod_id"),
		field.String("role"),
	}
}

func (UserMod) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user", User.Type).
			Unique().
			Required().
			Field("user_id"),
		edge.To("mod", Mod.Type).
			Unique().
			Required().
			Field("mod_id"),
	}
}
