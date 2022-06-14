package postgres

import (
	"context"

	"github.com/patrickmn/go-cache"
	"github.com/satisfactorymodding/smr-api/util"
)

func CreateModLink(ctx context.Context, modLink *ModLink) (*ModLink, error) {
	modLink.ID = util.GenerateUniqueID()
	DBCtx(ctx).Create(&modLink)
	return modLink, nil
}

func GetModLinkByID(ctx context.Context, modLinksID string) *ModLink {
	cacheKey := "GetModLinkByID_" + modLinksID

	if modLink, ok := dbCache.Get(cacheKey); ok {
		return modLink.(*ModLink)
	}

	var modLink ModLink
	DBCtx(ctx).Find(&modLink, "id = ?", modLinksID)

	if modLink.ID == "" {
		return nil
	}

	dbCache.Set(cacheKey, modLink, cache.DefaultExpiration)

	return &modLink
}

func GetModLinks(ctx context.Context) []ModLink {
	cacheKey := "GetModLinks"

	if modLink, ok := dbCache.Get(cacheKey); ok {
		return modLink.([]ModLink)
	}

	var modLink []ModLink
	DBCtx(ctx).Find(&modLink)

	dbCache.Set(cacheKey, modLink, cache.DefaultExpiration)

	return modLink
}

func GetLinksByMod(ctx context.Context, importance string) []ModLink {
	cacheKey := "GetLinksByMod_" + importance

	if modLink, ok := dbCache.Get(cacheKey); ok {
		return modLink.([]ModLink)
	}

	var modLink []ModLink
	DBCtx(ctx).Find(&modLink, "importance = ?", importance)

	dbCache.Set(cacheKey, modLink, cache.DefaultExpiration)

	return modLink
}

func GetModLink(ctx context.Context, versionID string, platform string) *ModLink {
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
