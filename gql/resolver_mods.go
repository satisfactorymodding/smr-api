package gql

import (
	"bytes"
	"context"
	"io/ioutil"
	"time"

	"github.com/satisfactorymodding/smr-api/dataloader"
	"github.com/satisfactorymodding/smr-api/db/postgres"
	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/integrations"
	"github.com/satisfactorymodding/smr-api/models"
	"github.com/satisfactorymodding/smr-api/redis"
	"github.com/satisfactorymodding/smr-api/storage"
	"github.com/satisfactorymodding/smr-api/util"
	"github.com/satisfactorymodding/smr-api/util/converter"

	"github.com/pkg/errors"

	"github.com/99designs/gqlgen/graphql"
	"github.com/dgraph-io/ristretto"
	"gopkg.in/go-playground/validator.v9"
)

func (r *mutationResolver) CreateMod(ctx context.Context, mod generated.NewMod) (*generated.Mod, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "createMod")
	defer wrapper.end()

	val := ctx.Value(util.ContextValidator{}).(*validator.Validate)
	if err := val.Struct(&mod); err != nil {
		return nil, errors.Wrap(err, "validation failed")
	}

	if postgres.GetModByReference(mod.ModReference, &newCtx) != nil {
		return nil, errors.New("mod with this mod reference already exists")
	}

	dbMod := &postgres.Mod{
		Name:             mod.Name,
		ShortDescription: mod.ShortDescription,
		Approved:         true,
		ModReference:     mod.ModReference,
	}

	SetStringINN(mod.SourceURL, &dbMod.SourceURL)
	SetStringINN(mod.FullDescription, &dbMod.FullDescription)
	SetBoolINN(mod.Hidden, &dbMod.Hidden)

	user := ctx.Value(postgres.UserKey{}).(*postgres.User)

	dbMod.CreatorID = user.ID

	var logoData []byte

	if mod.Logo != nil {
		file, err := ioutil.ReadAll(mod.Logo.File)

		if err != nil {
			return nil, errors.Wrap(err, "failed to read logo file")
		}

		logoData, err = converter.ConvertAnyImageToWebp(ctx, file)

		if err != nil {
			return nil, errors.Wrap(err, "failed to convert logo file")
		}
	} else {
		dbMod.Logo = ""
	}

	resultMod, err := postgres.CreateMod(dbMod, &newCtx)

	if err != nil {
		return nil, err
	}

	if logoData != nil {
		success, logoKey := storage.UploadModLogo(ctx, resultMod.ID, bytes.NewReader(logoData))
		if success {
			resultMod.Logo = storage.GenerateDownloadLink(logoKey)
			postgres.Save(&resultMod, &newCtx)
		}
	}

	err = postgres.SetModTags(resultMod.ID, mod.TagIDs, &newCtx)

	if err != nil {
		return nil, err
	}

	// Need to get the mod again to populate tags
	return DBModToGenerated(postgres.GetModByIDNoCache(resultMod.ID, &newCtx)), nil
}

