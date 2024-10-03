package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type VirustotalResult struct {
	ent.Schema
}

func (VirustotalResult) Mixin() []ent.Mixin {
	return []ent.Mixin{
		IDMixin{},
		TimeMixin{},
	}
}

func (VirustotalResult) Fields() []ent.Field {
	return []ent.Field{
		field.Bool("safe").Default(false),
		field.String("hash").MaxLen(64).NotEmpty(),
		field.String("file_name").NotEmpty(),
		field.String("version_id").MaxLen(14).NotEmpty(),
	}
}

func (VirustotalResult) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("safe"),
		index.Fields("hash", "version_id").Unique(),
		index.Fields("file_name"),
	}
}

func (VirustotalResult) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("version", Version.Type).
			Ref("virustotal_results").
			Field("version_id").
			Unique().
			Required(),
	}
}
