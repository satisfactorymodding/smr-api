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

func GetModLinkByID(ctx context.Context, modLinksId string) *ModLink {
	cacheKey := "GetModLinkByID_" + modLinksId

	if modLink, ok := dbCache.Get(cacheKey); ok {
		return modLink.(*ModLink)
	}

	var modLink ModLink
	DBCtx(ctx).Find(&modLink, "id = ?", modLinksId)

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
