package db

import (
	"context"
	"regexp"
	"strconv"

	"entgo.io/ent/dialect/sql"

	"github.com/satisfactorymodding/smr-api/generated/ent"
	"github.com/satisfactorymodding/smr-api/generated/ent/version"
)

var semverCheck = regexp.MustCompile(`^(<=|<|>|>=|\^)?(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)

func OrderToOrder(order string) sql.OrderTermOption {
	if order == "asc" || order == "ascending" {
		return sql.OrderAsc()
	}

	return sql.OrderDesc()
}

func GetModVersionsConstraint(ctx context.Context, modID string, constraint string) ([]*ent.Version, error) {
	matches := semverCheck.FindAllStringSubmatch(constraint, -1)
	if len(matches) == 0 {
		return nil, nil
	}

	major, err := strconv.Atoi(matches[0][2])
	if err != nil {
		return nil, nil
	}

	minor, err := strconv.Atoi(matches[0][3])
	if err != nil {
		return nil, nil
	}

	patch, err := strconv.Atoi(matches[0][4])
	if err != nil {
		return nil, nil
	}

	q := From(ctx).Version.Query().WithTargets().Where(version.ModID(modID))

	/*
		<=1.2.3
		major < 1
		major = 1, minor < 2
		major = 1, minor = 2, patch <= 3

		<1.2.3
		major < 1
		major = 1, minor < 2
		major = 1, minor = 2, patch < 3

		>1.2.3
		major > 1
		major = 1, minor > 2
		major = 1, minor = 2, patch > 3

		>=1.2.3
		major > 1
		major = 1, minor > 2
		major = 1, minor = 2, patch >= 3

		1.2.3
		major = 1, minor = 2, patch = 3

		^1.2.3 (>=1.2.3, <2.0.0)
		major = 1, minor > 2
		major = 1, minor = 2, patch >= 3

		^0.2.3 (>=0.2.3, <0.3.0)
		major = 0, minor = 2, patch >= 3

		^0.0.3 (>=0.0.3, <0.0.4)
		major = 0, minor = 0, patch = 3
	*/

	sign := matches[0][1]
	switch sign {
	case "<=":
		q = q.Where(
			version.Or(
				version.VersionMajorLT(major),
				version.And(version.VersionMajorEQ(major), version.VersionMinorLT(minor)),
				version.And(version.VersionMajorEQ(major), version.VersionMinorEQ(minor), version.VersionPatchLTE(patch)),
			),
		)
	case "<":
		q = q.Where(
			version.Or(
				version.VersionMajorLT(major),
				version.And(version.VersionMajorEQ(major), version.VersionMinorLT(minor)),
				version.And(version.VersionMajorEQ(major), version.VersionMinorEQ(minor), version.VersionPatchLT(patch)),
			),
		)
	case ">":
		q = q.Where(
			version.Or(
				version.VersionMajorGT(major),
				version.And(version.VersionMajorEQ(major), version.VersionMinorGT(minor)),
				version.And(version.VersionMajorEQ(major), version.VersionMinorEQ(minor), version.VersionPatchGT(patch)),
			),
		)
	case ">=":
		q = q.Where(
			version.Or(
				version.VersionMajorGT(major),
				version.And(version.VersionMajorEQ(major), version.VersionMinorGT(minor)),
				version.And(version.VersionMajorEQ(major), version.VersionMinorEQ(minor), version.VersionPatchGTE(patch)),
			),
		)
	case "^":
		if major != 0 {
			q = q.Where(
				version.Or(
					version.And(version.VersionMajorEQ(major), version.VersionMinorGT(minor)),
					version.And(version.VersionMajorEQ(major), version.VersionMinorEQ(minor), version.VersionPatchGTE(patch)),
				),
			)
		} else if minor != 0 {
			q = q.Where(
				version.And(version.VersionMajorEQ(major), version.VersionMinorEQ(minor), version.VersionPatchGTE(patch)),
			)
		} else {
			q = q.Where(
				version.And(version.VersionMajorEQ(major), version.VersionMinorEQ(minor), version.VersionPatchEQ(patch)),
			)
		}
	default:
		q = q.Where(
			version.And(version.VersionMajorEQ(major), version.VersionMinorEQ(minor), version.VersionPatchEQ(patch)),
		)
	}

	return q.All(ctx)
}
