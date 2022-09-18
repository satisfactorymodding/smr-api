package postgres

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/util"

	"github.com/finnbear/moderation"
	"github.com/mitchellh/hashstructure/v2"
	"github.com/patrickmn/go-cache"
)

func ValidateTagName(tag string) error {
	if len(tag) > 24 {
		return errors.New("Tag name is over 24 characters long")
	}
	if len(tag) < 3 {
		return errors.New("Tag name is under 3 characters long")
	}
	if moderation.IsInappropriate(tag) {
		return errors.New("Tag name is inapproriate")
	}
	return nil
}

func CreateTag(ctx context.Context, tag *Tag, ratelimit bool) (*Tag, error) {
	err := ValidateTagName(tag.Name)
	if err != nil {
		return nil, err
	}

	if GetTagByName(ctx, tag.Name) != nil {
		return nil, fmt.Errorf("Tag %v already exists", tag.Name)
	}

	if ratelimit {
		var tags []Tag
		DBCtx(ctx).Order("created_at asc").Find(&tags, "created_at > ?", time.Now().Add(time.Hour*24*-1))

		currentAvailable := float64(8)
		lastTagTime := time.Now()
		for _, tag := range tags {
			currentAvailable--
			if tag.CreatedAt.After(lastTagTime) {
				diff := tag.CreatedAt.Sub(lastTagTime)
				currentAvailable = math.Min(4, currentAvailable+diff.Hours()/6)
			}
			lastTagTime = tag.CreatedAt
		}

		if currentAvailable < 1 {
			timeToWait := time.Until(lastTagTime.Add(time.Hour * 6)).Minutes()
			return nil, fmt.Errorf("please wait %.0f minutes to create another tag", timeToWait)
		}
	}

	tag.ID = util.GenerateUniqueID()

	err = DBCtx(ctx).Create(&tag).Error

	if err != nil {
		return nil, fmt.Errorf("Could not create tag: %w", err)
	}

	return tag, nil
}

func GetTagByName(ctx context.Context, tagName string) *Tag {
	cacheKey := "GetTagByName_" + tagName

	if tag, ok := dbCache.Get(cacheKey); ok {
		return tag.(*Tag)
	}

	var tag Tag
	DBCtx(ctx).Find(&tag, "name = ?", tagName)

	if tag.Name == "" {
		return nil
	}

	dbCache.Set(cacheKey, &tag, cache.DefaultExpiration)

	return &tag
}

func GetTagByID(ctx context.Context, tagID string) *Tag {
	cacheKey := "GetTagById_" + tagID

	if tag, ok := dbCache.Get(cacheKey); ok {
		return tag.(*Tag)
	}

	var tag Tag
	DBCtx(ctx).Find(&tag, "id = ?", tagID)

	if tag.ID == "" {
		return nil
	}

	dbCache.Set(cacheKey, &tag, cache.DefaultExpiration)

	return &tag
}

func GetTags(ctx context.Context, filter *generated.TagFilter) []Tag {
	cacheKey := ""
	hash, err := hashstructure.Hash(filter, hashstructure.FormatV2, nil)
	if err == nil {
		cacheKey = "GetTags" + strconv.FormatUint(hash, 10)
		if tags, ok := dbCache.Get(cacheKey); ok {
			return tags.([]Tag)
		}
	}

	query := DBCtx(ctx)

	if filter != nil {
		if filter.Search != nil && *filter.Search != "" {
			cleanSearch := strings.Replace(strings.TrimSpace(*filter.Search), " ", " & ", -1)
			sub := DBCtx(ctx).Table("tags")
			sub = sub.Select("id, similarity(name, ?) as s", cleanSearch, cleanSearch, cleanSearch)

			query = query.Joins("INNER JOIN (?) AS t1 on t1.id = tags.id", sub)
			query = query.Where("t1.s > 0.2")
		}

		query = query.Limit(*filter.Limit).Offset(*filter.Offset)

		if filter.Ids != nil && len(filter.Ids) > 0 {
			query = query.Where("id in (?)", filter.Ids)
		}

		query = query.Order("name " + string(*filter.Order))
	}

	var tags []Tag
	query.Find(&tags)

	if cacheKey != "" {
		dbCache.Set(cacheKey, tags, cache.DefaultExpiration)
	}

	return tags
}
