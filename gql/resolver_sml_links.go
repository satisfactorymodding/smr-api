package gql

import (
	"context"

	"github.com/pkg/errors"
	"github.com/satisfactorymodding/smr-api/util"
	"gopkg.in/go-playground/validator.v9"

	"github.com/satisfactorymodding/smr-api/db/postgres"
	"github.com/satisfactorymodding/smr-api/generated"
)

func (r *mutationResolver) CreateSMLLink(ctx context.Context, smlLink generated.NewSMLLink) (*generated.SMLLink, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "createSMLLink")
	defer wrapper.end()

	val := ctx.Value(util.ContextValidator{}).(*validator.Validate)
	if err := val.Struct(&smlLink); err != nil {
		return nil, errors.Wrap(err, "validation failed")
	}

	dbSMLLinks := &postgres.SMLLink{
		ID:               string(util.GenerateUniqueID()),
		SMLVersionLinkID: string(smlLink.SMLVersionLinkID),
		Platform:         string(smlLink.Platform),
		Link:             string(smlLink.Link),
	}

	resultSMLLink, err := postgres.CreateSMLLink(newCtx, dbSMLLinks)

	if err != nil {
		return nil, err
	}

	return DBSMLLinkToGenerated(resultSMLLink), nil
}

func (r *mutationResolver) DeleteSMLLink(ctx context.Context, linksID string) (bool, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "deleteSMLLink")
	defer wrapper.end()

	dbSMLLink := postgres.GetSMLLinkByID(newCtx, linksID)

	if dbSMLLink == nil {
		return false, errors.New("SML Link not found")
	}

	postgres.Delete(newCtx, &dbSMLLink)

	return true, nil
}

func (r *mutationResolver) UpdateSMLLink(ctx context.Context, smlLinkID string, smlLink generated.UpdateSMLLink) (*generated.SMLLink, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "updateSMLLink")
	defer wrapper.end()
	val := ctx.Value(util.ContextValidator{}).(*validator.Validate)
	if err := val.Struct(&smlLink); err != nil {
		return nil, errors.Wrap(err, "validation failed")
	}

	dbSMLLink := postgres.GetSMLLink(newCtx, smlLinkID)

	if dbSMLLink == nil {
		return nil, errors.New("sml link not found")
	}

	SetStringINNOE((*string)(&smlLink.SMLVersionLinkID), &dbSMLLink.SMLVersionLinkID)
	SetStringINNOE((*string)(&smlLink.Platform), &dbSMLLink.Platform)
	SetStringINNOE((*string)(&smlLink.Link), &dbSMLLink.Link)

	postgres.Save(newCtx, &dbSMLLink)

	return DBSMLLinkToGenerated(dbSMLLink), nil
}

func (r *queryResolver) GetSMLLink(ctx context.Context, smlLinkID string) (*generated.SMLLink, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "getSMLLink")
	defer wrapper.end()

	smlLink := postgres.GetSMLLink(newCtx, smlLinkID)

	return DBSMLLinkToGenerated(smlLink), nil
}

func (r *queryResolver) GetSMLLinks(ctx context.Context, filter map[string]interface{}) (*generated.GetSMLLinks, error) {
	wrapper, _ := WrapQueryTrace(ctx, "getSMLLinks")
	defer wrapper.end()
	return &generated.GetSMLLinks{}, nil
}

type getSMLLinksResolver struct{ *Resolver }
