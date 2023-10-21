package gql

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"runtime/debug"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/Vilsol/slox"
	"github.com/dgraph-io/ristretto"
	"github.com/pkg/errors"

	"github.com/satisfactorymodding/smr-api/dataloader"
	"github.com/satisfactorymodding/smr-api/db/postgres"
	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/integrations"
	"github.com/satisfactorymodding/smr-api/models"
	"github.com/satisfactorymodding/smr-api/redis"
	"github.com/satisfactorymodding/smr-api/storage"
	"github.com/satisfactorymodding/smr-api/util"
)

func (r *mutationResolver) CreateVersion(ctx context.Context, modID string) (string, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "createVersion")
	defer wrapper.end()

	mod := postgres.GetModByID(newCtx, modID)

	if mod == nil {
		return "", errors.New("mod not found")
	}

	if !mod.Approved {
		return "", errors.New("mod is not validated")
	}

	if mod.ID == mod.ModReference {
		return "", errors.New("you must update your mod reference on the site to match your mod_reference in your data.json")
	}

	versionID := util.GenerateUniqueID()

	storage.StartUploadMultipartMod(ctx, mod.ID, mod.Name, versionID)

	return versionID, nil
}

func (r *mutationResolver) UploadVersionPart(ctx context.Context, modID string, versionID string, part int, file graphql.Upload) (bool, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "createVersion")
	defer wrapper.end()

	if part > 100 {
		return false, errors.New("files can consist of max 41 chunks")
	}

	mod := postgres.GetModByID(newCtx, modID)

	if mod == nil {
		return false, errors.New("mod not found")
	}

	if !mod.Approved {
		return false, errors.New("mod is not validated")
	}

	if mod.ID == mod.ModReference {
		return false, errors.New("you must update your mod reference on the site to match your mod_reference in your data.json")
	}

	// TODO Optimize
	fileData, err := io.ReadAll(file.File)
	if err != nil {
		return false, fmt.Errorf("failed to read file: %w", err)
	}

	success, _ := storage.UploadMultipartMod(ctx, mod.ID, mod.Name, versionID, int64(part), bytes.NewReader(fileData))

	return success, nil
}

func (r *mutationResolver) FinalizeCreateVersion(ctx context.Context, modID string, versionID string, version generated.NewVersion) (bool, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "finalizeCreateVersion")
	defer wrapper.end()

	mod := postgres.GetModByID(newCtx, modID)

	if mod == nil {
		return false, errors.New("mod not found")
	}

	if !mod.Approved {
		return false, errors.New("mod is not validated")
	}

	if mod.ID == mod.ModReference {
		return false, errors.New("you must update your mod reference on the site to match your mod_reference in your data.json")
	}

	newCtx = slox.With(ctx, slog.String("mod_id", mod.ID), slog.String("version_id", versionID))

	slox.Info(newCtx, "finalization gql call")

	go func(ctx context.Context, mod *postgres.Mod, versionID string, version generated.NewVersion) {
		defer func() {
			if r := recover(); r != nil {
				slox.Error(ctx, "recovered from version finalization", slog.Any("recover", r), slog.String("stack", string(debug.Stack())))

				if err := redis.StoreVersionUploadState(versionID, nil, errors.New("internal error, please try again, if it fails again, please report on discord")); err != nil {
					slox.Error(ctx, "failed to store version upload state", slog.Any("err", err))
				}
			}
		}()

		slox.Info(ctx, "calling FinalizeVersionUploadAsync")

		data, err := FinalizeVersionUploadAsync(ctx, mod, versionID, version)
		if err2 := redis.StoreVersionUploadState(versionID, data, err); err2 != nil {
			slox.Error(ctx, "error storing redis state", slog.Any("err", err))
			return
		}

		slox.Info(ctx, "finished FinalizeVersionUploadAsync")

		if err != nil {
			slox.Error(ctx, "error completing version upload", slog.Any("err", err))
		} else {
			slox.Info(ctx, "completed version upload")
		}
	}(util.ReWrapCtx(ctx), mod, versionID, version)

	return true, nil
}

