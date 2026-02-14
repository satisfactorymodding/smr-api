package gql

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"maps"
	"math"
	"slices"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/Vilsol/slox"
	"github.com/go-playground/validator/v10"
	resolver "github.com/satisfactorymodding/ficsit-resolver"

	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/generated/conv"
	"github.com/satisfactorymodding/smr-api/generated/ent"
	"github.com/satisfactorymodding/smr-api/generated/ent/mod"
	"github.com/satisfactorymodding/smr-api/generated/ent/modpack"
	"github.com/satisfactorymodding/smr-api/generated/ent/modpackrelease"
	"github.com/satisfactorymodding/smr-api/generated/ent/version"
	"github.com/satisfactorymodding/smr-api/models"
	"github.com/satisfactorymodding/smr-api/redis"
	"github.com/satisfactorymodding/smr-api/storage"
	"github.com/satisfactorymodding/smr-api/util"
	"github.com/satisfactorymodding/smr-api/util/converter"
)

func (r *queryResolver) GetModpack(ctx context.Context, modpackID string) (*generated.Modpack, error) {
	dbModpack, err := db.From(ctx).Modpack.Query().
		Where(modpack.ID(modpackID)).
		WithTags().
		WithTargets().
		WithReleases().
		WithModpackMods().
		WithParent().
		WithChildren().
		First(ctx)
	if err != nil {
		return nil, err
	}

	if dbModpack == nil {
		return nil, nil
	}

	if redis.CanIncrement(RealIP(ctx), "view", "modpack:"+modpackID, time.Hour*4) {
		if err := dbModpack.Update().AddViews(1).Exec(ctx); err != nil {
			return nil, err
		}
	}

	return (*conv.ModpackImpl)(nil).Convert(dbModpack), nil
}

func (r *queryResolver) GetModpacks(_ context.Context, _ *generated.ModpackFilter) (*generated.GetModpacks, error) {
	return &generated.GetModpacks{}, nil
}

func (r *mutationResolver) CreateModpack(ctx context.Context, newModpack generated.NewModpack) (*generated.Modpack, error) {
	val := ctx.Value(util.ContextValidator{}).(*validator.Validate)
	if err := val.Struct(&newModpack); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	user, _, err := db.UserFromGQLContext(ctx)
	if err != nil {
		return nil, err
	}

	if newModpack.ParentID != nil {
		parent, err := db.From(ctx).Modpack.Query().
			Where(modpack.ID(*newModpack.ParentID)).
			First(ctx)
		if err != nil {
			return nil, err
		}

		if parent.ParentID != "" {
			return nil, fmt.Errorf("cannot create a remix of a remix")
		}
	}

	var resultModpack *ent.Modpack

	if err := db.Tx(ctx, func(ctx context.Context, tx *ent.Tx) error {
		dbModpack := tx.Modpack.Create().
			SetName(newModpack.Name).
			SetShortDescription(newModpack.ShortDescription).
			SetCreatorID(user.ID)

		SetINNF(newModpack.FullDescription, dbModpack.SetFullDescription)
		SetINNF(newModpack.Hidden, dbModpack.SetHidden)
		SetINNF(newModpack.ParentID, dbModpack.SetParentID)

		// Add tags
		if len(newModpack.TagIDs) > 0 {
			dbModpack = dbModpack.AddTagIDs(newModpack.TagIDs...)
		}

		resultModpack, err = dbModpack.Save(ctx)

		// Create targets
		for _, target := range newModpack.Targets {
			targets, err := tx.ModpackTarget.Create().
				SetModpackID(resultModpack.ID).
				SetTargetName(target).Save(ctx)
			if err != nil {
				return err
			}

			dbModpack = dbModpack.AddTargets(targets)
		}

		// Create modpack mod relationships
		for _, modInput := range newModpack.Mods {
			if err := tx.ModpackMod.Create().
				SetModpackID(resultModpack.ID).
				SetModID(modInput.ModID).
				SetVersionConstraint(modInput.VersionConstraint).
				Exec(ctx); err != nil {
				return err
			}
		}

		return err
	}, nil); err != nil {
		return nil, err
	}

	// Handle logo upload
	if newModpack.Logo != nil {
		file, err := io.ReadAll(newModpack.Logo.File)
		if err != nil {
			return nil, fmt.Errorf("failed to read logo file: %w", err)
		}

		logoData, thumbHash, err := converter.ConvertAnyImageToWebp(ctx, file)
		if err != nil {
			return nil, fmt.Errorf("failed to convert logo file: %w", err)
		}

		logoKey, err := storage.UploadModpackLogo(ctx, resultModpack.ID, bytes.NewReader(logoData))
		if err == nil {
			resultModpack, err = resultModpack.Update().
				SetLogo(storage.GenerateDownloadLink(ctx, logoKey)).
				SetLogoThumbhash(thumbHash).
				Save(ctx)
			if err != nil {
				return nil, err
			}
		}
	}

	// Get the modpack again with all relationships
	resultModpack, err = db.From(ctx).Modpack.Query().
		WithTags().
		WithTargets().
		WithModpackMods().
		WithParent().
		Where(modpack.ID(resultModpack.ID)).
		First(ctx)
	if err != nil {
		return nil, err
	}

	return (*conv.ModpackImpl)(nil).Convert(resultModpack), nil
}

