package postgres

import (
	"context"
	"strings"

	"github.com/patrickmn/go-cache"
	"github.com/satisfactorymodding/smr-api/models"
	"github.com/satisfactorymodding/smr-api/util"
)

func CreateModLink(ctx context.Context, modLink *ModLink) (*ModLink, error) {
	modLink.ID = util.GenerateUniqueID()
	DBCtx(ctx).Create(&modLink)
	return modLink, nil
}

func GetModLink(ctx context.Context, modLinkID string) *ModLink {
	cacheKey := "GetModLink_" + modLinkID

	if modLink, ok := dbCache.Get(cacheKey); ok {
		return modLink.(*ModLink)
	}

	var modLink ModLink
	DBCtx(ctx).Find(&modLink, "id = ?", modLinkID)

	if modLink.ID == "" {
		return nil
	}

	dbCache.Set(cacheKey, modLink, cache.DefaultExpiration)

	return &modLink
}

func GetModLinks(ctx context.Context, filter *models.ModLinkFilter) []ModLink {
	var modLinks []ModLink
	query := DBCtx(ctx)

	if filter != nil {
		query = query.Limit(*filter.Limit).
			Offset(*filter.Offset).
			Order(string(*filter.OrderBy) + " " + string(*filter.Order))

		if filter.Search != nil && *filter.Search != "" {
			query = query.Where("to_tsvector(name) @@ to_tsquery(?)", strings.Replace(*filter.Search, " ", " & ", -1))
		}
	}

	query.Find(&modLinks)
	return modLinks
}

func GetModLinkByID(ctx context.Context, modLinkID string) *ModLink {
	cacheKey := "GetModLink_" + modLinkID

	if modLink, ok := dbCache.Get(cacheKey); ok {
		return modLink.(*ModLink)
	}

	var modLink ModLink
	DBCtx(ctx).Find(&modLink, "id = ?", modLinkID)

	if modLink.ID == "" {
		return nil
	}

	dbCache.Set(cacheKey, modLink, cache.DefaultExpiration)

	return &modLink
}

func GetModLinksByID(ctx context.Context, modLinkIds []string) []ModLink {
	var modLinks []ModLink

	DBCtx(ctx).Find(&modLinks, "id in (?)", modLinkIds)

	if len(modLinkIds) != len(modLinks) {
		return nil
	}

	return modLinks
}

func GetModLinkDownload(ctx context.Context, versionID string, platform string) *ModLink {
	cacheKey := "GetModLink_" + versionID + "_" + platform
	if modplatform, ok := dbCache.Get(cacheKey); ok {
		return modplatform.(*ModLink)
	}

	var modplatform ModLink
	DBCtx(ctx).First(&modplatform, "mod_version_link_id = ? AND platform = ?", versionID, platform)

	if modplatform.ModVersionLinkID == "" {
		return nil
	}

	dbCache.Set(cacheKey, &modplatform, cache.DefaultExpiration)

	return &modplatform
}
