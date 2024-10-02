// Code generated by ent, DO NOT EDIT.

package version

import (
	"fmt"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/satisfactorymodding/smr-api/util"
)

const (
	// Label holds the string label denoting the version type in the database.
	Label = "version"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// FieldDeletedAt holds the string denoting the deleted_at field in the database.
	FieldDeletedAt = "deleted_at"
	// FieldModID holds the string denoting the mod_id field in the database.
	FieldModID = "mod_id"
	// FieldVersion holds the string denoting the version field in the database.
	FieldVersion = "version"
	// FieldGameVersion holds the string denoting the game_version field in the database.
	FieldGameVersion = "game_version"
	// FieldRequiredOnRemote holds the string denoting the required_on_remote field in the database.
	FieldRequiredOnRemote = "required_on_remote"
	// FieldChangelog holds the string denoting the changelog field in the database.
	FieldChangelog = "changelog"
	// FieldDownloads holds the string denoting the downloads field in the database.
	FieldDownloads = "downloads"
	// FieldKey holds the string denoting the key field in the database.
	FieldKey = "key"
	// FieldStability holds the string denoting the stability field in the database.
	FieldStability = "stability"
	// FieldApproved holds the string denoting the approved field in the database.
	FieldApproved = "approved"
	// FieldHotness holds the string denoting the hotness field in the database.
	FieldHotness = "hotness"
	// FieldDenied holds the string denoting the denied field in the database.
	FieldDenied = "denied"
	// FieldMetadata holds the string denoting the metadata field in the database.
	FieldMetadata = "metadata"
	// FieldModReference holds the string denoting the mod_reference field in the database.
	FieldModReference = "mod_reference"
	// FieldVersionMajor holds the string denoting the version_major field in the database.
	FieldVersionMajor = "version_major"
	// FieldVersionMinor holds the string denoting the version_minor field in the database.
	FieldVersionMinor = "version_minor"
	// FieldVersionPatch holds the string denoting the version_patch field in the database.
	FieldVersionPatch = "version_patch"
	// FieldSize holds the string denoting the size field in the database.
	FieldSize = "size"
	// FieldHash holds the string denoting the hash field in the database.
	FieldHash = "hash"
	// EdgeMod holds the string denoting the mod edge name in mutations.
	EdgeMod = "mod"
	// EdgeDependencies holds the string denoting the dependencies edge name in mutations.
	EdgeDependencies = "dependencies"
	// EdgeTargets holds the string denoting the targets edge name in mutations.
	EdgeTargets = "targets"
	// EdgeVersionDependencies holds the string denoting the version_dependencies edge name in mutations.
	EdgeVersionDependencies = "version_dependencies"
	// Table holds the table name of the version in the database.
	Table = "versions"
	// ModTable is the table that holds the mod relation/edge.
	ModTable = "versions"
	// ModInverseTable is the table name for the Mod entity.
	// It exists in this package in order to avoid circular dependency with the "mod" package.
	ModInverseTable = "mods"
	// ModColumn is the table column denoting the mod relation/edge.
	ModColumn = "mod_id"
	// DependenciesTable is the table that holds the dependencies relation/edge. The primary key declared below.
	DependenciesTable = "version_dependencies"
	// DependenciesInverseTable is the table name for the Mod entity.
	// It exists in this package in order to avoid circular dependency with the "mod" package.
	DependenciesInverseTable = "mods"
	// TargetsTable is the table that holds the targets relation/edge.
	TargetsTable = "version_targets"
	// TargetsInverseTable is the table name for the VersionTarget entity.
	// It exists in this package in order to avoid circular dependency with the "versiontarget" package.
	TargetsInverseTable = "version_targets"
	// TargetsColumn is the table column denoting the targets relation/edge.
	TargetsColumn = "version_id"
	// VersionDependenciesTable is the table that holds the version_dependencies relation/edge.
	VersionDependenciesTable = "version_dependencies"
	// VersionDependenciesInverseTable is the table name for the VersionDependency entity.
	// It exists in this package in order to avoid circular dependency with the "versiondependency" package.
	VersionDependenciesInverseTable = "version_dependencies"
	// VersionDependenciesColumn is the table column denoting the version_dependencies relation/edge.
	VersionDependenciesColumn = "version_id"
)

// Columns holds all SQL columns for version fields.
var Columns = []string{
	FieldID,
	FieldCreatedAt,
	FieldUpdatedAt,
	FieldDeletedAt,
	FieldModID,
	FieldVersion,
	FieldGameVersion,
	FieldRequiredOnRemote,
	FieldChangelog,
	FieldDownloads,
	FieldKey,
	FieldStability,
	FieldApproved,
	FieldHotness,
	FieldDenied,
	FieldMetadata,
	FieldModReference,
	FieldVersionMajor,
	FieldVersionMinor,
	FieldVersionPatch,
	FieldSize,
	FieldHash,
}

var (
	// DependenciesPrimaryKey and DependenciesColumn2 are the table columns denoting the
	// primary key for the dependencies relation (M2M).
	DependenciesPrimaryKey = []string{"version_id", "mod_id"}
)

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

