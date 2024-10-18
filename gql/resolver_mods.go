package gql

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"math"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/99designs/gqlgen/graphql"
	"github.com/Vilsol/slox"
	"github.com/dgraph-io/ristretto"
	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"

	"github.com/satisfactorymodding/smr-api/dataloader"
	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/generated/conv"
	"github.com/satisfactorymodding/smr-api/generated/ent"
	"github.com/satisfactorymodding/smr-api/generated/ent/mod"
	"github.com/satisfactorymodding/smr-api/generated/ent/usermod"
	"github.com/satisfactorymodding/smr-api/generated/ent/version"
	"github.com/satisfactorymodding/smr-api/integrations"
	"github.com/satisfactorymodding/smr-api/models"
	"github.com/satisfactorymodding/smr-api/redis"
	"github.com/satisfactorymodding/smr-api/storage"
	"github.com/satisfactorymodding/smr-api/util"
	"github.com/satisfactorymodding/smr-api/util/converter"
)

var DisallowedModReferences = map[string]bool{
	"satisfactory":          true,
	"factorygame":           true,
	"sml":                   true,
	"satisfactorymodloader": true,
	"examplemod":            true,
	"docmod":                true,
}

func (r *mutationResolver) CreateMod(ctx context.Context, newMod generated.NewMod) (*generated.Mod, error) {
	val := ctx.Value(util.ContextValidator{}).(*validator.Validate)
	if err := val.Struct(&newMod); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	if DisallowedModReferences[strings.ToLower(newMod.ModReference)] {
		return nil, errors.New("using this mod reference is not allowed")
	}

	exist, err := db.From(ctx).Mod.Query().Where(mod.ModReference(newMod.ModReference)).Exist(ctx)
	if err != nil {
		return nil, err
	}

	if exist {
		return nil, errors.New("mod with this mod reference already exists")
	}

	dbMod := db.From(ctx).Mod.Create().
		SetName(newMod.Name).
		SetShortDescription(newMod.ShortDescription).
		SetApproved(true).
		SetModReference(newMod.ModReference)

	SetINNF(newMod.SourceURL, dbMod.SetSourceURL)
	SetINNF(newMod.FullDescription, dbMod.SetFullDescription)
	SetINNF(newMod.Hidden, dbMod.SetHidden)
	SetINNF(newMod.ToggleNetworkUse, dbMod.SetToggleNetworkUse)
	SetINNF(newMod.ToggleExplicitContent, dbMod.SetToggleExplicitContent)

	user, _, err := db.UserFromGQLContext(ctx)
	if err != nil {
		return nil, err
	}

	dbMod.SetCreatorID(user.ID)

	// Allow only new 4 mods per 24h
	existingMods, err := db.From(ctx).Mod.Query().
		Order(mod.ByCreatedAt(sql.OrderAsc())).
		Where(mod.CreatorID(user.ID), mod.CreatedAtGT(time.Now().Add(time.Hour*24*-1))).
		All(ctx)
	if err != nil {
		return nil, err
	}

	currentAvailable := float64(util.ModsPer24h)
	lastModTime := time.Now()
	for _, mod := range existingMods {
		currentAvailable--
		if mod.CreatedAt.After(lastModTime) {
			diff := mod.CreatedAt.Sub(lastModTime)
			currentAvailable = math.Min(float64(util.ModsPer24h), currentAvailable+diff.Hours()/6)
		}
		lastModTime = mod.CreatedAt
	}

	if currentAvailable < 1 {
		timeToWait := time.Until(lastModTime.Add(time.Hour * 6)).Minutes()
		return nil, fmt.Errorf("please wait %.0f minutes to post another mod", timeToWait)
	}

	// Create mod
	resultMod, err := dbMod.Save(ctx)
	if err != nil {
		return nil, err
	}

	if err := db.From(ctx).UserMod.Create().
		SetRole("creator").
		SetModID(resultMod.ID).
		SetUserID(user.ID).
		Exec(ctx); err != nil {
		return nil, err
	}

	if newMod.Logo != nil {
		file, err := io.ReadAll(newMod.Logo.File)
		if err != nil {
			return nil, fmt.Errorf("failed to read logo file: %w", err)
		}

		logoData, thumbHash, err := converter.ConvertAnyImageToWebp(ctx, file)
		if err != nil {
			return nil, fmt.Errorf("failed to convert logo file: %w", err)
		}

		logoKey, err := storage.UploadModLogo(ctx, resultMod.ID, bytes.NewReader(logoData))
		if err == nil {
			resultMod, err = resultMod.Update().
				SetLogo(storage.GenerateDownloadLink(ctx, logoKey)).
				SetLogoThumbhash(thumbHash).
				Save(ctx)
			if err != nil {
				return nil, err
			}
		}
	}

	if len(newMod.TagIDs) > 0 {
		if err := resultMod.Update().AddTagIDs(newMod.TagIDs...).Exec(ctx); err != nil {
			return nil, err
		}
	}

	// Need to get the mod again to populate tags

	resultMod, err = db.From(ctx).Mod.Query().WithTags().Where(mod.ID(resultMod.ID)).First(ctx)
	if err != nil {
		return nil, err
	}

	return (*conv.ModImpl)(nil).Convert(resultMod), nil
}

