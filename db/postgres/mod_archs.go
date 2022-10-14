package postgres

import (
	"context"
	"strings"

	"github.com/patrickmn/go-cache"

	"github.com/satisfactorymodding/smr-api/models"
	"github.com/satisfactorymodding/smr-api/util"
)

func CreateModArch(ctx context.Context, modArch *ModArch) (*ModArch, error) {
	modArch.ID = util.GenerateUniqueID()
	DBCtx(ctx).Create(&modArch)
	return modArch, nil
}

func GetModArch(ctx context.Context, modArchID string) *ModArch {
	cacheKey := "GetModArch_" + modArchID

	if modArch, ok := dbCache.Get(cacheKey); ok {
		return modArch.(*ModArch)
	}

	var modArch ModArch
	DBCtx(ctx).Find(&modArch, "id = ?", modArchID)

	if modArch.ID == "" {
		return nil
	}

	dbCache.Set(cacheKey, &modArch, cache.DefaultExpiration)

	return &modArch
}

func GetModArchs(ctx context.Context, filter *models.ModArchFilter) []ModArch {
	var modArchs []ModArch
	query := DBCtx(ctx)

	if filter != nil {
		query = query.Limit(*filter.Limit).
			Offset(*filter.Offset).
			Order(string(*filter.OrderBy) + " " + string(*filter.Order))

		if filter.Search != nil && *filter.Search != "" {
			query = query.Where("to_tsvector(name) @@ to_tsquery(?)", strings.ReplaceAll(*filter.Search, " ", " & "))
		}
	}

	query.Find(&modArchs)
	return modArchs
}

func GetVersionModArchs(ctx context.Context, versionID string) []ModArch {
	var modArchs []ModArch
	query := DBCtx(ctx).Find(&modArchs, "mod_version_arch_id = ?", versionID)

	query.Find(&modArchs)
	return modArchs
}

func GetModArchByID(ctx context.Context, modArchID string) *ModArch {
	cacheKey := "GetModArch_" + modArchID

	if modArch, ok := dbCache.Get(cacheKey); ok {
		return modArch.(*ModArch)
	}

	var modArch ModArch
	DBCtx(ctx).Find(&modArch, "id = ?", modArchID)

	if modArch.ID == "" {
		return nil
	}

	dbCache.Set(cacheKey, &modArch, cache.DefaultExpiration)

	return &modArch
}

func GetModArchsByID(ctx context.Context, modArchIds []string) []ModArch {
	var modArchs []ModArch

	DBCtx(ctx).Find(&modArchs, "id in (?)", modArchIds)

	if len(modArchIds) != len(modArchs) {
		return nil
	}

	return modArchs
}

func GetModArchByPlatform(ctx context.Context, versionID string, platform string) *ModArch {
	cacheKey := "GetModArch_" + versionID + "_" + platform
	if modplatform, ok := dbCache.Get(cacheKey); ok {
		return modplatform.(*ModArch)
	}

	var modplatform ModArch
	DBCtx(ctx).First(&modplatform, "mod_version_arch_id = ? AND platform = ?", versionID, platform)

	if modplatform.ModVersionID == "" {
		return nil
	}

	dbCache.Set(cacheKey, &modplatform, cache.DefaultExpiration)

	return &modplatform
}

func GetModArchDownload(ctx context.Context, versionID string, platform string) string {
	var modPlatform ModArch
	DBCtx(ctx).First(&modPlatform, "mod_version_arch_id = ? AND platform = ?", versionID, platform)

	if modPlatform.ModVersionID == "" {
		return ""
	}

	return modPlatform.Key
}