// Note that the variables below are initialized by the runtime
// package on the initialization of the application. Therefore,
// it should be imported in the main as follows:
//
//	import _ "github.com/satisfactorymodding/smr-api/generated/ent/runtime"
var (
	Hooks        [1]ent.Hook
	Interceptors [1]ent.Interceptor
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
	// DefaultUpdatedAt holds the default value on creation for the "updated_at" field.
	DefaultUpdatedAt func() time.Time
	// UpdateDefaultUpdatedAt holds the default value on update for the "updated_at" field.
	UpdateDefaultUpdatedAt func() time.Time
	// VersionValidator is a validator for the "version" field. It is called by the builders before save.
	VersionValidator func(string) error
	// DefaultDownloads holds the default value on creation for the "downloads" field.
	DefaultDownloads uint
	// DefaultApproved holds the default value on creation for the "approved" field.
	DefaultApproved bool
	// DefaultHotness holds the default value on creation for the "hotness" field.
	DefaultHotness uint
	// DefaultDenied holds the default value on creation for the "denied" field.
	DefaultDenied bool
	// ModReferenceValidator is a validator for the "mod_reference" field. It is called by the builders before save.
	ModReferenceValidator func(string) error
	// HashValidator is a validator for the "hash" field. It is called by the builders before save.
	HashValidator func(string) error
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() string
)

const DefaultStability util.Stability = "release"

// StabilityValidator is a validator for the "stability" field enum values. It is called by the builders before save.
func StabilityValidator(s util.Stability) error {
	switch s {
	case "release", "beta", "alpha":
		return nil
	default:
		return fmt.Errorf("version: invalid enum value for stability field: %q", s)
	}
}

// OrderOption defines the ordering options for the Version queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByCreatedAt orders the results by the created_at field.
func ByCreatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreatedAt, opts...).ToFunc()
}

// ByUpdatedAt orders the results by the updated_at field.
func ByUpdatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUpdatedAt, opts...).ToFunc()
}

// ByDeletedAt orders the results by the deleted_at field.
func ByDeletedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDeletedAt, opts...).ToFunc()
}

// ByModID orders the results by the mod_id field.
func ByModID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldModID, opts...).ToFunc()
}

// ByVersion orders the results by the version field.
func ByVersion(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldVersion, opts...).ToFunc()
}

// ByGameVersion orders the results by the game_version field.
func ByGameVersion(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldGameVersion, opts...).ToFunc()
}

// ByRequiredOnRemote orders the results by the required_on_remote field.
func ByRequiredOnRemote(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldRequiredOnRemote, opts...).ToFunc()
}

// ByChangelog orders the results by the changelog field.
func ByChangelog(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldChangelog, opts...).ToFunc()
}

// ByDownloads orders the results by the downloads field.
func ByDownloads(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDownloads, opts...).ToFunc()
}

// ByKey orders the results by the key field.
func ByKey(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldKey, opts...).ToFunc()
}

// ByStability orders the results by the stability field.
func ByStability(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldStability, opts...).ToFunc()
}

// ByApproved orders the results by the approved field.
func ByApproved(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldApproved, opts...).ToFunc()
}

// ByHotness orders the results by the hotness field.
func ByHotness(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldHotness, opts...).ToFunc()
}

// ByDenied orders the results by the denied field.
func ByDenied(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDenied, opts...).ToFunc()
}

// ByMetadata orders the results by the metadata field.
func ByMetadata(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldMetadata, opts...).ToFunc()
}

// ByModReference orders the results by the mod_reference field.
func ByModReference(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldModReference, opts...).ToFunc()
}

// ByVersionMajor orders the results by the version_major field.
func ByVersionMajor(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldVersionMajor, opts...).ToFunc()
}

// ByVersionMinor orders the results by the version_minor field.
func ByVersionMinor(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldVersionMinor, opts...).ToFunc()
}

// ByVersionPatch orders the results by the version_patch field.
func ByVersionPatch(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldVersionPatch, opts...).ToFunc()
}

// BySize orders the results by the size field.
func BySize(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSize, opts...).ToFunc()
}

// ByHash orders the results by the hash field.
func ByHash(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldHash, opts...).ToFunc()
}

// ByModField orders the results by mod field.
func ByModField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newModStep(), sql.OrderByField(field, opts...))
	}
}

// ByDependenciesCount orders the results by dependencies count.
func ByDependenciesCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newDependenciesStep(), opts...)
	}
}

// ByDependencies orders the results by dependencies terms.
func ByDependencies(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newDependenciesStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByTargetsCount orders the results by targets count.
func ByTargetsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newTargetsStep(), opts...)
	}
}

// ByTargets orders the results by targets terms.
func ByTargets(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newTargetsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByVersionDependenciesCount orders the results by version_dependencies count.
func ByVersionDependenciesCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newVersionDependenciesStep(), opts...)
	}
}

// ByVersionDependencies orders the results by version_dependencies terms.
func ByVersionDependencies(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newVersionDependenciesStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newModStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ModInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, ModTable, ModColumn),
	)
}
func newDependenciesStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(DependenciesInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2M, false, DependenciesTable, DependenciesPrimaryKey...),
	)
}
func newTargetsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(TargetsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, TargetsTable, TargetsColumn),
	)
}
func newVersionDependenciesStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(VersionDependenciesInverseTable, VersionDependenciesColumn),
		sqlgraph.Edge(sqlgraph.O2M, true, VersionDependenciesTable, VersionDependenciesColumn),
	)
}
