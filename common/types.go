package common

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/lib/pq"
)

// JSONNullTime represents a time.Time that may be null.
// JSONNullTime implements the sql.Scanner interface so it can
// be used as a scan destination, similar to sql.NullString. It
// also implements the json.Marshaler and json.Unmarshaler json interfaces.
type JSONNullTime struct {
	pq.NullTime
}

// MarshalJSON satisifies the json.Marshaler interface. This
// allows the wrapped time field to be directly returned during
// encoding.
func (v JSONNullTime) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Time)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON satisifies the json.Unmarshaler interface. This
// allows the wrapped time field to be directly returned during
// decoding.
func (v *JSONNullTime) UnmarshalJSON(data []byte) error {
	// Unmarshaling into a pointer will let us detect null
	var x *time.Time
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		v.Valid = true
		v.Time = *x
	} else {
		v.Valid = false
	}
	return nil
}

var _ sql.Scanner = &JSONNullTime{}
var _ json.Marshaler = JSONNullTime{}
var _ json.Unmarshaler = &JSONNullTime{}
