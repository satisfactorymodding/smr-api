package db

import (
	"context"
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/satisfactorymodding/smr-api/generated/ent/satisfactoryversion"
)

func GetEngineVersionForSatisfactoryVersion(ctx context.Context, satisfactoryVersionRange string) (string, error) {
	r, err := semver.NewConstraint(satisfactoryVersionRange)
	if err != nil {
		return "", fmt.Errorf("failed to parse version range: %w", err)
	}

	query := From(ctx).SatisfactoryVersion.Query().
		Order(satisfactoryversion.ByVersion())
	versions, err := query.All(ctx)
	if err != nil {
		return "", err
	}

	for _, v := range versions {
		versionSemver := semver.New(uint64(v.Version), 0, 0, "", "")
		if r.Check(versionSemver) {
			return v.EngineVersion, nil
		}
	}

	return "", fmt.Errorf("no engine version found for game version range %s", satisfactoryVersionRange)
}
