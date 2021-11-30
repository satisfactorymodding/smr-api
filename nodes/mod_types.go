package nodes

import (
	"time"

	"github.com/satisfactorymodding/smr-api/db/postgres"
)

type Mod struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	ShortDescription string    `json:"short_description"`
	FullDescription  string    `json:"full_description"`
	Logo             string    `json:"logo"`
	SourceURL        string    `json:"source_url"`
	CreatorID        string    `json:"creator_id"`
	Approved         bool      `json:"approved"`
	Views            uint      `json:"views"`
	Downloads        uint      `json:"downloads"`
	Hotness          uint      `json:"hotness"`
	Popularity       uint      `json:"popularity"`
	UpdatedAt        time.Time `json:"updated_at"`
	CreatedAt        time.Time `json:"created_at"`
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
	ID         string    `json:"id"`
	Version    string    `json:"version"`
	SMLVersion string    `json:"sml_version"`
	Changelog  string    `json:"changelog"`
	Downloads  uint      `json:"downloads"`
	Stability  string    `json:"stability"`
	ModID      string    `json:"mod_id"`
	Approved   bool      `json:"approved"`
	UpdatedAt  time.Time `json:"updated_at"`
	CreatedAt  time.Time `json:"created_at"`
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
	ID                  string    `json:"id"`
	Version             string    `json:"version"`
	SatisfactoryVersion int       `json:"satisfactory_version"`
	BootstrapVersion    *string   `json:"bootstrap_version"`
	Stability           string    `json:"stability"`
	Date                time.Time `json:"date"`
	Link                string    `json:"link"`
	Changelog           string    `json:"changelog"`
	UpdatedAt           time.Time `json:"updated_at"`
	CreatedAt           time.Time `json:"created_at"`
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