func (r *mutationResolver) UpdateModpack(ctx context.Context, modpackID string, updatedModpack generated.UpdateModpack) (*generated.Modpack, error) {
	val := ctx.Value(util.ContextValidator{}).(*validator.Validate)
	if err := val.Struct(&updatedModpack); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	user, _, err := db.UserFromGQLContext(ctx)
	if err != nil {
		return nil, err
	}

	dbModpack, err := db.From(ctx).Modpack.Query().
		Where(modpack.ID(modpackID), modpack.CreatorID(user.ID)).
		First(ctx)
	if err != nil {
		return nil, err
	}

	update := dbModpack.Update()
	SetINNF(updatedModpack.Name, update.SetName)
	SetINNF(updatedModpack.ShortDescription, update.SetShortDescription)
	SetINNF(updatedModpack.FullDescription, update.SetFullDescription)
	SetINNF(updatedModpack.Hidden, update.SetHidden)

	resultModpack, err := update.Save(ctx)
	if err != nil {
		return nil, err
	}

	if updatedModpack.Logo != nil {
		file, err := io.ReadAll(updatedModpack.Logo.File)
		if err != nil {
			return nil, fmt.Errorf("failed to read logo file: %w", err)
		}

		logoData, thumbHash, err := converter.ConvertAnyImageToWebp(ctx, file)
		if err != nil {
			return nil, fmt.Errorf("failed to convert logo file: %w", err)
		}

		logoKey, err := storage.UploadModpackLogo(ctx, resultModpack.ID, bytes.NewReader(logoData))
		if err == nil {
			resultModpack, err = resultModpack.Update().
				SetLogo(storage.GenerateDownloadLink(ctx, logoKey)).
				SetLogoThumbhash(thumbHash).
				Save(ctx)
			if err != nil {
				return nil, err
			}
		}
	}

	if updatedModpack.TagIDs != nil {
		if err := resultModpack.Update().ClearTags().AddTagIDs(updatedModpack.TagIDs...).Exec(ctx); err != nil {
			return nil, err
		}
	}

	resultModpack, err = db.From(ctx).Modpack.Query().
		WithTags().
		WithTargets().
		WithModpackMods().
		Where(modpack.ID(resultModpack.ID)).
		First(ctx)
	if err != nil {
		return nil, err
	}

	return (*conv.ModpackImpl)(nil).Convert(resultModpack), nil
}

func (r *mutationResolver) DeleteModpack(ctx context.Context, modpackID string) (bool, error) {
	user, _, err := db.UserFromGQLContext(ctx)
	if err != nil {
		return false, err
	}

	dbModpack, err := db.From(ctx).Modpack.Query().
		Where(modpack.ID(modpackID), modpack.CreatorID(user.ID)).
		First(ctx)
	if err != nil {
		return false, err
	}

	if err := db.From(ctx).Modpack.DeleteOne(dbModpack).Exec(ctx); err != nil {
		return false, err
	}

	return true, nil
}

func (r *queryResolver) GetModpackRelease(ctx context.Context, modpackID string, version string) (*generated.ModpackRelease, error) {
	dbRelease, err := db.From(ctx).ModpackRelease.Query().
		Where(modpackrelease.HasModpackWith(modpack.ID(modpackID)), modpackrelease.Version(version)).
		First(ctx)
	if err != nil {
		return nil, err
	}

	return (*conv.ModpackReleaseImpl)(nil).Convert(dbRelease), nil
}

