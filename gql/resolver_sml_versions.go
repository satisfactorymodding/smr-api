package gql

import (
	"context"
	"fmt"
	"strconv"

	"github.com/99designs/gqlgen/graphql"

	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/generated/conv"
	"github.com/satisfactorymodding/smr-api/generated/ent/mod"
	"github.com/satisfactorymodding/smr-api/generated/ent/version"
	"github.com/satisfactorymodding/smr-api/generated/ent/versiondependency"
	"github.com/satisfactorymodding/smr-api/models"
)

func (r *queryResolver) GetSMLVersion(ctx context.Context, smlVersionID string) (*generated.SMLVersion, error) {
	result, err := db.From(ctx).Version.
		Query().
		WithTargets().
		Where(version.ID(smlVersionID)).
		First(ctx)
	if err != nil {
		return nil, err
	}

	return (*conv.SMLVersionImpl)(nil).Convert(result), nil
}

func (r *queryResolver) GetSMLVersions(_ context.Context, _ map[string]interface{}) (*generated.GetSMLVersions, error) {
	return &generated.GetSMLVersions{}, nil
}

type getSMLVersionsResolver struct{ *Resolver }

func (r *getSMLVersionsResolver) SmlVersions(ctx context.Context, _ *generated.GetSMLVersions) ([]*generated.SMLVersion, error) {
	resolverContext := graphql.GetFieldContext(ctx)
	versionFilter, err := models.ProcessSMLVersionFilter(resolverContext.Parent.Args["filter"].(map[string]interface{}))
	if err != nil {
		return nil, fmt.Errorf("failed to process version filter: %w", err)
	}

	smlModQuery := db.From(ctx).Mod.Query().Where(mod.ModReference("SML"))
	smlMod, err := smlModQuery.First(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get SML mod: %w", err)
	}

	query := db.From(ctx).Version.Query().WithTargets().Where(version.ModID(smlMod.ID))
	query = convertVersionFilter(query, versionFilter, false)

	result, err := query.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get versions: %w", err)
	}

	return (*conv.SMLVersionImpl)(nil).ConvertSlice(result), nil
}

func (r *getSMLVersionsResolver) Count(ctx context.Context, _ *generated.GetSMLVersions) (int, error) {
	resolverContext := graphql.GetFieldContext(ctx)
	versionFilter, err := models.ProcessSMLVersionFilter(resolverContext.Parent.Args["filter"].(map[string]interface{}))
	if err != nil {
		return 0, err
	}

	smlModQuery := db.From(ctx).Mod.Query().Where(mod.ModReference("SML"))
	smlMod, err := smlModQuery.First(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get SML mod: %w", err)
	}

	query := db.From(ctx).Version.Query().WithTargets().Where(version.ModID(smlMod.ID))
	query = convertVersionFilter(query, versionFilter, false)

	result, err := query.Count(ctx)
	if err != nil {
		return 0, err
	}

	return result, nil
}

type smlVersionResolver struct{ *Resolver }

func (s *smlVersionResolver) SatisfactoryVersion(ctx context.Context, obj *generated.SMLVersion) (int, error) {
	query := db.From(ctx).VersionDependency.Query().Where(versiondependency.HasVersionWith(version.ID(obj.ID)))
	dependencies, err := query.All(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get dependencies: %w", err)
	}
	for _, dep := range dependencies {
		if dep.ModID == "FactoryGame" {
			ver := dep.Condition[2:] // Strip the >=
			versionNum, err := strconv.Atoi(ver)
			if err != nil {
				return 0, fmt.Errorf("failed to parse version number: %w", err)
			}
			return versionNum, nil
		}
	}
	return 0, fmt.Errorf("no satisfactory version found for SML version")
}

func (s *smlVersionResolver) Link(_ context.Context, obj *generated.SMLVersion) (string, error) {
	return "https://github.com/satisfactorymodding/SatisfactoryModLoader/releases/tag/v" + obj.Version, nil
}

func (s *smlVersionResolver) BootstrapVersion(ctx context.Context, obj *generated.SMLVersion) (*string, error) {
	query := db.From(ctx).VersionDependency.Query().Where(versiondependency.HasVersionWith(version.ID(obj.ID)))
	dependencies, err := query.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get dependencies: %w", err)
	}
	for _, dep := range dependencies {
		if dep.ModID == "bootstrap" {
			return &dep.Condition, nil
		}
	}
	return nil, nil
}

func (s *smlVersionResolver) EngineVersion(ctx context.Context, obj *generated.SMLVersion) (string, error) {
	satisfactoryVersion, err := s.SatisfactoryVersion(ctx, obj)
	if err != nil {
		return "", fmt.Errorf("failed to get satisfactory version: %w", err)
	}
	v, err := db.GetEngineVersionForSatisfactoryVersion(ctx, satisfactoryVersion)
	if err != nil {
		return "", fmt.Errorf("failed to get engine version: %w", err)
	}
	return v, nil
}

type smlVersionTargetResolver struct{ *Resolver }

func (s *smlVersionTargetResolver) Link(ctx context.Context, obj *generated.SMLVersionTarget) (string, error) {
	query := db.From(ctx).Version.Query().WithTargets().Where(version.ID(obj.VersionID))
	v, err := query.First(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get version: %w", err)
	}
	if len(v.Edges.Targets) > 1 {
		return fmt.Sprintf("https://github.com/satisfactorymodding/SatisfactoryModLoader/releases/download/v%s/SML-%s.zip", v.Version, obj.TargetName), nil
	}
	return fmt.Sprintf("https://github.com/satisfactorymodding/SatisfactoryModLoader/releases/download/v%s/SML.zip", v.Version), nil
}
