package gql

import (
	"time"

	"github.com/satisfactorymodding/smr-api/db/postgres"
	"github.com/satisfactorymodding/smr-api/generated"
)

func DBUserToGenerated(user *postgres.User) *generated.User {
	if user == nil {
		return nil
	}

	Email := (*user).Email
	Avatar := (*user).Avatar

	result := &generated.User{
		ID:         (*user).ID,
		Username:   (*user).Username,
		Email:      &Email,
		Avatar:     &Avatar,
		CreatedAt:  user.CreatedAt.Format(time.RFC3339Nano),
		GithubID:   (*user).GithubID,
		GoogleID:   (*user).GoogleID,
		FacebookID: (*user).FacebookID,
	}

	return result
}

func DBModToGenerated(mod *postgres.Mod) *generated.Mod {
	if mod == nil {
		return nil
	}

	Logo := (*mod).Logo
	SourceURL := (*mod).SourceURL
	FullDescription := (*mod).FullDescription

	var LastVersionDate string
	if (*mod).LastVersionDate != nil {
		LastVersionDate = (*mod).LastVersionDate.Format(time.RFC3339Nano)
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
		Approved:   version.Approved,
		UpdatedAt:  version.UpdatedAt.Format(time.RFC3339Nano),
		CreatedAt:  version.CreatedAt.Format(time.RFC3339Nano),
		ModID:      version.ModID,
		Metadata:   version.Metadata,
		Hash:       version.Hash,
		Size:       &size,
	}
}

func DBGuideToGenerated(guide *postgres.Guide) *generated.Guide {
	if guide == nil {
		return nil
	}

	return &generated.Guide{
		ID:               guide.ID,
		Name:             guide.Name,
		ShortDescription: guide.ShortDescription,
		Guide:            guide.Guide,
		UserID:           guide.UserID,
		Views:            int(guide.Views),
		UpdatedAt:        guide.UpdatedAt.Format(time.RFC3339Nano),
		CreatedAt:        guide.CreatedAt.Format(time.RFC3339Nano),
	}
}

func DBSMLVersionToGenerated(smlVersion *postgres.SMLVersion) *generated.SMLVersion {
	if smlVersion == nil {
		return nil
	}

	return &generated.SMLVersion{
		ID:                  smlVersion.ID,
		Version:             smlVersion.Version,
		SatisfactoryVersion: smlVersion.SatisfactoryVersion,
		BootstrapVersion:    smlVersion.BootstrapVersion,
		Stability:           generated.VersionStabilities(smlVersion.Stability),
		Link:                smlVersion.Link,
		Changelog:           smlVersion.Changelog,
		Date:                smlVersion.Date.Format(time.RFC3339Nano),
		UpdatedAt:           smlVersion.UpdatedAt.Format(time.RFC3339Nano),
		CreatedAt:           smlVersion.CreatedAt.Format(time.RFC3339Nano),
	}
}

func DBBootstrapVersionToGenerated(bootstrapVersion *postgres.BootstrapVersion) *generated.BootstrapVersion {
	if bootstrapVersion == nil {
		return nil
	}

	return &generated.BootstrapVersion{
		ID:                  bootstrapVersion.ID,
		Version:             bootstrapVersion.Version,
		SatisfactoryVersion: bootstrapVersion.SatisfactoryVersion,
		Stability:           generated.VersionStabilities(bootstrapVersion.Stability),
		Link:                bootstrapVersion.Link,
		Changelog:           bootstrapVersion.Changelog,
		Date:                bootstrapVersion.Date.Format(time.RFC3339Nano),
		UpdatedAt:           bootstrapVersion.UpdatedAt.Format(time.RFC3339Nano),
		CreatedAt:           bootstrapVersion.CreatedAt.Format(time.RFC3339Nano),
	}
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

func DBAnnouncementToGenerated(announcement *postgres.Announcement) *generated.Announcement {
	if announcement == nil {
		return nil
	}

	return &generated.Announcement{
		ID:         announcement.ID,
		Message:    announcement.Message,
		Importance: generated.AnnouncementImportance(announcement.Importance),
	}
}

func DBAnnouncementsToGeneratedSlice(announcements []postgres.Announcement) []*generated.Announcement {
	converted := make([]*generated.Announcement, len(announcements))
	for i, announcement := range announcements {
		converted[i] = DBAnnouncementToGenerated(&announcement)
	}
	return converted
}

func DBModTagToGenerated(modTag *postgres.ModTag) *generated.ModTag {
	if modTag == nil {
		return nil
	}

	return &generated.ModTag{
		Name: modTag.Name,
	}
}

func DBModTagsToGeneratedSlice(modTags []postgres.ModTag) []*generated.ModTag {
	converted := make([]*generated.ModTag, len(modTags))
	for i, modTag := range modTags {
		converted[i] = DBModTagToGenerated(&modTag)
	}
	return converted
}
