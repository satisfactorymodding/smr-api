package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type Announcement struct {
	ent.Schema
}

func (Announcement) Mixin() []ent.Mixin {
	return []ent.Mixin{
		IDMixin{},
		TimeMixin{},
		SoftDeleteMixin{},
	}
}

func (Announcement) Fields() []ent.Field {
	return []ent.Field{
		field.String("message"),
		field.String("importance"),
	}
}
