package models

import (
	"strings"
	"time"
)

const (
	DhanDateTimeFormat = "2006-01-02 15:04:05"
	DhanDateFormat     = "2006-01-02"
	DhanISOWithoutTZ   = "2006-01-02T15:04:05"
)

// BrokerTime is a custom timestamp wrapper to parse various DhanHQ timestamp formats
type BrokerTime struct {
	time.Time
}

// UnmarshalJSON unmarshals DhanHQ timestamps from JSON bytes
func (bt *BrokerTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" || s == "" {
		bt.Time = time.Time{}
		return nil
	}

	// 1. Try standard Dhan space-separated datetime: "2006-01-02 15:04:05"
	if t, err := time.Parse(DhanDateTimeFormat, s); err == nil {
		bt.Time = t
		return nil
	}

	// 2. Try RFC3339 / ISO8601: "2006-01-02T15:04:05+05:30"
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		bt.Time = t
		return nil
	}

	// 3. Try ISO without timezone: "2006-01-02T15:04:05"
	if t, err := time.Parse(DhanISOWithoutTZ, s); err == nil {
		bt.Time = t
		return nil
	}

	// 4. Try date-only format: "2006-01-02"
	if t, err := time.Parse(DhanDateFormat, s); err == nil {
		bt.Time = t
		return nil
	}

	// Fallback to RFC3339Nano
	if t, err := time.Parse(time.RFC3339Nano, s); err == nil {
		bt.Time = t
		return nil
	}

	return &time.ParseError{
		Layout:     DhanDateTimeFormat,
		Value:      s,
		LayoutElem: "2006",
		ValueElem:  s,
		Message:    ": cannot parse as Dhan timestamp",
	}
}

// MarshalJSON marshals BrokerTime into JSON string
func (bt BrokerTime) MarshalJSON() ([]byte, error) {
	if bt.IsZero() {
		return []byte("null"), nil
	}
	return []byte(`"` + bt.Format(DhanDateTimeFormat) + `"`), nil
}
