package gql

import (
	"context"
	"fmt"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"

	"github.com/satisfactorymodding/smr-api/db/postgres"
	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/models"
	"github.com/satisfactorymodding/smr-api/util"
)

func (r *mutationResolver) CreateBootstrapVersion(ctx context.Context, bootstrapVersion generated.NewBootstrapVersion) (*generated.BootstrapVersion, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "createBootstrapVersion")
	defer wrapper.end()

	val := ctx.Value(util.ContextValidator{}).(*validator.Validate)
	if err := val.Struct(&bootstrapVersion); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	date, err := time.Parse(time.RFC3339Nano, bootstrapVersion.Date)
	if err != nil {
		return nil, fmt.Errorf("failed to parse date: %w", err)
	}

	dbBootstrapVersion := &postgres.BootstrapVersion{
		Version:             bootstrapVersion.Version,
		SatisfactoryVersion: bootstrapVersion.SatisfactoryVersion,
		Stability:           string(bootstrapVersion.Stability),
		Link:                bootstrapVersion.Link,
		Changelog:           bootstrapVersion.Changelog,
		Date:                date,
	}

	resultBootstrapVersion, err := postgres.CreateBootstrapVersion(newCtx, dbBootstrapVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to create bootstrap version: %w", err)
	}

	return DBBootstrapVersionToGenerated(resultBootstrapVersion), nil
}

func (r *mutationResolver) UpdateBootstrapVersion(ctx context.Context, bootstrapVersionID string, bootstrapVersion generated.UpdateBootstrapVersion) (*generated.BootstrapVersion, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "updateBootstrapVersion")
	defer wrapper.end()

	val := ctx.Value(util.ContextValidator{}).(*validator.Validate)
	if err := val.Struct(&bootstrapVersion); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	dbBootstrapVersion := postgres.GetBootstrapVersionByID(newCtx, bootstrapVersionID)

	if dbBootstrapVersion == nil {
		return nil, errors.New("bootstrapVersion not found")
	}

	SetStringINNOE(bootstrapVersion.Version, &dbBootstrapVersion.Version)
	SetINN(bootstrapVersion.SatisfactoryVersion, &dbBootstrapVersion.SatisfactoryVersion)
	SetStabilityINN(bootstrapVersion.Stability, &dbBootstrapVersion.Stability)
	SetStringINNOE(bootstrapVersion.Link, &dbBootstrapVersion.Link)
	SetStringINNOE(bootstrapVersion.Changelog, &dbBootstrapVersion.Changelog)
	SetDateINN(bootstrapVersion.Date, &dbBootstrapVersion.Date)

	postgres.Save(newCtx, &dbBootstrapVersion)

	return DBBootstrapVersionToGenerated(dbBootstrapVersion), nil
}

func (r *mutationResolver) DeleteBootstrapVersion(ctx context.Context, bootstrapVersionID string) (bool, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "deleteBootstrapVersion")
	defer wrapper.end()

	dbBootstrapVersion := postgres.GetBootstrapVersionByID(newCtx, bootstrapVersionID)

	if dbBootstrapVersion == nil {
		return false, errors.New("bootstrapVersion not found")
	}

	postgres.Delete(newCtx, &dbBootstrapVersion)

	return true, nil
}

func (r *queryResolver) GetBootstrapVersion(ctx context.Context, bootstrapVersionID string) (*generated.BootstrapVersion, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "getBootstrapVersion")
	defer wrapper.end()

	return DBBootstrapVersionToGenerated(postgres.GetBootstrapVersionByID(newCtx, bootstrapVersionID)), nil
}

func (r *queryResolver) GetBootstrapVersions(ctx context.Context, _ map[string]interface{}) (*generated.GetBootstrapVersions, error) {
	wrapper, _ := WrapQueryTrace(ctx, "getBootstrapVersions")
	defer wrapper.end()
	return &generated.GetBootstrapVersions{}, nil
}

type getBootstrapVersionsResolver struct{ *Resolver }

func (r *getBootstrapVersionsResolver) BootstrapVersions(ctx context.Context, _ *generated.GetBootstrapVersions) ([]*generated.BootstrapVersion, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "GetBootstrapVersions.bootstrapVersions")
	defer wrapper.end()

	resolverContext := graphql.GetFieldContext(ctx)
	bootstrapVersionFilter, err := models.ProcessBootstrapVersionFilter(resolverContext.Parent.Args["filter"].(map[string]interface{}))
	if err != nil {
		return nil, err
	}

	var bootstrapVersions []postgres.BootstrapVersion

	if bootstrapVersionFilter.Ids == nil || len(bootstrapVersionFilter.Ids) == 0 {
		bootstrapVersions = postgres.GetBootstrapVersions(newCtx, bootstrapVersionFilter)
	} else {
		bootstrapVersions = postgres.GetBootstrapVersionsByID(newCtx, bootstrapVersionFilter.Ids)
	}

	if bootstrapVersions == nil {
		return nil, errors.New("bootstrap releases not found")
	}

	converted := make([]*generated.BootstrapVersion, len(bootstrapVersions))
	for k, v := range bootstrapVersions {
		converted[k] = DBBootstrapVersionToGenerated(&v)
	}

	return converted, nil
}

func (r *getBootstrapVersionsResolver) Count(ctx context.Context, _ *generated.GetBootstrapVersions) (int, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "GetBootstrapVersions.count")
	defer wrapper.end()

	resolverContext := graphql.GetFieldContext(ctx)
	bootstrapVersionFilter, err := models.ProcessBootstrapVersionFilter(resolverContext.Parent.Args["filter"].(map[string]interface{}))
	if err != nil {
		return 0, err
	}

	if bootstrapVersionFilter.Ids != nil && len(bootstrapVersionFilter.Ids) != 0 {
		return len(bootstrapVersionFilter.Ids), nil
	}

	return int(postgres.GetBootstrapVersionCount(newCtx, bootstrapVersionFilter)), nil
}
