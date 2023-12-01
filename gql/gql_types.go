package gql

import (
	"time"

	"github.com/satisfactorymodding/smr-api/db/postgres"
	"github.com/satisfactorymodding/smr-api/generated"
)

func DBModToGenerated(mod *postgres.Mod) *generated.Mod {
	if mod == nil {
		return nil
	}

	Logo := mod.Logo
	SourceURL := mod.SourceURL
	FullDescription := mod.FullDescription

	var LastVersionDate string
	if mod.LastVersionDate != nil {
		LastVersionDate = mod.LastVersionDate.Format(time.RFC3339Nano)
	}

	return &generated.Mod{
		ID:               mod.ID,
		Name:             mod.Name,
		ShortDescription: mod.ShortDescription,
		Logo:             &Logo,
		SourceURL:        &SourceURL,
		CreatorID:        mod.CreatorID,
		Approved:         mod.Approved,
		Views:            int(mod.Views),
		Downloads:        int(mod.Downloads),
		Hotness:          int(mod.Hotness),
		Popularity:       int(mod.Popularity),
		UpdatedAt:        mod.UpdatedAt.Format(time.RFC3339Nano),
		CreatedAt:        mod.CreatedAt.Format(time.RFC3339Nano),
		FullDescription:  &FullDescription,
		LastVersionDate:  &LastVersionDate,
		ModReference:     mod.ModReference,
		Hidden:           mod.Hidden,
		Versions:         DBVersionsToGeneratedSlice(mod.Versions),
		Tags:             DBTagsToGeneratedSlice(mod.Tags),
		Compatibility:    DBCompInfoToGenCompInfo(mod.Compatibility),
	}
}

func DBVersionToGenerated(version *postgres.Version) *generated.Version {
	if version == nil {
		return nil
	}

	size := 0

	if version.Size != nil {
		size = int(*version.Size)
	}

	return &generated.Version{
		ID:         version.ID,
		Version:    version.Version,
		SmlVersion: version.SMLVersion,
		Changelog:  version.Changelog,
		Downloads:  int(version.Downloads),
		Stability:  generated.VersionStabilities(version.Stability),
		Targets:    DBVersionTargetsToGeneratedSlice(version.Targets),
		Approved:   version.Approved,
		UpdatedAt:  version.UpdatedAt.Format(time.RFC3339Nano),
		CreatedAt:  version.CreatedAt.Format(time.RFC3339Nano),
		ModID:      version.ModID,
		Metadata:   version.Metadata,
		Hash:       version.Hash,
		Size:       &size,
	}
}

func DBVersionsToGeneratedSlice(versions []postgres.Version) []*generated.Version {
	converted := make([]*generated.Version, len(versions))
	for i, version := range versions {
		converted[i] = DBVersionToGenerated(&version)
	}
	return converted
}

func DBVersionDependencyToGenerated(versionDependency *postgres.VersionDependency) *generated.VersionDependency {
	if versionDependency == nil {
		return nil
	}

	return &generated.VersionDependency{
		VersionID: versionDependency.VersionID,
		ModID:     versionDependency.ModID,
		Condition: versionDependency.Condition,
		Optional:  versionDependency.Optional,
	}
}

func DBTagToGenerated(tag *postgres.Tag) *generated.Tag {
	if tag == nil {
		return nil
	}
	return &generated.Tag{
		Name:        tag.Name,
		ID:          tag.ID,
		Description: tag.Description,
	}
}

func DBTagsToGeneratedSlice(tags []postgres.Tag) []*generated.Tag {
	converted := make([]*generated.Tag, len(tags))
	for i, tag := range tags {
		converted[i] = DBTagToGenerated(&tag)
	}
	return converted
}

func DBVersionTargetToGenerated(versionTarget *postgres.VersionTarget) *generated.VersionTarget {
	if versionTarget == nil {
		return nil
	}

	hash := versionTarget.Hash
	size := int(versionTarget.Size)

	return &generated.VersionTarget{
		VersionID:  versionTarget.VersionID,
		TargetName: generated.TargetName(versionTarget.TargetName),
		Hash:       &hash,
		Size:       &size,
	}
}

func DBVersionTargetsToGeneratedSlice(versionTargets []postgres.VersionTarget) []*generated.VersionTarget {
	converted := make([]*generated.VersionTarget, len(versionTargets))
	for i, versionTarget := range versionTargets {
		converted[i] = DBVersionTargetToGenerated(&versionTarget)
	}
	return converted
}

func GenCompInfoToDBCompInfo(gen *generated.CompatibilityInfoInput) *postgres.CompatibilityInfo {
	if gen == nil {
		return nil
	}
	return &postgres.CompatibilityInfo{
		Ea:  GenCompToDBComp(gen.Ea),
		Exp: GenCompToDBComp(gen.Exp),
	}
}

func GenCompToDBComp(gen *generated.CompatibilityInput) postgres.Compatibility {
	r := postgres.Compatibility{
		State: string(gen.State),
	}
	SetINN(gen.Note, &r.Note)
	return r
}

func DBCompInfoToGenCompInfo(gen *postgres.CompatibilityInfo) *generated.CompatibilityInfo {
	if gen == nil {
		return nil
	}
	return &generated.CompatibilityInfo{
		Ea:  DBCompToGenComp(gen.Ea),
		Exp: DBCompToGenComp(gen.Exp),
	}
}

func DBCompToGenComp(db postgres.Compatibility) *generated.Compatibility {
	return &generated.Compatibility{
		State: generated.CompatibilityState(db.State),
		Note:  &db.Note,
	}
}
