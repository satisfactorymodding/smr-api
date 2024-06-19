package db

import (
	"context"

	"entgo.io/ent/dialect/sql"

	"github.com/satisfactorymodding/smr-api/generated/ent/satisfactoryversion"
)

func GetEngineVersionForSatisfactoryVersion(ctx context.Context, satisfactoryVersion int) (string, error) {
	query := From(ctx).SatisfactoryVersion.Query().
		Where(satisfactoryversion.VersionLTE(satisfactoryVersion)).
		Order(satisfactoryversion.ByVersion(sql.OrderDesc()))
	v, err := query.First(ctx)
	if err != nil {
		return "", err
	}
	return v.EngineVersion, nil
}
