package structs

import (
	"time"
	"database/sql/driver"
)

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