package gql

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"gopkg.in/go-playground/validator.v9"

	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/generated/conv"
	"github.com/satisfactorymodding/smr-api/generated/ent/satisfactoryversion"
	"github.com/satisfactorymodding/smr-api/util"
)

func (r *queryResolver) GetSatisfactoryVersions(ctx context.Context) ([]*generated.SatisfactoryVersion, error) {
	query := db.From(ctx).SatisfactoryVersion.Query().Order(satisfactoryversion.ByVersion(sql.OrderDesc()))
	versions, err := query.All(ctx)
	if err != nil {
		return nil, err
	}
	return (*conv.SatisfactoryVersionImpl)(nil).ConvertSlice(versions), nil
}

func (r *queryResolver) GetSatisfactoryVersion(ctx context.Context, id string) (*generated.SatisfactoryVersion, error) {
	query := db.From(ctx).SatisfactoryVersion.Query().Where(satisfactoryversion.ID(id))
	version, err := query.First(ctx)
	if err != nil {
		return nil, err
	}
	return (*conv.SatisfactoryVersionImpl)(nil).Convert(version), nil
}

func (r *mutationResolver) CreateSatisfactoryVersion(ctx context.Context, satisfactoryVersion generated.NewSatisfactoryVersion) (*generated.SatisfactoryVersion, error) {
	val := ctx.Value(util.ContextValidator{}).(*validator.Validate)
	if err := val.Struct(&satisfactoryVersion); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	result, err := db.From(ctx).SatisfactoryVersion.
		Create().
		SetVersion(satisfactoryVersion.Version).
		SetEngineVersion(satisfactoryVersion.EngineVersion).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return (*conv.SatisfactoryVersionImpl)(nil).Convert(result), nil
}

func (r *mutationResolver) UpdateSatisfactoryVersion(ctx context.Context, satisfactoryVersionID string, satisfactoryVersion generated.UpdateSatisfactoryVersion) (*generated.SatisfactoryVersion, error) {
	val := ctx.Value(util.ContextValidator{}).(*validator.Validate)
	if err := val.Struct(&satisfactoryVersion); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	update := db.From(ctx).SatisfactoryVersion.UpdateOneID(satisfactoryVersionID)
	SetINNOEF(satisfactoryVersion.Version, update.SetVersion)
	SetINNOEF(satisfactoryVersion.EngineVersion, update.SetEngineVersion)

	result, err := update.Save(ctx)
	if err != nil {
		return nil, err
	}

	return (*conv.SatisfactoryVersionImpl)(nil).Convert(result), nil
}

func (r *mutationResolver) DeleteSatisfactoryVersion(ctx context.Context, satisfactoryVersionID string) (bool, error) {
	if err := db.From(ctx).SatisfactoryVersion.DeleteOneID(satisfactoryVersionID).Exec(ctx); err != nil {
		return false, err
	}

	return true, nil
}
