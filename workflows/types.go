package workflows

import (
	"github.com/satisfactorymodding/smr-api/workflows/removemod"
	"github.com/satisfactorymodding/smr-api/workflows/statistics"
	"github.com/satisfactorymodding/smr-api/workflows/updatemodfromstorage"
	"github.com/satisfactorymodding/smr-api/workflows/versionupload"
)

var Workflows = struct {
	Statistics           *statistics.A
	UpdateModFromStorage *updatemodfromstorage.A
	VersionUpload        *versionupload.A
	RemoveMod            *removemod.A
}{
	Statistics:           statistics.Statistics,
	UpdateModFromStorage: updatemodfromstorage.UpdateModFromStorage,
	VersionUpload:        versionupload.VersionUpload,
	RemoveMod:            removemod.RemoveMod,
}