func (r *mutationResolver) UpdateMod(ctx context.Context, modID string, updateMod generated.UpdateMod) (*generated.Mod, error) {
	val := ctx.Value(util.ContextValidator{}).(*validator.Validate)
	if err := val.Struct(&updateMod); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	dbMod, err := db.From(ctx).Mod.Get(ctx, modID)
	if err != nil {
		return nil, err
	}

	if updateMod.ModReference != nil && *updateMod.ModReference != dbMod.ModReference && dbMod.ID != dbMod.ModReference {
		return nil, errors.New("this mod already has set a mod reference")
	}

	dbUpdate := dbMod.Update()
	if updateMod.TagIDs != nil {
		dbUpdate.ClearTags().AddTagIDs(updateMod.TagIDs...)
	}

	SetINNOEF(updateMod.Name, dbUpdate.SetName)
	SetINNOEF(updateMod.ShortDescription, dbUpdate.SetShortDescription)
	SetINNF(updateMod.SourceURL, dbUpdate.SetSourceURL)
	SetINNF(updateMod.FullDescription, dbUpdate.SetFullDescription)
	SetINNF(updateMod.ModReference, dbUpdate.SetModReference)
	SetINNF(updateMod.Hidden, dbUpdate.SetHidden)
	SetCompatibilityINNF(updateMod.Compatibility, dbUpdate.SetCompatibility)
	SetINNF(updateMod.ToggleNetworkUse, dbUpdate.SetToggleNetworkUse)
	SetINNF(updateMod.ToggleExplicitContent, dbUpdate.SetToggleExplicitContent)

	if updateMod.Logo != nil {
		file, err := io.ReadAll(updateMod.Logo.File)
		if err != nil {
			return nil, fmt.Errorf("failed to read logo file: %w", err)
		}

		logoData, thumbHash, err := converter.ConvertAnyImageToWebp(ctx, file)
		if err != nil {
			return nil, err
		}

		logoKey, err := storage.UploadModLogo(ctx, dbMod.ID, bytes.NewReader(logoData))
		if err == nil {
			dbUpdate.SetLogo(storage.GenerateDownloadLink(ctx, logoKey))
			dbUpdate.SetLogoThumbhash(thumbHash)
		} else {
			dbUpdate.ClearLogo()
			dbUpdate.ClearLogoThumbhash()
		}
	}

	dbMod, err = dbUpdate.Save(ctx)
	if err != nil {
		return nil, err
	}

	if updateMod.Authors != nil {
		authors, err := dataloader.For(ctx).UserModsByModID.Load(ctx, modID)()
		if err != nil {
			return nil, err
		}

		for _, author := range authors {
			// Creators cannot be deleted
			if author.Role == "creator" {
				continue
			}

			found := false
			for _, userMod := range updateMod.Authors {
				if userMod.UserID == author.UserID {
					found = true
					break
				}
			}

			if !found {
				if _, err := db.From(ctx).UserMod.Delete().
					Where(usermod.UserID(author.UserID), usermod.ModID(author.ModID)).
					Exec(ctx); err != nil {
					return nil, err
				}
			}
		}

		for _, userMod := range updateMod.Authors {
			role := "creator"

			if userMod.Role == "editor" {
				role = "editor"
			}

			var existing *ent.UserMod
			for _, author := range authors {
				if author.UserID == userMod.UserID {
					existing = author
					break
				}
			}

			if existing != nil {
				if err := db.From(ctx).UserMod.UpdateOne(existing).
					SetRole(role).
					Exec(ctx); err != nil {
					return nil, err
				}
			} else {
				if err := db.From(ctx).UserMod.Create().
					SetUserID(userMod.UserID).
					SetModID(modID).
					SetRole(role).
					Exec(ctx); err != nil {
					return nil, err
				}
			}
		}
	}

	return (*conv.ModImpl)(nil).Convert(dbMod), nil
}

