package gql

import (
	"context"
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/99designs/gqlgen/graphql"
	"gopkg.in/go-playground/validator.v9"

	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/generated/conv"
	"github.com/satisfactorymodding/smr-api/generated/ent"
	"github.com/satisfactorymodding/smr-api/generated/ent/smlversion"
	"github.com/satisfactorymodding/smr-api/generated/ent/smlversiontarget"
	"github.com/satisfactorymodding/smr-api/models"
	"github.com/satisfactorymodding/smr-api/util"
)

func (r *mutationResolver) CreateSMLVersion(ctx context.Context, smlVersion generated.NewSMLVersion) (*generated.SMLVersion, error) {
	wrapper, ctx := WrapMutationTrace(ctx, "createSMLVersion")
	defer wrapper.end()

	val := ctx.Value(util.ContextValidator{}).(*validator.Validate)
	if err := val.Struct(&smlVersion); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	date, err := time.Parse(time.RFC3339Nano, smlVersion.Date)
	if err != nil {
		return nil, fmt.Errorf("failed to parse date: %w", err)
	}

	var result *ent.SmlVersion
	if err := db.Tx(ctx, func(ctx context.Context, tx *ent.Tx) error {
		result, err = db.From(ctx).SmlVersion.
			Create().
			SetVersion(smlVersion.Version).
			SetSatisfactoryVersion(smlVersion.SatisfactoryVersion).
			SetNillableBootstrapVersion(smlVersion.BootstrapVersion).
			SetStability(smlversion.Stability(smlVersion.Stability)).
			SetLink(smlVersion.Link).
			SetChangelog(smlVersion.Changelog).
			SetDate(date).
			SetEngineVersion(smlVersion.EngineVersion).
			Save(ctx)
		if err != nil {
			return err
		}

		println(result.ID)

		for _, smlVersionTarget := range smlVersion.Targets {
			if _, err := db.From(ctx).SmlVersionTarget.
				Create().
				SetVersionID(result.ID).
				SetTargetName(string(smlVersionTarget.TargetName)).
				SetLink(smlVersionTarget.Link).
				Save(ctx); err != nil {
				println("HERE")
				return err
			}
		}

		return nil
	}, nil); err != nil {
		return nil, err
	}

	result, err = db.From(ctx).SmlVersion.
		Query().
		WithTargets().
		Where(smlversion.ID(result.ID)).
		First(ctx)
	if err != nil {
		return nil, err
	}

	return (*conv.SMLVersionImpl)(nil).Convert(result), nil
}

func (r *mutationResolver) UpdateSMLVersion(ctx context.Context, smlVersionID string, smlVersion generated.UpdateSMLVersion) (*generated.SMLVersion, error) {
	wrapper, ctx := WrapMutationTrace(ctx, "updateSMLVersion")
	defer wrapper.end()

	val := ctx.Value(util.ContextValidator{}).(*validator.Validate)
	if err := val.Struct(&smlVersion); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	err := db.Tx(ctx, func(ctx context.Context, tx *ent.Tx) error {
		dbSMLVersion, err := tx.SmlVersion.Get(ctx, smlVersionID)
		if err != nil {
			return err
		}

		update := dbSMLVersion.Update()

		SetINNOEF(smlVersion.Version, update.SetVersion)
		SetINNF(smlVersion.SatisfactoryVersion, update.SetSatisfactoryVersion)
		SetINNOEF(smlVersion.BootstrapVersion, update.SetBootstrapVersion)
		SetStabilityINNF(smlVersion.Stability, update.SetStability)
		SetINNOEF(smlVersion.Link, update.SetLink)
		SetINNOEF(smlVersion.Changelog, update.SetChangelog)
		SetDateINNF(smlVersion.Date, update.SetDate)
		SetINNOEF(smlVersion.EngineVersion, update.SetEngineVersion)

		if err := update.Exec(ctx); err != nil {
			return err
		}

		dbSMLTargets, err := dbSMLVersion.QueryTargets().All(ctx)
		if err != nil {
			return err
		}

		for _, dbSMLTarget := range dbSMLTargets {
			found := false

			for _, smlTarget := range smlVersion.Targets {
				if dbSMLTarget.TargetName == string(smlTarget.TargetName) {
					found = true
				}
			}

			if !found {
				tx.SmlVersionTarget.DeleteOneID(dbSMLTarget.ID)
			}
		}

		for _, smlTarget := range smlVersion.Targets {
			if err := tx.SmlVersionTarget.Update().Where(
				smlversiontarget.VersionID(dbSMLVersion.ID),
				smlversiontarget.TargetName(string(smlTarget.TargetName)),
			).SetLink(smlTarget.Link).Exec(ctx); err != nil {
				return err
			}
		}

		return nil
	}, nil)
	if err != nil {
		return nil, err
	}

	result, err := db.From(ctx).SmlVersion.
		Query().
		WithTargets().
		Where(smlversion.ID(smlVersionID)).
		First(ctx)
	if err != nil {
		return nil, err
	}

	return (*conv.SMLVersionImpl)(nil).Convert(result), nil
}

