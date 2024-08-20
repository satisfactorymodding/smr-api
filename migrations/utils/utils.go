package utils

import (
	"context"
	"fmt"
	"log/slog"

	"entgo.io/ent/dialect/sql"
	"github.com/Vilsol/slox"
	"go.temporal.io/sdk/client"

	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/db/schema"
	"github.com/satisfactorymodding/smr-api/generated/ent"
	"github.com/satisfactorymodding/smr-api/generated/ent/mod"
	"github.com/satisfactorymodding/smr-api/generated/ent/version"
	"github.com/satisfactorymodding/smr-api/workflows"
)

func ReindexAllModFiles(ctx context.Context, withMetadata bool, modFilter func(*ent.Mod) bool, versionFilter func(version *ent.Version) bool) error {
	return ExecuteOnVersions(ctx, false, modFilter, versionFilter, func(m *ent.Mod, v *ent.Version) {
		if _, err := workflows.Client(ctx).ExecuteWorkflow(ctx, client.StartWorkflowOptions{
			ID:        fmt.Sprintf("update-mod-data-from-storage-%s-%s", m.ID, v.ID),
			TaskQueue: workflows.RepoTaskQueue,
		}, workflows.UpdateModDataFromStorageWorkflow, m.ID, v.ID, withMetadata); err != nil {
			slox.Error(ctx, "failed to start finalization workflow", slog.Any("err", err))
		}
	})
}

func ExecuteOnVersions(ctx context.Context, withDeleted bool, modFilter func(*ent.Mod) bool, versionFilter func(version *ent.Version) bool, f func(mod *ent.Mod, version *ent.Version)) error {
	offset := 0

	if withDeleted {
		ctx = schema.SkipSoftDelete(ctx)
	}

	for {
		q := db.From(ctx).Mod.Query().Limit(100).Offset(offset).Order(mod.ByCreatedAt(sql.OrderDesc()))
		mods, err := q.All(ctx)
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

						f(m, v)
					}
				} else {
					break
				}
			}
		}
	}

	return nil
}
