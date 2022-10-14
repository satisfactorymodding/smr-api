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
	Arch       []ModArch `gorm:"foreignKey:ModVersionID;preload:true"`
	Hotness    uint
	Downloads  uint
	Denied     bool `gorm:"default:false;not null"`
	Approved   bool `gorm:"default:false;not null"`
}

type Guide struct {
	SMRModel
	Name             string `gorm:"type:varchar(50)"`
	ShortDescription string `gorm:"type:varchar(128)"`
	Guide            string
	UserID           string
	Tags             []Tag `gorm:"many2many:guide_tags"`
	User             User
	Views            uint
}

type UserGroup struct {
	SMRDates

	UserID  string `gorm:"primary_key"`
	GroupID string `gorm:"primary_key"`
}

type SMLVersion struct {
	Date             time.Time
	BootstrapVersion *string
	SMRModel
	Version             string `gorm:"type:varchar(32);unique_index"`
	Stability           string `sql:"type:version_stability"`
	Link                string
	Changelog           string
	Arch                []SMLArch `gorm:"foreignKey:SMLVersionID;preload:true"`
	SatisfactoryVersion int
}

type VersionDependency struct {
	SMRDates

	VersionID string `gorm:"primary_key;type:varchar(14)"`
	ModID     string `gorm:"primary_key;type:varchar(14)"`

	Condition string `gorm:"type:varchar(64)"`
	Optional  bool
}

type BootstrapVersion struct {
	Date time.Time
	SMRModel
	Version             string `gorm:"type:varchar(32);unique_index"`
	Stability           string `sql:"type:version_stability"`
	Link                string
	Changelog           string
	SatisfactoryVersion int
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

type ModArch struct {
	ID           string `gorm:"primary_key;type:varchar(16)"`
	ModVersionID string `gorm:"column:mod_version_arch_id"`
	Platform     string
	Key          string
	Hash         string
	Size         int64
}

func (ModArch) TableName() string {
	return "mod_archs"
}

type SMLArch struct {
	ID           string `gorm:"primary_key;type:varchar(14)"`
	SMLVersionID string `gorm:"column:sml_version_arch_id"`
	Platform     string
	Link         string
}

func (SMLArch) TableName() string {
	return "sml_archs"
}
