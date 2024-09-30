package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// YourTableName holds the schema definition for the YourTableName entity.
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
		field.String("url").
			NotEmpty(),
		field.String("hash").Unique().NotEmpty(),
		field.String("file_name").NotEmpty(),
		field.String("version_id").NotEmpty(),
	}
}

func (VirustotalResult) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "virustotal_results"},
	}
}

func (VirustotalResult) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("version", Version.Type).
			Ref("virustotalResults").
			Field("version_id").
			Unique().
			Required(),
	}
}