func (r *mutationResolver) UpdateMod(ctx context.Context, modID string, mod generated.UpdateMod) (*generated.Mod, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "updateMod")
	defer wrapper.end()

	val := ctx.Value(util.ContextValidator{}).(*validator.Validate)
	if err := val.Struct(&mod); err != nil {
		return nil, errors.Wrap(err, "validation failed")
	}

	err := postgres.ResetModTags(modID, mod.TagIDs, &newCtx)

	if err != nil {
		return nil, err
	}

	dbMod := postgres.GetModByIDNoCache(modID, &newCtx)

	if dbMod == nil {
		return nil, errors.New("mod not found")
	}

	if mod.ModReference != nil && *mod.ModReference != dbMod.ModReference && dbMod.ID != dbMod.ModReference {
		return nil, errors.New("this mod already has set a mod reference")
	}

	SetStringINNOE(mod.Name, &dbMod.Name)
	SetStringINNOE(mod.ShortDescription, &dbMod.ShortDescription)
	SetStringINN(mod.SourceURL, &dbMod.SourceURL)
	SetStringINN(mod.FullDescription, &dbMod.FullDescription)
	SetStringINN(mod.ModReference, &dbMod.ModReference)
	SetBoolINN(mod.Hidden, &dbMod.Hidden)

	if mod.Logo != nil {
		file, err := ioutil.ReadAll(mod.Logo.File)

		if err != nil {
			return nil, errors.Wrap(err, "failed to read logo file")
		}

		logoData, err := converter.ConvertAnyImageToWebp(ctx, file)

		if err != nil {
			return nil, err
		}

		success, logoKey := storage.UploadModLogo(ctx, dbMod.ID, bytes.NewReader(logoData))
		if success {
			dbMod.Logo = storage.GenerateDownloadLink(logoKey)
		} else {
			dbMod.Logo = ""
		}
	}

	postgres.Save(&dbMod, &newCtx)

	if mod.Authors != nil {
		authors, err := dataloader.For(ctx).UserModsByModID.Load(modID)

		if err != nil {
			return nil, err
		}

		for _, author := range authors {
			// Creators cannot be deleted
			if author.Role == "creator" {
				continue
			}

			found := false
			for _, userMod := range mod.Authors {
				if userMod.UserID == author.UserID {
					found = true
					break
				}
			}

			if !found {
				postgres.Delete(author, &newCtx)
			}
		}

		for _, userMod := range mod.Authors {
			role := "creator"

			if userMod.Role == "editor" {
				role = "editor"
			}

			postgres.Save(&postgres.UserMod{
				UserID: userMod.UserID,
				ModID:  modID,
				Role:   role,
			}, &newCtx)
		}
	}

	return DBModToGenerated(dbMod), nil
}

func (r *mutationResolver) DeleteMod(ctx context.Context, modID string) (bool, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "deleteMod")
	defer wrapper.end()

	dbMod := postgres.GetModByID(modID, &newCtx)

	if dbMod == nil {
		return false, errors.New("mod not found")
	}

	postgres.Delete(&dbMod, &newCtx)

	return true, nil
}

func (r *mutationResolver) ApproveMod(ctx context.Context, modID string) (bool, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "approveMod")
	defer wrapper.end()

	dbMod := postgres.GetModByID(modID, &newCtx)

	if dbMod == nil {
		return false, errors.New("mod not found")
	}

	dbMod.Approved = true

	postgres.Save(&dbMod, &newCtx)

	go integrations.NewMod(util.ReWrapCtx(ctx), dbMod)

	return true, nil
}

func (r *mutationResolver) DenyMod(ctx context.Context, modID string) (bool, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "denyMod")
	defer wrapper.end()

	dbMod := postgres.GetModByID(modID, &newCtx)

	if dbMod == nil {
		return false, errors.New("mod not found")
	}

	dbMod.Denied = true

	postgres.Save(&dbMod, &newCtx)
	postgres.Delete(&dbMod, &newCtx)

	return true, nil
}

func (r *queryResolver) GetMod(ctx context.Context, modID string) (*generated.Mod, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "getMod")
	defer wrapper.end()

	mod := postgres.GetModByID(modID, &newCtx)

	if mod != nil {
		if redis.CanIncrement(RealIP(ctx), "view", "mod:"+modID, time.Hour*4) {
			postgres.IncrementModViews(mod, &newCtx)
		}
	}

	return DBModToGenerated(mod), nil
}

func (r *queryResolver) GetModByReference(ctx context.Context, modReference string) (*generated.Mod, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "getModByReference")
	defer wrapper.end()

	mod := postgres.GetModByReference(modReference, &newCtx)

	if mod != nil {
		if redis.CanIncrement(RealIP(ctx), "view", "mod:"+mod.ID, time.Hour*4) {
			postgres.IncrementModViews(mod, &newCtx)
		}
	}

	return DBModToGenerated(mod), nil
}

func (r *queryResolver) GetMods(ctx context.Context, filter map[string]interface{}) (*generated.GetMods, error) {
	wrapper, _ := WrapQueryTrace(ctx, "getMods")
	defer wrapper.end()
	return &generated.GetMods{}, nil
}

