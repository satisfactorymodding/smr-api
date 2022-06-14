package gql

import (
	"context"

	"github.com/pkg/errors"
	"github.com/satisfactorymodding/smr-api/util"
	"gopkg.in/go-playground/validator.v9"

	"github.com/satisfactorymodding/smr-api/db/postgres"
	"github.com/satisfactorymodding/smr-api/generated"
)

func (r *mutationResolver) CreateModLink(ctx context.Context, modLink generated.NewModLink) (*generated.ModLink, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "createModLink")
	defer wrapper.end()

	val := ctx.Value(util.ContextValidator{}).(*validator.Validate)
	if err := val.Struct(&modLink); err != nil {
		return nil, errors.Wrap(err, "validation failed")
	}

	dbModLink := &postgres.ModLink{
		Platform: string(modLink.Platform),
		//Side:     string(modLink.Side),
		Link: string(modLink.Link),
	}

	resultModLink, err := postgres.CreateModLink(newCtx, dbModLink)
	if err != nil {
		return nil, err
	}
	return DBModLinkToGenerated(resultModLink), nil
}

func (r *mutationResolver) DeleteModLink(ctx context.Context, modLinksID string) (bool, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "deleteModLink")
	defer wrapper.end()

	dbModLink := postgres.GetModLinkByID(ctx, modLinksID)

	if dbModLink == nil {
		return false, errors.New("Mod Link not found")
	}

	postgres.Delete(newCtx, &dbModLink)

	return true, nil
}

func (r *mutationResolver) UpdateModLink(ctx context.Context, modLinkID string, modLink generated.UpdateModLink) (*generated.ModLink, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "updateModLink")
	defer wrapper.end()

	val := ctx.Value(util.ContextValidator{}).(*validator.Validate)
	if err := val.Struct(&modLink); err != nil {
		return nil, errors.Wrap(err, "validation failed")
	}

	dbModLink := postgres.GetModLinkByID(newCtx, modLinkID)

	if dbModLink == nil {
		return nil, errors.New("guide not found")
	}

	SetStringINNOE((*string)(&modLink.Platform), &dbModLink.Platform)
	SetStringINNOE((*string)(&modLink.Link), &dbModLink.Link)

	postgres.Save(newCtx, &dbModLink)

	return DBModLinkToGenerated(dbModLink), nil
}

func (r *queryResolver) GetModLink(ctx context.Context, modLinkID string) (*generated.ModLink, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "getModLink")
	defer wrapper.end()

	modLink := postgres.GetModLinkByID(newCtx, modLinkID)

	return DBModLinkToGenerated(modLink), nil
}
