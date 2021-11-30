package utils

import (
	"context"

	"github.com/satisfactorymodding/smr-api/db/postgres"
	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/models"
	"github.com/satisfactorymodding/smr-api/redis/jobs"
)

func ReindexAllModFiles(ctx context.Context, withMetadata bool, modFilter func(postgres.Mod) bool, versionFilter func(version postgres.Version) bool) {
	offset := 0

	limit := 100
	createdAt := generated.VersionFieldsCreatedAt
	orderDesc := generated.OrderDesc

	for {
		mods := postgres.GetMods(100, offset, "created_at", "asc", "", false, &ctx)
		offset += 100

		if len(mods) == 0 {
			break
		}

		for _, mod := range mods {
			versionOffset := 0

			if modFilter != nil {
				if !modFilter(mod) {
					continue
				}
			}

			for {
				versions := postgres.GetModVersionsNew(mod.ID, &models.VersionFilter{
					Limit:   &limit,
					Offset:  &versionOffset,
					OrderBy: &createdAt,
					Order:   &orderDesc,
				}, false, &ctx)

				versionOffset += len(versions)

				if len(versions) > 0 {
					for _, version := range versions {
						if versionFilter != nil {
							if !versionFilter(version) {
								continue
							}
						}

						if withMetadata {
							jobs.SubmitJobUpdateDBFromModVersionFileTask(ctx, mod.ID, version.ID)
						} else {
							jobs.SubmitJobUpdateDBFromModVersionJSONFileTask(ctx, mod.ID, version.ID)
						}
					}
				} else {
					break
				}
			}
		}
	}
}
