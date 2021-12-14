package gql

import (
	"context"
	"time"

	"github.com/satisfactorymodding/smr-api/db/postgres"
	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/models"
	"github.com/satisfactorymodding/smr-api/redis"
	"github.com/satisfactorymodding/smr-api/util"

	"github.com/pkg/errors"

	"github.com/99designs/gqlgen/graphql"
	"gopkg.in/go-playground/validator.v9"
)

func (r *mutationResolver) CreateGuide(ctx context.Context, guide generated.NewGuide) (*generated.Guide, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "createGuide")
	defer wrapper.end()

	val := ctx.Value(util.ContextValidator{}).(*validator.Validate)
	if err := val.Struct(&guide); err != nil {
		return nil, errors.Wrap(err, "validation failed")
	}

	dbGuide := &postgres.Guide{
		Name:             guide.Name,
		ShortDescription: guide.ShortDescription,
		Guide:            guide.Guide,
	}

	user := ctx.Value(postgres.UserKey{}).(*postgres.User)

	dbGuide.UserID = user.ID

	resultGuide, err := postgres.CreateGuide(newCtx, dbGuide)

	if err != nil {
		return nil, err
	}

	err = postgres.SetGuideTags(resultGuide.ID, guide.TagIDs, &newCtx)

	if err != nil {
		return nil, err
	}

	// Need to get the guide again to populate tags
	return DBGuideToGenerated(postgres.GetGuideByIDNoCache(resultGuide.ID, &newCtx)), nil
}

func (r *mutationResolver) UpdateGuide(ctx context.Context, guideID string, guide generated.UpdateGuide) (*generated.Guide, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "updateGuide")
	defer wrapper.end()

	val := ctx.Value(util.ContextValidator{}).(*validator.Validate)
	if err := val.Struct(&guide); err != nil {
		return nil, errors.Wrap(err, "validation failed")
	}

	err := postgres.ResetGuideTags(guideID, guide.TagIDs, &newCtx)
	if err != nil {
		return nil, err
	}

	dbGuide := postgres.GetGuideByIDNoCache(guideID, &newCtx)

	if dbGuide == nil {
		return nil, errors.New("guide not found")
	}

	SetStringINNOE(guide.Name, &dbGuide.Name)
	SetStringINNOE(guide.ShortDescription, &dbGuide.ShortDescription)
	SetStringINNOE(guide.Guide, &dbGuide.Guide)

	postgres.Save(newCtx, &dbGuide)

	return DBGuideToGenerated(dbGuide), nil
}

func (r *mutationResolver) DeleteGuide(ctx context.Context, guideID string) (bool, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "deleteGuide")
	defer wrapper.end()

	dbGuide := postgres.GetGuideByID(newCtx, guideID)

	if dbGuide == nil {
		return false, errors.New("guide not found")
	}

	postgres.Delete(newCtx, &dbGuide)

	return true, nil
}

func (r *queryResolver) GetGuide(ctx context.Context, guideID string) (*generated.Guide, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "getGuide")
	defer wrapper.end()

	guide := postgres.GetGuideByID(newCtx, guideID)

	if guide != nil {
		if redis.CanIncrement(RealIP(ctx), "view", "guide:"+guideID, time.Hour*4) {
			postgres.IncrementGuideViews(newCtx, guide)
		}
	}

	return DBGuideToGenerated(guide), nil
}

func (r *queryResolver) GetGuides(ctx context.Context, filter map[string]interface{}) (*generated.GetGuides, error) {
	wrapper, _ := WrapQueryTrace(ctx, "getGuides")
	defer wrapper.end()
	return &generated.GetGuides{}, nil
}

type getGuidesResolver struct{ *Resolver }

func (r *getGuidesResolver) Guides(ctx context.Context, obj *generated.GetGuides) ([]*generated.Guide, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "GetGuides.guides")
	defer wrapper.end()

	resolverContext := graphql.GetFieldContext(ctx)
	guideFilter, err := models.ProcessGuideFilter(resolverContext.Parent.Args["filter"].(map[string]interface{}))

	if err != nil {
		return nil, err
	}

	var guides []postgres.Guide

	if guideFilter.Ids == nil || len(guideFilter.Ids) == 0 {
		guides = postgres.GetGuides(newCtx, guideFilter)
	} else {
		guides = postgres.GetGuidesByID(newCtx, guideFilter.Ids)
	}

	if guides == nil {
		return nil, errors.New("guides not found")
	}

	converted := make([]*generated.Guide, len(guides))
	for k, v := range guides {
		converted[k] = DBGuideToGenerated(&v)
	}

	return converted, nil
}

func (r *getGuidesResolver) Count(ctx context.Context, obj *generated.GetGuides) (int, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "GetGuides.count")
	defer wrapper.end()

	resolverContext := graphql.GetFieldContext(ctx)
	guideFilter, err := models.ProcessGuideFilter(resolverContext.Parent.Args["filter"].(map[string]interface{}))

	if err != nil {
		return 0, err
	}

	if guideFilter.Ids != nil && len(guideFilter.Ids) != 0 {
		return len(guideFilter.Ids), nil
	}

	return int(postgres.GetGuideCount(newCtx, guideFilter)), nil
}

type guideResolver struct{ *Resolver }

func (r *guideResolver) User(ctx context.Context, obj *generated.Guide) (*generated.User, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "Guide.user")
	defer wrapper.end()

	user := postgres.GetUserByID(newCtx, obj.UserID)

	if user == nil {
		return nil, errors.New("user not found")
	}

	return DBUserToGenerated(user), nil
}