func (r *queryResolver) GetUnapprovedMods(ctx context.Context, filter map[string]interface{}) (*generated.GetMods, error) {
	wrapper, _ := WrapQueryTrace(ctx, "getUnapprovedMods")
	defer wrapper.end()
	return &generated.GetMods{}, nil
}

func (r *queryResolver) GetMyMods(ctx context.Context, filter map[string]interface{}) (*generated.GetMyMods, error) {
	wrapper, _ := WrapQueryTrace(ctx, "getMyMods")
	defer wrapper.end()
	return &generated.GetMyMods{}, nil
}

func (r *queryResolver) GetMyUnapprovedMods(ctx context.Context, filter map[string]interface{}) (*generated.GetMyMods, error) {
	wrapper, _ := WrapQueryTrace(ctx, "getMyUnapprovedMods")
	defer wrapper.end()
	return &generated.GetMyMods{}, nil
}

type getModsResolver struct{ *Resolver }

func (r *getModsResolver) Mods(ctx context.Context, obj *generated.GetMods) ([]*generated.Mod, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "GetMods.mods")
	defer wrapper.end()

	resolverContext := graphql.GetFieldContext(ctx)
	unapproved := resolverContext.Parent.Field.Field.Name == "getUnapprovedMods"

	modFilter, err := models.ProcessModFilter(resolverContext.Parent.Args["filter"].(map[string]interface{}))

	if err != nil {
		return nil, err
	}

	for _, field := range graphql.CollectFieldsCtx(ctx, nil) {
		modFilter.AddField(field.Name)
	}

	mods := postgres.GetModsNew(modFilter, unapproved, &newCtx)

	if mods == nil {
		return nil, errors.New("mods not found")
	}

	converted := make([]*generated.Mod, len(mods))
	for k, v := range mods {
		converted[k] = DBModToGenerated(&v)
	}

	return converted, nil
}

func (r *getModsResolver) Count(ctx context.Context, obj *generated.GetMods) (int, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "GetMods.count")
	defer wrapper.end()

	resolverContext := graphql.GetFieldContext(ctx)
	unapproved := resolverContext.Parent.Field.Field.Name == "getUnapprovedMods"

	modFilter, err := models.ProcessModFilter(resolverContext.Parent.Args["filter"].(map[string]interface{}))

	if err != nil {
		return 0, err
	}

	if modFilter.Ids != nil && len(modFilter.Ids) != 0 {
		return len(modFilter.Ids), nil
	}

	return int(postgres.GetModCountNew(modFilter, unapproved, &newCtx)), nil
}

type getMyModsResolver struct{ *Resolver }

func (r *getMyModsResolver) Mods(ctx context.Context, obj *generated.GetMyMods) ([]*generated.Mod, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "GetMyMods.mods")
	defer wrapper.end()

	resolverContext := graphql.GetFieldContext(ctx)
	unapproved := resolverContext.Parent.Field.Field.Name == "getMyUnapprovedMods"

	modFilter, err := models.ProcessModFilter(resolverContext.Parent.Args["filter"].(map[string]interface{}))

	if err != nil {
		return nil, err
	}

	for _, field := range graphql.CollectFieldsCtx(ctx, nil) {
		modFilter.AddField(field.Name)
	}

	var mods []postgres.Mod

	if modFilter.Ids == nil || len(modFilter.Ids) == 0 {
		mods = postgres.GetModsNew(modFilter, unapproved, &newCtx)
	} else {
		mods = postgres.GetModsByID(modFilter.Ids, &newCtx)
	}

	if mods == nil {
		return nil, errors.New("mods not found")
	}

	converted := make([]*generated.Mod, len(mods))
	for k, v := range mods {
		converted[k] = DBModToGenerated(&v)
	}

	return converted, nil
}

