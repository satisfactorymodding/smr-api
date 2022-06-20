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

func GetSMLLink(ctx context.Context, smlLinkID string) *SMLLink {
	cacheKey := "GetSMLLink_" + smlLinkID

	if smlLink, ok := dbCache.Get(cacheKey); ok {
		return smlLink.(*SMLLink)
	}

	var smlLink SMLLink
	DBCtx(ctx).Find(&smlLink, "id = ?", smlLinkID)

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

func GetSMLLinkByID(ctx context.Context, smlLinkID string) []SMLLink {
	var smlLinks []SMLLink

	DBCtx(ctx).Find(&smlLinks, "id in ?", smlLinkID)

	if len(smlLinks) != 0 {
		return nil
	}

	return smlLinks
}

func GetSMLLinksByID(ctx context.Context, smlLinkIds []string) []SMLLink {
	var smlLinks []SMLLink

	DBCtx(ctx).Find(&smlLinks, "id in (?)", smlLinkIds)

	if len(smlLinkIds) != len(smlLinks) {
		return nil
	}

	return smlLinks
}

func GetSMLLinkDownload(ctx context.Context, smlVersionID string, platform string) string {
	var smlPlatform SMLLink
	DBCtx(ctx).First(&smlPlatform, "sml_version_link_id = ? AND platform = ?", smlVersionID, platform)

	if smlPlatform.SMLVersionLinkID == "" {
		return ""
	}

	return smlPlatform.Link
}
