package main

import "time"

type HabiticaData struct {
	Success bool           `json:"success"`
	Data    []HabiticaItem `json:"data"`
}

type HabiticaItem struct {
	Repeat    map[string]bool `json:"repeat"`
	Frequency string          `json:"frequency"`
	Type      string          `json:"type"`
	Notes     string          `json:"notes"`
	Checklist []interface{}   `json:"checklist"`
	Reminders []Remider       `json:"reminders,omitempty"`
	Text      string          `json:"text"`
	IsDue     bool            `json:"isDue"`
}

type Remider struct {
	Time time.Time `json:"time"`
}
