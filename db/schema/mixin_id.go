package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"

	"github.com/satisfactorymodding/smr-api/util"
)

type IDMixin struct {
	mixin.Schema
}

func (IDMixin) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").DefaultFunc(util.GenerateUniqueID),
	}
}

func (IDMixin) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id"),
	}
}
