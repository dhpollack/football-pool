package apiespn

import (
	"encoding/json"
	"strings"
	"time"
)

// ESPNDateTime is a custom time type that handles ESPN's time format
// ESPN returns times like "2025-09-05T00:20Z" but Go expects "2025-09-05T00:20:00Z"
type ESPNDateTime struct {
	time.Time
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (e *ESPNDateTime) UnmarshalJSON(data []byte) error {
	// Remove quotes from the JSON string
	str := strings.Trim(string(data), `"`)

	// Try parsing with standard format first
	t, err := time.Parse(time.RFC3339, str)
	if err == nil {
		e.Time = t
		return nil
	}

	// If that fails, try parsing ESPN format (without seconds)
	t, err = time.Parse("2006-01-02T15:04Z", str)
	if err != nil {
		return err
	}

	e.Time = t
	return nil
}

// MarshalJSON implements the json.Marshaler interface
func (e ESPNDateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Time.Format(time.RFC3339))
}

// String returns the time in RFC3339 format
func (e ESPNDateTime) String() string {
	return e.Time.Format(time.RFC3339)
}
