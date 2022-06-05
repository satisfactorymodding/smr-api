package gql

import (
	"context"
	"time"

	"github.com/satisfactorymodding/smr-api/db/postgres"
	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/models"
	"github.com/satisfactorymodding/smr-api/util"

	"github.com/pkg/errors"

	"github.com/99designs/gqlgen/graphql"
	"gopkg.in/go-playground/validator.v9"
)

func (r *mutationResolver) CreateSMLVersion(ctx context.Context, smlVersion generated.NewSMLVersion) (*generated.SMLVersion, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "createSMLVersion")
	defer wrapper.end()

	val := ctx.Value(util.ContextValidator{}).(*validator.Validate)
	if err := val.Struct(&smlVersion); err != nil {
		return nil, errors.Wrap(err, "validation failed")
	}

	date, err := time.Parse(time.RFC3339Nano, smlVersion.Date)

	if err != nil {
		return nil, errors.Wrap(err, "failed to parse date")
	}

	dbSMLVersion := &postgres.SMLVersion{
		Version:             smlVersion.Version,
		SatisfactoryVersion: smlVersion.SatisfactoryVersion,
		BootstrapVersion:    smlVersion.BootstrapVersion,
		Stability:           string(smlVersion.Stability),
		Link:                smlVersion.Link,
		Changelog:           smlVersion.Changelog,
		Date:                date,
	}

	resultSMLVersion, err := postgres.CreateSMLVersion(newCtx, dbSMLVersion)

	for _, smlLink := range smlVersion.Links {
		dbSMLLinks := &postgres.SMLLink{
			ID:               util.GenerateUniqueID(),
			SMLVersionLinkID: string(resultSMLVersion.ID),
			Platform:         string(smlLink.Platform),
			//Side:             string(smlLink.Side),
			Link: string(smlLink.Link),
		}

		resultSMLLink, err := postgres.CreateSMLLink(newCtx, dbSMLLinks)

		if err != nil {
			return nil, err
		}

		DBSMLLinkToGenerated(resultSMLLink)

	}

	if err != nil {
		return nil, err
	}

	return DBSMLVersionToGenerated(resultSMLVersion), nil
}

func (r *mutationResolver) UpdateSMLVersion(ctx context.Context, smlVersionID string, smlVersion generated.UpdateSMLVersion) (*generated.SMLVersion, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "updateSMLVersion")
	defer wrapper.end()

	val := ctx.Value(util.ContextValidator{}).(*validator.Validate)
	if err := val.Struct(&smlVersion); err != nil {
		return nil, errors.Wrap(err, "validation failed")
	}

	dbSMLVersion := postgres.GetSMLVersionByID(newCtx, smlVersionID)

	if dbSMLVersion == nil {
		return nil, errors.New("smlVersion not found")
	}

	SetStringINNOE(smlVersion.Version, &dbSMLVersion.Version)
	SetIntINN(smlVersion.SatisfactoryVersion, &dbSMLVersion.SatisfactoryVersion)
	SetStringINNOE(smlVersion.BootstrapVersion, dbSMLVersion.BootstrapVersion)
	SetStabilityINN(smlVersion.Stability, &dbSMLVersion.Stability)
	SetStringINNOE(smlVersion.Link, &dbSMLVersion.Link)
	SetStringINNOE(smlVersion.Changelog, &dbSMLVersion.Changelog)
	SetDateINN(smlVersion.Date, &dbSMLVersion.Date)

	postgres.Save(newCtx, &dbSMLVersion)

	return DBSMLVersionToGenerated(dbSMLVersion), nil
}

func (r *mutationResolver) DeleteSMLVersion(ctx context.Context, smlVersionID string) (bool, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "deleteSMLVersion")
	defer wrapper.end()

	dbSMLVersion := postgres.GetSMLVersionByID(newCtx, smlVersionID)

	if dbSMLVersion == nil {
		return false, errors.New("smlVersion not found")
	}

	postgres.Delete(newCtx, &dbSMLVersion)

	return true, nil
}

func (r *queryResolver) GetSMLVersion(ctx context.Context, smlVersionID string) (*generated.SMLVersion, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "getSMLVersion")
	defer wrapper.end()

	return DBSMLVersionToGenerated(postgres.GetSMLVersionByID(newCtx, smlVersionID)), nil
}

func (r *queryResolver) GetSMLVersions(ctx context.Context, filter map[string]interface{}) (*generated.GetSMLVersions, error) {
	wrapper, _ := WrapQueryTrace(ctx, "getSMLVersions")
	defer wrapper.end()
	return &generated.GetSMLVersions{}, nil
}

type getSMLVersionsResolver struct{ *Resolver }

func (r *getSMLVersionsResolver) SmlVersions(ctx context.Context, obj *generated.GetSMLVersions) ([]*generated.SMLVersion, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "GetSMLVersions.smlVersions")
	defer wrapper.end()

	resolverContext := graphql.GetFieldContext(ctx)
	smlVersionFilter, err := models.ProcessSMLVersionFilter(resolverContext.Parent.Args["filter"].(map[string]interface{}))

	if err != nil {
		return nil, err
	}

	var smlVersions []postgres.SMLVersion

	if smlVersionFilter.Ids == nil || len(smlVersionFilter.Ids) == 0 {
		smlVersions = postgres.GetSMLVersions(newCtx, smlVersionFilter)
	} else {
		smlVersions = postgres.GetSMLVersionsByID(newCtx, smlVersionFilter.Ids)
	}

	if smlVersions == nil {
		return nil, errors.New("sml releases not found")
	}

	converted := make([]*generated.SMLVersion, len(smlVersions))
	for k, v := range smlVersions {
		converted[k] = DBSMLVersionToGenerated(&v)
	}

	return converted, nil
}

func (r *getSMLVersionsResolver) Count(ctx context.Context, obj *generated.GetSMLVersions) (int, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "GetSMLVersions.count")
	defer wrapper.end()

	resolverContext := graphql.GetFieldContext(ctx)
	smlVersionFilter, err := models.ProcessSMLVersionFilter(resolverContext.Parent.Args["filter"].(map[string]interface{}))

	if err != nil {
		return 0, err
	}

	if smlVersionFilter.Ids != nil && len(smlVersionFilter.Ids) != 0 {
		return len(smlVersionFilter.Ids), nil
	}

	return int(postgres.GetSMLVersionCount(newCtx, smlVersionFilter)), nil
}
