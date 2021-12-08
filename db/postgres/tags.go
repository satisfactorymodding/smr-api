package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/satisfactorymodding/smr-api/util"

	"github.com/finnbear/moderation"

	"github.com/patrickmn/go-cache"
)

func ValidateTagName(tag string) error {
	if len(tag) > 20 {
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
		return nil, fmt.Errorf("Tag '%v' already exists", tag.Name)
	}
	tag.ID = util.GenerateUniqueID()

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

func GetTagByID(tagId string, ctx *context.Context) *Tag {
	cacheKey := "GetTagById_" + tagId

	if tag, ok := dbCache.Get(cacheKey); ok {
		return tag.(*Tag)
	}

	var tag Tag
	DBCtx(ctx).Find(&tag, "id = ?", tagId)

	if tag.ID == "" {
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
