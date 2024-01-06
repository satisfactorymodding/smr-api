package dataloader

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/patrickmn/go-cache"

	"github.com/satisfactorymodding/smr-api/db/postgres"
)

type loadersKey struct{}

type Loaders struct {
	VersionDependenciesByVersionID VersionDependencyLoader
	UserModsByModID                UserModLoader
	VersionsByModID                VersionLoader
	VersionsByModIDNoMeta          VersionLoaderNoMeta
	UserByID                       UserLoader
}

func Middleware() func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
	dbCache := cache.New(time.Second*5, time.Second*10)

	return func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := context.WithValue(c.Request().Context(), loadersKey{}, &Loaders{
				VersionDependenciesByVersionID: VersionDependencyLoader{
					maxBatch: 100,
					wait:     time.Millisecond,
					fetch: func(ids []string) ([][]postgres.VersionDependency, []error) {
						fetchIds := make([]string, 0)
						byID := map[string][]postgres.VersionDependency{}
						for _, id := range ids {
							if dependencies, ok := dbCache.Get("VersionDependenciesByVersionID_" + id); ok {
								byID[id] = dependencies.([]postgres.VersionDependency)
							} else {
								fetchIds = append(fetchIds, id)
							}
						}

						var entities []postgres.VersionDependency
						reqCtx := c.Request().Context()
						postgres.DBCtx(reqCtx).Where("version_id IN ?", fetchIds).Find(&entities)

						for _, entity := range entities {
							byID[entity.VersionID] = append(byID[entity.VersionID], entity)
						}

						results := make([][]postgres.VersionDependency, len(ids))
						for i, id := range ids {
							results[i] = byID[id]
							dbCache.Set("VersionDependenciesByVersionID_"+id, results[i], cache.DefaultExpiration)
						}

						return results, nil
					},
				},
				UserModsByModID: UserModLoader{
					maxBatch: 100,
					wait:     time.Millisecond,
					fetch: func(ids []string) ([][]postgres.UserMod, []error) {
						fetchIds := make([]string, 0)
						byID := map[string][]postgres.UserMod{}
						for _, id := range ids {
							if mods, ok := dbCache.Get("UserModsByModID_" + id); ok {
								byID[id] = mods.([]postgres.UserMod)
							} else {
								fetchIds = append(fetchIds, id)
							}
						}

						var entities []postgres.UserMod
						reqCtx := c.Request().Context()
						postgres.DBCtx(reqCtx).Where("mod_id IN ?", fetchIds).Find(&entities)

						for _, entity := range entities {
							byID[entity.ModID] = append(byID[entity.ModID], entity)
						}

						results := make([][]postgres.UserMod, len(ids))
						for i, id := range ids {
							results[i] = byID[id]

							if results[i] == nil {
								results[i] = make([]postgres.UserMod, 0)
							}

							dbCache.Set("UserModsByModID_"+id, results[i], cache.DefaultExpiration)
						}

						return results, nil
					},
				},
				VersionsByModID: VersionLoader{
					maxBatch: 100,
					wait:     time.Millisecond,
					fetch: func(ids []string) ([][]postgres.Version, []error) {
						fetchIds := make([]string, 0)
						byID := map[string][]postgres.Version{}
						for _, id := range ids {
							if versions, ok := dbCache.Get("VersionsByModID_" + id); ok {
								byID[id] = versions.([]postgres.Version)
							} else {
								fetchIds = append(fetchIds, id)
							}
						}

						var entities []postgres.Version
						reqCtx := c.Request().Context()
						postgres.DBCtx(reqCtx).Preload("Targets").Where("approved = ? AND denied = ? AND mod_id IN ?", true, false, fetchIds).Order("created_at desc").Find(&entities)

						for _, entity := range entities {
							byID[entity.ModID] = append(byID[entity.ModID], entity)
						}

						results := make([][]postgres.Version, len(ids))
						for i, id := range ids {
							results[i] = byID[id]

							if results[i] == nil {
								results[i] = make([]postgres.Version, 0)
							}

							dbCache.Set("VersionsByModID_"+id, results[i], cache.DefaultExpiration)
						}

						return results, nil
					},
				},
				VersionsByModIDNoMeta: VersionLoaderNoMeta{
					maxBatch: 100,
					wait:     time.Millisecond,
					fetch: func(ids []string) ([][]postgres.Version, []error) {
						fetchIds := make([]string, 0)
						byID := map[string][]postgres.Version{}
						for _, id := range ids {
							if versions, ok := dbCache.Get("VersionsByModIDNoMeta_" + id); ok {
								byID[id] = versions.([]postgres.Version)
							} else {
								fetchIds = append(fetchIds, id)
							}
						}

						var entities []postgres.Version
						reqCtx := c.Request().Context()
						postgres.DBCtx(reqCtx).Preload("Targets").Select(
							"id",
							"created_at",
							"updated_at",
							"deleted_at",
							"mod_id",
							"version",
							"sml_version",
							"changelog",
							"downloads",
							"key",
							"stability",
							"approved",
							"hotness",
							"denied",
							"mod_reference",
							"version_major",
							"version_minor",
							"version_patch",
							"size",
							"hash",
						).Where("approved = ? AND denied = ? AND mod_id IN ?", true, false, fetchIds).Order("created_at desc").Find(&entities)

						for _, entity := range entities {
							byID[entity.ModID] = append(byID[entity.ModID], entity)
						}

						results := make([][]postgres.Version, len(ids))
						for i, id := range ids {
							results[i] = byID[id]

							if results[i] == nil {
								results[i] = make([]postgres.Version, 0)
							}

							dbCache.Set("VersionsByModIDNoMeta_"+id, results[i], cache.DefaultExpiration)
						}

						return results, nil
					},
				},
				UserByID: UserLoader{
					maxBatch: 100,
					wait:     time.Millisecond,
					fetch: func(ids []string) ([]*postgres.User, []error) {
						fetchIds := make([]string, 0)
						byID := map[string]*postgres.User{}
						for _, id := range ids {
							if versions, ok := dbCache.Get("UserByID_" + id); ok {
								byID[id] = versions.(*postgres.User)
							} else {
								fetchIds = append(fetchIds, id)
							}
						}

						var entities []postgres.User
						reqCtx := c.Request().Context()
						postgres.DBCtx(reqCtx).Where("id IN ?", fetchIds).Find(&entities)

						for _, entity := range entities {
							tempEntity := entity
							byID[entity.ID] = &tempEntity
						}

						results := make([]*postgres.User, len(ids))
						for i, id := range ids {
							results[i] = byID[id]

							dbCache.Set("UserByID_"+id, results[i], cache.DefaultExpiration)
						}

						return results, nil
					},
				},
			})

			c.SetRequest(c.Request().WithContext(ctx))

			return handlerFunc(c)
		}
	}
}

func For(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey{}).(*Loaders)
}
