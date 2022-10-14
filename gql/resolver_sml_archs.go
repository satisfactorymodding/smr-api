package gql

import (
	"context"

	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"

	"github.com/satisfactorymodding/smr-api/util"

	"github.com/satisfactorymodding/smr-api/db/postgres"
	"github.com/satisfactorymodding/smr-api/generated"
)

func (r *mutationResolver) CreateSMLArch(ctx context.Context, smlArch generated.NewSMLArch) (*generated.SMLArch, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "createSMLArch")
	defer wrapper.end()

	val := ctx.Value(util.ContextValidator{}).(*validator.Validate)
	if err := val.Struct(&smlArch); err != nil {
		return nil, errors.Wrap(err, "validation failed")
	}

	dbSMLArchs := &postgres.SMLArch{
		ID:       util.GenerateUniqueID(),
		Platform: smlArch.Platform,
		Link:     smlArch.Link,
	}

	resultSMLArch, err := postgres.CreateSMLArch(newCtx, dbSMLArchs)

	if err != nil {
		return nil, err
	}

	return DBSMLArchToGenerated(resultSMLArch), nil
}

func (r *mutationResolver) DeleteSMLArch(ctx context.Context, linksID string) (bool, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "deleteSMLArch")
	defer wrapper.end()

	dbSMLArch := postgres.GetSMLArchByID(newCtx, linksID)

	if dbSMLArch == nil {
		return false, errors.New("SML Link not found")
	}

	postgres.Delete(newCtx, &dbSMLArch)

	return true, nil
}

func (r *mutationResolver) UpdateSMLArch(ctx context.Context, smlLinkID string, smlArch generated.UpdateSMLArch) (*generated.SMLArch, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "updateSMLArch")
	defer wrapper.end()
	val := ctx.Value(util.ContextValidator{}).(*validator.Validate)
	if err := val.Struct(&smlArch); err != nil {
		return nil, errors.Wrap(err, "validation failed")
	}

	dbSMLArch := postgres.GetSMLArch(newCtx, smlLinkID)

	if dbSMLArch == nil {
		return nil, errors.New("sml link not found")
	}

	SetStringINNOE(&smlArch.Platform, &dbSMLArch.Platform)
	SetStringINNOE(&smlArch.Link, &dbSMLArch.Link)

	postgres.Save(newCtx, &dbSMLArch)

	return DBSMLArchToGenerated(dbSMLArch), nil
}

func (r *queryResolver) GetSMLArch(ctx context.Context, smlLinkID string) (*generated.SMLArch, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "getSMLArch")
	defer wrapper.end()

	smlArch := postgres.GetSMLArch(newCtx, smlLinkID)

	return DBSMLArchToGenerated(smlArch), nil
}

func (r *queryResolver) GetSMLArchs(ctx context.Context, filter map[string]interface{}) (*generated.GetSMLArchs, error) {
	wrapper, _ := WrapQueryTrace(ctx, "getSMLArchs")
	defer wrapper.end()
	return &generated.GetSMLArchs{}, nil
}

func (r *queryResolver) GetSMLArchBySMLID(ctx context.Context, smlVersionID string) ([]postgres.SMLArch, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "GetSMLArchBySMLID")
	defer wrapper.end()
	smlArch := postgres.GetSMLArchBySMLID(newCtx, smlVersionID)
	return smlArch, nil
}

func (r *queryResolver) GetSMLDownload(ctx context.Context, smlVersionID string, platform string) (string, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "getSMLDownload")
	defer wrapper.end()
	smlArch := postgres.GetSMLArchDownload(newCtx, smlVersionID, platform)
	return smlArch, nil
}
