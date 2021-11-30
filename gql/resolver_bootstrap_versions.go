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

func (r *mutationResolver) CreateBootstrapVersion(ctx context.Context, bootstrapVersion generated.NewBootstrapVersion) (*generated.BootstrapVersion, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "createBootstrapVersion")
	defer wrapper.end()

	val := ctx.Value(util.ContextValidator{}).(*validator.Validate)
	if err := val.Struct(&bootstrapVersion); err != nil {
		return nil, errors.Wrap(err, "validation failed")
	}

	date, err := time.Parse(time.RFC3339Nano, bootstrapVersion.Date)

	if err != nil {
		return nil, errors.Wrap(err, "failed to parse date")
	}

	dbBootstrapVersion := &postgres.BootstrapVersion{
		Version:             bootstrapVersion.Version,
		SatisfactoryVersion: bootstrapVersion.SatisfactoryVersion,
		Stability:           string(bootstrapVersion.Stability),
		Link:                bootstrapVersion.Link,
		Changelog:           bootstrapVersion.Changelog,
		Date:                date,
	}

	resultBootstrapVersion, err := postgres.CreateBootstrapVersion(dbBootstrapVersion, &newCtx)

	if err != nil {
		return nil, errors.Wrap(err, "failed to create bootstrap version")
	}

	return DBBootstrapVersionToGenerated(resultBootstrapVersion), nil
}

func (r *mutationResolver) UpdateBootstrapVersion(ctx context.Context, bootstrapVersionID string, bootstrapVersion generated.UpdateBootstrapVersion) (*generated.BootstrapVersion, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "updateBootstrapVersion")
	defer wrapper.end()

	val := ctx.Value(util.ContextValidator{}).(*validator.Validate)
	if err := val.Struct(&bootstrapVersion); err != nil {
		return nil, errors.Wrap(err, "validation failed")
	}

	dbBootstrapVersion := postgres.GetBootstrapVersionByID(bootstrapVersionID, &newCtx)

	if dbBootstrapVersion == nil {
		return nil, errors.New("bootstrapVersion not found")
	}

	SetStringINNOE(bootstrapVersion.Version, &dbBootstrapVersion.Version)
	SetIntINN(bootstrapVersion.SatisfactoryVersion, &dbBootstrapVersion.SatisfactoryVersion)
	SetStabilityINN(bootstrapVersion.Stability, &dbBootstrapVersion.Stability)
	SetStringINNOE(bootstrapVersion.Link, &dbBootstrapVersion.Link)
	SetStringINNOE(bootstrapVersion.Changelog, &dbBootstrapVersion.Changelog)
	SetDateINN(bootstrapVersion.Date, &dbBootstrapVersion.Date)

	postgres.Save(&dbBootstrapVersion, &newCtx)

	return DBBootstrapVersionToGenerated(dbBootstrapVersion), nil
}

func (r *mutationResolver) DeleteBootstrapVersion(ctx context.Context, bootstrapVersionID string) (bool, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "deleteBootstrapVersion")
	defer wrapper.end()

	dbBootstrapVersion := postgres.GetBootstrapVersionByID(bootstrapVersionID, &newCtx)

	if dbBootstrapVersion == nil {
		return false, errors.New("bootstrapVersion not found")
	}

	postgres.Delete(&dbBootstrapVersion, &newCtx)

	return true, nil
}

func (r *queryResolver) GetBootstrapVersion(ctx context.Context, bootstrapVersionID string) (*generated.BootstrapVersion, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "getBootstrapVersion")
	defer wrapper.end()

	return DBBootstrapVersionToGenerated(postgres.GetBootstrapVersionByID(bootstrapVersionID, &newCtx)), nil
}

func (r *queryResolver) GetBootstrapVersions(ctx context.Context, filter map[string]interface{}) (*generated.GetBootstrapVersions, error) {
	wrapper, _ := WrapQueryTrace(ctx, "getBootstrapVersions")
	defer wrapper.end()
	return &generated.GetBootstrapVersions{}, nil
}

type getBootstrapVersionsResolver struct{ *Resolver }

func (r *getBootstrapVersionsResolver) BootstrapVersions(ctx context.Context, obj *generated.GetBootstrapVersions) ([]*generated.BootstrapVersion, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "GetBootstrapVersions.bootstrapVersions")
	defer wrapper.end()

	resolverContext := graphql.GetFieldContext(ctx)
	bootstrapVersionFilter, err := models.ProcessBootstrapVersionFilter(resolverContext.Parent.Args["filter"].(map[string]interface{}))

	if err != nil {
		return nil, err
	}

	var bootstrapVersions []postgres.BootstrapVersion

	if bootstrapVersionFilter.Ids == nil || len(bootstrapVersionFilter.Ids) == 0 {
		bootstrapVersions = postgres.GetBootstrapVersions(bootstrapVersionFilter, &newCtx)
	} else {
		bootstrapVersions = postgres.GetBootstrapVersionsByID(bootstrapVersionFilter.Ids, &newCtx)
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

func (r *getBootstrapVersionsResolver) Count(ctx context.Context, obj *generated.GetBootstrapVersions) (int, error) {
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

	return int(postgres.GetBootstrapVersionCount(bootstrapVersionFilter, &newCtx)), nil
}
