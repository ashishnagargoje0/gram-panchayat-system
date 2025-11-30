package models

import "time"

// Complaint is a minimal model used by user relations.
// Add fields as your original schema requires.
type Complaint struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	Subject   string    `json:"subject,omitempty"`
	Body      string    `json:"body,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
