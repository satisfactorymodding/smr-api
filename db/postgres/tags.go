package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/finnbear/moderation"

	"github.com/patrickmn/go-cache"
)

func ValidateTagName(tag string) error {
	if len(tag) < 21 {
		return errors.New("Tag name is over 20 characters long")
	}
	if moderation.IsInappropriate(tag) {
		return errors.New("Tag name is inapproriate")
	}
	return nil
}

func CreateTag(tag *Tag, ctx *context.Context, ratelimit bool) (*Tag, error) {
	err := ValidateTagName(tag.Name)
	if err != nil {
		return nil, err
	}

	exists := GetTagByName(tag.Name, ctx)
	if exists != nil {
		return nil, errors.New(fmt.Sprintf("Tag '%v' already exists", tag.Name))
	}

	DBCtx(ctx).Create(&tag)
	return tag, nil
}

func GetTagByName(tagName string, ctx *context.Context) *Tag {
	cacheKey := "GetTagByName_" + tagName

	if tag, ok := dbCache.Get(cacheKey); ok {
		return tag.(*Tag)
	}

	var tag Tag
	DBCtx(ctx).Find(&tag, "name = ?", tagName)

	if tag.Name == "" {
		return nil
	}

	dbCache.Set(cacheKey, tag, cache.DefaultExpiration)

	return &tag
}

func GetTags(ctx *context.Context) []Tag {
	cacheKey := "GetTags"

	if tags, ok := dbCache.Get(cacheKey); ok {
		return tags.([]Tag)
	}

	var tags []Tag
	DBCtx(ctx).Find(&tags)

	dbCache.Set(cacheKey, tags, cache.DefaultExpiration)

	return tags
}