func (r *mutationResolver) UpdateVersion(ctx context.Context, versionID string, version generated.UpdateVersion) (*generated.Version, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "updateVersion")
	defer wrapper.end()

	dbVersion := postgres.GetVersion(newCtx, versionID)

	if dbVersion == nil {
		return nil, errors.New("version not found")
	}

	SetStringINNOE(version.Changelog, &dbVersion.Changelog)
	SetStabilityINN(version.Stability, &dbVersion.Stability)

	postgres.Save(newCtx, &dbVersion)

	return DBVersionToGenerated(dbVersion), nil
}

func (r *mutationResolver) DeleteVersion(ctx context.Context, versionID string) (bool, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "deleteVersion")
	defer wrapper.end()

	dbVersion := postgres.GetVersion(newCtx, versionID)

	if dbVersion == nil {
		return false, errors.New("version not found")
	}

	postgres.Delete(newCtx, &dbVersion)

	return true, nil
}

func (r *mutationResolver) ApproveVersion(ctx context.Context, versionID string) (bool, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "approveVersion")
	defer wrapper.end()

	dbVersion := postgres.GetVersion(newCtx, versionID)

	if dbVersion == nil {
		return false, errors.New("version not found")
	}

	dbVersion.Approved = true

	postgres.Save(newCtx, &dbVersion)

	mod := postgres.GetModByID(newCtx, dbVersion.ModID)
	now := time.Now()
	mod.LastVersionDate = &now
	postgres.Save(newCtx, &mod)

	go integrations.NewVersion(util.ReWrapCtx(ctx), dbVersion)

	return true, nil
}

func (r *mutationResolver) DenyVersion(ctx context.Context, versionID string) (bool, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "denyVersion")
	defer wrapper.end()

	dbVersion := postgres.GetVersion(newCtx, versionID)

	if dbVersion == nil {
		return false, errors.New("version not found")
	}

	dbVersion.Denied = true

	postgres.Save(newCtx, &dbVersion)
	postgres.Delete(newCtx, &dbVersion)

	mod := postgres.GetModByID(newCtx, dbVersion.ModID)
	postgres.Save(newCtx, &mod)

	return true, nil
}

func (r *queryResolver) GetVersion(ctx context.Context, versionID string) (*generated.Version, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "getVersion")
	defer wrapper.end()
	return DBVersionToGenerated(postgres.GetVersion(newCtx, versionID)), nil
}

func (r *queryResolver) GetVersions(ctx context.Context, _ map[string]interface{}) (*generated.GetVersions, error) {
	wrapper, _ := WrapQueryTrace(ctx, "getVersions")
	defer wrapper.end()
	return &generated.GetVersions{}, nil
}

func (r *queryResolver) GetUnapprovedVersions(ctx context.Context, _ map[string]interface{}) (*generated.GetVersions, error) {
	wrapper, _ := WrapQueryTrace(ctx, "getUnapprovedVersions")
	defer wrapper.end()
	return &generated.GetVersions{}, nil
}

func (r *queryResolver) GetMyVersions(ctx context.Context, _ map[string]interface{}) (*generated.GetMyVersions, error) {
	wrapper, _ := WrapQueryTrace(ctx, "getMyVersions")
	defer wrapper.end()
	return &generated.GetMyVersions{}, nil
}

func (r *queryResolver) GetMyUnapprovedVersions(ctx context.Context, _ map[string]interface{}) (*generated.GetMyVersions, error) {
	wrapper, _ := WrapQueryTrace(ctx, "getMyUnapprovedVersions")
	defer wrapper.end()
	return &generated.GetMyVersions{}, nil
}

