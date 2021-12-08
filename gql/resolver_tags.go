package gql

import (
	"context"

	"github.com/pkg/errors"
	"github.com/satisfactorymodding/smr-api/db/postgres"
	"github.com/satisfactorymodding/smr-api/generated"
)

func (r *mutationResolver) CreateTag(ctx context.Context, tagName string) (*generated.Tag, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "createTag")
	defer wrapper.end()

	dbTag := &postgres.Tag{
		Name: tagName,
	}

	resultTag, err := postgres.CreateTag(dbTag, &newCtx, true)
	if err != nil {
		return nil, err
	}
	return DBTagToGenerated(resultTag), nil
}

func (r *mutationResolver) CreateMultipleTags(ctx context.Context, tagNames []string) ([]*generated.Tag, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "createMultipleTags")
	defer wrapper.end()

	resultTags := make([]postgres.Tag, len(tagNames))

	for i, tagName := range tagNames {
		dbTag := &postgres.Tag{
			Name: tagName,
		}

		resultTag, err := postgres.CreateTag(dbTag, &newCtx, false)
		if err != nil {
			return nil, err
		}

		resultTags[i] = *resultTag
	}

	return DBTagsToGeneratedSlice(resultTags), nil
}

func (r *mutationResolver) DeleteTag(ctx context.Context, id string) (bool, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "deleteTag")
	defer wrapper.end()

	dbTag := postgres.GetTagByID(id, &newCtx)

	if dbTag == nil {
		return false, errors.New("Tag not found")
	}

	postgres.Delete(&dbTag, &newCtx)

	return true, nil
}

func (r *mutationResolver) UpdateTag(ctx context.Context, id string, newName string) (*generated.Tag, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "updateTag")
	defer wrapper.end()

	dbTag := postgres.GetTagByID(id, &newCtx)

	if dbTag == nil {
		return nil, errors.New("Tag not found")
	}

	err := postgres.ValidateTagName(dbTag.Name)
	if err != nil {
		return nil, err
	}

	SetStringINNOE(&newName, &dbTag.Name)

	postgres.Save(&dbTag, &newCtx)

	return DBTagToGenerated(dbTag), nil
}

func (r *queryResolver) GetTag(ctx context.Context, id string) (*generated.Tag, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "getTag_"+id)
	defer wrapper.end()

	tag := postgres.GetTagByID(id, &newCtx)

	return DBTagToGenerated(tag), nil
}

func (r *queryResolver) GetTags(ctx context.Context, filter *generated.TagFilter) ([]*generated.Tag, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "getTags")
	defer wrapper.end()

	tags := postgres.GetTags(&newCtx)

	return DBTagsToGeneratedSlice(tags), nil
}
