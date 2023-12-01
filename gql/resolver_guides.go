package gql

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/99designs/gqlgen/graphql"
	"github.com/Vilsol/slox"
	"gopkg.in/go-playground/validator.v9"

	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/generated/conv"
	"github.com/satisfactorymodding/smr-api/generated/ent"
	"github.com/satisfactorymodding/smr-api/generated/ent/guide"
	"github.com/satisfactorymodding/smr-api/generated/ent/guidetag"
	"github.com/satisfactorymodding/smr-api/generated/ent/tag"
	"github.com/satisfactorymodding/smr-api/models"
	"github.com/satisfactorymodding/smr-api/redis"
	"github.com/satisfactorymodding/smr-api/util"
)

func (r *mutationResolver) CreateGuide(ctx context.Context, g generated.NewGuide) (*generated.Guide, error) {
	wrapper, ctx := WrapMutationTrace(ctx, "createGuide")
	defer wrapper.end()

	val := ctx.Value(util.ContextValidator{}).(*validator.Validate)
	if err := val.Struct(&g); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	user, err := db.UserFromGQLContext(ctx)
	if err != nil {
		return nil, err
	}

	// Allow only 8 new guides per 24h
	guides, err := db.From(ctx).Guide.Query().Where(
		guide.UserID(user.ID),
		guide.CreatedAtGT(time.Now().Add(time.Hour*24*-1)),
	).All(ctx)
	if err != nil {
		return nil, err
	}

	currentAvailable := float64(8)
	lastGuideTime := time.Now()
	for _, existingGuide := range guides {
		currentAvailable--
		if existingGuide.CreatedAt.After(lastGuideTime) {
			diff := existingGuide.CreatedAt.Sub(lastGuideTime)
			currentAvailable = math.Min(8, currentAvailable+diff.Hours()/3)
		}
		lastGuideTime = existingGuide.CreatedAt
	}

	if currentAvailable < 1 {
		timeToWait := time.Until(lastGuideTime.Add(time.Hour * 6)).Minutes()
		return nil, fmt.Errorf("please wait %.0f minutes to post another guide", timeToWait)
	}

	result, err := db.From(ctx).Guide.
		Create().
		SetName(g.Name).
		SetShortDescription(g.ShortDescription).
		SetGuide(g.Guide).
		SetUser(user).
		AddTagIDs(g.TagIDs...).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return (*conv.GuideImpl)(nil).Convert(result), nil
}

func (r *mutationResolver) UpdateGuide(ctx context.Context, guideID string, g generated.UpdateGuide) (*generated.Guide, error) {
	wrapper, ctx := WrapMutationTrace(ctx, "updateGuide")
	defer wrapper.end()

	val := ctx.Value(util.ContextValidator{}).(*validator.Validate)
	if err := val.Struct(&g); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	if err := db.Tx(ctx, func(ctx context.Context, tx *ent.Tx) error {
		update := tx.Guide.UpdateOneID(guideID)

		SetINNOEF(g.Name, update.SetName)
		SetINNOEF(g.ShortDescription, update.SetShortDescription)
		SetINNOEF(g.Guide, update.SetGuide)

		if err := update.Exec(ctx); err != nil {
			return err
		}

		if _, err := tx.GuideTag.Delete().Where(
			guidetag.GuideID(guideID),
			guidetag.TagIDNotIn(g.TagIDs...),
		).Exec(ctx); err != nil {
			return err
		}

		return tx.GuideTag.MapCreateBulk(g.TagIDs, func(create *ent.GuideTagCreate, i int) {
			create.SetGuideID(guideID).SetTagID(g.TagIDs[i])
		}).Exec(ctx)
	}, nil); err != nil {
		return nil, err
	}

	result, err := db.From(ctx).Guide.Query().WithTags().Where(guide.ID(guideID)).First(ctx)
	if err != nil {
		return nil, err
	}

	return (*conv.GuideImpl)(nil).Convert(result), nil
}

