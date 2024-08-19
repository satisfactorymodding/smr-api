package gql

import (
	"context"

	"github.com/satisfactorymodding/smr-api/generated"
)

type versionDependencyResolver struct{ *Resolver }

func (r *versionDependencyResolver) Mod(ctx context.Context, obj *generated.VersionDependency) (*generated.Mod, error) {
	return r.Query().GetModByReference(ctx, obj.ModID)
}
