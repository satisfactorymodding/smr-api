package postgres

import (
	"context"
	"strings"

	"github.com/patrickmn/go-cache"
	"github.com/satisfactorymodding/smr-api/models"
	"github.com/satisfactorymodding/smr-api/util"
)

func CreateSMLArch(ctx context.Context, smlLink *SMLArch) (*SMLArch, error) {
	smlLink.ID = util.GenerateUniqueID()

	DBCtx(ctx).Create(&smlLink)

	return smlLink, nil
}

func GetSMLArch(ctx context.Context, smlLinkID string) *SMLArch {
	cacheKey := "GetSMLArch_" + smlLinkID

	if smlLink, ok := dbCache.Get(cacheKey); ok {
		return smlLink.(*SMLArch)
	}

	var smlLink SMLArch
	DBCtx(ctx).Find(&smlLink, "id = ?", smlLinkID)

	if smlLink.ID == "" {
		return nil
	}

	dbCache.Set(cacheKey, smlLink, cache.DefaultExpiration)

	return &smlLink
}

func GetSMLArchs(ctx context.Context, filter *models.SMLArchFilter) []SMLArch {
	var smlLinks []SMLArch
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

func GetSMLArchByID(ctx context.Context, smlLinkID string) []SMLArch {
	var smlLinks []SMLArch

	DBCtx(ctx).Find(&smlLinks, "id in ?", smlLinkID)

	if len(smlLinks) != 0 {
		return nil
	}

	return smlLinks
}

func GetSMLArchsByID(ctx context.Context, smlLinkIds []string) []SMLArch {
	var smlLinks []SMLArch

	DBCtx(ctx).Find(&smlLinks, "id in (?)", smlLinkIds)

	if len(smlLinkIds) != len(smlLinks) {
		return nil
	}

	return smlLinks
}

func GetSMLArchDownload(ctx context.Context, smlVersionID string, platform string) string {
	var smlPlatform SMLArch
	DBCtx(ctx).First(&smlPlatform, "sml_version_arch_id = ? AND platform = ?", smlVersionID, platform)

	if smlPlatform.SMLVersionArchID == "" {
		return ""
	}

	return smlPlatform.Link
}
