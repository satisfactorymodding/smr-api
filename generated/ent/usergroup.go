// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/satisfactorymodding/smr-api/generated/ent/user"
	"github.com/satisfactorymodding/smr-api/generated/ent/usergroup"
)

// UserGroup is the model entity for the UserGroup schema.
type UserGroup struct {
	config `json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// DeletedAt holds the value of the "deleted_at" field.
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	// UserID holds the value of the "user_id" field.
	UserID string `json:"user_id,omitempty"`
	// GroupID holds the value of the "group_id" field.
	GroupID string `json:"group_id,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the UserGroupQuery when eager-loading is set.
	Edges        UserGroupEdges `json:"edges"`
	selectValues sql.SelectValues
}

// UserGroupEdges holds the relations/edges for other nodes in the graph.
type UserGroupEdges struct {
	// User holds the value of the user edge.
	User *User `json:"user,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// UserOrErr returns the User value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e UserGroupEdges) UserOrErr() (*User, error) {
	if e.User != nil {
		return e.User, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: user.Label}
	}
	return nil, &NotLoadedError{edge: "user"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*UserGroup) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case usergroup.FieldID, usergroup.FieldUserID, usergroup.FieldGroupID:
			values[i] = new(sql.NullString)
		case usergroup.FieldCreatedAt, usergroup.FieldUpdatedAt, usergroup.FieldDeletedAt:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the UserGroup fields.
func (ug *UserGroup) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case usergroup.FieldID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value.Valid {
				ug.ID = value.String
			}
		case usergroup.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				ug.CreatedAt = value.Time
			}
		case usergroup.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				ug.UpdatedAt = value.Time
			}
		case usergroup.FieldDeletedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field deleted_at", values[i])
			} else if value.Valid {
				ug.DeletedAt = value.Time
			}
		case usergroup.FieldUserID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field user_id", values[i])
			} else if value.Valid {
				ug.UserID = value.String
			}
		case usergroup.FieldGroupID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field group_id", values[i])
			} else if value.Valid {
				ug.GroupID = value.String
			}
		default:
			ug.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the UserGroup.
// This includes values selected through modifiers, order, etc.
func (ug *UserGroup) Value(name string) (ent.Value, error) {
	return ug.selectValues.Get(name)
}

// QueryUser queries the "user" edge of the UserGroup entity.
func (ug *UserGroup) QueryUser() *UserQuery {
	return NewUserGroupClient(ug.config).QueryUser(ug)
}

// Update returns a builder for updating this UserGroup.
// Note that you need to call UserGroup.Unwrap() before calling this method if this UserGroup
// was returned from a transaction, and the transaction was committed or rolled back.
func (ug *UserGroup) Update() *UserGroupUpdateOne {
	return NewUserGroupClient(ug.config).UpdateOne(ug)
}

// Unwrap unwraps the UserGroup entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (ug *UserGroup) Unwrap() *UserGroup {
	_tx, ok := ug.config.driver.(*txDriver)
	if !ok {
		panic("ent: UserGroup is not a transactional entity")
	}
	ug.config.driver = _tx.drv
	return ug
}

// String implements the fmt.Stringer.
func (ug *UserGroup) String() string {
	var builder strings.Builder
	builder.WriteString("UserGroup(")
	builder.WriteString(fmt.Sprintf("id=%v, ", ug.ID))
	builder.WriteString("created_at=")
	builder.WriteString(ug.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(ug.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("deleted_at=")
	builder.WriteString(ug.DeletedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("user_id=")
	builder.WriteString(ug.UserID)
	builder.WriteString(", ")
	builder.WriteString("group_id=")
	builder.WriteString(ug.GroupID)
	builder.WriteByte(')')
	return builder.String()
}

// UserGroups is a parsable slice of UserGroup.
type UserGroups []*UserGroup
