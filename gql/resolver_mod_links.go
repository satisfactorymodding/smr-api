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

	dbModLinks := &postgres.ModLink{
		ID:               string(util.GenerateUniqueID()),
		ModVersionLinkID: string(modLink.ModVersionLinkID),
		Platform:         string(modLink.Platform),
		Link:             string(modLink.Link),
		Hash:             *modLink.Hash,
		Size:             int64(*modLink.Size),
	}

	resultModLink, err := postgres.CreateModLink(newCtx, dbModLinks)

	if err != nil {
		return nil, err
	}

	return DBModLinkToGenerated(resultModLink), nil
}

func (r *mutationResolver) DeleteModLink(ctx context.Context, linksID string) (bool, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "deleteModLink")
	defer wrapper.end()

	dbModLink := postgres.GetModLinkByID(newCtx, linksID)

	if dbModLink == nil {
		return false, errors.New("Mod Link not found")
	}

	postgres.Delete(newCtx, &dbModLink)

	return true, nil
}

func (r *queryResolver) GetModLink(ctx context.Context, modLinkID string) (*generated.ModLink, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "getModLink")
	defer wrapper.end()

	modLink := postgres.GetModLink(newCtx, modLinkID)

	return DBModLinkToGenerated(modLink), nil
}

func (r *queryResolver) GetModLinks(ctx context.Context, filter map[string]interface{}) (*generated.GetModLinks, error) {
	wrapper, _ := WrapQueryTrace(ctx, "getModLinks")
	defer wrapper.end()
	return &generated.GetModLinks{}, nil
}

func (r *queryResolver) GetModLinkByID(ctx context.Context, modLinkID string) (*generated.ModLink, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "getModLink")
	defer wrapper.end()

	modLink := postgres.GetModLinkByID(newCtx, modLinkID)

	return DBModLinkToGenerated(modLink), nil
}

type getModLinksResolver struct{ *Resolver }
