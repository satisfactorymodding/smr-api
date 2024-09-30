package gql

import (
	"context"
	"log/slog"

	"github.com/Vilsol/slox"
	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/generated/conv"
	"github.com/satisfactorymodding/smr-api/generated/ent/virustotalresult"
)

func (r *queryResolver) GetVirustotalResult(ctx context.Context, virusTotalHash string) (*generated.VirustotalResult, error) {
	result, err := db.From(ctx).VirustotalResult.
		Query().
		Where(virustotalresult.Hash(virusTotalHash)).
		First(ctx)
	if err != nil {
		return nil, err
	}

	return (*conv.VirustotalResultImpl)(nil).Convert(result), nil
}

func (r *queryResolver) GetVirustotalResults(_ context.Context) (*generated.GetVirustotalResults, error) {
	return &generated.GetVirustotalResults{}, nil
}

type getVirustotalResultsResolver struct{ *Resolver }

func (r *getVirustotalResultsResolver) VirustotalResults(ctx context.Context, _ *generated.GetVirustotalResults) ([]*generated.VirustotalResult, error) {

	query := db.From(ctx).VirustotalResult.Query()

	result, err := query.All(ctx)
	if err != nil {
		return nil, err
	}
	slox.Info(ctx, "Getting all virustotal results", slog.String("results", "a"))
	return (*conv.VirustotalResultImpl)(nil).ConvertSlice(result), nil
}

func (r *getVirustotalResultsResolver) Count(ctx context.Context, _ *generated.GetVirustotalResults) (int, error) {
	query := db.From(ctx).VirustotalResult.Query()

	count, err := query.Count(ctx)
	if err != nil {
		return 0, err
	}
	return count, nil
}