func (r *mutationResolver) CreateModpackRelease(ctx context.Context, modpackID string, release generated.NewModpackRelease) (*generated.ModpackRelease, error) {
	val := ctx.Value(util.ContextValidator{}).(*validator.Validate)
	if err := val.Struct(&release); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	user, _, err := db.UserFromGQLContext(ctx)
	if err != nil {
		return nil, err
	}

	dbModpack, err := db.From(ctx).Modpack.Query().
		WithTargets().
		Where(modpack.ID(modpackID), modpack.CreatorID(user.ID)).
		First(ctx)
	if err != nil {
		return nil, err
	}

	targetNames := make([]resolver.TargetName, len(dbModpack.Edges.Targets))
	for i, target := range dbModpack.Edges.Targets {
		targetNames[i] = resolver.TargetName(target.TargetName)
	}

	lockfile, err := resolveModpackToLockfile(ctx, modpackID, targetNames)
	if err != nil {
		return nil, err
	}

	dbRelease := db.From(ctx).ModpackRelease.Create().
		SetModpackID(dbModpack.ID).
		SetVersion(release.Version).
		SetChangelog(release.Changelog).
		SetLockfile(lockfile)

	resultRelease, err := dbRelease.Save(ctx)
	if err != nil {
		return nil, err
	}

	resultRelease, err = db.From(ctx).ModpackRelease.Query().
		Where(modpackrelease.ID(resultRelease.ID)).
		First(ctx)
	if err != nil {
		return nil, err
	}

	return (*conv.ModpackReleaseImpl)(nil).Convert(resultRelease), nil
}

func (r *mutationResolver) DeleteModpackRelease(ctx context.Context, modpackID string, version string) (bool, error) {
	user, _, err := db.UserFromGQLContext(ctx)
	if err != nil {
		return false, err
	}

	dbRelease, err := db.From(ctx).ModpackRelease.Query().
		Where(
			modpackrelease.HasModpackWith(modpack.ID(modpackID), modpack.CreatorID(user.ID)),
			modpackrelease.Version(version),
		).
		First(ctx)
	if err != nil {
		return false, err
	}

	if err := db.From(ctx).ModpackRelease.DeleteOne(dbRelease).Exec(ctx); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) ResolveModpack(ctx context.Context, modpackID string, targets []string) (*string, error) {
	targetNames := make([]resolver.TargetName, len(targets))
	for i, target := range targets {
		targetNames[i] = resolver.TargetName(target)
	}

	lockfile, err := resolveModpackToLockfile(ctx, modpackID, targetNames)
	if err != nil {
		return nil, err
	}

	return &lockfile, nil
}

// GetModpacks resolver type
type getModpacksResolver struct{ *Resolver }

func (r *getModpacksResolver) Modpacks(ctx context.Context, _ *generated.GetModpacks) ([]*generated.Modpack, error) {
	resolverContext := graphql.GetFieldContext(ctx)
	modpackFilter, err := models.ProcessModpackFilter(resolverContext.Parent.Args["filter"].(*generated.ModpackFilter))
	if err != nil {
		return nil, err
	}

	query := db.From(ctx).Modpack.Query().WithTags()
	query = db.ConvertModpackFilter(query, modpackFilter, false)

	result, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	return (*conv.ModpackImpl)(nil).ConvertSlice(result), nil
}

func (r *getModpacksResolver) Count(ctx context.Context, _ *generated.GetModpacks) (int, error) {
	resolverContext := graphql.GetFieldContext(ctx)
	modpackFilter, err := models.ProcessModpackFilter(resolverContext.Parent.Args["filter"].(*generated.ModpackFilter))
	if err != nil {
		return 0, err
	}

	query := db.From(ctx).Modpack.Query().WithTags()
	query = db.ConvertModpackFilter(query, modpackFilter, true)

	result, err := query.Count(ctx)
	if err != nil {
		return 0, err
	}

	return result, nil
}

// GetModpacks resolver type
type modpackResolver struct{ *Resolver }

