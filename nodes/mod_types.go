package nodes

import (
	"time"

	"github.com/satisfactorymodding/smr-api/db/postgres"
)

type Mod struct {
	UpdatedAt        time.Time `json:"updated_at"`
	CreatedAt        time.Time `json:"created_at"`
	CreatorID        string    `json:"creator_id"`
	FullDescription  string    `json:"full_description"`
	Logo             string    `json:"logo"`
	SourceURL        string    `json:"source_url"`
	ID               string    `json:"id"`
	ShortDescription string    `json:"short_description"`
	Name             string    `json:"name"`
	Views            uint      `json:"views"`
	Downloads        uint      `json:"downloads"`
	Hotness          uint      `json:"hotness"`
	Popularity       uint      `json:"popularity"`
	Approved         bool      `json:"approved"`
}

func ModToMod(mod *postgres.Mod, short bool) *Mod {
	result := Mod{
		ID:               mod.ID,
		Name:             mod.Name,
		ShortDescription: mod.ShortDescription,
		Logo:             mod.Logo,
		SourceURL:        mod.SourceURL,
		CreatorID:        mod.CreatorID,
		Approved:         mod.Approved,
		Views:            mod.Views,
		Downloads:        mod.Downloads,
		Hotness:          mod.Hotness,
		Popularity:       mod.Popularity,
		UpdatedAt:        mod.UpdatedAt,
		CreatedAt:        mod.CreatedAt,
	}

	if !short {
		result.FullDescription = mod.FullDescription
	}

	return &result
}

type Version struct {
	UpdatedAt    time.Time           `json:"updated_at,omitempty"`
	CreatedAt    time.Time           `json:"created_at,omitempty"`
	ID           string              `json:"id,omitempty"`
	Version      string              `json:"version,omitempty"`
	SMLVersion   string              `json:"sml_version,omitempty"`
	Changelog    string              `json:"changelog,omitempty"`
	Stability    string              `json:"stability,omitempty"`
	ModID        string              `json:"mod_id,omitempty"`
	Downloads    uint                `json:"downloads,omitempty"`
	Approved     bool                `json:"approved,omitempty"`
	Dependencies []VersionDependency `json:"dependencies,omitempty"`
	Arch         []VersionTarget     `json:"arch,omitempty"`
}

type VersionDependency struct {
	ModID     string `json:"mod_id"`
	Condition string `json:"condition"`
	Optional  bool   `json:"optional"`
}

type VersionTarget struct {
	VersionID  string `json:"version_id"`
	TargetName string `json:"target_name"`
	Key        string `json:"key"`
	Hash       string `json:"hash"`
	Size       int64  `json:"size"`
}

func TinyVersionToVersion(version *postgres.TinyVersion) *Version {
	var dependencies []VersionDependency
	if version.Dependencies != nil {
		dependencies = make([]VersionDependency, len(version.Dependencies))
		for i, v := range version.Dependencies {
			dependencies[i] = VersionDependencyToVersionDependency(v)
		}
	}

	var archs []VersionTarget
	if version.Arch != nil {
		archs = make([]VersionTarget, len(version.Arch))
		for i, v := range version.Arch {
			archs[i] = VersionArchToVersionArch(v)
		}
	}

	return &Version{
		UpdatedAt:    version.UpdatedAt,
		CreatedAt:    version.CreatedAt,
		ID:           version.ID,
		Version:      version.Version,
		SMLVersion:   version.SMLVersion,
		Dependencies: dependencies,
		Arch:         archs,
	}
}

func VersionToVersion(version *postgres.Version) *Version {
	return &Version{
		ID:         version.ID,
		Version:    version.Version,
		SMLVersion: version.SMLVersion,
		Changelog:  version.Changelog,
		Downloads:  version.Downloads,
		Stability:  version.Stability,
		Approved:   version.Approved,
		UpdatedAt:  version.UpdatedAt,
		CreatedAt:  version.CreatedAt,
		ModID:      version.ModID,
	}
}

func VersionDependencyToVersionDependency(version postgres.VersionDependency) VersionDependency {
	return VersionDependency{
		ModID:     version.ModID,
		Condition: version.Condition,
		Optional:  version.Optional,
	}
}

func VersionArchToVersionArch(version postgres.VersionTarget) VersionTarget {
	return VersionTarget{
		VersionID:  version.VersionID,
		TargetName: version.TargetName,
		Key:        version.Key,
		Hash:       version.Hash,
		Size:       version.Size,
	}
}

type ModUser struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

func ModUserToModUser(userMod *postgres.UserMod) *ModUser {
	return &ModUser{
		UserID: userMod.UserID,
		Role:   userMod.Role,
	}
}

type SMLVersion struct {
	Date                time.Time `json:"date"`
	UpdatedAt           time.Time `json:"updated_at"`
	CreatedAt           time.Time `json:"created_at"`
	BootstrapVersion    *string   `json:"bootstrap_version"`
	ID                  string    `json:"id"`
	Version             string    `json:"version"`
	Stability           string    `json:"stability"`
	Link                string    `json:"link"`
	Changelog           string    `json:"changelog"`
	SatisfactoryVersion int       `json:"satisfactory_version"`
}

func SMLVersionToSMLVersion(version *postgres.SMLVersion) *SMLVersion {
	return &SMLVersion{
		ID:                  version.ID,
		Version:             version.Version,
		SatisfactoryVersion: version.SatisfactoryVersion,
		BootstrapVersion:    version.BootstrapVersion,
		Stability:           version.Stability,
		Date:                version.Date,
		Link:                version.Link,
		Changelog:           version.Changelog,
		UpdatedAt:           version.UpdatedAt,
		CreatedAt:           version.CreatedAt,
	}
}
