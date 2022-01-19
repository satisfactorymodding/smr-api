package gql

import (
	"context"
	"fmt"

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

	resultTag, err := postgres.CreateTag(newCtx, dbTag, true)
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

		resultTag, err := postgres.CreateTag(newCtx, dbTag, false)
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

	dbTag := postgres.GetTagByID(newCtx, id)

	if dbTag == nil {
		return false, errors.New("Tag not found")
	}

	postgres.Delete(newCtx, &dbTag)

	return true, nil
}

func (r *mutationResolver) UpdateTag(ctx context.Context, id string, newName string) (*generated.Tag, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "updateTag")
	defer wrapper.end()

	dbTag := postgres.GetTagByID(newCtx, id)

	if dbTag == nil {
		return nil, errors.New("Tag not found")
	}

	err := postgres.ValidateTagName(dbTag.Name)
	if err != nil {
		return nil, err
	}

	SetStringINNOE(&newName, &dbTag.Name)

	postgres.Save(newCtx, &dbTag)

	return DBTagToGenerated(dbTag), nil
}

func (r *queryResolver) GetTag(ctx context.Context, id string) (*generated.Tag, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "getTag_"+id)
	defer wrapper.end()

	tag := postgres.GetTagByID(newCtx, id)

	return DBTagToGenerated(tag), nil
}

func (r *queryResolver) GetTags(ctx context.Context, filter *generated.TagFilter) ([]*generated.Tag, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "getTags")
	defer wrapper.end()

	insertFilterDefaults(filter)
	tags := postgres.GetTags(newCtx, filter)

	return DBTagsToGeneratedSlice(tags), nil
}

func insertFilterDefaults(filter *generated.TagFilter) {
	fmt.Printf("%#v", filter)
	offset := 0
	limit := 10
	order := generated.OrderDesc
	if filter.Offset == nil {
		filter.Offset = &offset
	}
	if filter.Limit == nil {
		filter.Limit = &limit
	}
	if filter.Order == nil {
		filter.Order = &order
	}
}
