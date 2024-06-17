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
	Ea  Compatibility `gorm:"type:compatibility" json:"EA"`
	Exp Compatibility `gorm:"type:compatibility" json:"EXP"`
}

type Compatibility struct {
	State string
	Note  string
}