func (r *getMyModsResolver) Count(ctx context.Context, obj *generated.GetMyMods) (int, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "GetMyMods.count")
	defer wrapper.end()

	resolverContext := graphql.GetFieldContext(ctx)
	unapproved := resolverContext.Parent.Field.Field.Name == "getMyUnapprovedMods"

	modFilter, err := models.ProcessModFilter(resolverContext.Parent.Args["filter"].(map[string]interface{}))

	if err != nil {
		return 0, err
	}

	if modFilter.Ids != nil && len(modFilter.Ids) != 0 {
		return len(modFilter.Ids), nil
	}

	return int(postgres.GetModCountNew(modFilter, unapproved, &newCtx)), nil
}

type modResolver struct{ *Resolver }

func (r *modResolver) Authors(ctx context.Context, obj *generated.Mod) ([]*generated.UserMod, error) {
	wrapper, _ := WrapQueryTrace(ctx, "Mod.authors")
	defer wrapper.end()

	authors, err := dataloader.For(ctx).UserModsByModID.Load(obj.ID)

	if err != nil {
		return nil, err
	}

	if authors == nil {
		return nil, errors.New("authors not found")
	}

	converted := make([]*generated.UserMod, len(authors))
	for k, v := range authors {
		converted[k] = &generated.UserMod{
			UserID: v.UserID,
			ModID:  v.ModID,
			Role:   v.Role,
		}
	}

	return converted, nil
}

func (r *modResolver) Version(ctx context.Context, obj *generated.Mod, version string) (*generated.Version, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "Mod.version")
	defer wrapper.end()
	return DBVersionToGenerated(postgres.GetModVersionByName(obj.ID, version, &newCtx)), nil
}

var versionNoMetaCache, _ = ristretto.NewCache(&ristretto.Config{
	NumCounters: 1e6, // number of keys to track frequency of (1M).
	MaxCost:     1e6, // maximum cost of cache (1M).
	BufferItems: 64,  // number of keys per Get buffer.
})

const versionNoMetaCacheTTL = time.Second * 30

func (r *modResolver) Versions(ctx context.Context, obj *generated.Mod, filter map[string]interface{}) ([]*generated.Version, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "Mod.versions")
	defer wrapper.end()

	versionFilter, err := models.ProcessVersionFilter(filter)

	if err != nil {
		return nil, err
	}

	hasMetadata := false
	for _, field := range graphql.CollectFieldsCtx(ctx, nil) {
		versionFilter.AddField(field.Name)

		if field.Name == "metadata" {
			hasMetadata = true
		}
	}

	var versions []postgres.Version

	if versionFilter == nil || versionFilter.IsDefault(true) {
		if hasMetadata {
			versions, err = dataloader.For(ctx).VersionsByModID.Load(obj.ID)
		} else {
			if cacheVersions, ok := versionNoMetaCache.Get(obj.ID); ok {
				versions = cacheVersions.([]postgres.Version)
			}

			if versions == nil {
				versions, err = dataloader.For(ctx).VersionsByModIDNoMeta.Load(obj.ID)
				if err == nil && versions != nil {
					versionNoMetaCache.SetWithTTL(obj.ID, versions, int64(len(versions)), versionNoMetaCacheTTL)
				}
			}
		}

		if err != nil {
			return nil, err
		}

		if versionFilter.Limit != nil && *versionFilter.Limit < len(versions) {
			versions = versions[:*versionFilter.Limit]
		}
	} else {
		versions = postgres.GetModVersionsNew(obj.ID, versionFilter, false, &newCtx)
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

func (r *modResolver) LatestVersions(ctx context.Context, obj *generated.Mod) (*generated.LatestVersions, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "Mod.latestVersions")
	defer wrapper.end()

	versions := postgres.GetModLatestVersions(obj.ID, false, &newCtx)

	if versions == nil {
		return nil, errors.New("versions not found")
	}

	versionsD := *versions

	converted := generated.LatestVersions{}
	for _, v := range versionsD {
		switch v.Stability {
		case string(generated.VersionStabilitiesAlpha):
			converted.Alpha = DBVersionToGenerated(&v)
		case string(generated.VersionStabilitiesBeta):
			converted.Beta = DBVersionToGenerated(&v)
		case string(generated.VersionStabilitiesRelease):
			converted.Release = DBVersionToGenerated(&v)
		}
	}

	return &converted, nil
}
