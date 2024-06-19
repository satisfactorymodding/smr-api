package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type SatisfactoryVersion struct {
	ent.Schema
}

func (SatisfactoryVersion) Mixin() []ent.Mixin {
	return []ent.Mixin{
		IDMixin{},
		TimeMixin{},
		SoftDeleteMixin{},
	}
}

func (SatisfactoryVersion) Fields() []ent.Field {
	return []ent.Field{
		field.Int("version").Unique(),
		field.String("engine_version").MaxLen(16).Default("4.26"),
	}
}
