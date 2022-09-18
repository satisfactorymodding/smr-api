package gql

import (
	"context"

	"github.com/satisfactorymodding/smr-api/db/postgres"
	"github.com/satisfactorymodding/smr-api/generated"
)

type modlinkResolver struct{ *Resolver }

func (r *modlinkResolver) Asset(_ context.Context, obj *generated.ModArch) (string, error) {
	return "/v1/version/" + obj.ModVersionArchID + "/" + obj.Platform + "/download", nil
}

func (r *queryResolver) GetModArch(ctx context.Context, modArchID string) (*generated.ModArch, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "getModArch")
	defer wrapper.end()
	modArch := postgres.GetModArch(newCtx, modArchID)
	return DBModArchToGenerated(modArch), nil
}

func (r *queryResolver) GetModArchs(ctx context.Context, filter map[string]interface{}) (*generated.GetModArchs, error) {
	wrapper, _ := WrapQueryTrace(ctx, "getModArchs")
	defer wrapper.end()
	return &generated.GetModArchs{}, nil
}

func (r *queryResolver) GetModArchByID(ctx context.Context, modArchID string) (*generated.ModArch, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "getModArchByID")
	defer wrapper.end()
	modArch := postgres.GetModArchByID(newCtx, modArchID)
	return DBModArchToGenerated(modArch), nil
}
