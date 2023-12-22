package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"github.com/satisfactorymodding/smr-api/util"
)

type SmlVersion struct {
	ent.Schema
}

func (SmlVersion) Mixin() []ent.Mixin {
	return []ent.Mixin{
		IDMixin{},
		TimeMixin{},
		SoftDeleteMixin{},
	}
}

func (SmlVersion) Fields() []ent.Field {
	return []ent.Field{
		field.String("version").MaxLen(32).Unique(),
		field.Int("satisfactory_version"),
		field.Enum("stability").GoType(util.Stability("")),
		field.Time("date"),
		field.String("link"),
		field.String("changelog"),
		field.String("bootstrap_version").MaxLen(14).Optional(),
		field.String("engine_version").MaxLen(16).Default("4.26"),
	}
}

func (SmlVersion) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("targets", SmlVersionTarget.Type),
	}
}
