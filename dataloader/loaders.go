package dataloader

import (
	"context"
	"errors"

	"entgo.io/ent/dialect/sql"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/labstack/echo/v4"

	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/generated/ent"
	"github.com/satisfactorymodding/smr-api/generated/ent/user"
	"github.com/satisfactorymodding/smr-api/generated/ent/usermod"
	"github.com/satisfactorymodding/smr-api/generated/ent/version"
	"github.com/satisfactorymodding/smr-api/generated/ent/versiondependency"
)

type loadersKey struct{}

type Loaders struct {
	UserByID                       *dataloader.Loader[string, *ent.User]
	VersionDependenciesByVersionID *dataloader.Loader[string, []*ent.VersionDependency]
	UserModsByModID                *dataloader.Loader[string, []*ent.UserMod]
	VersionsByModID                *dataloader.Loader[string, []*ent.Version]
	VersionsByModIDNoMeta          *dataloader.Loader[string, []*ent.Version]
}

func Middleware() func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
	return func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := context.WithValue(c.Request().Context(), loadersKey{}, &Loaders{
				VersionDependenciesByVersionID: dataloader.NewBatchedLoader(func(ctx context.Context, ids []string) []*dataloader.Result[[]*ent.VersionDependency] {
					// TODO Query only selected fields from context
					entities, err := db.From(ctx).VersionDependency.Query().Where(versiondependency.VersionIDIn(ids...)).All(ctx)
					if err != nil {
						return nil
					}

					byID := map[string][]*ent.VersionDependency{}
					for _, entity := range entities {
						byID[entity.VersionID] = append(byID[entity.VersionID], entity)
					}

					results := make([]*dataloader.Result[[]*ent.VersionDependency], len(ids))
					for i, id := range ids {
						if u, ok := byID[id]; ok {
							results[i] = &dataloader.Result[[]*ent.VersionDependency]{Data: u}
						} else {
							results[i] = &dataloader.Result[[]*ent.VersionDependency]{Error: errors.New("version not found")}
						}
					}

					return results
				}, dataloader.WithCache[string, []*ent.VersionDependency](&dataloader.NoCache[string, []*ent.VersionDependency]{})),
				UserModsByModID: dataloader.NewBatchedLoader(func(ctx context.Context, ids []string) []*dataloader.Result[[]*ent.UserMod] {
					// TODO Query only selected fields from context
					entities, err := db.From(ctx).UserMod.Query().Where(usermod.ModIDIn(ids...)).All(ctx)
					if err != nil {
						return nil
					}

					byID := map[string][]*ent.UserMod{}
					for _, entity := range entities {
						byID[entity.ModID] = append(byID[entity.ModID], entity)
					}

					results := make([]*dataloader.Result[[]*ent.UserMod], len(ids))
					for i, id := range ids {
						if u, ok := byID[id]; ok {
							results[i] = &dataloader.Result[[]*ent.UserMod]{Data: u}
						} else {
							results[i] = &dataloader.Result[[]*ent.UserMod]{Error: errors.New("version not found")}
						}
					}

					return results
				}, dataloader.WithCache[string, []*ent.UserMod](&dataloader.NoCache[string, []*ent.UserMod]{})),
				VersionsByModID: dataloader.NewBatchedLoader(func(ctx context.Context, ids []string) []*dataloader.Result[[]*ent.Version] {
					// TODO Query only selected fields from context
					entities, err := db.From(ctx).Version.Query().WithTargets().Where(
						version.ModIDIn(ids...),
						version.Approved(true),
						version.Denied(false),
					).Order(version.ByCreatedAt(sql.OrderDesc())).All(ctx)
					if err != nil {
						return nil
					}

					byID := map[string][]*ent.Version{}
					for _, entity := range entities {
						byID[entity.ModID] = append(byID[entity.ModID], entity)
					}

					results := make([]*dataloader.Result[[]*ent.Version], len(ids))
					for i, id := range ids {
						if u, ok := byID[id]; ok {
							results[i] = &dataloader.Result[[]*ent.Version]{Data: u}
						} else {
							results[i] = &dataloader.Result[[]*ent.Version]{Error: errors.New("version not found")}
						}
					}

					return results
				}, dataloader.WithCache[string, []*ent.Version](&dataloader.NoCache[string, []*ent.Version]{})),
				VersionsByModIDNoMeta: dataloader.NewBatchedLoader(func(ctx context.Context, ids []string) []*dataloader.Result[[]*ent.Version] {
					// TODO Query only selected fields from context
					entities, err := db.From(ctx).Version.Query().Select(
						version.FieldID,
						version.FieldCreatedAt,
						version.FieldUpdatedAt,
						version.FieldDeletedAt,
						version.FieldModID,
						version.FieldVersion,
						version.FieldSmlVersion,
						version.FieldChangelog,
						version.FieldDownloads,
						version.FieldKey,
						version.FieldStability,
						version.FieldApproved,
						version.FieldHotness,
						version.FieldDenied,
						version.FieldModReference,
						version.FieldVersionMajor,
						version.FieldVersionMinor,
						version.FieldVersionPatch,
						version.FieldSize,
						version.FieldHash,
					).WithTargets().Where(
						version.ModIDIn(ids...),
						version.Approved(true),
						version.Denied(false),
					).Order(version.ByCreatedAt(sql.OrderDesc())).All(ctx)
					if err != nil {
						return nil
					}

					byID := map[string][]*ent.Version{}
					for _, entity := range entities {
						byID[entity.ModID] = append(byID[entity.ModID], entity)
					}

					results := make([]*dataloader.Result[[]*ent.Version], len(ids))
					for i, id := range ids {
						if u, ok := byID[id]; ok {
							results[i] = &dataloader.Result[[]*ent.Version]{Data: u}
						} else {
							results[i] = &dataloader.Result[[]*ent.Version]{Error: errors.New("version not found")}
						}
					}

					return results
				}, dataloader.WithCache[string, []*ent.Version](&dataloader.NoCache[string, []*ent.Version]{})),
				UserByID: dataloader.NewBatchedLoader(func(ctx context.Context, ids []string) []*dataloader.Result[*ent.User] {
					// TODO Query only selected fields from context
					entities, err := db.From(ctx).User.Query().Where(user.IDIn(ids...)).All(ctx)
					if err != nil {
						return nil
					}

					byID := map[string]*ent.User{}
					for _, entity := range entities {
						byID[entity.ID] = entity
					}

					results := make([]*dataloader.Result[*ent.User], len(ids))
					for i, id := range ids {
						if u, ok := byID[id]; ok {
							results[i] = &dataloader.Result[*ent.User]{Data: u}
						} else {
							results[i] = &dataloader.Result[*ent.User]{Error: errors.New("user not found")}
						}
					}

					return results
				}, dataloader.WithCache[string, *ent.User](&dataloader.NoCache[string, *ent.User]{})),
			})

			c.SetRequest(c.Request().WithContext(ctx))

			return handlerFunc(c)
		}
	}
}

func For(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey{}).(*Loaders)
}
