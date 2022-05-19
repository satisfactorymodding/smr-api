package postgres

import (
	"context"
	"strings"

	"github.com/patrickmn/go-cache"
	"github.com/satisfactorymodding/smr-api/models"
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

func GetSMLLinks(ctx context.Context, filter *models.SMLLinkFilter) []SMLLink {
	var smlLinks []SMLLink
	query := DBCtx(ctx)

	if filter != nil {
		query = query.Limit(*filter.Limit).
			Offset(*filter.Offset).
			Order(string(*filter.OrderBy) + " " + string(*filter.Order))

		if filter.Search != nil && *filter.Search != "" {
			query = query.Where("to_tsvector(name) @@ to_tsquery(?)", strings.Replace(*filter.Search, " ", " & ", -1))
		}
	}

	query.Find(&smlLinks)
	return smlLinks
}
