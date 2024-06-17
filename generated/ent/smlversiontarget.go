// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/satisfactorymodding/smr-api/generated/ent/smlversion"
	"github.com/satisfactorymodding/smr-api/generated/ent/smlversiontarget"
)

// SmlVersionTarget is the model entity for the SmlVersionTarget schema.
type SmlVersionTarget struct {
	config `json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// VersionID holds the value of the "version_id" field.
	VersionID string `json:"version_id,omitempty"`
	// TargetName holds the value of the "target_name" field.
	TargetName string `json:"target_name,omitempty"`
	// Link holds the value of the "link" field.
	Link string `json:"link,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the SmlVersionTargetQuery when eager-loading is set.
	Edges        SmlVersionTargetEdges `json:"edges"`
	selectValues sql.SelectValues
}

// SmlVersionTargetEdges holds the relations/edges for other nodes in the graph.
type SmlVersionTargetEdges struct {
	// SmlVersion holds the value of the sml_version edge.
	SmlVersion *SmlVersion `json:"sml_version,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// SmlVersionOrErr returns the SmlVersion value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e SmlVersionTargetEdges) SmlVersionOrErr() (*SmlVersion, error) {
	if e.SmlVersion != nil {
		return e.SmlVersion, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: smlversion.Label}
	}
	return nil, &NotLoadedError{edge: "sml_version"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*SmlVersionTarget) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case smlversiontarget.FieldID, smlversiontarget.FieldVersionID, smlversiontarget.FieldTargetName, smlversiontarget.FieldLink:
			values[i] = new(sql.NullString)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the SmlVersionTarget fields.
func (svt *SmlVersionTarget) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case smlversiontarget.FieldID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value.Valid {
				svt.ID = value.String
			}
		case smlversiontarget.FieldVersionID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field version_id", values[i])
			} else if value.Valid {
				svt.VersionID = value.String
			}
		case smlversiontarget.FieldTargetName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field target_name", values[i])
			} else if value.Valid {
				svt.TargetName = value.String
			}
		case smlversiontarget.FieldLink:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field link", values[i])
			} else if value.Valid {
				svt.Link = value.String
			}
		default:
			svt.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the SmlVersionTarget.
// This includes values selected through modifiers, order, etc.
func (svt *SmlVersionTarget) Value(name string) (ent.Value, error) {
	return svt.selectValues.Get(name)
}

// QuerySmlVersion queries the "sml_version" edge of the SmlVersionTarget entity.
func (svt *SmlVersionTarget) QuerySmlVersion() *SmlVersionQuery {
	return NewSmlVersionTargetClient(svt.config).QuerySmlVersion(svt)
}

// Update returns a builder for updating this SmlVersionTarget.
// Note that you need to call SmlVersionTarget.Unwrap() before calling this method if this SmlVersionTarget
// was returned from a transaction, and the transaction was committed or rolled back.
func (svt *SmlVersionTarget) Update() *SmlVersionTargetUpdateOne {
	return NewSmlVersionTargetClient(svt.config).UpdateOne(svt)
}

// Unwrap unwraps the SmlVersionTarget entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (svt *SmlVersionTarget) Unwrap() *SmlVersionTarget {
	_tx, ok := svt.config.driver.(*txDriver)
	if !ok {
		panic("ent: SmlVersionTarget is not a transactional entity")
	}
	svt.config.driver = _tx.drv
	return svt
}

// String implements the fmt.Stringer.
func (svt *SmlVersionTarget) String() string {
	var builder strings.Builder
	builder.WriteString("SmlVersionTarget(")
	builder.WriteString(fmt.Sprintf("id=%v, ", svt.ID))
	builder.WriteString("version_id=")
	builder.WriteString(svt.VersionID)
	builder.WriteString(", ")
	builder.WriteString("target_name=")
	builder.WriteString(svt.TargetName)
	builder.WriteString(", ")
	builder.WriteString("link=")
	builder.WriteString(svt.Link)
	builder.WriteByte(')')
	return builder.String()
}

// SmlVersionTargets is a parsable slice of SmlVersionTarget.
type SmlVersionTargets []*SmlVersionTarget
