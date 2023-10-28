package conversion

import (
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
	Convert(source *ent.SmlVersion) *generated.SMLVersion
	ConvertSlice(source []*ent.SmlVersion) []*generated.SMLVersion
}

// goverter:converter
// goverter:output:file ../generated/conv/user.go
// goverter:output:package conv
// goverter:extend TimeToString
type User interface {
	// goverter:ignore Roles Groups Mods Guides
	Convert(source *ent.User) *generated.User
}

func TimeToString(i time.Time) string {
	return i.Format(time.RFC3339)
}
