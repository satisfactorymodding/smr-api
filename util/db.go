package util

type Stability string

const (
	StabilityRelease = "release"
	StabilityBeta    = "beta"
	StabilityAlpha   = "alpha"
)

func (s Stability) Values() []string {
	return []string{
		StabilityRelease,
		StabilityBeta,
		StabilityAlpha,
	}
}

type CompatibilityInfo struct {
	Ea         Compatibility `gorm:"type:compatibility" json:"EA"`
	Exp        Compatibility `gorm:"type:compatibility" json:"EXP"`
	Controller Compatibility `gorm:"type:compatibility" json:"Controller"`
}

type Compatibility struct {
	State string
	Note  string
}

type AIUseDisclosureInfo struct {
	DisclosureType   string `gorm:"type:string" json:"DisclosureType"`
	DisclosureString string `gorm:"type:string" json:"DisclosureString"`
}
