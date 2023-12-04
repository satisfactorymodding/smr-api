package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/patrickmn/go-cache"

	"github.com/satisfactorymodding/smr-api/models"
)

func GetModsLatestVersions(ctx context.Context, modIds []string, unapproved bool) *[]Version {
	cacheKey := "GetModsLatestVersions_" + strings.Join(modIds, ":") + "_" + fmt.Sprint(unapproved)
	if versions, ok := dbCache.Get(cacheKey); ok {
		return versions.(*[]Version)
	}

	var versions []Version

	DBCtx(ctx).Preload("Targets").Select("distinct on (mod_id, stability) *").
		Where("mod_id in (?)", modIds).
		Where("approved = ? AND denied = ?", !unapproved, false).
		Order("mod_id, stability, created_at desc").
		Find(&versions)

	dbCache.Set(cacheKey, &versions, cache.DefaultExpiration)

	return &versions
}

func GetModVersions(ctx context.Context, modID string, limit int, offset int, orderBy string, order string, unapproved bool) []Version {
	cacheKey := "GetModVersions_" + modID + "_" + fmt.Sprint(limit) + "_" + fmt.Sprint(offset) + "_" + orderBy + "_" + order + "_" + fmt.Sprint(unapproved)
	if versions, ok := dbCache.Get(cacheKey); ok {
		return versions.([]Version)
	}

	var versions []Version
	DBCtx(ctx).Preload("Targets").Limit(limit).Offset(offset).Order(orderBy+" "+order).Where("approved = ? AND denied = ?", !unapproved, false).Find(&versions, "mod_id = ?", modID)

	dbCache.Set(cacheKey, versions, cache.DefaultExpiration)

	return versions
}

func GetAllModVersionsWithDependencies(ctx context.Context, modID string) []TinyVersion {
	cacheKey := "GetAllModVersionsWithDependencies_" + modID
	if versions, ok := dbCache.Get(cacheKey); ok {
		return versions.([]TinyVersion)
	}

	var versions []TinyVersion
	DBCtx(ctx).
		Preload("Dependencies").
		Preload("Targets").
		Where("approved = ? AND denied = ?", true, false).
		Find(&versions, "mod_id = ?", modID)

	dbCache.Set(cacheKey, versions, cache.DefaultExpiration)

	return versions
}

func GetModVersionsNew(ctx context.Context, modID string, filter *models.VersionFilter, unapproved bool) []Version {
	hash, err := filter.Hash()
	cacheKey := ""
	if err == nil {
		cacheKey = "GetModVersionsNew_" + modID + "_" + hash + "_" + fmt.Sprint(unapproved)
		if versions, ok := dbCache.Get(cacheKey); ok {
			return versions.([]Version)
		}
	}

	var versions []Version
	query := DBCtx(ctx).Preload("Targets")

	if filter != nil {
		query = query.Limit(*filter.Limit).
			Offset(*filter.Offset).
			Order(string(*filter.OrderBy) + " " + string(*filter.Order))
	}

	query.Preload("Targets").Where("approved = ? AND denied = ?", !unapproved, false).Find(&versions, "mod_id = ?", modID)

	if cacheKey != "" {
		dbCache.Set(cacheKey, versions, cache.DefaultExpiration)
	}

	return versions
}

func GetModVersion(ctx context.Context, modID string, versionID string) *Version {
	cacheKey := "GetModVersion_" + modID + "_" + versionID
	if version, ok := dbCache.Get(cacheKey); ok {
		return version.(*Version)
	}

	var version Version
	DBCtx(ctx).Preload("Targets").First(&version, "mod_id = ? AND id = ?", modID, versionID)

	if version.ID == "" {
		return nil
	}

	dbCache.Set(cacheKey, &version, cache.DefaultExpiration)

	return &version
}

func IncrementVersionDownloads(ctx context.Context, version *Version) {
	DBCtx(ctx).Model(version).Update("downloads", version.Downloads+1)
}

func GetVersion(ctx context.Context, versionID string) *Version {
	cacheKey := "GetVersion_" + versionID
	if version, ok := dbCache.Get(cacheKey); ok {
		return version.(*Version)
	}

	var version Version
	DBCtx(ctx).Preload("Targets").First(&version, "id = ?", versionID)

	if version.ID == "" {
		return nil
	}

	dbCache.Set(cacheKey, &version, cache.DefaultExpiration)

	return &version
}

func GetVersionTarget(ctx context.Context, versionID string, target string) *VersionTarget {
	cacheKey := "GetVersionTarget_" + versionID + "_" + target
	if versionTarget, ok := dbCache.Get(cacheKey); ok {
		return versionTarget.(*VersionTarget)
	}

	var versionTarget VersionTarget
	DBCtx(ctx).First(&versionTarget, "version_id = ? AND target_name = ?", versionID, target)

	if versionTarget.VersionID == "" {
		return nil
	}

	dbCache.Set(cacheKey, &versionTarget, cache.DefaultExpiration)

	return &versionTarget
}
