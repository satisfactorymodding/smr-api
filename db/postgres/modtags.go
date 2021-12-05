package postgres

import (
	"context"
	"errors"

	"github.com/patrickmn/go-cache"
)

func CreateModTag(modTag *ModTag, ctx *context.Context) (*ModTag, error) {
	if len(modTag.Name) > 20 {
		return nil, errors.New("Mod tag name is more than 20 characters")
	}

	exists := GetModTagByName(modTag.Name, ctx)
	if exists != nil {
		return nil, errors.New("Mod tag already exists")
	}

	DBCtx(ctx).Create(&modTag)
	return modTag, nil
}

func GetModTagByName(modTagName string, ctx *context.Context) *ModTag {
	cacheKey := "GetModTagByName_" + modTagName

	if modTag, ok := dbCache.Get(cacheKey); ok {
		return modTag.(*ModTag)
	}

	var modTag ModTag
	DBCtx(ctx).Find(&modTag, "name = ?", modTagName)

	if modTag.Name == "" {
		return nil
	}

	dbCache.Set(cacheKey, modTag, cache.DefaultExpiration)

	return &modTag
}

func GetModTags(ctx *context.Context) []ModTag {
	cacheKey := "GetModTags"

	if modTags, ok := dbCache.Get(cacheKey); ok {
		return modTags.([]ModTag)
	}

	var modTags []ModTag
	DBCtx(ctx).Find(&modTags)

	dbCache.Set(cacheKey, modTags, cache.DefaultExpiration)

	return modTags
}
