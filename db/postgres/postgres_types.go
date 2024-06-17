package postgres

import (
	"time"

	"gorm.io/gorm"
)

type Tabler interface {
	TableName() string
}

type SMRDates struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type SMRModel struct {
	ID string `gorm:"primary_key;type:varchar(14)"`
	SMRDates
}

type User struct {
	GithubID   *string
	GoogleID   *string
	FacebookID *string
	SMRModel
	Email      string `gorm:"type:varchar(256);unique_index"`
	Username   string `gorm:"type:varchar(32)"`
	Avatar     string
	JoinedFrom string
	Mods       []Mod `gorm:"many2many:user_mods;"`
	Banned     bool  `gorm:"default:false;not null"`
}

type UserSession struct {
	SMRModel
	UserID    string
	Token     string `gorm:"type:varchar(256);unique_index"`
	UserAgent string
	User      User
}

type Mod struct {
	LastVersionDate *time.Time
	Compatibility   *CompatibilityInfo `gorm:"serializer:json"`
	SMRModel
	CreatorID        string
	Logo             string
	SourceURL        string
	FullDescription  string
	ShortDescription string `gorm:"type:varchar(128)"`
	Name             string `gorm:"type:varchar(32)"`
	ModReference     string
	Versions         []Version
	Tags             []Tag  `gorm:"many2many:mod_tags"`
	Users            []User `gorm:"many2many:user_mods;"`
	Downloads        uint
	Popularity       uint
	Hotness          uint
	Views            uint
	Hidden           bool
	Denied           bool `gorm:"default:false;not null"`
	Approved         bool `gorm:"default:false;not null"`
}

type UserMod struct {
	UserID string `gorm:"primary_key"`
	ModID  string `gorm:"primary_key"`
	Role   string
}

// If updated, update dataloader
type Version struct {
	Metadata     *string
	Hash         *string
	Size         *int64
	VersionPatch *int
	VersionMinor *int
	VersionMajor *int
	ModReference *string
	SMRModel
	Changelog  string
	Stability  string `gorm:"default:'alpha'" sql:"type:version_stability"`
	Key        string
	SMLVersion string `gorm:"type:varchar(16)"`
	Version    string `gorm:"type:varchar(16)"`
	ModID      string
	Targets    []VersionTarget `gorm:"foreignKey:VersionID"`
	Hotness    uint
	Downloads  uint
	Denied     bool `gorm:"default:false;not null"`
	Approved   bool `gorm:"default:false;not null"`
}

type TinyVersion struct {
	Hash *string
	Size *int64
	SMRModel
	SMLVersion   string              `gorm:"type:varchar(16)"`
	Version      string              `gorm:"type:varchar(16)"`
	Targets      []VersionTarget     `gorm:"foreignKey:VersionID;preload:true"`
	Dependencies []VersionDependency `gorm:"foreignKey:VersionID"`
}

func (TinyVersion) TableName() string {
	return "versions"
}

type SMLVersion struct {
	Date             time.Time
	BootstrapVersion *string
	SMRModel
	Version             string `gorm:"type:varchar(32);unique_index"`
	Stability           string `sql:"type:version_stability"`
	Link                string
	Changelog           string
	EngineVersion       string
	Targets             []SMLVersionTarget `gorm:"foreignKey:VersionID"`
	SatisfactoryVersion int
}

type VersionDependency struct {
	SMRDates

	VersionID string `gorm:"primary_key;type:varchar(14)"`
	ModID     string `gorm:"primary_key;type:varchar(14)"`

	Condition string `gorm:"type:varchar(64)"`
	Optional  bool
}

type Tag struct {
	SMRModel

	Name        string `gorm:"type:varchar(24)"`
	Description string `gorm:"type:varchar(512)"`

	Mods []Mod `gorm:"many2many:mod_tags"`
}

type CompatibilityInfo struct {
	Ea  Compatibility `gorm:"type:compatibility" json:"EA"`
	Exp Compatibility `gorm:"type:compatibility" json:"EXP"`
}

type Compatibility struct {
	State string
	Note  string
}

type VersionTarget struct {
	VersionID  string `gorm:"primary_key;type:varchar(14)"`
	TargetName string `gorm:"primary_key;type:varchar(16)"`
	Key        string
	Hash       string
	Size       int64
}

type SMLVersionTarget struct {
	VersionID  string `gorm:"primary_key;type:varchar(14)"`
	TargetName string `gorm:"primary_key;type:varchar(16)"`
	Link       string
}
