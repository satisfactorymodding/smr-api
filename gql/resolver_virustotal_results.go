package gql

import (
	"context"
	"log/slog"

	// "github.com/99designs/gqlgen/graphql"

	"github.com/Vilsol/slox"
	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/generated/conv"
	"github.com/satisfactorymodding/smr-api/generated/ent/virustotalresult"
	// "github.com/satisfactorymodding/smr-api/models"
)

func (r *queryResolver) GetVirustotalResult(ctx context.Context, virusTotalHash string) (*generated.VirustotalResult, error) {
	result, err := db.From(ctx).VirustotalResult.
		Query().
		// WithTargets().
		Where(virustotalresult.Hash(virusTotalHash)).
		First(ctx)
	if err != nil {
		return nil, err
	}
	// return nil, nil
	return (*conv.VirustotalResultImpl)(nil).Convert(result), nil
}

func (r *queryResolver) GetVirustotalResults(_ context.Context) (*generated.GetVirustotalResults, error) {
	return &generated.GetVirustotalResults{}, nil
}

type getVirustotalResultsResolver struct{ *Resolver }

func (r *getVirustotalResultsResolver) VirustotalResults(ctx context.Context, _ *generated.GetVirustotalResults) ([]*generated.VirustotalResult, error) {
	// resolverContext := graphql.GetFieldContext(ctx)
	// unapproved := resolverContext.Parent.Field.Field.Name == "getUnapprovedVersions"

	// versionFilter, err := models.ProcessVersionFilter(resolverContext.Parent.Args["filter"].(map[string]interface{}))
	// if err != nil {
	// 	return nil, err
	// }

	// for _, field := range graphql.CollectFieldsCtx(ctx, nil) {
	// 	versionFilter.AddField(field.Name)
	// }

	query := db.From(ctx).VirustotalResult.Query()
	// query = convertVersionFilter(query, versionFilter, unapproved)

	result, err := query.All(ctx)
	if err != nil {
		return nil, err
	}
	slox.Info(ctx, "Getting all virustotal results", slog.String("results", "a"))
	// return nil, nil
	return (*conv.VirustotalResultImpl)(nil).ConvertSlice(result), nil
}

func (r *getVirustotalResultsResolver) Count(ctx context.Context, _ *generated.GetVirustotalResults) (int, error) {
	// resolverContext := graphql.GetFieldContext(ctx)
	// unapproved := resolverContext.Parent.Field.Field.Name == "getUnapprovedVersions"

	// versionFilter, err := models.ProcessVersionFilter(resolverContext.Parent.Args["filter"].(map[string]interface{}))
	// if err != nil {
	// 	return 0, err
	// }

	query := db.From(ctx).VirustotalResult.Query()
	// query = convertVersionFilter(query, versionFilter, unapproved)

	count, err := query.Count(ctx)
	if err != nil {
		return 0, err
	}
	// return 0, nil
	return count, nil
}
