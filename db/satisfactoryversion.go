package db

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"github.com/Masterminds/semver/v3"

	"github.com/satisfactorymodding/smr-api/generated/ent/satisfactoryversion"
)

func GetEngineVersionForSatisfactoryVersion(ctx context.Context, satisfactoryVersionRange string) (string, error) {
	r, err := semver.NewConstraint(satisfactoryVersionRange)
	if err != nil {
		return "", fmt.Errorf("failed to parse version range: %w", err)
	}

	// Each entry's engine version represents that
	// the engine version was in use for satisfactory versions
	// >= (entry version) < (next entry version)
	// So we need to find the first version where this range intersects with the given range,
	// which is equivalent to finding the highest version that is <= the min version of the given range

	minSatisfactoryVersion, err := r.MinVersion()
	if err != nil {
		return "", fmt.Errorf("failed to get min version: %w", err)
	}

	query := From(ctx).SatisfactoryVersion.Query().
		Where(satisfactoryversion.VersionLTE(int(minSatisfactoryVersion.Major()))).
		Order(satisfactoryversion.ByVersion(sql.OrderDesc()))
	v, err := query.First(ctx)
	if err != nil {
		return "", fmt.Errorf("no engine version found for game version range %s: %w", satisfactoryVersionRange, err)
	}

	return v.EngineVersion, nil
}
