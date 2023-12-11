package postgres

import (
	"context"
)

func GetSMLLatestVersions(ctx context.Context) *[]SMLVersion {
	var smlVersions []SMLVersion

	DBCtx(ctx).Preload("Targets").Select("distinct on (stability) *").
		Order("stability, created_at desc").
		Find(&smlVersions)

	return &smlVersions
}
