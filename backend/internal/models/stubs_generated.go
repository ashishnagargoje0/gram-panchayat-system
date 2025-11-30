package models

import "time"

// Meeting & related minimal models (used to satisfy compile-time refs)
type Meeting struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	MeetingType string    `json:"meeting_type,omitempty"`
	ScheduledAt time.Time `json:"scheduled_at,omitempty"`
	Location    string    `json:"location,omitempty"`
	Agenda      string    `json:"agenda,omitempty"`
	Status      string    `json:"status,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type MeetingMinutes struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	MeetingID  uint      `json:"meeting_id"`
	Attendees  string    `json:"attendees,omitempty"`
	Decisions  string    `json:"decisions,omitempty"`
	RecordedBy uint      `json:"recorded_by,omitempty"` // changed to uint to match service usage
	Content    string    `json:"content,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}
