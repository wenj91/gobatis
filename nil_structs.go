package gobatis

import (
	"database/sql"
	"database/sql/driver"
	"time"
)

type NullBool = sql.NullBool
type NullFloat64 = sql.NullFloat64
type NullInt64 = sql.NullInt64
type NullString = sql.NullString

// There is implementation of this in lib/pg
type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}
// Scan implements the Scanner interface.
func (nt *NullTime) Scan(value interface{}) error {
	nt.Time, nt.Valid = value.(time.Time)
	return nil
}
// Value implements the driver Valuer interface.
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}