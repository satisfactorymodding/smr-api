package postgres

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"

	"github.com/satisfactorymodding/smr-api/models"
	"github.com/satisfactorymodding/smr-api/util"
)

var semverCheck = regexp.MustCompile(`^(<=|<|>|>=|\^)?(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)

func GetVersionsByID(ctx context.Context, versionIds []string) []Version {
	cacheKey := "GetVersionsById_" + strings.Join(versionIds, ":")
	if versions, ok := dbCache.Get(cacheKey); ok {
		return versions.([]Version)
	}

	var versions []Version
	DBCtx(ctx).Preload("Targets").Find(&versions, "id in (?)", versionIds)

	if len(versionIds) != len(versions) {
		return nil
	}

	dbCache.Set(cacheKey, versions, cache.DefaultExpiration)

	return versions
}

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

func GetModVersionByName(ctx context.Context, modID string, versionName string) *Version {
	cacheKey := "GetModVersionByName_" + modID + "_" + versionName
	if version, ok := dbCache.Get(cacheKey); ok {
		return version.(*Version)
	}

	var version Version
	DBCtx(ctx).Preload("Targets").First(&version, "mod_id = ? AND version = ?", modID, versionName)

	if version.ID == "" {
		return nil
	}

	dbCache.Set(cacheKey, &version, cache.DefaultExpiration)

	return &version
}

func CreateVersion(ctx context.Context, version *Version) error {
	var versionCount int64
	DBCtx(ctx).Model(Version{}).Where("mod_id = ? AND version = ?", version.ModID, version.Version).Count(&versionCount)

	if versionCount > 0 {
		return errors.New("this mod already has a version with this name")
	}

	// Allow only new 5 versions per 24h

	var versions []Version
	DBCtx(ctx).Order("created_at asc").Find(&versions, "mod_id = ? AND created_at > ?", version.ModID, time.Now().Add(time.Hour*24*-1))

	if len(versions) >= 5 {
		timeToWait := time.Until(versions[0].CreatedAt.Add(time.Hour * 24)).Minutes()
		return fmt.Errorf("please wait %.0f minutes to post another version", timeToWait)
	}

	version.ID = util.GenerateUniqueID()
	DBCtx(ctx).Create(&version)

	return nil
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

func GetVersionsNew(ctx context.Context, filter *models.VersionFilter, unapproved bool) []Version {
	hash, err := filter.Hash()
	cacheKey := ""
	if err == nil {
		cacheKey = "GetVersionsNew_" + hash + "_" + fmt.Sprint(unapproved)
		if versions, ok := dbCache.Get(cacheKey); ok {
			return versions.([]Version)
		}
	}

	var versions []Version
	query := DBCtx(ctx).Preload("Targets").Where("approved = ? AND denied = ?", !unapproved, false)

	if filter != nil {
		query = query.Limit(*filter.Limit).
			Offset(*filter.Offset).
			Order(string(*filter.OrderBy) + " " + string(*filter.Order))

		if filter.Search != nil && *filter.Search != "" {
			query = query.Where("to_tsvector(version) @@ to_tsquery(?)", strings.ReplaceAll(*filter.Search, " ", " & "))
		}

		if filter.Fields != nil && len(filter.Fields) > 0 {
			query = query.Select(filter.Fields)
		}
	}

	query.Preload("Targets").Find(&versions)

	if cacheKey != "" {
		dbCache.Set(cacheKey, versions, cache.DefaultExpiration)
	}

	return versions
}

func GetVersionCountNew(ctx context.Context, filter *models.VersionFilter, unapproved bool) int64 {
	hash, err := filter.Hash()
	cacheKey := ""
	if err == nil {
		cacheKey = "GetVersionCountNew_" + hash + "_" + fmt.Sprint(unapproved)
		if versionCount, ok := dbCache.Get(cacheKey); ok {
			return versionCount.(int64)
		}
	}

	var versionCount int64
	query := DBCtx(ctx).Model(Version{}).Where("approved = ? AND denied = ?", !unapproved, false)

	if filter != nil {
		if filter.Search != nil && *filter.Search != "" {
			query = query.Where("to_tsvector(version) @@ to_tsquery(?)", strings.ReplaceAll(*filter.Search, " ", " & "))
		}
	}

	query.Count(&versionCount)

	if cacheKey != "" {
		dbCache.Set(cacheKey, versionCount, cache.DefaultExpiration)
	}

	return versionCount
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

func GetModVersionsConstraint(ctx context.Context, modID string, constraint string) []Version {
	matches := semverCheck.FindAllStringSubmatch(constraint, -1)
	if len(matches) == 0 {
		return nil
	}

	major, err := strconv.Atoi(matches[0][2])
	if err != nil {
		return nil
	}

	minor, err := strconv.Atoi(matches[0][3])
	if err != nil {
		return nil
	}

	patch, err := strconv.Atoi(matches[0][4])
	if err != nil {
		return nil
	}

	query := DBCtx(ctx).Preload("Targets").Where("mod_id", modID)

	/*
		<=1.2.3
		major < 1
		major = 1, minor < 2
		major = 1, minor = 2, patch <= 3

		<1.2.3
		major < 1
		major = 1, minor < 2
		major = 1, minor = 2, patch < 3

		>1.2.3
		major > 1
		major = 1, minor > 2
		major = 1, minor = 2, patch > 3

		>=1.2.3
		major > 1
		major = 1, minor > 2
		major = 1, minor = 2, patch >= 3

		1.2.3
		major = 1, minor = 2, patch = 3

		^1.2.3 (>=1.2.3, <2.0.0)
		major = 1, minor > 2
		major = 1, minor = 2, patch >= 3

		^0.2.3 (>=0.2.3, <0.3.0)
		major = 0, minor = 2, patch >= 3

		^0.0.3 (>=0.0.3, <0.0.4)
		major = 0, minor = 0, patch = 3
	*/

	sign := matches[0][1]
	switch sign {
	case "<=":
		query = query.Where(db.Or("version_major < ?", major).
			Or("version_major = ? AND version_minor < ?", major, minor).
			Or("version_major = ? AND version_minor = ? AND version_patch <= ?", major, minor, patch))
	case "<":
		query = query.Where(db.Or("version_major < ?", major).
			Or("version_major = ? AND version_minor < ?", major, minor).
			Or("version_major = ? AND version_minor = ? AND version_patch < ?", major, minor, patch))
	case ">":
		query = query.Where(db.Or("version_major > ?", major).
			Or("version_major = ? AND version_minor > ?", major, minor).
			Or("version_major = ? AND version_minor = ? AND version_patch > ?", major, minor, patch))
	case ">=":
		query = query.Where(db.Or("version_major > ?", major).
			Or("version_major = ? AND version_minor > ?", major, minor).
			Or("version_major = ? AND version_minor = ? AND version_patch >= ?", major, minor, patch))
	case "^":
		if major != 0 {
			query = query.Where(db.Or("version_major = ? AND version_minor > ?", major, minor).
				Or("version_major = ? AND version_minor = ? AND version_patch >= ?", major, minor, patch))
		} else if minor != 0 {
			query = query.Where("version_major = ? AND version_minor = ? AND version_patch >= ?", major, minor, patch)
		} else {
			query = query.Where("version_major = ? AND version_minor = ? AND version_patch = ?", major, minor, patch)
		}
	default:
		query = query.Where("version_major = ? AND version_minor = ? AND version_patch = ?", major, minor, patch)
	}

	var versions []Version
	query.Preload("Targets").Find(&versions)
	return versions
}
