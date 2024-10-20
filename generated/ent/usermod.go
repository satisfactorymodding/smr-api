// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/satisfactorymodding/smr-api/generated/ent/mod"
	"github.com/satisfactorymodding/smr-api/generated/ent/user"
	"github.com/satisfactorymodding/smr-api/generated/ent/usermod"
)

// UserMod is the model entity for the UserMod schema.
type UserMod struct {
	config `json:"-"`
	// UserID holds the value of the "user_id" field.
	UserID string `json:"user_id,omitempty"`
	// ModID holds the value of the "mod_id" field.
	ModID string `json:"mod_id,omitempty"`
	// Role holds the value of the "role" field.
	Role string `json:"role,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the UserModQuery when eager-loading is set.
	Edges        UserModEdges `json:"edges"`
	selectValues sql.SelectValues
}

// UserModEdges holds the relations/edges for other nodes in the graph.
type UserModEdges struct {
	// User holds the value of the user edge.
	User *User `json:"user,omitempty"`
	// Mod holds the value of the mod edge.
	Mod *Mod `json:"mod,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// UserOrErr returns the User value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e UserModEdges) UserOrErr() (*User, error) {
	if e.User != nil {
		return e.User, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: user.Label}
	}
	return nil, &NotLoadedError{edge: "user"}
}

// ModOrErr returns the Mod value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e UserModEdges) ModOrErr() (*Mod, error) {
	if e.Mod != nil {
		return e.Mod, nil
	} else if e.loadedTypes[1] {
		return nil, &NotFoundError{label: mod.Label}
	}
	return nil, &NotLoadedError{edge: "mod"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*UserMod) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case usermod.FieldUserID, usermod.FieldModID, usermod.FieldRole:
			values[i] = new(sql.NullString)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the UserMod fields.
func (um *UserMod) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case usermod.FieldUserID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field user_id", values[i])
			} else if value.Valid {
				um.UserID = value.String
			}
		case usermod.FieldModID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field mod_id", values[i])
			} else if value.Valid {
				um.ModID = value.String
			}
		case usermod.FieldRole:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field role", values[i])
			} else if value.Valid {
				um.Role = value.String
			}
		default:
			um.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the UserMod.
// This includes values selected through modifiers, order, etc.
func (um *UserMod) Value(name string) (ent.Value, error) {
	return um.selectValues.Get(name)
}

// QueryUser queries the "user" edge of the UserMod entity.
func (um *UserMod) QueryUser() *UserQuery {
	return NewUserModClient(um.config).QueryUser(um)
}

// QueryMod queries the "mod" edge of the UserMod entity.
func (um *UserMod) QueryMod() *ModQuery {
	return NewUserModClient(um.config).QueryMod(um)
}

// Update returns a builder for updating this UserMod.
// Note that you need to call UserMod.Unwrap() before calling this method if this UserMod
// was returned from a transaction, and the transaction was committed or rolled back.
func (um *UserMod) Update() *UserModUpdateOne {
	return NewUserModClient(um.config).UpdateOne(um)
}

// Unwrap unwraps the UserMod entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (um *UserMod) Unwrap() *UserMod {
	_tx, ok := um.config.driver.(*txDriver)
	if !ok {
		panic("ent: UserMod is not a transactional entity")
	}
	um.config.driver = _tx.drv
	return um
}

// String implements the fmt.Stringer.
func (um *UserMod) String() string {
	var builder strings.Builder
	builder.WriteString("UserMod(")
	builder.WriteString("user_id=")
	builder.WriteString(um.UserID)
	builder.WriteString(", ")
	builder.WriteString("mod_id=")
	builder.WriteString(um.ModID)
	builder.WriteString(", ")
	builder.WriteString("role=")
	builder.WriteString(um.Role)
	builder.WriteByte(')')
	return builder.String()
}

// UserMods is a parsable slice of UserMod.
type UserMods []*UserMod
