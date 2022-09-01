package postgres

import (
	"time"

	"gorm.io/gorm"
)

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
	SMRModel

	Email      string `gorm:"type:varchar(256);unique_index"`
	Username   string `gorm:"type:varchar(32)"`
	Avatar     string
	JoinedFrom string
	Banned     bool `gorm:"default:false;not null"`

	GithubID   *string
	GoogleID   *string
	FacebookID *string

	Mods []Mod `gorm:"many2many:user_mods;"`
}

type UserSession struct {
	SMRModel

	UserID string
	User   User

	Token     string `gorm:"type:varchar(256);unique_index"`
	UserAgent string
}

type Mod struct {
	SMRModel

	Name             string `gorm:"type:varchar(32)"`
	ShortDescription string `gorm:"type:varchar(128)"`
	FullDescription  string
	Logo             string
	SourceURL        string
	CreatorID        string
	Approved         bool `gorm:"default:false;not null"`
	Denied           bool `gorm:"default:false;not null"`
	Views            uint
	Downloads        uint
	Hotness          uint
	Popularity       uint
	LastVersionDate  *time.Time
	ModReference     string
	Hidden           bool
	Compatibility    *CompatibilityInfo `gorm:"serializer:json"`

	Users []User `gorm:"many2many:user_mods;"`

	Tags []Tag `gorm:"many2many:mod_tags"`

	Versions []Version
}

type UserMod struct {
	UserID string `gorm:"primary_key"`
	ModID  string `gorm:"primary_key"`
	Role   string
}

// If updated, update dataloader
type Version struct {
	SMRModel

	ModID string

	Version      string `gorm:"type:varchar(16)"`
	SMLVersion   string `gorm:"type:varchar(16)"`
	Changelog    string
	Downloads    uint
	Key          string
	Stability    string `gorm:"default:'alpha'" sql:"type:version_stability"`
	Approved     bool   `gorm:"default:false;not null"`
	Denied       bool   `gorm:"default:false;not null"`
	Hotness      uint
	Metadata     *string
	ModReference *string
	VersionMajor *int
	VersionMinor *int
	VersionPatch *int
	Size         *int64
	Hash         *string
}

type Guide struct {
	SMRModel

	Name             string `gorm:"type:varchar(50)"`
	ShortDescription string `gorm:"type:varchar(128)"`
	Guide            string
	Views            uint
	Tags             []Tag `gorm:"many2many:guide_tags"`

	UserID string
	User   User
}

type UserGroup struct {
	SMRDates

	UserID  string `gorm:"primary_key"`
	GroupID string `gorm:"primary_key"`
}

type SMLVersion struct {
	SMRModel

	Version             string `gorm:"type:varchar(32);unique_index"`
	SatisfactoryVersion int
	Stability           string `sql:"type:version_stability"`
	Date                time.Time
	Link                string
	Changelog           string
	BootstrapVersion    *string
}

type VersionDependency struct {
	SMRDates

	VersionID string `gorm:"primary_key;type:varchar(14)"`
	ModID     string `gorm:"primary_key;type:varchar(14)"`

	Condition string `gorm:"type:varchar(64)"`
	Optional  bool
}

type BootstrapVersion struct {
	SMRModel

	Version             string `gorm:"type:varchar(32);unique_index"`
	SatisfactoryVersion int
	Stability           string `sql:"type:version_stability"`
	Date                time.Time
	Link                string
	Changelog           string
}

type Announcement struct {
	SMRModel

	Message    string
	Importance string
}

type Tag struct {
	SMRModel

	Name string `gorm:"type:varchar(24)"`

	Mods []Mod `gorm:"many2many:mod_tags"`
}

type ModTag struct {
	TagID string `gorm:"primary_key;type:varchar(24)"`
	ModID string `gorm:"primary_key;type:varchar(16)"`
}

type GuideTag struct {
	TagID   string `gorm:"primary_key;type:varchar(24)"`
	GuideID string `gorm:"primary_key;type:varchar(16)"`
}

type CompatibilityInfo struct {
	EA  Compatibility `gorm:"type:compatibility"`
	EXP Compatibility `gorm:"type:compatibility"`
}

type Compatibility struct {
	State string
	Note  string
}