func (r *mutationResolver) DeleteSMLVersion(ctx context.Context, smlVersionID string) (bool, error) {
	wrapper, ctx := WrapMutationTrace(ctx, "deleteSMLVersion")
	defer wrapper.end()

	err := db.Tx(ctx, func(ctx context.Context, tx *ent.Tx) error {
		if _, err := tx.SmlVersionTarget.Delete().Where(smlversiontarget.VersionID(smlVersionID)).Exec(ctx); err != nil {
			return err
		}

		return tx.SmlVersion.DeleteOneID(smlVersionID).Exec(ctx)
	}, nil)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *queryResolver) GetSMLVersion(ctx context.Context, smlVersionID string) (*generated.SMLVersion, error) {
	wrapper, ctx := WrapQueryTrace(ctx, "getSMLVersion")
	defer wrapper.end()

	result, err := db.From(ctx).SmlVersion.
		Query().
		WithTargets().
		Where(smlversion.ID(smlVersionID)).
		First(ctx)
	if err != nil {
		return nil, err
	}

	return (*conv.SMLVersionImpl)(nil).Convert(result), nil
}

func (r *queryResolver) GetSMLVersions(ctx context.Context, _ map[string]interface{}) (*generated.GetSMLVersions, error) {
	wrapper, _ := WrapQueryTrace(ctx, "getSMLVersions")
	defer wrapper.end()
	return &generated.GetSMLVersions{}, nil
}

type getSMLVersionsResolver struct{ *Resolver }

func (r *getSMLVersionsResolver) SmlVersions(ctx context.Context, _ *generated.GetSMLVersions) ([]*generated.SMLVersion, error) {
	wrapper, ctx := WrapQueryTrace(ctx, "GetSMLVersions.smlVersions")
	defer wrapper.end()

	resolverContext := graphql.GetFieldContext(ctx)
	smlVersionFilter, err := models.ProcessSMLVersionFilter(resolverContext.Parent.Args["filter"].(map[string]interface{}))
	if err != nil {
		return nil, err
	}

	query := db.From(ctx).SmlVersion.Query().WithTargets()
	query = convertSMLVersionFilter(query, smlVersionFilter)

	result, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	return (*conv.SMLVersionImpl)(nil).ConvertSlice(result), nil
}

func (r *getSMLVersionsResolver) Count(ctx context.Context, _ *generated.GetSMLVersions) (int, error) {
	wrapper, ctx := WrapQueryTrace(ctx, "GetSMLVersions.count")
	defer wrapper.end()

	resolverContext := graphql.GetFieldContext(ctx)
	smlVersionFilter, err := models.ProcessSMLVersionFilter(resolverContext.Parent.Args["filter"].(map[string]interface{}))
	if err != nil {
		return 0, err
	}

	query := db.From(ctx).SmlVersion.Query().WithTargets()
	query = convertSMLVersionFilter(query, smlVersionFilter)

	result, err := query.Count(ctx)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func convertSMLVersionFilter(query *ent.SmlVersionQuery, filter *models.SMLVersionFilter) *ent.SmlVersionQuery {
	if len(filter.Ids) > 0 {
		query = query.Where(smlversion.IDIn(filter.Ids...))
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
	}
	return query
}
