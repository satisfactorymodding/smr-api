package postgres

import (
	"context"

	"github.com/patrickmn/go-cache"
	"github.com/satisfactorymodding/smr-api/util"
)

func CreateSMLLink(ctx context.Context, smlLink *SMLLink) (*SMLLink, error) {
	smlLink.ID = util.GenerateUniqueID()
	DBCtx(ctx).Create(&smlLink)
	return smlLink, nil
}

func GetSMLLinkByID(ctx context.Context, smlLinksId string) *SMLLink {
	cacheKey := "GetSMLLinkByID_" + smlLinksId

	if smlLink, ok := dbCache.Get(cacheKey); ok {
		return smlLink.(*SMLLink)
	}

	var smlLink SMLLink
	DBCtx(ctx).Find(&smlLink, "id = ?", smlLinksId)

	if smlLink.ID == "" {
		return nil
	}

	dbCache.Set(cacheKey, smlLink, cache.DefaultExpiration)

	return &smlLink
}

func GetSMLLinks(ctx context.Context) []SMLLink {
	cacheKey := "GetSMLLinks"

	if smlLink, ok := dbCache.Get(cacheKey); ok {
		return smlLink.([]SMLLink)
	}

	var smlLink []SMLLink
	DBCtx(ctx).Find(&smlLink)

	dbCache.Set(cacheKey, smlLink, cache.DefaultExpiration)

	return smlLink
}
