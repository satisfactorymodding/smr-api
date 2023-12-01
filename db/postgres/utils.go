package postgres

import (
	"context"

	"gorm.io/gorm"
)

type (
	ContextDB struct{}
)

func DBFromContext(ctx context.Context) *gorm.DB {
	value := ctx.Value(ContextDB{})

	if value == nil {
		return nil
	}

	return value.(*gorm.DB)
}
