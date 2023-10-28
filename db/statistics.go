package db

import (
	"context"
	"log/slog"
	"regexp"
	"time"

	"github.com/Vilsol/slox"

	"github.com/satisfactorymodding/smr-api/generated/ent"
	"github.com/satisfactorymodding/smr-api/generated/ent/version"
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
						err := Tx(ctx, func(ctx context.Context, tx *ent.Tx) error {
							switch entityType {
							case "mod":
								if action == "view" {
									mod, err := tx.Mod.Get(ctx, entityID)
									if err != nil {
										return err
									}

									if mod != nil {
										currentHotness := mod.Hotness
										if currentHotness > 4 {
											// Preserve some of the hotness
											currentHotness /= 4
										}

										return mod.Update().SetHotness(currentHotness + count).Exec(ctx)
									}
								}
							case "version":
								if action == "download" {
									version, err := tx.Version.Get(ctx, entityID)
									if err != nil {
										return err
									}

									if version != nil {
										currentHotness := version.Hotness
										if currentHotness > 4 {
											// Preserve some of the popularity
											currentHotness /= 4
										}
										return version.Update().SetHotness(currentHotness + count).Exec(ctx)
									}
								}
							}

							return nil
						}, nil)
						if err != nil {
							slox.From(ctx).Error("failed updating statistics", slog.Any("err", err))
						}
					}
				}
			}

			type Result struct {
				ModID     string `json:"mod_id"`
				Hotness   uint   `json:"hotness"`
				Downloads uint   `json:"downloads"`
			}

			var resultRows []Result

			err := From(ctx).Version.
				Query().
				GroupBy("mod_id").
				Aggregate(
					ent.As(ent.Sum(version.FieldHotness), "hotness"),
					ent.As(ent.Sum(version.FieldDownloads), "downloads"),
				).
				Scan(ctx, &resultRows)
			if err != nil {
				slox.From(ctx).Error("failed summing version data", slog.Any("err", err))
				continue
			}

			for _, row := range resultRows {
				err := Tx(ctx, func(ctx context.Context, tx *ent.Tx) error {
					mod, err := tx.Mod.Get(ctx, row.ModID)
					if err != nil {
						return err
					}

					if mod != nil {
						currentPopularity := mod.Popularity
						if currentPopularity > 4 {
							// Preserve some of the popularity
							currentPopularity /= 4
						}
						return mod.Update().SetPopularity(currentPopularity + row.Hotness).SetDownloads(row.Downloads).Exec(ctx)
					}

					return nil
				}, nil)
				if err != nil {
					slox.From(ctx).Error("failed updating mod data", slog.Any("err", err))
					continue
				}
			}

			slox.Info(ctx, "statistics updated", slog.Duration("took", time.Since(start)))
			time.Sleep(time.Minute)
		}
	}()
}
