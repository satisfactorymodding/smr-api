// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/satisfactorymodding/smr-api/generated/ent/announcement"
)

// Announcement is the model entity for the Announcement schema.
type Announcement struct {
	config `json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// DeletedAt holds the value of the "deleted_at" field.
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	// Message holds the value of the "message" field.
	Message string `json:"message,omitempty"`
	// Importance holds the value of the "importance" field.
	Importance   string `json:"importance,omitempty"`
	selectValues sql.SelectValues
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Announcement) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case announcement.FieldID, announcement.FieldMessage, announcement.FieldImportance:
			values[i] = new(sql.NullString)
		case announcement.FieldCreatedAt, announcement.FieldUpdatedAt, announcement.FieldDeletedAt:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Announcement fields.
func (a *Announcement) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case announcement.FieldID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value.Valid {
				a.ID = value.String
			}
		case announcement.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				a.CreatedAt = value.Time
			}
		case announcement.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				a.UpdatedAt = value.Time
			}
		case announcement.FieldDeletedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field deleted_at", values[i])
			} else if value.Valid {
				a.DeletedAt = value.Time
			}
		case announcement.FieldMessage:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field message", values[i])
			} else if value.Valid {
				a.Message = value.String
			}
		case announcement.FieldImportance:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field importance", values[i])
			} else if value.Valid {
				a.Importance = value.String
			}
		default:
			a.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Announcement.
// This includes values selected through modifiers, order, etc.
func (a *Announcement) Value(name string) (ent.Value, error) {
	return a.selectValues.Get(name)
}

// Update returns a builder for updating this Announcement.
// Note that you need to call Announcement.Unwrap() before calling this method if this Announcement
// was returned from a transaction, and the transaction was committed or rolled back.
func (a *Announcement) Update() *AnnouncementUpdateOne {
	return NewAnnouncementClient(a.config).UpdateOne(a)
}

// Unwrap unwraps the Announcement entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (a *Announcement) Unwrap() *Announcement {
	_tx, ok := a.config.driver.(*txDriver)
	if !ok {
		panic("ent: Announcement is not a transactional entity")
	}
	a.config.driver = _tx.drv
	return a
}

// String implements the fmt.Stringer.
func (a *Announcement) String() string {
	var builder strings.Builder
	builder.WriteString("Announcement(")
	builder.WriteString(fmt.Sprintf("id=%v, ", a.ID))
	builder.WriteString("created_at=")
	builder.WriteString(a.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(a.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("deleted_at=")
	builder.WriteString(a.DeletedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("message=")
	builder.WriteString(a.Message)
	builder.WriteString(", ")
	builder.WriteString("importance=")
	builder.WriteString(a.Importance)
	builder.WriteByte(')')
	return builder.String()
}

// Announcements is a parsable slice of Announcement.
type Announcements []*Announcement
