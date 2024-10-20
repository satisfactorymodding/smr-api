package types

type ModAllVersionsVersion struct {
	ID               string                             `json:"id"`
	Version          string                             `json:"version"`
	GameVersion      string                             `json:"game_version"`
	RequiredOnRemote bool                               `json:"required_on_remote"`
	Targets          []*ModAllVersionsVersionTarget     `json:"targets"`
	Dependencies     []*ModAllVersionsVersionDependency `json:"dependencies"`
}

type ModAllVersionsVersionTarget struct {
	VersionID  string `json:"version_id"`
	TargetName string `json:"target_name"`
	Link       string `json:"link"`
	Hash       string `json:"hash"`
	Size       int    `json:"size"`
}

type ModAllVersionsVersionDependency struct {
	ModID        string `json:"mod_id"`
	ModReference string `json:"mod_reference"`
	Condition    string `json:"condition"`
	Optional     bool   `json:"optional"`
}
