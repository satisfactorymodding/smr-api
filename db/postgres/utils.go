package postgres

import (
	"context"
	"gorm.io/gorm"
)

type (
	ContextDB struct{}
)

func ContextWithDB(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, ContextDB{}, db)
}

func DBFromContext(ctx context.Context) *gorm.DB {
	value := ctx.Value(ContextDB{})

	if value == nil {
		return nil
	}

	return value.(*gorm.DB)
}
