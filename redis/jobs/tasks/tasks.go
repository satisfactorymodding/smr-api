package tasks

import "github.com/vmihailenco/taskq/v3"

var (
	UpdateDBFromModVersionFileTask     *taskq.Task
	UpdateDBFromModVersionJSONFileTask *taskq.Task
	CopyObjectFromOldBucketTask        *taskq.Task
	ScanModOnVirusTotalTask            *taskq.Task
)

type UpdateDBFromModVersionFileData struct {
	ModID     string `json:"mod_id"`
	VersionID string `json:"version_id"`
}

type UpdateDBFromModVersionJSONFileData struct {
	ModID     string `json:"mod_id"`
	VersionID string `json:"version_id"`
}

type CopyObjectFromOldBucketData struct {
	Key string `json:"key"`
}

type CopyObjectToOldBucketData struct {
	Key string `json:"key"`
}

type ScanModOnVirusTotalData struct {
	ModID        string `json:"mod_id"`
	VersionID    string `json:"version_id"`
	ApproveAfter bool   `json:"approve_after"`
}