func (r *queryResolver) CheckVersionUploadState(ctx context.Context, modID string, versionID string) (*generated.CreateVersionResponse, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "checkVersionUploadState")
	defer wrapper.end()

	mod := postgres.GetModByID(newCtx, modID)

	if mod == nil {
		return nil, errors.New("mod not found")
	}

	if !mod.Approved {
		return nil, errors.New("mod is not validated")
	}

	if mod.ID == mod.ModReference {
		return nil, errors.New("you must update your mod reference on the site to match your mod_reference in your data.json")
	}

	return redis.GetVersionUploadState(versionID)
}

type getVersionsResolver struct{ *Resolver }

func (r *getVersionsResolver) Versions(ctx context.Context, _ *generated.GetVersions) ([]*generated.Version, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "GetVersions.versions")
	defer wrapper.end()

	resolverContext := graphql.GetFieldContext(ctx)
	unapproved := resolverContext.Parent.Field.Field.Name == "getUnapprovedVersions"

	versionFilter, err := models.ProcessVersionFilter(resolverContext.Parent.Args["filter"].(map[string]interface{}))
	if err != nil {
		return nil, err
	}

	for _, field := range graphql.CollectFieldsCtx(ctx, nil) {
		versionFilter.AddField(field.Name)
	}

	var versions []postgres.Version

	if versionFilter.Ids == nil || len(versionFilter.Ids) == 0 {
		versions = postgres.GetVersionsNew(newCtx, versionFilter, unapproved)
	} else {
		versions = postgres.GetVersionsByID(newCtx, versionFilter.Ids)
	}

	if versions == nil {
		return nil, errors.New("versions not found")
	}

	converted := make([]*generated.Version, len(versions))
	for k, v := range versions {
		converted[k] = DBVersionToGenerated(&v)
	}

	return converted, nil
}

func (r *getVersionsResolver) Count(ctx context.Context, _ *generated.GetVersions) (int, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "GetVersions.count")
	defer wrapper.end()

	resolverContext := graphql.GetFieldContext(ctx)
	unapproved := resolverContext.Parent.Field.Field.Name == "getUnapprovedVersions"

	versionFilter, err := models.ProcessVersionFilter(resolverContext.Parent.Args["filter"].(map[string]interface{}))
	if err != nil {
		return 0, err
	}

	if versionFilter.Ids != nil && len(versionFilter.Ids) != 0 {
		return len(versionFilter.Ids), nil
	}

	return int(postgres.GetVersionCountNew(newCtx, versionFilter, unapproved)), nil
}

type versionResolver struct{ *Resolver }

func findWindowsTarget(obj *generated.Version) *generated.VersionTarget {
	var windowsTarget *generated.VersionTarget
	for _, target := range obj.Targets {
		if target.TargetName == "Windows" {
			windowsTarget = target
			break
		}
	}
	return windowsTarget
}

func (r *versionResolver) Link(ctx context.Context, obj *generated.Version) (string, error) {
	wrapper, _ := WrapQueryTrace(ctx, "Version.link")
	defer wrapper.end()

	windowsTarget := findWindowsTarget(obj)
	if windowsTarget != nil {
		link, _ := r.VersionTarget().Link(ctx, windowsTarget)
		return link, nil
	}

	return "/v1/version/" + obj.ID + "/download", nil
}

func (r *versionResolver) Mod(ctx context.Context, obj *generated.Version) (*generated.Mod, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "Version.mod")
	defer wrapper.end()

	return DBModToGenerated(postgres.GetModByID(newCtx, obj.ModID)), nil
}

func (r *versionResolver) Hash(ctx context.Context, obj *generated.Version) (*string, error) {
	wrapper, _ := WrapQueryTrace(ctx, "Version.hash")
	defer wrapper.end()

	hash := ""

	windowsTarget := findWindowsTarget(obj)
	if windowsTarget == nil {
		if obj.Hash == nil {
			return nil, nil
		}
		hash = *obj.Hash
	} else {
		if windowsTarget.Hash == nil {
			return nil, nil
		}
		hash = *windowsTarget.Hash
	}

	return &hash, nil
}

