package utils

import (
	"context"

	"entgo.io/ent/dialect/sql"

	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/generated/ent"
	"github.com/satisfactorymodding/smr-api/generated/ent/mod"
	"github.com/satisfactorymodding/smr-api/generated/ent/version"
	"github.com/satisfactorymodding/smr-api/redis/jobs"
)

func ReindexAllModFiles(ctx context.Context, withMetadata bool, modFilter func(*ent.Mod) bool, versionFilter func(version *ent.Version) bool) error {
	offset := 0

	for {
		mods, err := db.From(ctx).Mod.Query().Limit(100).Offset(offset).Order(mod.ByCreatedAt(sql.OrderDesc())).All(ctx)
		if err != nil {
			return err
		}
		offset += len(mods)

		if len(mods) == 0 {
			break
		}

		for _, m := range mods {
			versionOffset := 0

			if modFilter != nil {
				if !modFilter(m) {
					continue
				}
			}

			for {
				versions, err := m.QueryVersions().Limit(100).Offset(versionOffset).Order(version.ByCreatedAt(sql.OrderDesc())).All(ctx)
				if err != nil {
					return err
				}

				versionOffset += len(versions)

				if len(versions) > 0 {
					for _, v := range versions {
						if versionFilter != nil {
							if !versionFilter(v) {
								continue
							}
						}

						if withMetadata {
							jobs.SubmitJobUpdateDBFromModVersionFileTask(ctx, m.ID, v.ID)
						} else {
							jobs.SubmitJobUpdateDBFromModVersionJSONFileTask(ctx, m.ID, v.ID)
						}
					}
				} else {
					break
				}
			}
		}
	}

	return nil
}
