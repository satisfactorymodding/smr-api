package postgres

import (
	"context"
	"strings"

	"github.com/satisfactorymodding/smr-api/models"
)

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

func GetSMLLatestVersions(ctx context.Context) *[]SMLVersion {
	var smlVersions []SMLVersion

	DBCtx(ctx).Preload("Targets").Select("distinct on (stability) *").
		Order("stability, created_at desc").
		Find(&smlVersions)

	return &smlVersions
}