func (r *versionResolver) Size(ctx context.Context, obj *generated.Version) (*int, error) {
	wrapper, _ := WrapQueryTrace(ctx, "Version.size")
	defer wrapper.end()

	size := 0

	windowsTarget := findWindowsTarget(obj)
	if windowsTarget == nil {
		if obj.Size == nil {
			return nil, nil
		}
		size = *obj.Size
	} else {
		if windowsTarget.Size == nil {
			return nil, nil
		}
		size = *windowsTarget.Size
	}

	return &size, nil
}

var versionDependencyCache, _ = ristretto.NewCache(&ristretto.Config{
	NumCounters: 1e6, // number of keys to track frequency of (1M).
	MaxCost:     1e6, // maximum cost of cache (1M).
	BufferItems: 64,  // number of keys per Get buffer.
})

const versionDependencyCacheTTL = time.Minute * 10

func (r *versionResolver) Dependencies(ctx context.Context, obj *generated.Version) ([]*generated.VersionDependency, error) {
	wrapper, _ := WrapQueryTrace(ctx, "Version.dependencies")
	defer wrapper.end()

	var dependencies []postgres.VersionDependency

	if cacheVersions, ok := versionDependencyCache.Get(obj.ID); ok {
		dependencies = cacheVersions.([]postgres.VersionDependency)
	}

	if dependencies == nil {
		var err error
		dependencies, err = dataloader.For(ctx).VersionDependenciesByVersionID.Load(obj.ID)

		if err != nil {
			return nil, err
		}

		versionDependencyCache.SetWithTTL(obj.ID, dependencies, int64(len(dependencies)), versionDependencyCacheTTL)
	}

	converted := make([]*generated.VersionDependency, len(dependencies))
	for k, v := range dependencies {
		converted[k] = DBVersionDependencyToGenerated(&v)
	}

	return converted, nil
}

type versionTargetResolver struct{ *Resolver }

func (r *versionTargetResolver) Link(_ context.Context, obj *generated.VersionTarget) (string, error) {
	return "/v1/version/" + obj.VersionID + "/" + string(obj.TargetName) + "/download", nil
}

type getMyVersionsResolver struct{ *Resolver }

func (r *getMyVersionsResolver) Versions(ctx context.Context, _ *generated.GetMyVersions) ([]*generated.Version, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "GetMyVersions.versions")
	defer wrapper.end()

	resolverContext := graphql.GetFieldContext(ctx)
	unapproved := resolverContext.Parent.Field.Field.Name == "getMyUnapprovedVersions"

	versionFilter, err := models.ProcessVersionFilter(resolverContext.Parent.Args["filter"].(map[string]interface{}))
	if err != nil {
		return nil, err
	}

	for _, field := range graphql.CollectFieldsCtx(ctx, nil) {
		versionFilter.AddField(field.Name)
	}

	var versions []postgres.Version

	if versionFilter.Ids == nil || len(versionFilter.Ids) == 0 {
		versions = postgres.GetVersionsNew(newCtx, versionFilter, unapproved)
	} else {
		versions = postgres.GetVersionsByID(newCtx, versionFilter.Ids)
	}

	if versions == nil {
		return nil, errors.New("versions not found")
	}

	converted := make([]*generated.Version, len(versions))
	for k, v := range versions {
		converted[k] = DBVersionToGenerated(&v)
	}

	return converted, nil
}

func (r *getMyVersionsResolver) Count(ctx context.Context, _ *generated.GetMyVersions) (int, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "GetMyVersions.count")
	defer wrapper.end()

	resolverContext := graphql.GetFieldContext(ctx)
	unapproved := resolverContext.Parent.Field.Field.Name == "getMyUnapprovedVersions"

	versionFilter, err := models.ProcessVersionFilter(resolverContext.Parent.Args["filter"].(map[string]interface{}))
	if err != nil {
		return 0, err
	}

	if versionFilter.Ids != nil && len(versionFilter.Ids) != 0 {
		return len(versionFilter.Ids), nil
	}

	return int(postgres.GetVersionCountNew(newCtx, versionFilter, unapproved)), nil
}
