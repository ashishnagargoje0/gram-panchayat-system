package models

import "time"

// Property is a minimal model representing a tax/property record.
type Property struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	OwnerID    uint      `json:"owner_id"`
	Address    string    `json:"address,omitempty"`
	AssessedAt time.Time `json:"assessed_at,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
