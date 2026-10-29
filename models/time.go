package models

import (
	"strconv"
	"strings"
	"time"
)

const (
	DhanDateTimeFormat = "2006-01-02 15:04:05"
	DhanDateFormat     = "2006-01-02"
	DhanISOWithoutTZ   = "2006-01-02T15:04:05"
)

var brokerTimeLayouts = []string{
	// ISO / RFC3339 formats with timezone
	time.RFC3339Nano,
	time.RFC3339,

	// ISO without timezone (with / without fractional seconds)
	"2006-01-02T15:04:05.999999999",
	"2006-01-02T15:04:05",

	// Space-separated formats (with timezone)
	"2006-01-02 15:04:05.999999999-07:00",
	"2006-01-02 15:04:05.999999999Z07:00",
	"2006-01-02 15:04:05-07:00",
	"2006-01-02 15:04:05Z07:00",

	// Space-separated formats (without timezone)
	"2006-01-02 15:04:05.999999999",
	"2006-01-02 15:04:05",

	// Slash-separated formats
	"2006/01/02 15:04:05.999999999",
	"2006/01/02 15:04:05",
	"2006/01/02",

	// Indian DD-MM-YYYY formats
	"02-01-2006 15:04:05.999999999",
	"02-01-2006 15:04:05",
	"02-01-2006",
	"02/01/2006 15:04:05.999999999",
	"02/01/2006 15:04:05",
	"02/01/2006",

	// Date only
	"2006-01-02",
}

// BrokerTime is a custom timestamp wrapper to parse various DhanHQ timestamp formats
type BrokerTime struct {
	time.Time
}

// UnmarshalJSON unmarshals DhanHQ timestamps from JSON bytes
func (bt *BrokerTime) UnmarshalJSON(b []byte) error {
	s := strings.TrimSpace(strings.Trim(strings.TrimSpace(string(b)), "\""))
	if s == "null" || s == "" {
		bt.Time = time.Time{}
		return nil
	}

	// 1. Try numeric epoch timestamp (seconds or milliseconds)
	if epoch, err := strconv.ParseInt(s, 10, 64); err == nil {
		if epoch > 1e11 {
			bt.Time = time.UnixMilli(epoch)
		} else {
			bt.Time = time.Unix(epoch, 0)
		}
		return nil
	}

	// 2. Try candidate layouts
	for _, layout := range brokerTimeLayouts {
		if t, err := time.Parse(layout, s); err == nil {
			bt.Time = t
			return nil
		}
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
