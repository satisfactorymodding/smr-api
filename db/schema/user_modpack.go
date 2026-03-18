package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type UserModpack struct {
	ent.Schema
}

func (UserModpack) Annotations() []schema.Annotation {
	return []schema.Annotation{
		field.ID("user_id", "modpack_id"),
	}
}

func (UserModpack) Fields() []ent.Field {
	return []ent.Field{
		field.String("user_id"),
		field.String("modpack_id"),
		field.String("role"),
	}
}

func (UserModpack) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("user", User.Type).
            Unique().
            Required().
            Field("user_id"),

        edge.To("modpack", Modpack.Type).
            Unique().
            Required().
            Field("modpack_id"),
    }
}