func (r *modpackResolver) Children(ctx context.Context, m *generated.Modpack) ([]*generated.Modpack, error) {
	packs, err := db.From(ctx).Modpack.Query().
		WithTags().
		WithTargets().
		WithModpackMods().
		Where(modpack.ParentID(m.ID)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return (*conv.ModpackImpl)(nil).ConvertSlice(packs), nil
}

type lockfileResolver struct {
	Context context.Context
}

func (l lockfileResolver) ModVersionsWithDependencies(_ context.Context, modID string) ([]resolver.ModVersion, error) {
	//nolint:contextcheck
	slox.Info(l.Context, "resolver looking up mod", slog.String("id", modID))

	//nolint:contextcheck
	all, err := db.From(l.Context).Version.Query().
		WithVersionDependencies().
		WithTargets().
		Where(version.HasModWith(mod.Or(mod.ID(modID), mod.ModReference(modID)))).
		All(l.Context) //nolint:contextcheck
	if err != nil {
		return nil, fmt.Errorf("failed to query versions: %w", err)
	}

	result := make([]resolver.ModVersion, len(all))
	for i, v := range all {
		dependencies := make([]resolver.Dependency, len(v.Edges.VersionDependencies))
		for j, d := range v.Edges.VersionDependencies {
			dependencies[j] = resolver.Dependency{
				ModID:     d.ModID,
				Condition: d.Condition,
				Optional:  d.Optional,
			}
		}

		targets := make([]resolver.Target, len(v.Edges.Targets))
		for j, t := range v.Edges.Targets {
			targets[j] = resolver.Target{
				TargetName: resolver.TargetName(t.TargetName),
				Link:       "/v1/version/" + t.VersionID + "/" + t.TargetName + "/download",
				Hash:       t.Hash,
				Size:       t.Size,
			}
		}

		result[i] = resolver.ModVersion{
			Version:          v.Version,
			GameVersion:      v.GameVersion,
			Dependencies:     dependencies,
			Targets:          targets,
			RequiredOnRemote: v.RequiredOnRemote,
		}
	}

	return result, nil
}

func (l lockfileResolver) GetModName(_ context.Context, modReference string) (*resolver.ModName, error) {
	//nolint:contextcheck
	slox.Info(l.Context, "resolver looking up mod name", slog.String("reference", modReference))

	//nolint:contextcheck
	m, err := db.From(l.Context).Mod.Query().
		Where(mod.ModReference(modReference)).
		Select(mod.FieldID, mod.FieldModReference, mod.FieldName).
		First(l.Context) //nolint:contextcheck
	if err != nil {
		return nil, err
	}

	return &resolver.ModName{
		ID:           m.ID,
		ModReference: m.ModReference,
		Name:         m.Name,
	}, nil
}

func resolveModpackToLockfile(ctx context.Context, modpackID string, targets []resolver.TargetName) (string, error) {
	pack, err := db.From(ctx).Modpack.Query().
		WithModpackMods().
		WithParent(func(query *ent.ModpackQuery) {
			query.WithModpackMods()
		}).
		Where(modpack.ID(modpackID)).
		Only(ctx)
	if err != nil {
		return "", err
	}

	constraints := make(map[string]string)
	if pack.ParentID != "" {
		for _, m := range pack.Edges.Parent.Edges.ModpackMods {
			constraints[m.ModID] = m.VersionConstraint
		}
	}

	for _, m := range pack.Edges.ModpackMods {
		constraints[m.ModID] = m.VersionConstraint
	}

	modReferences, err := db.From(ctx).Mod.Query().
		Where(mod.IDIn(slices.Collect(maps.Keys(constraints))...)).
		Select(mod.FieldID, mod.FieldModReference).
		All(ctx)
	if err != nil {
		return "", err
	}

	referenceConstraints := make(map[string]string, len(constraints))
	for _, reference := range modReferences {
		referenceConstraints[reference.ModReference] = constraints[reference.ID]
	}

	dependencyResolver := resolver.NewDependencyResolver(lockfileResolver{
		Context: ctx,
	})

	lockfile, err := dependencyResolver.ResolveModDependencies(referenceConstraints, nil, math.MaxInt, targets)
	if err != nil {
		return "", fmt.Errorf("failed to resolve dependencies: %w", err)
	}

	b, err := json.Marshal(lockfile)
	if err != nil {
		return "", fmt.Errorf("failed to marshal lockfile: %w", err)
	}

	return string(b), nil
}
