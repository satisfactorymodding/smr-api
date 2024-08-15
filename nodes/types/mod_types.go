package types

type ModAllVersionsVersion struct {
	ID           string                             `json:"id"`
	Version      string                             `json:"version"`
	GameVersion  string                             `json:"game_version"`
	Targets      []*ModAllVersionsVersionTarget     `json:"targets"`
	Dependencies []*ModAllVersionsVersionDependency `json:"dependencies"`
}

type ModAllVersionsVersionTarget struct {
	VersionID  string `json:"version_id"`
	TargetName string `json:"target_name"`
	Link       string `json:"link"`
	Hash       string `json:"hash"`
	Size       int    `json:"size"`
}

type ModAllVersionsVersionDependency struct {
	ModID     string `json:"mod_id"`
	Condition string `json:"condition"`
	Optional  bool   `json:"optional"`
}
