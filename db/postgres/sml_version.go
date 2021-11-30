package postgres

import (
	"context"
	"strings"

	"github.com/satisfactorymodding/smr-api/models"
	"github.com/satisfactorymodding/smr-api/util"
)

func CreateSMLVersion(smlVersion *SMLVersion, ctx *context.Context) (*SMLVersion, error) {
	smlVersion.ID = util.GenerateUniqueID()

	DBCtx(ctx).Create(&smlVersion)

	return smlVersion, nil
}

func GetSMLVersionByID(smlVersionID string, ctx *context.Context) *SMLVersion {
	var smlVersion SMLVersion
	DBCtx(ctx).Find(&smlVersion, "id = ?", smlVersionID)

	if smlVersion.ID == "" {
		return nil
	}

	return &smlVersion
}

func GetSMLVersions(filter *models.SMLVersionFilter, ctx *context.Context) []SMLVersion {
	var smlVersions []SMLVersion
	query := DBCtx(ctx)

	if filter != nil {
		query = query.Limit(*filter.Limit).
			Offset(*filter.Offset).
			Order(string(*filter.OrderBy) + " " + string(*filter.Order))

		if filter.Search != nil && *filter.Search != "" {
			query = query.Where("to_tsvector(name) @@ to_tsquery(?)", strings.Replace(*filter.Search, " ", " & ", -1))
		}
	}

	query.Find(&smlVersions)
	return smlVersions
}

func GetSMLVersionsByID(smlVersionIds []string, ctx *context.Context) []SMLVersion {
	var smlVersions []SMLVersion
	DBCtx(ctx).Find(&smlVersions, "id in (?)", smlVersionIds)

	if len(smlVersionIds) != len(smlVersions) {
		return nil
	}

	return smlVersions
}

func GetSMLVersionCount(filter *models.SMLVersionFilter, ctx *context.Context) int64 {
	var smlVersionCount int64
	query := DBCtx(ctx).Model(SMLVersion{})

	if filter != nil {
		if filter.Search != nil && *filter.Search != "" {
			query = query.Where("to_tsvector(name) @@ to_tsquery(?)", strings.Replace(*filter.Search, " ", " & ", -1))
		}
	}

	query.Count(&smlVersionCount)
	return smlVersionCount
}

func GetSMLLatestVersions(ctx *context.Context) *[]SMLVersion {
	var smlVersions []SMLVersion

	DBCtx(ctx).Select("distinct on (stability) *").
		Order("stability, created_at desc").
		Find(&smlVersions)

	return &smlVersions
}
