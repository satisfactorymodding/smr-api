// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/satisfactorymodding/smr-api/generated/ent/satisfactoryversion"
)

// SatisfactoryVersion is the model entity for the SatisfactoryVersion schema.
type SatisfactoryVersion struct {
	config `json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// DeletedAt holds the value of the "deleted_at" field.
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	// Version holds the value of the "version" field.
	Version int `json:"version,omitempty"`
	// EngineVersion holds the value of the "engine_version" field.
	EngineVersion string `json:"engine_version,omitempty"`
	selectValues  sql.SelectValues
}

// scanValues returns the types for scanning values from sql.Rows.
func (*SatisfactoryVersion) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case satisfactoryversion.FieldVersion:
			values[i] = new(sql.NullInt64)
		case satisfactoryversion.FieldID, satisfactoryversion.FieldEngineVersion:
			values[i] = new(sql.NullString)
		case satisfactoryversion.FieldCreatedAt, satisfactoryversion.FieldUpdatedAt, satisfactoryversion.FieldDeletedAt:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the SatisfactoryVersion fields.
func (sv *SatisfactoryVersion) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case satisfactoryversion.FieldID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value.Valid {
				sv.ID = value.String
			}
		case satisfactoryversion.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				sv.CreatedAt = value.Time
			}
		case satisfactoryversion.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				sv.UpdatedAt = value.Time
			}
		case satisfactoryversion.FieldDeletedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field deleted_at", values[i])
			} else if value.Valid {
				sv.DeletedAt = value.Time
			}
		case satisfactoryversion.FieldVersion:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field version", values[i])
			} else if value.Valid {
				sv.Version = int(value.Int64)
			}
		case satisfactoryversion.FieldEngineVersion:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field engine_version", values[i])
			} else if value.Valid {
				sv.EngineVersion = value.String
			}
		default:
			sv.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the SatisfactoryVersion.
// This includes values selected through modifiers, order, etc.
func (sv *SatisfactoryVersion) Value(name string) (ent.Value, error) {
	return sv.selectValues.Get(name)
}

// Update returns a builder for updating this SatisfactoryVersion.
// Note that you need to call SatisfactoryVersion.Unwrap() before calling this method if this SatisfactoryVersion
// was returned from a transaction, and the transaction was committed or rolled back.
func (sv *SatisfactoryVersion) Update() *SatisfactoryVersionUpdateOne {
	return NewSatisfactoryVersionClient(sv.config).UpdateOne(sv)
}

// Unwrap unwraps the SatisfactoryVersion entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (sv *SatisfactoryVersion) Unwrap() *SatisfactoryVersion {
	_tx, ok := sv.config.driver.(*txDriver)
	if !ok {
		panic("ent: SatisfactoryVersion is not a transactional entity")
	}
	sv.config.driver = _tx.drv
	return sv
}

// String implements the fmt.Stringer.
func (sv *SatisfactoryVersion) String() string {
	var builder strings.Builder
	builder.WriteString("SatisfactoryVersion(")
	builder.WriteString(fmt.Sprintf("id=%v, ", sv.ID))
	builder.WriteString("created_at=")
	builder.WriteString(sv.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(sv.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("deleted_at=")
	builder.WriteString(sv.DeletedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("version=")
	builder.WriteString(fmt.Sprintf("%v", sv.Version))
	builder.WriteString(", ")
	builder.WriteString("engine_version=")
	builder.WriteString(sv.EngineVersion)
	builder.WriteByte(')')
	return builder.String()
}

// SatisfactoryVersions is a parsable slice of SatisfactoryVersion.
type SatisfactoryVersions []*SatisfactoryVersion
