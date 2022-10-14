package postgres

import (
	"context"
	"strings"

	"github.com/patrickmn/go-cache"

	"github.com/satisfactorymodding/smr-api/models"
	"github.com/satisfactorymodding/smr-api/util"
)

func CreateSMLArch(ctx context.Context, smlArch *SMLArch) (*SMLArch, error) {
	smlArch.ID = util.GenerateUniqueID()

	DBCtx(ctx).Create(&smlArch)

	return smlArch, nil
}

func GetSMLArch(ctx context.Context, smlLinkID string) *SMLArch {
	cacheKey := "GetSMLArch_" + smlLinkID

	if smlArch, ok := dbCache.Get(cacheKey); ok {
		return smlArch.(*SMLArch)
	}

	var smlArch SMLArch
	DBCtx(ctx).Find(&smlArch, "id = ?", smlLinkID)

	if smlArch.ID == "" {
		return nil
	}

	dbCache.Set(cacheKey, &smlArch, cache.DefaultExpiration)

	return &smlArch
}

func GetSMLArchs(ctx context.Context, filter *models.SMLArchFilter) []SMLArch {
	var smlLinks []SMLArch
	query := DBCtx(ctx)

	if filter != nil {
		query = query.Limit(*filter.Limit).
			Offset(*filter.Offset).
			Order(string(*filter.OrderBy) + " " + string(*filter.Order))

		if filter.Search != nil && *filter.Search != "" {
			query = query.Where("to_tsvector(name) @@ to_tsquery(?)", strings.ReplaceAll(*filter.Search, " ", " & "))
		}
	}

	query.Find(&smlLinks)
	return smlLinks
}

func GetSMLArchByID(ctx context.Context, smlLinkID string) []SMLArch {
	var smlArchs []SMLArch

	DBCtx(ctx).Find(&smlArchs, "id in ?", smlLinkID)

	if len(smlArchs) != 0 {
		return nil
	}

	return smlArchs
}

func GetSMLArchsByID(ctx context.Context, smlArchIds []string) []SMLArch {
	var smlArchs []SMLArch

	DBCtx(ctx).Find(&smlArchs, "id in (?)", smlArchIds)

	if len(smlArchIds) != len(smlArchs) {
		return nil
	}

	return smlArchs
}

func GetSMLArchBySMLID(ctx context.Context, smlVersionID string) []SMLArch {
	var smlArchs []SMLArch

	DBCtx(ctx).Find(&smlArchs, "sml_version_arch_id = ?", smlVersionID)

	return smlArchs
}

func GetSMLArchDownload(ctx context.Context, smlVersionID string, platform string) string {
	var smlPlatform SMLArch
	DBCtx(ctx).First(&smlPlatform, "sml_version_arch_id = ? AND platform = ?", smlVersionID, platform)

	if smlPlatform.ID == "" {
		return ""
	}

	return smlPlatform.Link
}
