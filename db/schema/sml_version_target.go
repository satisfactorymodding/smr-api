package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type SmlVersionTarget struct {
	ent.Schema
}

func (SmlVersionTarget) Mixin() []ent.Mixin {
	return []ent.Mixin{
		IDMixin{},
	}
}

func (SmlVersionTarget) Fields() []ent.Field {
	return []ent.Field{
		field.String("version_id"),
		field.String("target_name"),
		field.String("link"),
	}
}

func (SmlVersionTarget) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("sml_version", SmlVersion.Type).
			Ref("targets").
			Field("version_id").
			Unique().
			Required(),
	}
}

func (SmlVersionTarget) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("version_id", "target_name").Unique(),
	}
}
