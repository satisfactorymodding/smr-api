package gql

import (
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/finnbear/moderation"
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/generated/conv"
	"github.com/satisfactorymodding/smr-api/generated/ent"
	"github.com/satisfactorymodding/smr-api/generated/ent/tag"
)

func (r *mutationResolver) CreateTag(ctx context.Context, tagName string, description string) (*generated.Tag, error) {
	if err := ValidateTagName(tagName); err != nil {
		return nil, err
	}

	result, err := db.From(ctx).Tag.Create().SetName(tagName).SetDescription(description).Save(ctx)
	if err != nil {
		return nil, err
	}

	return (*conv.TagImpl)(nil).Convert(result), nil
}

func (r *mutationResolver) CreateMultipleTags(ctx context.Context, tags []*generated.NewTag) ([]*generated.Tag, error) {
	for _, t := range tags {
		if err := ValidateTagName(t.Name); err != nil {
			return nil, err
		}
	}

	var result []*ent.Tag
	if err := db.Tx(ctx, func(ctx context.Context, tx *ent.Tx) error {
		for _, newTag := range tags {
			t, err := tx.Tag.Create().SetName(newTag.Name).SetDescription(newTag.Description).Save(ctx)
			if err != nil {
				return err
			}
			result = append(result, t)
		}
		return nil
	}, nil); err != nil {
		return nil, err
	}

	return (*conv.TagImpl)(nil).ConvertSlice(result), nil
}

func (r *mutationResolver) DeleteTag(ctx context.Context, id string) (bool, error) {
	if err := db.From(ctx).Tag.DeleteOneID(id).Exec(ctx); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) UpdateTag(ctx context.Context, id string, newName string, description string) (*generated.Tag, error) {
	if err := ValidateTagName(newName); err != nil {
		return nil, err
	}

	update := db.From(ctx).Tag.UpdateOneID(id)

	result, err := update.SetName(newName).SetDescription(description).Save(ctx)
	if err != nil {
		return nil, err
	}

	return (*conv.TagImpl)(nil).Convert(result), nil
}

func (r *queryResolver) GetTag(ctx context.Context, id string) (*generated.Tag, error) {
	result, err := db.From(ctx).Tag.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return (*conv.TagImpl)(nil).Convert(result), nil
}

func (r *queryResolver) GetTags(ctx context.Context, filter *generated.TagFilter) ([]*generated.Tag, error) {
	insertFilterDefaults(&filter)

	query := db.From(ctx).Tag.Query()

	if filter != nil {
		query = query.
			Limit(*filter.Limit).
			Offset(*filter.Offset).
			Order(sql.OrderByField(
				tag.FieldName,
				db.OrderToOrder(filter.Order.String()),
			).ToFunc())

		if filter.Search != nil && *filter.Search != "" {
			cleanSearch := strings.ReplaceAll(strings.TrimSpace(*filter.Search), " ", " & ")

			query = query.Modify(func(s *sql.Selector) {
				s.Where(sql.P(func(builder *sql.Builder) {
					builder.WriteString("similarity(name, ").Arg(cleanSearch).WriteString(") > 0.2")
				}))
			}).TagQuery
		}

		if filter.Ids != nil && len(filter.Ids) > 0 {
			query.Where(tag.IDIn(filter.Ids...))
		}
	}

	result, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	return (*conv.TagImpl)(nil).ConvertSlice(result), nil
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

func ValidateTagName(tag string) error {
	if len(tag) > 24 {
		return errors.New("Tag name is over 24 characters long")
	}
	if len(tag) < 3 {
		return errors.New("Tag name is under 3 characters long")
	}
	if moderation.IsInappropriate(tag) {
		return errors.New("Tag name is inapproriate")
	}
	return nil
}
