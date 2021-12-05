package gql

import (
	"context"

	"github.com/pkg/errors"
	"github.com/satisfactorymodding/smr-api/db/postgres"
	"github.com/satisfactorymodding/smr-api/generated"
)

func (r *mutationResolver) CreateModTag(ctx context.Context, modTagName string) (*generated.ModTag, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "createModTag")
	defer wrapper.end()

	dbModTag := &postgres.ModTag{
		Name: modTagName,
	}

	resultModTag, err := postgres.CreateModTag(dbModTag, &newCtx)
	if err != nil {
		return nil, err
	}
	return DBModTagToGenerated(resultModTag), nil
}

func (r *mutationResolver) DeleteModTag(ctx context.Context, modTagName string) (bool, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "deleteModTag")
	defer wrapper.end()

	dbModTag := postgres.GetModTagByName(modTagName, &newCtx)

	if dbModTag == nil {
		return false, errors.New("ModTag not found")
	}

	postgres.Delete(&dbModTag, &newCtx)

	return true, nil
}

func (r *mutationResolver) UpdateModTag(ctx context.Context, modTagName string, newName string) (*generated.ModTag, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "updateModTag")
	defer wrapper.end()

	dbModTag := postgres.GetModTagByName(modTagName, &newCtx)

	if dbModTag == nil {
		return nil, errors.New("ModTag not found")
	}

	//TODO Use validation in here too

	SetStringINNOE(&newName, &dbModTag.Name)

	postgres.Save(&dbModTag, &newCtx)

	return DBModTagToGenerated(dbModTag), nil
}

func (r *queryResolver) GetModTag(ctx context.Context, modTagName string) (*generated.ModTag, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "getModTag")
	defer wrapper.end()

	modTag := postgres.GetModTagByName(modTagName, &newCtx)

	return DBModTagToGenerated(modTag), nil
}

func (r *queryResolver) GetModTags(ctx context.Context, filter *generated.ModTagFilter) ([]*generated.ModTag, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "getModTags")
	defer wrapper.end()

	modTags := postgres.GetModTags(&newCtx)

	return DBModTagsToGeneratedSlice(modTags), nil
}
