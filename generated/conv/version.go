// Code generated by github.com/jmattheis/goverter, DO NOT EDIT.

package conv

import (
	conversion "github.com/satisfactorymodding/smr-api/conversion"
	generated "github.com/satisfactorymodding/smr-api/generated"
	ent "github.com/satisfactorymodding/smr-api/generated/ent"
)

type VersionImpl struct{}

func (c *VersionImpl) Convert(source *ent.Version) *generated.Version {
	var pGeneratedVersion *generated.Version
	if source != nil {
		var generatedVersion generated.Version
		generatedVersion.ID = (*source).ID
		generatedVersion.ModID = (*source).ModID
		generatedVersion.Version = (*source).Version
		generatedVersion.SmlVersion = (*source).SmlVersion
		generatedVersion.Changelog = (*source).Changelog
		generatedVersion.Downloads = conversion.UIntToInt((*source).Downloads)
		generatedVersion.Stability = generated.VersionStabilities((*source).Stability)
		generatedVersion.Approved = (*source).Approved
		generatedVersion.UpdatedAt = conversion.TimeToString((*source).UpdatedAt)
		generatedVersion.CreatedAt = conversion.TimeToString((*source).CreatedAt)
		var pGeneratedVersionTargetList []*generated.VersionTarget
		if (*source).Edges.Targets != nil {
			pGeneratedVersionTargetList = make([]*generated.VersionTarget, len((*source).Edges.Targets))
			for i := 0; i < len((*source).Edges.Targets); i++ {
				pGeneratedVersionTargetList[i] = c.ConvertTarget((*source).Edges.Targets[i])
			}
		}
		generatedVersion.Targets = pGeneratedVersionTargetList
		pString := (*source).Metadata
		generatedVersion.Metadata = &pString
		pGeneratedVersion = &generatedVersion
	}
	return pGeneratedVersion
}
func (c *VersionImpl) ConvertSlice(source []*ent.Version) []*generated.Version {
	var pGeneratedVersionList []*generated.Version
	if source != nil {
		pGeneratedVersionList = make([]*generated.Version, len(source))
		for i := 0; i < len(source); i++ {
			pGeneratedVersionList[i] = c.Convert(source[i])
		}
	}
	return pGeneratedVersionList
}
func (c *VersionImpl) ConvertTarget(source *ent.VersionTarget) *generated.VersionTarget {
	var pGeneratedVersionTarget *generated.VersionTarget
	if source != nil {
		var generatedVersionTarget generated.VersionTarget
		generatedVersionTarget.VersionID = (*source).VersionID
		generatedVersionTarget.TargetName = generated.TargetName((*source).TargetName)
		pInt := conversion.Int64ToInt((*source).Size)
		generatedVersionTarget.Size = &pInt
		pString := (*source).Hash
		generatedVersionTarget.Hash = &pString
		pGeneratedVersionTarget = &generatedVersionTarget
	}
	return pGeneratedVersionTarget
}
