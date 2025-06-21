package util

type CompatibilityInfo struct {
	Ea  Compatibility `gorm:"type:compatibility" json:"EA"`
	Exp Compatibility `gorm:"type:compatibility" json:"EXP"`
}

type Compatibility struct {
	State string
	Note  string
}
