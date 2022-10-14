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
	UpdatedAt  time.Time `json:"updated_at"`
	CreatedAt  time.Time `json:"created_at"`
	ID         string    `json:"id"`
	Version    string    `json:"version"`
	SMLVersion string    `json:"sml_version"`
	Changelog  string    `json:"changelog"`
	Stability  string    `json:"stability"`
	ModID      string    `json:"mod_id"`
	Downloads  uint      `json:"downloads"`
	Approved   bool      `json:"approved"`
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
