package postgres

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/satisfactorymodding/smr-api/models"
	"github.com/satisfactorymodding/smr-api/util"

	"github.com/patrickmn/go-cache"
)

func CreateGuide(guide *Guide, ctx *context.Context) (*Guide, error) {
	// Allow only 8 new guides per 24h

	guide.ID = util.GenerateUniqueID()

	var guides []Guide
	DBCtx(ctx).Order("created_at asc").Find(&guides, "user_id = ? AND created_at > ?", guide.UserID, time.Now().Add(time.Hour*24*-1))

	currentAvailable := float64(8)
	lastGuideTime := time.Now()
	for _, guide := range guides {
		currentAvailable--
		if guide.CreatedAt.After(lastGuideTime) {
			diff := guide.CreatedAt.Sub(lastGuideTime)
			currentAvailable = math.Min(8, currentAvailable+diff.Hours()/3)
		}
		lastGuideTime = guide.CreatedAt
	}

	if currentAvailable < 1 {
		timeToWait := time.Until(lastGuideTime.Add(time.Hour * 6)).Minutes()
		return nil, fmt.Errorf("please wait %.0f minutes to post another guide", timeToWait)
	}

	DBCtx(ctx).Create(&guide)

	return guide, nil
}

func GetGuideByID(guideID string, ctx *context.Context) *Guide {
	cacheKey := "GetGuideById_" + guideID

	if guide, ok := dbCache.Get(cacheKey); ok {
		return guide.(*Guide)
	}

	return GetGuideByIDNoCache(guideID, ctx)
}

func GetGuideByIDNoCache(guideID string, ctx *context.Context) *Guide {
	var guide Guide
	DBCtx(ctx).Preload("Tags").Find(&guide, "id = ?", guideID)

	if guide.ID == "" {
		return nil
	}

	cacheKey := "GetGuideById_" + guideID
	dbCache.Set(cacheKey, &guide, cache.DefaultExpiration)

	return &guide
}

func GetGuides(filter *models.GuideFilter, ctx *context.Context) []Guide {
	hash, err := filter.Hash()
	cacheKey := ""
	if err == nil {
		cacheKey = "GetGuides_" + hash
		if guides, ok := dbCache.Get(cacheKey); ok {
			return guides.([]Guide)
		}
	}

	var guides []Guide
	query := DBCtx(ctx).Preload("Tags")

	if filter != nil {
		query = query.Limit(*filter.Limit).
			Offset(*filter.Offset).
			Order(string(*filter.OrderBy) + " " + string(*filter.Order))

		if filter.Search != nil && *filter.Search != "" {
			query = query.Where("to_tsvector(name) @@ to_tsquery(?)", strings.Replace(*filter.Search, " ", " & ", -1))
		}

		if filter.TagIDs != nil && len(filter.TagIDs) > 0 {
			query.Joins("INNER JOIN guide_tags on guide_tags.tag_id in ? AND guide_tags.guide_id = guides.id", filter.TagIDs)
		}
	}

	query.Find(&guides)

	if cacheKey != "" {
		dbCache.Set(cacheKey, guides, cache.DefaultExpiration)
	}

	return guides
}

func GetGuidesByID(guideIds []string, ctx *context.Context) []Guide {
	cacheKey := "GetGuidesById_" + strings.Join(guideIds, ":")

	if guides, ok := dbCache.Get(cacheKey); ok {
		return guides.([]Guide)
	}

	var guides []Guide
	DBCtx(ctx).Preload("Tags").Find(&guides, "id in (?)", guideIds)

	if len(guideIds) != len(guides) {
		return nil
	}

	dbCache.Set(cacheKey, guides, cache.DefaultExpiration)

	return guides
}

func GetGuideCount(filter *models.GuideFilter, ctx *context.Context) int64 {
	hash, err := filter.Hash()
	cacheKey := ""
	if err == nil {
		cacheKey = "GetGuideCount_" + hash
		if count, ok := dbCache.Get(cacheKey); ok {
			return count.(int64)
		}
	}

	var guideCount int64
	query := DBCtx(ctx).Model(Guide{})

	if filter != nil {
		if filter.Search != nil && *filter.Search != "" {
			query = query.Where("to_tsvector(name) @@ to_tsquery(?)", strings.Replace(*filter.Search, " ", " & ", -1))
		}
	}

	query.Count(&guideCount)

	if cacheKey != "" {
		dbCache.Set(cacheKey, guideCount, cache.DefaultExpiration)
	}

	return guideCount
}

func IncrementGuideViews(guide *Guide, ctx *context.Context) {
	DBCtx(ctx).Model(guide).Update("views", guide.Views+1)
}

func GetUserGuides(userID string, ctx *context.Context) []Guide {
	var guides []Guide
	DBCtx(ctx).Preload("Tags").Find(&guides, "user_id = ?", userID)
	return guides
}

func ClearGuideTags(guideID string, ctx *context.Context) error {
	r := DBCtx(ctx).Where("guide_id = ?", guideID).Delete(&GuideTag{})
	return r.Error
}

func SetGuideTags(guideID string, tagIDs []string, ctx *context.Context) error {
	for _, tag := range tagIDs {
		err := AddGuideTag(guideID, tag, ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func ResetGuideTags(guideID string, tagIDs []string, ctx *context.Context) error {
	err := ClearGuideTags(guideID, ctx)
	if err != nil {
		return err
	}
	err = SetGuideTags(guideID, tagIDs, ctx)
	if err != nil {
		return err
	}
	return nil
}

func AddGuideTag(guideID string, tagID string, ctx *context.Context) error {
	r := DBCtx(ctx).Create(&GuideTag{GuideID: guideID, TagID: tagID})
	return r.Error
}

func RemoveGuideTag(guideID string, tagID string, ctx *context.Context) error {
	r := DBCtx(ctx).Delete(&GuideTag{GuideID: guideID, TagID: tagID})
	return r.Error
}
