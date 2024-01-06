package gql

import (
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/satisfactorymodding/smr-api/db/postgres"
	"github.com/satisfactorymodding/smr-api/generated"
)

func (r *mutationResolver) CreateTag(ctx context.Context, tagName string, description string) (*generated.Tag, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "createTag")
	defer wrapper.end()

	dbTag := &postgres.Tag{
		Name:        tagName,
		Description: description,
	}

	resultTag, err := postgres.CreateTag(newCtx, dbTag, true)
	if err != nil {
		return nil, err
	}
	return DBTagToGenerated(resultTag), nil
}

func (r *mutationResolver) CreateMultipleTags(ctx context.Context, tags []*generated.NewTag) ([]*generated.Tag, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "createMultipleTags")
	defer wrapper.end()

	resultTags := make([]postgres.Tag, len(tags))

	for i, tag := range tags {
		dbTag := &postgres.Tag{
			Name:        tag.Name,
			Description: tag.Description,
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

func (r *mutationResolver) UpdateTag(ctx context.Context, id string, newName string, description string) (*generated.Tag, error) {
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
	SetStringINNOE(&description, &dbTag.Description)

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

	insertFilterDefaults(&filter)
	tags := postgres.GetTags(newCtx, filter)

	return DBTagsToGeneratedSlice(tags), nil
}

func insertFilterDefaults(filter **generated.TagFilter) {
	Offset := 0
	Limit := 10
	Order := generated.OrderDesc
	if *filter == nil {
		*filter = &generated.TagFilter{}
	}
	inner := *filter
	if inner.Offset == nil {
		inner.Offset = &Offset
	}
	if inner.Limit == nil {
		inner.Limit = &Limit
	}
	if inner.Order == nil {
		inner.Order = &Order
	}
}
