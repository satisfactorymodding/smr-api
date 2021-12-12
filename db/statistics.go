package db

import (
	"context"
	"regexp"
	"time"

	"github.com/satisfactorymodding/smr-api/db/postgres"
	"github.com/satisfactorymodding/smr-api/redis"

	"github.com/rs/zerolog/log"
)

var keyRegex = regexp.MustCompile(`^([^:]+):([^:]+):([^:]+):([^:]+)$`)

func RunAsyncStatisticLoop(ctx context.Context) {
	go func() {
		for {
			start := time.Now()
			keys := redis.GetAllKeys()
			log.Ctx(ctx).Info().Msgf("Fetched: %d keys in %s", len(keys), time.Since(start).String())
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

					resultMap[entityType][action][entityID] = resultMap[entityType][action][entityID] + 1
				}
			}

			updateTx := postgres.DBCtx(ctx).Begin()

			for entityType, entityValue := range resultMap {
				for action, actionValue := range entityValue {
					for entityID, count := range actionValue {
						switch entityType {
						case "mod":
							switch action {
							case "view":
								mod := postgres.GetModByID(ctx, entityID)
								if mod != nil {
									currentHotness := mod.Hotness
									if currentHotness > 4 {
										// Preserve some of the hotness
										currentHotness = currentHotness / 4
									}
									updateTx.Model(&mod).UpdateColumns(postgres.Mod{Hotness: currentHotness + count})
								}
							}
						case "version":
							switch action {
							case "download":
								version := postgres.GetVersion(ctx, entityID)
								if version != nil {
									currentHotness := version.Hotness
									if currentHotness > 4 {
										// Preserve some of the popularity
										currentHotness = currentHotness / 4
									}
									updateTx.Model(&version).UpdateColumns(postgres.Version{Hotness: currentHotness + count})
								}
							}
						}
					}
				}
			}

			updateTx.Commit()
			updateTx = postgres.DBCtx(ctx).Begin()

			type Result struct {
				ModID     string
				Hotness   uint
				Downloads uint
			}

			var resultRows []Result

			postgres.DBCtx(ctx).Raw("SELECT mod_id, SUM(hotness) AS hotness, SUM(downloads) AS downloads FROM versions GROUP BY mod_id").Scan(&resultRows)

			for _, row := range resultRows {
				mod := postgres.GetModByID(ctx, row.ModID)
				if mod != nil {
					currentPopularity := mod.Popularity
					if currentPopularity > 4 {
						// Preserve some of the popularity
						currentPopularity = currentPopularity / 4
					}
					updateTx.Model(&mod).UpdateColumns(postgres.Mod{
						Popularity: currentPopularity + row.Hotness,
						Downloads:  row.Downloads,
					})
				}
			}

			updateTx.Commit()

			log.Ctx(ctx).Info().Msgf("Statistics Updated! Took %s", time.Since(start).String())
			time.Sleep(time.Minute)
		}
	}()
}
