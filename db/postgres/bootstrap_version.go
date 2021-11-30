package postgres

import (
	"context"
	"strings"

	"github.com/satisfactorymodding/smr-api/models"
	"github.com/satisfactorymodding/smr-api/util"
)

func CreateBootstrapVersion(bootstrapVersion *BootstrapVersion, ctx *context.Context) (*BootstrapVersion, error) {
	bootstrapVersion.ID = util.GenerateUniqueID()

	DBCtx(ctx).Create(&bootstrapVersion)

	return bootstrapVersion, nil
}

func GetBootstrapVersionByID(bootstrapVersionID string, ctx *context.Context) *BootstrapVersion {
	var bootstrapVersion BootstrapVersion
	DBCtx(ctx).Find(&bootstrapVersion, "id = ?", bootstrapVersionID)

	if bootstrapVersion.ID == "" {
		return nil
	}

	return &bootstrapVersion
}

func GetBootstrapVersions(filter *models.BootstrapVersionFilter, ctx *context.Context) []BootstrapVersion {
	var bootstrapVersions []BootstrapVersion
	query := DBCtx(ctx)

	if filter != nil {
		query = query.Limit(*filter.Limit).
			Offset(*filter.Offset).
			Order(string(*filter.OrderBy) + " " + string(*filter.Order))

		if filter.Search != nil && *filter.Search != "" {
			query = query.Where("to_tsvector(name) @@ to_tsquery(?)", strings.Replace(*filter.Search, " ", " & ", -1))
		}
	}

	query.Find(&bootstrapVersions)
	return bootstrapVersions
}

func GetBootstrapVersionsByID(bootstrapVersionIds []string, ctx *context.Context) []BootstrapVersion {
	var bootstrapVersions []BootstrapVersion
	DBCtx(ctx).Find(&bootstrapVersions, "id in (?)", bootstrapVersionIds)

	if len(bootstrapVersionIds) != len(bootstrapVersions) {
		return nil
	}

	return bootstrapVersions
}

func GetBootstrapVersionCount(filter *models.BootstrapVersionFilter, ctx *context.Context) int64 {
	var bootstrapVersionCount int64
	query := DBCtx(ctx).Model(BootstrapVersion{})

	if filter != nil {
		if filter.Search != nil && *filter.Search != "" {
			query = query.Where("to_tsvector(name) @@ to_tsquery(?)", strings.Replace(*filter.Search, " ", " & ", -1))
		}
	}

	query.Count(&bootstrapVersionCount)
	return bootstrapVersionCount
}

func GetBootstrapLatestVersions(ctx *context.Context) *[]BootstrapVersion {
	var bootstrapVersions []BootstrapVersion

	DBCtx(ctx).Select("distinct on (stability) *").
		Order("stability, created_at desc").
		Find(&bootstrapVersions)

	return &bootstrapVersions
}