func (r *mutationResolver) UpdateModCompatibility(ctx context.Context, modID string, compatibility generated.CompatibilityInfoInput) (bool, error) {
	updateMod := generated.UpdateMod{
		Compatibility: &compatibility,
	}
	_, err := r.UpdateMod(ctx, modID, updateMod)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *mutationResolver) UpdateMultipleModCompatibilities(ctx context.Context, modIDs []string, compatibility generated.CompatibilityInfoInput) (bool, error) {
	for _, modID := range modIDs {
		_, err := r.UpdateModCompatibility(ctx, modID, compatibility)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

func (r *mutationResolver) DeleteMod(ctx context.Context, modID string) (bool, error) {
	if err := db.From(ctx).Mod.DeleteOneID(modID).Exec(ctx); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) ApproveMod(ctx context.Context, modID string) (bool, error) {
	dbMod, err := db.From(ctx).Mod.Get(ctx, modID)
	if err != nil {
		return false, err
	}

	if dbMod == nil {
		return false, errors.New("mod not found")
	}

	if err := dbMod.Update().SetApproved(true).Exec(ctx); err != nil {
		return false, err
	}

	go integrations.NewMod(db.ReWrapCtx(ctx), dbMod)

	return true, nil
}

func (r *mutationResolver) DenyMod(ctx context.Context, modID string) (bool, error) {
	dbMod, err := db.From(ctx).Mod.Get(ctx, modID)
	if err != nil {
		return false, err
	}

	if dbMod == nil {
		return false, errors.New("mod not found")
	}

	if err := dbMod.Update().SetDenied(true).Exec(ctx); err != nil {
		return false, err
	}

	if err := db.From(ctx).Mod.DeleteOneID(modID).Exec(ctx); err != nil {
		return false, err
	}

	return true, nil
}

func (r *queryResolver) GetMod(ctx context.Context, modID string) (*generated.Mod, error) {
	dbMod, err := db.From(ctx).Mod.Query().Where(mod.ID(modID)).WithTags().First(ctx)
	if err != nil {
		return nil, err
	}

	if dbMod != nil {
		if redis.CanIncrement(RealIP(ctx), "view", "mod:"+modID, time.Hour*4) {
			if err := dbMod.Update().AddViews(1).Exec(ctx); err != nil {
				return nil, err
			}
		}
	}

	return (*conv.ModImpl)(nil).Convert(dbMod), nil
}

func (r *queryResolver) GetModByReference(ctx context.Context, modReference string) (*generated.Mod, error) {
	dbMod, err := db.From(ctx).Mod.Query().Where(mod.ModReference(modReference)).WithTags().First(ctx)
	if err != nil {
		return nil, err
	}

	if dbMod != nil {
		if redis.CanIncrement(RealIP(ctx), "view", "mod:"+dbMod.ID, time.Hour*4) {
			if err := dbMod.Update().AddViews(1).Exec(ctx); err != nil {
				return nil, err
			}
		}
	}

	return (*conv.ModImpl)(nil).Convert(dbMod), nil
}

func (r *queryResolver) GetMods(_ context.Context, _ map[string]interface{}) (*generated.GetMods, error) {
	return &generated.GetMods{}, nil
}

func (r *queryResolver) GetUnapprovedMods(_ context.Context, _ map[string]interface{}) (*generated.GetMods, error) {
	return &generated.GetMods{}, nil
}

func (r *queryResolver) GetMyMods(_ context.Context, _ map[string]interface{}) (*generated.GetMyMods, error) {
	return &generated.GetMyMods{}, nil
}

func (r *queryResolver) GetMyUnapprovedMods(_ context.Context, _ map[string]interface{}) (*generated.GetMyMods, error) {
	return &generated.GetMyMods{}, nil
}

type getModsResolver struct{ *Resolver }

func (r *getModsResolver) Mods(ctx context.Context, _ *generated.GetMods) ([]*generated.Mod, error) {
	resolverContext := graphql.GetFieldContext(ctx)
	unapproved := resolverContext.Parent.Field.Field.Name == "getUnapprovedMods"

	modFilter, err := models.ProcessModFilter(resolverContext.Parent.Args["filter"].(map[string]interface{}))
	if err != nil {
		return nil, err
	}

	for _, field := range graphql.CollectFieldsCtx(ctx, nil) {
		modFilter.AddField(field.Name)
	}

	query := db.From(ctx).Mod.Query()
	query = db.ConvertModFilter(query, modFilter, false, unapproved)

	result, err := query.All(ctx)
	if err != nil {
		slox.Error(ctx, "failed querying mods", slog.Any("err", err))
		return nil, err
	}

	return (*conv.ModImpl)(nil).ConvertSlice(result), nil
}

func (r *getModsResolver) Count(ctx context.Context, _ *generated.GetMods) (int, error) {
	resolverContext := graphql.GetFieldContext(ctx)
	unapproved := resolverContext.Parent.Field.Field.Name == "getUnapprovedMods"

	modFilter, err := models.ProcessModFilter(resolverContext.Parent.Args["filter"].(map[string]interface{}))
	if err != nil {
		return 0, err
	}

	query := db.From(ctx).Mod.Query()
	query = db.ConvertModFilter(query, modFilter, true, unapproved)

	result, err := query.Count(ctx)
	if err != nil {
		slox.Error(ctx, "failed querying mod count", slog.Any("err", err))
		return 0, err
	}

	return result, nil
}

type getMyModsResolver struct{ *Resolver }

func (r *getMyModsResolver) Mods(ctx context.Context, _ *generated.GetMyMods) ([]*generated.Mod, error) {
	resolverContext := graphql.GetFieldContext(ctx)
	unapproved := resolverContext.Parent.Field.Field.Name == "getMyUnapprovedMods"

	modFilter, err := models.ProcessModFilter(resolverContext.Parent.Args["filter"].(map[string]interface{}))
	if err != nil {
		return nil, err
	}

	for _, field := range graphql.CollectFieldsCtx(ctx, nil) {
		modFilter.AddField(field.Name)
	}

	query := db.From(ctx).Mod.Query()
	query = db.ConvertModFilter(query, modFilter, false, unapproved)

	result, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	return (*conv.ModImpl)(nil).ConvertSlice(result), nil
}

func (r *getMyModsResolver) Count(ctx context.Context, _ *generated.GetMyMods) (int, error) {
	resolverContext := graphql.GetFieldContext(ctx)
	unapproved := resolverContext.Parent.Field.Field.Name == "getMyUnapprovedMods"

	modFilter, err := models.ProcessModFilter(resolverContext.Parent.Args["filter"].(map[string]interface{}))
	if err != nil {
		return 0, err
	}

	query := db.From(ctx).Mod.Query()
	query = db.ConvertModFilter(query, modFilter, false, unapproved)

	result, err := query.Count(ctx)
	if err != nil {
		return 0, err
	}

	return result, nil
}

type modResolver struct{ *Resolver }

func (r *modResolver) Authors(ctx context.Context, obj *generated.Mod) ([]*generated.UserMod, error) {
	authors, err := dataloader.For(ctx).UserModsByModID.Load(ctx, obj.ID)()
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

func (r *modResolver) Version(ctx context.Context, obj *generated.Mod, versionName string) (*generated.Version, error) {
	dbVersion, err := db.From(ctx).Version.Query().
		WithTargets().
		Where(version.Version(versionName), version.ModID(obj.ID)).
		First(ctx)
	if err != nil {
		return nil, err
	}

	return (*conv.VersionImpl)(nil).Convert(dbVersion), nil
}

var versionNoMetaCache, _ = ristretto.NewCache(&ristretto.Config{
	NumCounters: 1e6, // number of keys to track frequency of (1M).
	MaxCost:     1e6, // maximum cost of cache (1M).
	BufferItems: 64,  // number of keys per Get buffer.
})

const versionNoMetaCacheTTL = time.Second * 30

func (r *modResolver) Versions(ctx context.Context, obj *generated.Mod, filter map[string]interface{}) ([]*generated.Version, error) {
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

	var versions []*ent.Version

	if versionFilter == nil || versionFilter.IsDefault(true) {
		if hasMetadata {
			versions, err = dataloader.For(ctx).VersionsByModID.Load(ctx, obj.ID)()
		} else {
			if cacheVersions, ok := versionNoMetaCache.Get(obj.ID); ok {
				versions = cacheVersions.([]*ent.Version)
			}

			if versions == nil {
				versions, err = dataloader.For(ctx).VersionsByModIDNoMeta.Load(ctx, obj.ID)()
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
		query := db.From(ctx).Version.Query().WithTargets().Where(
			version.Approved(true),
			version.Denied(false),
			version.ModID(obj.ID),
		)

		if filter != nil {
			query = query.Limit(*versionFilter.Limit).
				Offset(*versionFilter.Offset).
				Order(sql.OrderByField(
					versionFilter.OrderBy.String(),
					db.OrderToOrder(versionFilter.Order.String()),
				).ToFunc())
		}

		versions, err = query.All(ctx)
		if err != nil {
			return nil, err
		}
	}

	if versions == nil {
		return nil, errors.New("versions not found")
	}

	return (*conv.VersionImpl)(nil).ConvertSlice(versions), nil
}

func (r *modResolver) LatestVersions(ctx context.Context, obj *generated.Mod) (*generated.LatestVersions, error) {
	versions, err := db.From(ctx).Version.
		Query().
		WithTargets().
		Where(
			version.ModID(obj.ID),
			version.Approved(true),
			version.Denied(false),
		).
		Order(
			version.ByModID(),
			version.ByStability(),
			version.ByCreatedAt(sql.OrderDesc()),
		).
		Modify(func(s *sql.Selector) {
			s.SelectExpr(sql.Expr("DISTINCT on (mod_id, stability) *"))
		}).
		All(ctx)
	if err != nil {
		return nil, err
	}

	converted := generated.LatestVersions{}
	for _, v := range versions {
		switch v.Stability {
		case util.StabilityAlpha:
			converted.Alpha = (*conv.VersionImpl)(nil).Convert(v)
		case util.StabilityBeta:
			converted.Beta = (*conv.VersionImpl)(nil).Convert(v)
		case util.StabilityRelease:
			converted.Release = (*conv.VersionImpl)(nil).Convert(v)
		}
	}

	return &converted, nil
}

func (r *queryResolver) GetModByIDOrReference(ctx context.Context, modIDOrReference string) (*generated.Mod, error) {
	m, err := db.From(ctx).Mod.Query().WithTags().Where(mod.Or(
		mod.ID(modIDOrReference),
		mod.ModReference(modIDOrReference),
	)).First(ctx)
	if err != nil {
		return nil, err
	}

	if m != nil {
		if redis.CanIncrement(RealIP(ctx), "view", "mod:"+m.ID, time.Hour*4) {
			if err := m.Update().AddViews(1).Exec(ctx); err != nil {
				slox.Error(ctx, "failed incrementing mod views", slog.Any("err", err))
			}
		}
	}

	return (*conv.ModImpl)(nil).Convert(m), nil
}

func (r *queryResolver) ResolveModVersions(ctx context.Context, filter []*generated.ModVersionConstraint) ([]*generated.ModVersion, error) {
	constraintMapping := make(map[string]string)
	modIDOrReferences := make([]string, len(filter))
	for i, constraint := range filter {
		modIDOrReferences[i] = constraint.ModIDOrReference
		constraintMapping[constraint.ModIDOrReference] = constraint.Version
	}

	mods, err := db.From(ctx).Mod.Query().
		WithTags().
		Where(mod.Or(mod.IDIn(modIDOrReferences...), mod.ModReferenceIn(modIDOrReferences...))).
		All(ctx)
	if err != nil {
		return nil, err
	}

	modVersions := make([]*generated.ModVersion, len(mods))
	for i, m := range mods {
		constraint, ok := constraintMapping[m.ID]
		if !ok {
			constraint = constraintMapping[m.ModReference]
		}

		versions, err := db.GetModVersionsConstraint(ctx, m.ID, constraint)
		if err != nil {
			return nil, err
		}

		modVersions[i] = &generated.ModVersion{
			ID:           m.ID,
			ModReference: m.ModReference,
			Versions:     (*conv.VersionImpl)(nil).ConvertSlice(versions),
		}
	}

	return modVersions, nil
}

func (r *queryResolver) GetModAssetList(ctx context.Context, modReference string) ([]string, error) {
	list := redis.GetModAssetList(modReference)
	if list != nil {
		return list, nil
	}

	assets, err := storage.ListModAssets(ctx, modReference)
	if err != nil {
		return nil, err
	}

	redis.StoreModAssetList(modReference, assets)

	return assets, nil
}
