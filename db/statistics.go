package db

import (
	"context"
	"log/slog"
	"regexp"
	"time"

	"github.com/Vilsol/slox"

	"github.com/satisfactorymodding/smr-api/db/postgres"
	"github.com/satisfactorymodding/smr-api/redis"
)

var keyRegex = regexp.MustCompile(`^([^:]+):([^:]+):([^:]+):([^:]+)$`)

func RunAsyncStatisticLoop(ctx context.Context) {
	go func() {
		for {
			start := time.Now()
			keys := redis.GetAllKeys()
			slox.Info(ctx, "statistics fetched", slog.Int("keys", len(keys)), slog.Duration("took", time.Since(start)))
			resultMap := make(map[string]map[string]map[string]uint)
			for _, key := range keys {
				if matches := keyRegex.FindStringSubmatch(key); matches != nil {
					entityType := matches[1]
					entityID := matches[2]
					action := matches[3]

					if _, ok := resultMap[entityType]; !ok {
						resultMap[entityType] = make(map[string]map[string]uint)
					}

					if _, ok := resultMap[entityType][action]; !ok {
						resultMap[entityType][action] = make(map[string]uint)
					}

					resultMap[entityType][action][entityID]++
				}
			}

			for entityType, entityValue := range resultMap {
				for action, actionValue := range entityValue {
					for entityID, count := range actionValue {
						updateTx := postgres.DBCtx(ctx).Begin()
						ctxWithTx := postgres.ContextWithDB(ctx, updateTx)
						switch entityType {
						case "mod":
							if action == "view" {
								mod := postgres.GetModByID(ctxWithTx, entityID)
								if mod != nil {
									currentHotness := mod.Hotness
									if currentHotness > 4 {
										// Preserve some of the hotness
										currentHotness /= 4
									}
									updateTx.Model(&mod).UpdateColumns(postgres.Mod{Hotness: currentHotness + count})
								}
							}
						case "version":
							if action == "download" {
								version := postgres.GetVersion(ctxWithTx, entityID)
								if version != nil {
									currentHotness := version.Hotness
									if currentHotness > 4 {
										// Preserve some of the popularity
										currentHotness /= 4
									}
									updateTx.Model(&version).UpdateColumns(postgres.Version{Hotness: currentHotness + count})
								}
							}
						}
						updateTx.Commit()
					}
				}
			}

			type Result struct {
				ModID     string
				Hotness   uint
				Downloads uint
			}

			var resultRows []Result

			postgres.DBCtx(ctx).Raw("SELECT mod_id, SUM(hotness) AS hotness, SUM(downloads) AS downloads FROM versions GROUP BY mod_id").Scan(&resultRows)

			for _, row := range resultRows {
				updateTx := postgres.DBCtx(ctx).Begin()
				ctxWithTx := postgres.ContextWithDB(ctx, updateTx)
				mod := postgres.GetModByID(ctxWithTx, row.ModID)
				if mod != nil {
					currentPopularity := mod.Popularity
					if currentPopularity > 4 {
						// Preserve some of the popularity
						currentPopularity /= 4
					}
					updateTx.Model(&mod).UpdateColumns(postgres.Mod{
						Popularity: currentPopularity + row.Hotness,
						Downloads:  row.Downloads,
					})
				}
				updateTx.Commit()
			}

			slox.Info(ctx, "statistics updated", slog.Duration("took", time.Since(start)))
			time.Sleep(time.Minute)
		}
	}()
}
