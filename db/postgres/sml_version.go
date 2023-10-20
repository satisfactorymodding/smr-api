package postgres

import (
	"context"
	"strings"

	"github.com/satisfactorymodding/smr-api/models"
	"github.com/satisfactorymodding/smr-api/util"
)

func CreateSMLVersion(ctx context.Context, smlVersion *SMLVersion) (*SMLVersion, error) {
	smlVersion.ID = util.GenerateUniqueID()

	DBCtx(ctx).Create(&smlVersion)

	return smlVersion, nil
}

func GetSMLVersionByID(ctx context.Context, smlVersionID string) *SMLVersion {
	var smlVersion SMLVersion
	DBCtx(ctx).Preload("Targets").Find(&smlVersion, "id in (?)", smlVersionID)

	if smlVersion.ID == "" {
		return nil
	}

	return &smlVersion
}

func GetSMLVersions(ctx context.Context, filter *models.SMLVersionFilter) []SMLVersion {
	var smlVersions []SMLVersion
	query := DBCtx(ctx)

	if filter != nil {
		query = query.Limit(*filter.Limit).
			Offset(*filter.Offset).
			Order(string(*filter.OrderBy) + " " + string(*filter.Order))

		if filter.Search != nil && *filter.Search != "" {
			query = query.Where("to_tsvector(name) @@ to_tsquery(?)", strings.ReplaceAll(*filter.Search, " ", " & "))
		}
	}

	query.Preload("Targets").Find(&smlVersions)

	return smlVersions
}

func GetSMLVersionsByID(ctx context.Context, smlVersionIds []string) []SMLVersion {
	var smlVersions []SMLVersion
	DBCtx(ctx).Preload("Targets").Find(&smlVersions, "id in (?)", smlVersionIds)

	if len(smlVersionIds) != len(smlVersions) {
		return nil
	}

	return smlVersions
}

func GetSMLVersionCount(ctx context.Context, filter *models.SMLVersionFilter) int64 {
	var smlVersionCount int64
	query := DBCtx(ctx).Model(SMLVersion{})

	if filter != nil {
		if filter.Search != nil && *filter.Search != "" {
			query = query.Where("to_tsvector(name) @@ to_tsquery(?)", strings.ReplaceAll(*filter.Search, " ", " & "))
		}
	}

	query.Count(&smlVersionCount)
	return smlVersionCount
}

func GetSMLLatestVersions(ctx context.Context) *[]SMLVersion {
	var smlVersions []SMLVersion

	DBCtx(ctx).Preload("Targets").Select("distinct on (stability) *").
		Order("stability, created_at desc").
		Find(&smlVersions)

	return &smlVersions
}

func GetSMLVersionTargets(ctx context.Context, smlVersionID string) []SMLVersionTarget {
	var smlVersionTargets []SMLVersionTarget

	DBCtx(ctx).Find(&smlVersionTargets, "version_id = ?", smlVersionID)

	return smlVersionTargets
}
