package models_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/DragoHex/godhan/models"
)

type sampleTimePayload struct {
	Timestamp models.BrokerTime `json:"timestamp"`
}

func TestBrokerTime_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name       string
		jsonInput  string
		wantYear   int
		wantMonth  time.Month
		wantDay    int
		wantHour   int
		wantMinute int
		wantZero   bool
		wantErr    bool
	}{
		{
			name:       "Standard SQL-like Datetime",
			jsonInput:  `{"timestamp":"2023-10-15 14:30:45"}`,
			wantYear:   2023,
			wantMonth:  time.October,
			wantDay:    15,
			wantHour:   14,
			wantMinute: 30,
		},
		{
			name:       "ISO 8601 Datetime with Timezone",
			jsonInput:  `{"timestamp":"2023-10-15T14:30:45+05:30"}`,
			wantYear:   2023,
			wantMonth:  time.October,
			wantDay:    15,
		},
		{
			name:      "Date Only",
			jsonInput: `{"timestamp":"2023-10-15"}`,
			wantYear:  2023,
			wantMonth: time.October,
			wantDay:   15,
		},
		{
			name:      "Null Value",
			jsonInput: `{"timestamp":null}`,
			wantZero:  true,
		},
		{
			name:      "Empty String",
			jsonInput: `{"timestamp":""}`,
			wantZero:  true,
		},
		{
			name:      "Invalid Corrupt Date",
			jsonInput: `{"timestamp":"invalid-date-format"}`,
			wantErr:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var payload sampleTimePayload
			err := json.Unmarshal([]byte(tc.jsonInput), &payload)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected unmarshal error: %v", err)
			}
			if tc.wantZero {
				if !payload.Timestamp.IsZero() {
					t.Errorf("expected zero time, got %v", payload.Timestamp)
				}
				return
			}
			if payload.Timestamp.Year() != tc.wantYear || payload.Timestamp.Month() != tc.wantMonth || payload.Timestamp.Day() != tc.wantDay {
				t.Errorf("date mismatch: got %v", payload.Timestamp)
			}
		})
	}
}
