package gql

import (
	"context"

	"github.com/pkg/errors"
	"github.com/satisfactorymodding/smr-api/util"
	"gopkg.in/go-playground/validator.v9"

	"github.com/satisfactorymodding/smr-api/db/postgres"
	"github.com/satisfactorymodding/smr-api/generated"
)

func (r *mutationResolver) CreateSMLArch(ctx context.Context, smlLink generated.NewSMLArch) (*generated.SMLArch, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "createSMLArch")
	defer wrapper.end()

	val := ctx.Value(util.ContextValidator{}).(*validator.Validate)
	if err := val.Struct(&smlLink); err != nil {
		return nil, errors.Wrap(err, "validation failed")
	}

	dbSMLArchs := &postgres.SMLArch{
		ID:               string(util.GenerateUniqueID()),
		SMLVersionArchID: string(smlLink.SMLVersionArchID),
		Platform:         string(smlLink.Platform),
		Link:             string(smlLink.Link),
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

func (r *mutationResolver) UpdateSMLArch(ctx context.Context, smlLinkID string, smlLink generated.UpdateSMLArch) (*generated.SMLArch, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "updateSMLArch")
	defer wrapper.end()
	val := ctx.Value(util.ContextValidator{}).(*validator.Validate)
	if err := val.Struct(&smlLink); err != nil {
		return nil, errors.Wrap(err, "validation failed")
	}

	dbSMLArch := postgres.GetSMLArch(newCtx, smlLinkID)

	if dbSMLArch == nil {
		return nil, errors.New("sml link not found")
	}

	SetStringINNOE((*string)(&smlLink.SMLVersionArchID), &dbSMLArch.SMLVersionArchID)
	SetStringINNOE((*string)(&smlLink.Platform), &dbSMLArch.Platform)
	SetStringINNOE((*string)(&smlLink.Link), &dbSMLArch.Link)

	postgres.Save(newCtx, &dbSMLArch)

	return DBSMLArchToGenerated(dbSMLArch), nil
}

func (r *queryResolver) GetSMLArch(ctx context.Context, smlLinkID string) (*generated.SMLArch, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "getSMLArch")
	defer wrapper.end()

	smlLink := postgres.GetSMLArch(newCtx, smlLinkID)

	return DBSMLArchToGenerated(smlLink), nil
}

func (r *queryResolver) GetSMLArchs(ctx context.Context, filter map[string]interface{}) (*generated.GetSMLArchs, error) {
	wrapper, _ := WrapQueryTrace(ctx, "getSMLArchs")
	defer wrapper.end()
	return &generated.GetSMLArchs{}, nil
}

func (r *queryResolver) GetSMLDownload(ctx context.Context, smlVersionID string, platform string) (string, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "getSMLDownload")
	defer wrapper.end()
	smlLink := postgres.GetSMLArchDownload(newCtx, smlVersionID, platform)
	return smlLink, nil
}
