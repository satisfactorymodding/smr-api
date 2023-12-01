package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type VersionTarget struct {
	ent.Schema
}

func (VersionTarget) Mixin() []ent.Mixin {
	return []ent.Mixin{
		IDMixin{},
	}
}

func (VersionTarget) Fields() []ent.Field {
	return []ent.Field{
		field.String("version_id"),
		field.String("target_name"),
		field.String("key"),
		field.String("hash"),
		field.Int64("size"),
	}
}

func (VersionTarget) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("sml_version", Version.Type).
			Ref("targets").
			Field("version_id").
			Unique().
			Required(),
	}
}

func (VersionTarget) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("version_id", "target_name").Unique(),
	}
}
