package conversion

import (
	"strconv"
	"strings"
	"time"

	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/generated/ent"
)

// goverter:converter
// goverter:output:file ../generated/conv/announcement.go
// goverter:output:package conv
// goverter:extend TimeToString
type Announcement interface {
	Convert(source *ent.Announcement) *generated.Announcement
	ConvertSlice(source []*ent.Announcement) []*generated.Announcement
}

// goverter:converter
// goverter:output:file ../generated/conv/sml_version.go
// goverter:output:package conv
// goverter:extend TimeToString
type SMLVersion interface {
	// goverter:map Edges.Targets Targets
	// goverter:map CreatedAt Date
	// goverter:map GameVersion SatisfactoryVersion | SMLSatisfactoryVersion
	// goverter:ignore EngineVersion BootstrapVersion Link
	Convert(source *ent.Version) *generated.SMLVersion
	ConvertSlice(source []*ent.Version) []*generated.SMLVersion
	// goverter:ignore Link
	ConvertTarget(source ent.VersionTarget) generated.SMLVersionTarget
}

// goverter:converter
// goverter:output:file ../generated/conv/satisfactory_version.go
// goverter:output:package conv
// goverter:extend TimeToString
type SatisfactoryVersion interface {
	Convert(source *ent.SatisfactoryVersion) *generated.SatisfactoryVersion
	ConvertSlice(source []*ent.SatisfactoryVersion) []*generated.SatisfactoryVersion
}

// goverter:converter
// goverter:output:file ../generated/conv/user.go
// goverter:output:package conv
// goverter:extend TimeToString
type User interface {
	// goverter:ignore Roles Groups Mods Guides
	Convert(source *ent.User) *generated.User
	ConvertSlice(source []*ent.User) []*generated.User
}

// goverter:converter
// goverter:output:file ../generated/conv/guide.go
// goverter:output:package conv
// goverter:extend TimeToString
type Guide interface {
	// goverter:ignore User
	// goverter:map Edges.Tags Tags
	Convert(source *ent.Guide) *generated.Guide
	ConvertSlice(source []*ent.Guide) []*generated.Guide
}

// goverter:converter
// goverter:output:file ../generated/conv/tag.go
// goverter:output:package conv
// goverter:extend TimeToString
type Tag interface {
	Convert(source *ent.Tag) *generated.Tag
	ConvertSlice(source []*ent.Tag) []*generated.Tag
}

// goverter:converter
// goverter:output:file ../generated/conv/user_mod.go
// goverter:output:package conv
// goverter:extend TimeToString
type UserMod interface {
	// goverter:ignore User Mod
	Convert(source *ent.UserMod) *generated.UserMod
	ConvertSlice(source []*ent.UserMod) []*generated.UserMod
}

// goverter:converter
// goverter:output:file ../generated/conv/mod.go
// goverter:output:package conv
// goverter:extend TimeToString UIntToInt
type Mod interface {
	// goverter:map Edges.Tags Tags
	// goverter:ignore Authors Version Versions LatestVersions
	Convert(source *ent.Mod) *generated.Mod
	ConvertSlice(source []*ent.Mod) []*generated.Mod
}

// goverter:converter
// goverter:output:file ../generated/conv/version.go
// goverter:output:package conv
// goverter:extend TimeToString UIntToInt Int64ToInt
type Version interface {
	// goverter:map Edges.Targets Targets
	// goverter:ignore Link Mod Dependencies
	Convert(source *ent.Version) *generated.Version
	ConvertSlice(source []*ent.Version) []*generated.Version

	// goverter:ignore Link
	ConvertTarget(source *ent.VersionTarget) *generated.VersionTarget
}

// goverter:converter
// goverter:output:file ../generated/conv/version_dependency.go
// goverter:output:package conv
// goverter:extend TimeToString UIntToInt Int64ToInt
type VersionDependency interface {
	// goverter:ignore Mod Version
	Convert(source *ent.VersionDependency) *generated.VersionDependency
	ConvertSlice(source []*ent.VersionDependency) []*generated.VersionDependency
}

func TimeToString(i time.Time) string {
	return i.Format(time.RFC3339)
}

func UIntToInt(i uint) int {
	return int(i)
}

func Int64ToInt(i int64) int {
	return int(i)
}

func SMLSatisfactoryVersion(gameVersion string) int {
	if !strings.HasPrefix(gameVersion, ">=") {
		// the semver library doesn't have a constraint.MinVersion(),
		// so we only handle >= for now, which is what will be used 99% of the time
		return 0
	}
	res, err := strconv.Atoi(gameVersion[2:]) // Strip the >=
	if err != nil {
		return 0
	}
	return res
}