func (r *mutationResolver) DeleteGuide(ctx context.Context, guideID string) (bool, error) {
	wrapper, ctx := WrapMutationTrace(ctx, "deleteGuide")
	defer wrapper.end()

	if err := db.From(ctx).Guide.DeleteOneID(guideID).Exec(ctx); err != nil {
		return false, err
	}

	return true, nil
}

func (r *queryResolver) GetGuide(ctx context.Context, guideID string) (*generated.Guide, error) {
	wrapper, ctx := WrapQueryTrace(ctx, "getGuide")
	defer wrapper.end()

	result, err := db.From(ctx).Guide.Get(ctx, guideID)
	if err != nil {
		return nil, err
	}

	if result != nil {
		if redis.CanIncrement(RealIP(ctx), "view", "guide:"+guideID, time.Hour*4) {
			err := result.Update().AddViews(1).Exec(ctx)
			if err != nil {
				slox.Error(ctx, "failed incrementing views", slog.Any("err", err))
			}
		}
	}

	return (*conv.GuideImpl)(nil).Convert(result), nil
}

func (r *queryResolver) GetGuides(ctx context.Context, _ map[string]interface{}) (*generated.GetGuides, error) {
	wrapper, _ := WrapQueryTrace(ctx, "getGuides")
	defer wrapper.end()
	return &generated.GetGuides{}, nil
}

type getGuidesResolver struct{ *Resolver }

func (r *getGuidesResolver) Guides(ctx context.Context, _ *generated.GetGuides) ([]*generated.Guide, error) {
	wrapper, ctx := WrapQueryTrace(ctx, "GetGuides.guides")
	defer wrapper.end()

	resolverContext := graphql.GetFieldContext(ctx)
	guideFilter, err := models.ProcessGuideFilter(resolverContext.Parent.Args["filter"].(map[string]interface{}))
	if err != nil {
		return nil, err
	}

	query := db.From(ctx).Guide.Query().WithTags()
	query = convertGuideFilter(query, guideFilter)

	result, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	return (*conv.GuideImpl)(nil).ConvertSlice(result), nil
}

func (r *getGuidesResolver) Count(ctx context.Context, _ *generated.GetGuides) (int, error) {
	wrapper, ctx := WrapQueryTrace(ctx, "GetGuides.count")
	defer wrapper.end()

	resolverContext := graphql.GetFieldContext(ctx)
	guideFilter, err := models.ProcessGuideFilter(resolverContext.Parent.Args["filter"].(map[string]interface{}))
	if err != nil {
		return 0, err
	}

	query := db.From(ctx).Guide.Query().WithTags()
	query = convertGuideFilter(query, guideFilter)

	result, err := query.Count(ctx)
	if err != nil {
		return 0, err
	}

	return result, nil
}

type guideResolver struct{ *Resolver }

func (r *guideResolver) User(ctx context.Context, obj *generated.Guide) (*generated.User, error) {
	wrapper, ctx := WrapQueryTrace(ctx, "Guide.user")
	defer wrapper.end()

	result, err := db.From(ctx).User.Get(ctx, obj.UserID)
	if err != nil {
		return nil, err
	}

	return (*conv.UserImpl)(nil).Convert(result), nil
}

func convertGuideFilter(query *ent.GuideQuery, filter *models.GuideFilter) *ent.GuideQuery {
	if len(filter.Ids) > 0 {
		query = query.Where(guide.IDIn(filter.Ids...))
	} else if filter != nil {
		query = query.
			Limit(*filter.Limit).
			Offset(*filter.Offset).
			Order(sql.OrderByField(
				filter.OrderBy.String(),
				db.OrderToOrder(filter.Order.String()),
			).ToFunc())

		if filter.Search != nil && *filter.Search != "" {
			query = query.Modify(func(s *sql.Selector) {
				s.Where(sql.ExprP("to_tsvector(name) @@ to_tsquery(?)", strings.ReplaceAll(*filter.Search, " ", " & ")))
			}).Clone()
		}

		if filter.TagIDs != nil && len(filter.TagIDs) > 0 {
			query = query.Where(guide.HasTagsWith(tag.IDIn(filter.TagIDs...)))
		}
	}
	return query
}
