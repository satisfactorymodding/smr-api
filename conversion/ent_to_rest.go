package conversion

import (
	"github.com/satisfactorymodding/smr-api/generated/ent"
	"github.com/satisfactorymodding/smr-api/nodes/types"
)

// goverter:converter
// goverter:output:file ../generated/conv/version.go
// goverter:output:package conv
// goverter:extend TimeToString UIntToInt Int64ToInt
type ModAllVersions interface {
	// goverter:map Edges.Targets Targets
	// goverter:map Edges.VersionDependencies Dependencies
	Convert(source *ent.Version) *types.ModAllVersionsVersion
	ConvertSlice(source []*ent.Version) []*types.ModAllVersionsVersion
}
